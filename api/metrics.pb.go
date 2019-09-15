// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/metrics.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type QueryMetricsRequest struct {
	MetricName           string               `protobuf:"bytes,1,opt,name=metricName,proto3" json:"metricName,omitempty"`
	Labels               map[string]string    `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Start                *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	End                  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	Step                 *duration.Duration   `protobuf:"bytes,5,opt,name=step,proto3" json:"step,omitempty"`
	ChunkSize            *duration.Duration   `protobuf:"bytes,6,opt,name=chunkSize,proto3" json:"chunkSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *QueryMetricsRequest) Reset()         { *m = QueryMetricsRequest{} }
func (m *QueryMetricsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryMetricsRequest) ProtoMessage()    {}
func (*QueryMetricsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82f4d3bdff1c5591, []int{0}
}

func (m *QueryMetricsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryMetricsRequest.Unmarshal(m, b)
}
func (m *QueryMetricsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryMetricsRequest.Marshal(b, m, deterministic)
}
func (m *QueryMetricsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMetricsRequest.Merge(m, src)
}
func (m *QueryMetricsRequest) XXX_Size() int {
	return xxx_messageInfo_QueryMetricsRequest.Size(m)
}
func (m *QueryMetricsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMetricsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMetricsRequest proto.InternalMessageInfo

func (m *QueryMetricsRequest) GetMetricName() string {
	if m != nil {
		return m.MetricName
	}
	return ""
}

func (m *QueryMetricsRequest) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *QueryMetricsRequest) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *QueryMetricsRequest) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *QueryMetricsRequest) GetStep() *duration.Duration {
	if m != nil {
		return m.Step
	}
	return nil
}

func (m *QueryMetricsRequest) GetChunkSize() *duration.Duration {
	if m != nil {
		return m.ChunkSize
	}
	return nil
}

type QueryMetricsResponse struct {
	Metrics              []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *QueryMetricsResponse) Reset()         { *m = QueryMetricsResponse{} }
func (m *QueryMetricsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMetricsResponse) ProtoMessage()    {}
func (*QueryMetricsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82f4d3bdff1c5591, []int{1}
}

func (m *QueryMetricsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryMetricsResponse.Unmarshal(m, b)
}
func (m *QueryMetricsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryMetricsResponse.Marshal(b, m, deterministic)
}
func (m *QueryMetricsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMetricsResponse.Merge(m, src)
}
func (m *QueryMetricsResponse) XXX_Size() int {
	return xxx_messageInfo_QueryMetricsResponse.Size(m)
}
func (m *QueryMetricsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMetricsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMetricsResponse proto.InternalMessageInfo

func (m *QueryMetricsResponse) GetMetrics() []*Metric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

type Metric struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels               map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Values               []*SamplePair     `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_82f4d3bdff1c5591, []int{2}
}

func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (m *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(m, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metric) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Metric) GetValues() []*SamplePair {
	if m != nil {
		return m.Values
	}
	return nil
}

type SamplePair struct {
	Time                 *timestamp.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Value                float64              `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SamplePair) Reset()         { *m = SamplePair{} }
func (m *SamplePair) String() string { return proto.CompactTextString(m) }
func (*SamplePair) ProtoMessage()    {}
func (*SamplePair) Descriptor() ([]byte, []int) {
	return fileDescriptor_82f4d3bdff1c5591, []int{3}
}

func (m *SamplePair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SamplePair.Unmarshal(m, b)
}
func (m *SamplePair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SamplePair.Marshal(b, m, deterministic)
}
func (m *SamplePair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SamplePair.Merge(m, src)
}
func (m *SamplePair) XXX_Size() int {
	return xxx_messageInfo_SamplePair.Size(m)
}
func (m *SamplePair) XXX_DiscardUnknown() {
	xxx_messageInfo_SamplePair.DiscardUnknown(m)
}

var xxx_messageInfo_SamplePair proto.InternalMessageInfo

func (m *SamplePair) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *SamplePair) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type WriteMetricsRequest struct {
	Metrics              []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *WriteMetricsRequest) Reset()         { *m = WriteMetricsRequest{} }
func (m *WriteMetricsRequest) String() string { return proto.CompactTextString(m) }
func (*WriteMetricsRequest) ProtoMessage()    {}
func (*WriteMetricsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82f4d3bdff1c5591, []int{4}
}

