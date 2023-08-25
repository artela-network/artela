// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/types/v1/dynamic_fee.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

// ExtensionOptionDynamicFeeTx is an extension option that specifies the maxPrioPrice for cosmos txs
type ExtensionOptionDynamicFeeTx struct {
	// max_priority_price is the same as `max_priority_fee_per_gas` in eip-1559 spec
	MaxPriorityPrice github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=max_priority_price,json=maxPriorityPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"max_priority_price"`
}

func (m *ExtensionOptionDynamicFeeTx) Reset()         { *m = ExtensionOptionDynamicFeeTx{} }
func (m *ExtensionOptionDynamicFeeTx) String() string { return proto.CompactTextString(m) }
func (*ExtensionOptionDynamicFeeTx) ProtoMessage()    {}
func (*ExtensionOptionDynamicFeeTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_cfd40f81b1d1654b, []int{0}
}
func (m *ExtensionOptionDynamicFeeTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExtensionOptionDynamicFeeTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExtensionOptionDynamicFeeTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExtensionOptionDynamicFeeTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtensionOptionDynamicFeeTx.Merge(m, src)
}
func (m *ExtensionOptionDynamicFeeTx) XXX_Size() int {
	return m.Size()
}
func (m *ExtensionOptionDynamicFeeTx) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtensionOptionDynamicFeeTx.DiscardUnknown(m)
}

var xxx_messageInfo_ExtensionOptionDynamicFeeTx proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ExtensionOptionDynamicFeeTx)(nil), "artela.types.v1.ExtensionOptionDynamicFeeTx")
}

func init() { proto.RegisterFile("artela/types/v1/dynamic_fee.proto", fileDescriptor_cfd40f81b1d1654b) }

var fileDescriptor_cfd40f81b1d1654b = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x2c, 0x2a, 0x49,
	0xcd, 0x49, 0xd4, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x2f, 0x33, 0xd4, 0x4f, 0xa9, 0xcc, 0x4b,
	0xcc, 0xcd, 0x4c, 0x8e, 0x4f, 0x4b, 0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x87,
	0x28, 0xd1, 0x03, 0x2b, 0xd1, 0x2b, 0x33, 0x94, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0xcb, 0xe9,
	0x83, 0x58, 0x10, 0x65, 0x4a, 0xd5, 0x5c, 0xd2, 0xae, 0x15, 0x25, 0xa9, 0x79, 0xc5, 0x99, 0xf9,
	0x79, 0xfe, 0x05, 0x25, 0x99, 0xf9, 0x79, 0x2e, 0x10, 0xa3, 0xdc, 0x52, 0x53, 0x43, 0x2a, 0x84,
	0x62, 0xb8, 0x84, 0x72, 0x13, 0x2b, 0xe2, 0x0b, 0x8a, 0x32, 0xf3, 0x8b, 0x32, 0x4b, 0x2a, 0x41,
	0x8c, 0xe4, 0x54, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x4e, 0x27, 0xbd, 0x13, 0xf7, 0xe4, 0x19, 0x6e,
	0xdd, 0x93, 0x57, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce,
	0x2f, 0xce, 0xcd, 0x2f, 0x86, 0x52, 0xba, 0xc5, 0x29, 0xd9, 0x10, 0x27, 0xea, 0x79, 0xe6, 0x95,
	0x04, 0x09, 0xe4, 0x26, 0x56, 0x04, 0x40, 0x0d, 0x0a, 0x00, 0x99, 0xe3, 0xe4, 0x7c, 0xe2, 0x91,
	0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1,
	0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x9a, 0x48, 0x66, 0x42, 0x3c, 0xa2, 0x9b, 0x97,
	0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x0d, 0xe5, 0x82, 0x3c, 0x0d, 0x36, 0x3a, 0x89, 0x0d, 0xec, 0x11,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd4, 0x91, 0x20, 0xe2, 0x14, 0x01, 0x00, 0x00,
}

func (m *ExtensionOptionDynamicFeeTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtensionOptionDynamicFeeTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExtensionOptionDynamicFeeTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MaxPriorityPrice.Size()
		i -= size
		if _, err := m.MaxPriorityPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDynamicFee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintDynamicFee(dAtA []byte, offset int, v uint64) int {
	offset -= sovDynamicFee(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ExtensionOptionDynamicFeeTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MaxPriorityPrice.Size()
	n += 1 + l + sovDynamicFee(uint64(l))
	return n
}

func sovDynamicFee(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDynamicFee(x uint64) (n int) {
	return sovDynamicFee(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExtensionOptionDynamicFeeTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDynamicFee
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
			return fmt.Errorf("proto: ExtensionOptionDynamicFeeTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtensionOptionDynamicFeeTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxPriorityPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicFee
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
				return ErrInvalidLengthDynamicFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDynamicFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxPriorityPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDynamicFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDynamicFee
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
func skipDynamicFee(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDynamicFee
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
					return 0, ErrIntOverflowDynamicFee
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
					return 0, ErrIntOverflowDynamicFee
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
				return 0, ErrInvalidLengthDynamicFee
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDynamicFee
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDynamicFee
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDynamicFee        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDynamicFee          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDynamicFee = fmt.Errorf("proto: unexpected end of group")
)
