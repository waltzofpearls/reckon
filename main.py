# from model.train import TrainModel

import api.metrics_pb2
import api.metrics_pb2_grpc
import grpc
import os


def main():
    try:
        root_ca = os.environ.get("TLS_ROOT_CA")
        client_key = os.environ.get("TLS_CLIENT_KEY")
        client_cert = os.environ.get("TLS_CLIENT_CERT")
        server_address = os.environ.get("GRPC_SERVER_ADDRESS")

        credentials = grpc.ssl_channel_credentials(
            root_ca.encode(), client_key.encode(), client_cert.encode())
        channel = grpc.secure_channel(server_address, credentials)
        stub = api.metrics_pb2_grpc.MetricsStub(channel)

        req = api.metrics_pb2.QueryMetricsRequest()
        resp = stub.Query(req)
        print(resp)
    except Exception as e:
        print(e)
        return e


if __name__ == "__main__":
    main()