func (m *WriteMetricsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteMetricsRequest.Unmarshal(m, b)
}
func (m *WriteMetricsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteMetricsRequest.Marshal(b, m, deterministic)
}
func (m *WriteMetricsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteMetricsRequest.Merge(m, src)
}
func (m *WriteMetricsRequest) XXX_Size() int {
	return xxx_messageInfo_WriteMetricsRequest.Size(m)
}
func (m *WriteMetricsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteMetricsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WriteMetricsRequest proto.InternalMessageInfo

func (m *WriteMetricsRequest) GetMetrics() []*Metric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryMetricsRequest)(nil), "api.QueryMetricsRequest")
	proto.RegisterMapType((map[string]string)(nil), "api.QueryMetricsRequest.LabelsEntry")
	proto.RegisterType((*QueryMetricsResponse)(nil), "api.QueryMetricsResponse")
	proto.RegisterType((*Metric)(nil), "api.Metric")
	proto.RegisterMapType((map[string]string)(nil), "api.Metric.LabelsEntry")
	proto.RegisterType((*SamplePair)(nil), "api.SamplePair")
	proto.RegisterType((*WriteMetricsRequest)(nil), "api.WriteMetricsRequest")
}

func init() { proto.RegisterFile("api/metrics.proto", fileDescriptor_82f4d3bdff1c5591) }

