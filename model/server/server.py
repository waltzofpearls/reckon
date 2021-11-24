from concurrent import futures
from forecaster.prophet import Forecaster as ProphetForecaster
from multiprocessing import Event, Process, cpu_count
from pythonjsonlogger import jsonlogger
import contextlib
import grpc
import logging
import model.api.forecast_pb2_grpc as grpc_lib
import os
import signal
import socket
import sys
import time

class ForecastServicer(ProphetForecaster):
    def __init__(self, logger):
        self.logger = logger

    def pretty_timedelta(self, seconds):
        seconds = int(seconds)
        days, seconds = divmod(seconds, 86400)
        hours, seconds = divmod(seconds, 3600)
        minutes, seconds = divmod(seconds, 60)
        if days > 0:
            return '{:d}d{:d}h{:d}m{:d}s'.format(days, hours, minutes, seconds)
        elif hours > 0:
            return '{:d}h{:d}m{:d}s'.format(hours, minutes, seconds)
        elif minutes > 0:
            return '{:d}m{:d}s'.format(minutes, seconds)
        else:
            return '{:d}s'.format(seconds)

class GracefulShutdown:
    def __init__(self, logger):
        self.logger = logger
        self.event = Event()
        signal.signal(signal.SIGINT, self.handler('SIGINT'))
        signal.signal(signal.SIGTERM, self.handler('SIGTERM'))
        signal.signal(signal.SIGHUP, self.handler('SIGHUP'))

    def handler(self, signal_name):
        def fn(signal_received, frame):
            self.logger.info('signal received', extra={'signal': signal_name})
            self.event.set()
        return fn

class Config(object):
    def __init__(self):
        self.grpc_server_address = os.getenv('GRPC_SERVER_ADDRESS', '')
        self.grpc_server_key = str.encode(os.getenv('GRPC_SERVER_KEY', ''))
        self.grpc_server_cert = str.encode(os.getenv('GRPC_SERVER_CERT', ''))
        self.grpc_root_ca = str.encode(os.getenv('GRPC_ROOT_CA', ''))
        self.gprc_server_process_num = int(os.getenv('GPRC_SERVER_PROCESS_NUM', cpu_count()))
        self.grpc_server_thread_num = int(os.getenv('GRPC_SERVER_THREAD_NUM', 1))
        self.grpc_server_grace_period_in_secs = int(os.getenv('GRPC_SERVER_GRACE_PERIOD_IN_SECS', 2))
        self.grpc_server_kill_period_in_secs = int(os.getenv('GRPC_SERVER_KILL_PERIOD_IN_SECS', 5))

class Server(object):
    def __init__(self, config, logger):
        self.config = config
        self.logger = logger

    @contextlib.contextmanager
    def _reserve_port(self):
        """Find and reserve a port for all subprocesses to use"""
        sock = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
        if sock.getsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT) == 0:
            raise RuntimeError('failed to set SO_REUSEPORT.')
        _, port = self.config.grpc_server_address.split(':')
        sock.bind(('', int(port)))
        try:
            yield sock.getsockname()[1]
        finally:
            sock.close()

    def _run_server(self, shutdown_event):
        server_credentials = grpc.ssl_server_credentials(
            [(self.config.grpc_server_key, self.config.grpc_server_cert)],
            root_certificates=self.config.grpc_root_ca,
            require_client_auth=True
        )
        server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=self.config.grpc_server_thread_num),
            options=[
                ("grpc.so_reuseport", 1),
                ("grpc.use_local_subchannel_pool", 1),
            ],
        )
        grpc_lib.add_ForecastServicer_to_server(ForecastServicer(self.logger), server)
        server.add_secure_port(self.config.grpc_server_address, server_credentials)
        self.logger.info('starting python gRPC server...')
        server.start()

        while not shutdown_event.is_set():
            time.sleep(1)

        server.stop(5).wait()
        self.logger.info('python gRPC server stopped')

    def serve(self):
        with self._reserve_port():
            procs = []
            shutdown = GracefulShutdown(self.logger)
            for _ in range(self.config.gprc_server_process_num):
                proc = Process(target=self._run_server, args=(shutdown.event,))
                procs.append(proc)
                proc.start()
            while not shutdown.event.is_set():
                time.sleep(1)

            t = time.time()
            grace_period = self.config.grpc_server_grace_period_in_secs
            kill_period = self.config.grpc_server_kill_period_in_secs
            while True:
                # Send SIGINT if process doesn't exit quickly enough, and kill it as last resort
                # .is_alive() also implicitly joins the process (good practice in linux)
                alive_procs = [proc for proc in procs if proc.is_alive()]
                if len(alive_procs) == 0:
                    break
                elapsed = time.time() - t
                if elapsed >= grace_period and elapsed < kill_period:
                    for proc in alive_procs:
                        proc.terminate()
                        self.logger.info("sending SIGTERM to subprocess", extra={'proc': proc})
                elif elapsed >= kill_period:
                    for proc in alive_procs:
                        self.logger.warning("sending SIGKILL to subprocess", extra={'proc': proc})
                        # Queues and other inter-process communication primitives can break when
                        # process is killed, but we don't care here
                        proc.kill()
                time.sleep(1)

            time.sleep(1)
            for proc in procs:
                self.logger.info("subprocess terminated", extra={'proc': proc})

def json_logger():
    logger = logging.getLogger()
    log_handler = logging.StreamHandler(sys.stdout)
    formatter = jsonlogger.JsonFormatter(fmt='%(asctime)s %(name)s %(levelname)s %(message)s')
    log_handler.setFormatter(formatter)
    log_handler.flush = sys.stdout.flush
    logger.setLevel(logging.INFO)
    logger.addHandler(log_handler)
    return logger
