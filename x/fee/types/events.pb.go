// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/fee/v1/events.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this support file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// EventFee is the event type for the fee module
type EventFee struct {
	// base_fee for EIP-1559 blocks
	BaseFee string `protobuf:"bytes,1,opt,name=base_fee,json=baseFee,proto3" json:"base_fee,omitempty"`
}

func (m *EventFee) Reset()         { *m = EventFee{} }
func (m *EventFee) String() string { return proto.CompactTextString(m) }
func (*EventFee) ProtoMessage()    {}
func (*EventFee) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1be03613ec3d65c, []int{0}
}
func (m *EventFee) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventFee) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventFee.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventFee) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventFee.Merge(m, src)
}
func (m *EventFee) XXX_Size() int {
	return m.Size()
}
func (m *EventFee) XXX_DiscardUnknown() {
	xxx_messageInfo_EventFee.DiscardUnknown(m)
}

var xxx_messageInfo_EventFee proto.InternalMessageInfo

func (m *EventFee) GetBaseFee() string {
	if m != nil {
		return m.BaseFee
	}
	return ""
}

// EventBlockGas defines an Ethereum block gas event
type EventBlockGas struct {
	// height of the block
	Height string `protobuf:"bytes,1,opt,name=height,proto3" json:"height,omitempty"`
	// amount of gas wanted by the block
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventBlockGas) Reset()         { *m = EventBlockGas{} }
func (m *EventBlockGas) String() string { return proto.CompactTextString(m) }
func (*EventBlockGas) ProtoMessage()    {}
func (*EventBlockGas) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1be03613ec3d65c, []int{1}
}
func (m *EventBlockGas) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventBlockGas) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventBlockGas.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventBlockGas) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventBlockGas.Merge(m, src)
}
func (m *EventBlockGas) XXX_Size() int {
	return m.Size()
}
func (m *EventBlockGas) XXX_DiscardUnknown() {
	xxx_messageInfo_EventBlockGas.DiscardUnknown(m)
}

var xxx_messageInfo_EventBlockGas proto.InternalMessageInfo

func (m *EventBlockGas) GetHeight() string {
	if m != nil {
		return m.Height
	}
	return ""
}

func (m *EventBlockGas) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func init() {
	proto.RegisterType((*EventFee)(nil), "artela.fee.v1.EventFee")
	proto.RegisterType((*EventBlockGas)(nil), "artela.fee.v1.EventBlockGas")
}

func init() { proto.RegisterFile("artela/fee/v1/events.proto", fileDescriptor_e1be03613ec3d65c) }

var fileDescriptor_e1be03613ec3d65c = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4a, 0x2c, 0x2a, 0x49,
	0xcd, 0x49, 0xd4, 0x4f, 0x4b, 0x4d, 0xd5, 0x2f, 0x33, 0xd4, 0x4f, 0x2d, 0x4b, 0xcd, 0x2b, 0x29,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x85, 0xc8, 0xe9, 0xa5, 0xa5, 0xa6, 0xea, 0x95,
	0x19, 0x2a, 0xa9, 0x72, 0x71, 0xb8, 0x82, 0xa4, 0xdd, 0x52, 0x53, 0x85, 0x24, 0xb9, 0x38, 0x92,
	0x12, 0x8b, 0x53, 0xe3, 0xd3, 0x52, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xd8, 0x41,
	0x7c, 0xb7, 0xd4, 0x54, 0x25, 0x7b, 0x2e, 0x5e, 0xb0, 0x32, 0xa7, 0x9c, 0xfc, 0xe4, 0x6c, 0xf7,
	0xc4, 0x62, 0x21, 0x31, 0x2e, 0xb6, 0x8c, 0xd4, 0xcc, 0xf4, 0x8c, 0x12, 0xa8, 0x4a, 0x28, 0x0f,
	0x24, 0x9e, 0x98, 0x9b, 0x5f, 0x9a, 0x57, 0x22, 0xc1, 0x04, 0x11, 0x87, 0xf0, 0x9c, 0x3c, 0x4f,
	0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18,
	0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x3f, 0x3d, 0xb3, 0x24, 0xa3, 0x34,
	0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0xe2, 0x36, 0xdd, 0xbc, 0xd4, 0x92, 0xf2, 0xfc, 0xa2, 0x6c,
	0x28, 0x17, 0xe4, 0x85, 0x0a, 0xb0, 0x5f, 0x4a, 0x2a, 0x0b, 0x52, 0x8b, 0x93, 0xd8, 0xc0, 0x1e,
	0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x17, 0x88, 0xce, 0xe6, 0x00, 0x00, 0x00,
}

func (m *EventFee) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventFee) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventFee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BaseFee) > 0 {
		i -= len(m.BaseFee)
		copy(dAtA[i:], m.BaseFee)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.BaseFee)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventBlockGas) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventBlockGas) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventBlockGas) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Height) > 0 {
		i -= len(m.Height)
		copy(dAtA[i:], m.Height)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Height)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventFee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BaseFee)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *EventBlockGas) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Height)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventFee) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EventFee: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventFee: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseFee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EventBlockGas) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EventBlockGas: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventBlockGas: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Height = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