var fileDescriptor_82f4d3bdff1c5591 = []byte{
	// 440 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe5, 0x38, 0x76, 0xd5, 0xc9, 0x01, 0x98, 0x56, 0xb0, 0x35, 0x52, 0x89, 0x22, 0x10,
	0x39, 0x80, 0x83, 0xc2, 0x01, 0x8a, 0xc2, 0x8d, 0xde, 0x00, 0x81, 0x8b, 0xc4, 0x79, 0xd3, 0x0e,
	0x65, 0x55, 0xff, 0x59, 0x76, 0xd7, 0x48, 0xe1, 0xc8, 0x1b, 0xf1, 0x72, 0x9c, 0x51, 0x66, 0x37,
	0x4a, 0x9a, 0x1a, 0x8a, 0xd4, 0x9b, 0x77, 0xe6, 0xf7, 0x4d, 0x66, 0xe6, 0x9b, 0xc0, 0x1d, 0xa9,
	0xd5, 0xa4, 0x22, 0x67, 0xd4, 0xa9, 0xcd, 0xb5, 0x69, 0x5c, 0x83, 0xb1, 0xd4, 0x2a, 0xbb, 0x7f,
	0xde, 0x34, 0xe7, 0x25, 0x4d, 0x38, 0x34, 0x6f, 0xbf, 0x4c, 0xa8, 0xd2, 0x6e, 0xe1, 0x89, 0xec,
	0xc1, 0x76, 0xd2, 0xa9, 0x8a, 0xac, 0x93, 0x95, 0x0e, 0xc0, 0xe1, 0x36, 0x70, 0xd6, 0x1a, 0xe9,
	0x54, 0x53, 0xfb, 0xfc, 0xe8, 0x77, 0x0f, 0xf6, 0x3e, 0xb6, 0x64, 0x16, 0xef, 0xfc, 0x2f, 0x17,
	0xf4, 0xad, 0x25, 0xeb, 0xf0, 0x10, 0xc0, 0xf7, 0xf2, 0x5e, 0x56, 0x24, 0xa2, 0x61, 0x34, 0xde,
	0x2d, 0x36, 0x22, 0x38, 0x83, 0xb4, 0x94, 0x73, 0x2a, 0xad, 0xe8, 0x0d, 0xe3, 0xf1, 0x60, 0xfa,
	0x30, 0x97, 0x5a, 0xe5, 0x1d, 0x95, 0xf2, 0xb7, 0x8c, 0x1d, 0xd7, 0xce, 0x2c, 0x8a, 0xa0, 0xc1,
	0x67, 0x90, 0x58, 0x27, 0x8d, 0x13, 0xf1, 0x30, 0x1a, 0x0f, 0xa6, 0x59, 0xee, 0xbb, 0xcc, 0x57,
	0x5d, 0xe6, 0x9f, 0x56, 0x63, 0x14, 0x1e, 0xc4, 0x27, 0x10, 0x53, 0x7d, 0x26, 0xfa, 0xd7, 0xf2,
	0x4b, 0x0c, 0x9f, 0x42, 0xdf, 0x3a, 0xd2, 0x22, 0x61, 0xfc, 0xe0, 0x0a, 0xfe, 0x26, 0x2c, 0xa1,
	0x60, 0x0c, 0x5f, 0xc0, 0xee, 0xe9, 0xd7, 0xb6, 0xbe, 0x38, 0x51, 0x3f, 0x48, 0xa4, 0xd7, 0x69,
	0xd6, 0x6c, 0x76, 0x04, 0x83, 0x8d, 0xf1, 0xf0, 0x36, 0xc4, 0x17, 0xb4, 0x08, 0xdb, 0x5a, 0x7e,
	0xe2, 0x3e, 0x24, 0xdf, 0x65, 0xd9, 0x92, 0xe8, 0x71, 0xcc, 0x3f, 0x5e, 0xf5, 0x5e, 0x46, 0xa3,
	0xd7, 0xb0, 0x7f, 0x79, 0x5b, 0x56, 0x37, 0xb5, 0x25, 0x7c, 0x04, 0x3b, 0xe1, 0x08, 0x44, 0xc4,
	0x9b, 0x1d, 0xf0, 0x66, 0x3d, 0x56, 0xac, 0x72, 0xa3, 0x5f, 0x11, 0xa4, 0x3e, 0x86, 0x08, 0xfd,
	0x7a, 0x6d, 0x12, 0x7f, 0xe3, 0x64, 0xcb, 0x9e, 0x7b, 0x1b, 0x45, 0x3a, 0x1d, 0x79, 0x0c, 0x29,
	0xf7, 0x66, 0x45, 0xcc, 0x82, 0x5b, 0x2c, 0x38, 0x91, 0x95, 0x2e, 0xe9, 0x83, 0x54, 0xa6, 0x08,
	0xe9, 0x9b, 0x8c, 0x5c, 0x00, 0xac, 0x0b, 0x62, 0x0e, 0xfd, 0xe5, 0xb1, 0xb2, 0xf4, 0xdf, 0x96,
	0x32, 0x77, 0xb9, 0x6e, 0x14, 0xea, 0x8e, 0x66, 0xb0, 0xf7, 0xd9, 0x28, 0x47, 0x5b, 0xe7, 0xfb,
	0x7f, 0x5b, 0x9c, 0xfe, 0x8c, 0x60, 0x27, 0x28, 0x71, 0x06, 0x09, 0x1b, 0x82, 0xe2, 0x6f, 0xa7,
	0x9c, 0x1d, 0x74, 0x64, 0x82, 0x6d, 0x47, 0x90, 0x70, 0x1f, 0x41, 0xdd, 0xd1, 0x53, 0x76, 0xf7,
	0xca, 0x88, 0xc7, 0xcb, 0x7f, 0xf2, 0x3c, 0xe5, 0xf7, 0xf3, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x4f, 0xcc, 0x84, 0xbe, 0x01, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MetricsClient interface {
	Query(ctx context.Context, in *QueryMetricsRequest, opts ...grpc.CallOption) (*QueryMetricsResponse, error)
	Write(ctx context.Context, in *WriteMetricsRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type metricsClient struct {
	cc *grpc.ClientConn
}

func NewMetricsClient(cc *grpc.ClientConn) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) Query(ctx context.Context, in *QueryMetricsRequest, opts ...grpc.CallOption) (*QueryMetricsResponse, error) {
	out := new(QueryMetricsResponse)
	err := c.cc.Invoke(ctx, "/api.Metrics/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) Write(ctx context.Context, in *WriteMetricsRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/api.Metrics/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
type MetricsServer interface {
	Query(context.Context, *QueryMetricsRequest) (*QueryMetricsResponse, error)
	Write(context.Context, *WriteMetricsRequest) (*empty.Empty, error)
}

// UnimplementedMetricsServer can be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (*UnimplementedMetricsServer) Query(ctx context.Context, req *QueryMetricsRequest) (*QueryMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (*UnimplementedMetricsServer) Write(ctx context.Context, req *WriteMetricsRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}

func RegisterMetricsServer(s *grpc.Server, srv MetricsServer) {
	s.RegisterService(&_Metrics_serviceDesc, srv)
}

func _Metrics_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Metrics/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).Query(ctx, req.(*QueryMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Metrics/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).Write(ctx, req.(*WriteMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Metrics_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Query",
			Handler:    _Metrics_Query_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _Metrics_Write_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/metrics.proto",
}
