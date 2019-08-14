# from model.train import TrainModel

import grpc
import api.metrics_pb2
import api.metrics_pb2_grpc


def main():
    try:
        address = "localhost:3003"
        channel = grpc.insecure_channel(address)
        stub = api.metrics_pb2_grpc.MetricsStub(channel)
        req = api.metrics_pb2.QueryMetricsRequest()
        resp = stub.Query(req)
        print(resp)
    except Exception as e:
        print(e)
        return e


if __name__ == "__main__":
    main()
