# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: model/api/forecast.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='model/api/forecast.proto',
  package='api',
  syntax='proto3',
  serialized_options=b'Z\tmodel/api',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x18model/api/forecast.proto\x12\x03\x61pi\"C\n\x0eProphetRequest\x12\x1f\n\x06values\x18\x01 \x03(\x0b\x32\x0f.api.SamplePair\x12\x10\n\x08\x64uration\x18\x02 \x01(\x01\".\n\nSamplePair\x12\x11\n\ttimestamp\x18\x01 \x01(\x03\x12\r\n\x05value\x18\x02 \x01(\x01\"/\n\x0cProphetReply\x12\x1f\n\x06values\x18\x01 \x03(\x0b\x32\x0f.api.Forecasted\"S\n\nForecasted\x12\x11\n\ttimestamp\x18\x01 \x01(\x01\x12\x0c\n\x04yhat\x18\x02 \x01(\x01\x12\x11\n\tyhatLower\x18\x03 \x01(\x01\x12\x11\n\tyhatUpper\x18\x04 \x01(\x01\x32=\n\x08\x46orecast\x12\x31\n\x07Prophet\x12\x13.api.ProphetRequest\x1a\x11.api.ProphetReplyB\x0bZ\tmodel/apib\x06proto3'
)




_PROPHETREQUEST = _descriptor.Descriptor(
  name='ProphetRequest',
  full_name='api.ProphetRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='values', full_name='api.ProphetRequest.values', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='duration', full_name='api.ProphetRequest.duration', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=33,
  serialized_end=100,
)


_SAMPLEPAIR = _descriptor.Descriptor(
  name='SamplePair',
  full_name='api.SamplePair',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='api.SamplePair.timestamp', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='api.SamplePair.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=102,
  serialized_end=148,
)


_PROPHETREPLY = _descriptor.Descriptor(
  name='ProphetReply',
  full_name='api.ProphetReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='values', full_name='api.ProphetReply.values', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=150,
  serialized_end=197,
)


_FORECASTED = _descriptor.Descriptor(
  name='Forecasted',
  full_name='api.Forecasted',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='api.Forecasted.timestamp', index=0,
      number=1, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='yhat', full_name='api.Forecasted.yhat', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='yhatLower', full_name='api.Forecasted.yhatLower', index=2,
      number=3, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='yhatUpper', full_name='api.Forecasted.yhatUpper', index=3,
      number=4, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=199,
  serialized_end=282,
)

_PROPHETREQUEST.fields_by_name['values'].message_type = _SAMPLEPAIR
_PROPHETREPLY.fields_by_name['values'].message_type = _FORECASTED
DESCRIPTOR.message_types_by_name['ProphetRequest'] = _PROPHETREQUEST
DESCRIPTOR.message_types_by_name['SamplePair'] = _SAMPLEPAIR
DESCRIPTOR.message_types_by_name['ProphetReply'] = _PROPHETREPLY
DESCRIPTOR.message_types_by_name['Forecasted'] = _FORECASTED
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ProphetRequest = _reflection.GeneratedProtocolMessageType('ProphetRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROPHETREQUEST,
  '__module__' : 'model.api.forecast_pb2'
  # @@protoc_insertion_point(class_scope:api.ProphetRequest)
  })
_sym_db.RegisterMessage(ProphetRequest)

SamplePair = _reflection.GeneratedProtocolMessageType('SamplePair', (_message.Message,), {
  'DESCRIPTOR' : _SAMPLEPAIR,
  '__module__' : 'model.api.forecast_pb2'
  # @@protoc_insertion_point(class_scope:api.SamplePair)
  })
_sym_db.RegisterMessage(SamplePair)

ProphetReply = _reflection.GeneratedProtocolMessageType('ProphetReply', (_message.Message,), {
  'DESCRIPTOR' : _PROPHETREPLY,
  '__module__' : 'model.api.forecast_pb2'
  # @@protoc_insertion_point(class_scope:api.ProphetReply)
  })
_sym_db.RegisterMessage(ProphetReply)

Forecasted = _reflection.GeneratedProtocolMessageType('Forecasted', (_message.Message,), {
  'DESCRIPTOR' : _FORECASTED,
  '__module__' : 'model.api.forecast_pb2'
  # @@protoc_insertion_point(class_scope:api.Forecasted)
  })
_sym_db.RegisterMessage(Forecasted)


DESCRIPTOR._options = None

_FORECAST = _descriptor.ServiceDescriptor(
  name='Forecast',
  full_name='api.Forecast',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=284,
  serialized_end=345,
  methods=[
  _descriptor.MethodDescriptor(
    name='Prophet',
    full_name='api.Forecast.Prophet',
    index=0,
    containing_service=None,
    input_type=_PROPHETREQUEST,
    output_type=_PROPHETREPLY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_FORECAST)

DESCRIPTOR.services_by_name['Forecast'] = _FORECAST

# @@protoc_insertion_point(module_scope)