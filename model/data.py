import api.metrics_pb2
import api.metrics_pb2_grpc
import grpc
import os
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime, timedelta


class Data:
    def __init__(self):
        self.root_ca = os.environ.get('TLS_ROOT_CA', '').encode()
        self.tls_key = os.environ.get('TLS_CLIENT_KEY', '').encode()
        self.tls_cert = os.environ.get('TLS_CLIENT_CERT', '').encode()
        self.address = os.environ.get('GRPC_SERVER_ADDRESS', '')

    def fetch(self):
        try:
            credentials = grpc.ssl_channel_credentials(
                self.root_ca, self.tls_key, self.tls_cert)
            channel = grpc.secure_channel(self.address, credentials)
            stub = api.metrics_pb2_grpc.MetricsStub(channel)

            end = datetime.now()
            end_seconds = int(end.strftime('%s'))
            end_nanos = 0

            start = end - timedelta(hours=24)
            start_seconds = int(start.strftime('%s'))
            start_nanos = 0

            req = api.metrics_pb2.QueryMetricsRequest(
                metricName='steps',
                startTime=Timestamp(seconds=start_seconds, nanos=start_nanos),
                endTime=Timestamp(seconds=end_seconds, nanos=end_nanos),
            )
            resp = stub.Query(req)
            print(resp)
            return resp
        except Exception as e:
            print(e)
            return e

    def deliver(self, data):
        pass
