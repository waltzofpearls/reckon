from concurrent import futures
from forecaster.prophet import Forecaster as ProphetForecaster
from pythonjsonlogger import jsonlogger
import grpc
import logging
import model.api.forecast_pb2_grpc as grpc_lib
import os
import signal
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
    received = False
    def __init__(self, logger):
        self.logger = logger
        signal.signal(signal.SIGINT, self.sigint_received)
        signal.signal(signal.SIGTERM, self.sigterm_received)
        signal.signal(signal.SIGHUP, self.sighup_received)

    def sigint_received(self, *args):
        self.logger.info('SIGINT received')
        self.received = True

    def sigterm_received(self, *args):
        self.logger.info('SIGTERM received')
        self.received = True

    def sighup_received(self, *args):
        self.logger.info('SIGHUP received')
        self.received = True

class Config(object):
    def __init__(self):
        self.grpc_server_address = os.getenv('GRPC_SERVER_ADDRESS', '')
        self.grpc_server_key = str.encode(os.getenv('GRPC_SERVER_KEY', ''))
        self.grpc_server_cert = str.encode(os.getenv('GRPC_SERVER_CERT', ''))
        self.grpc_root_ca = str.encode(os.getenv('GRPC_ROOT_CA', ''))

class Server(object):
    def __init__(self, config, logger):
        self.config = config
        self.logger = logger

    def serve(self):
        server_credentials = grpc.ssl_server_credentials(
            [(self.config.grpc_server_key, self.config.grpc_server_cert)],
            root_certificates=self.config.grpc_root_ca,
            require_client_auth=True
        )
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        grpc_lib.add_ForecastServicer_to_server(ForecastServicer(self.logger), server)
        server.add_secure_port(self.config.grpc_server_address, server_credentials)
        self.logger.info('starting python gRPC server...')
        server.start()

        # wait for shutdown signals
        shutdown = GracefulShutdown(self.logger)
        while not shutdown.received:
            time.sleep(1)

        self.logger.info('stopping python gRPC server...')
        # wait for 5 seconds and then stop grpc server
        server.stop(5).wait()
        self.logger.info('python gRPC server gracefully stopped')

def json_logger():
    logger = logging.getLogger()
    log_handler = logging.StreamHandler(sys.stdout)
    formatter = jsonlogger.JsonFormatter(fmt='%(asctime)s %(name)s %(levelname)s %(message)s')
    log_handler.setFormatter(formatter)
    log_handler.flush = sys.stdout.flush
    logger.setLevel(logging.INFO)
    logger.addHandler(log_handler)
    return logger
