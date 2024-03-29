# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from model.api import forecast_pb2 as model_dot_api_dot_forecast__pb2


class ForecastStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Prophet = channel.unary_unary(
                '/api.Forecast/Prophet',
                request_serializer=model_dot_api_dot_forecast__pb2.ProphetRequest.SerializeToString,
                response_deserializer=model_dot_api_dot_forecast__pb2.ProphetReply.FromString,
                )


class ForecastServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Prophet(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ForecastServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Prophet': grpc.unary_unary_rpc_method_handler(
                    servicer.Prophet,
                    request_deserializer=model_dot_api_dot_forecast__pb2.ProphetRequest.FromString,
                    response_serializer=model_dot_api_dot_forecast__pb2.ProphetReply.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'api.Forecast', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Forecast(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Prophet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/api.Forecast/Prophet',
            model_dot_api_dot_forecast__pb2.ProphetRequest.SerializeToString,
            model_dot_api_dot_forecast__pb2.ProphetReply.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
