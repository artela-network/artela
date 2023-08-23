// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/fee/v1/fee.proto

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

// Params defines the EVM module parameters
type Params struct {
	// no_base_fee forces the EIP-1559 base fee to 0 (needed for 0 price calls)
	NoBaseFee bool `protobuf:"varint,1,opt,name=no_base_fee,json=noBaseFee,proto3" json:"no_base_fee,omitempty"`
	// base_fee_change_denominator bounds the amount the base fee can change
	// between blocks.
	BaseFeeChangeDenominator uint32 `protobuf:"varint,2,opt,name=base_fee_change_denominator,json=baseFeeChangeDenominator,proto3" json:"base_fee_change_denominator,omitempty"`
	// elasticity_multiplier bounds the maximum gas limit an EIP-1559 block may
	// have.
	ElasticityMultiplier uint32 `protobuf:"varint,3,opt,name=elasticity_multiplier,json=elasticityMultiplier,proto3" json:"elasticity_multiplier,omitempty"`
	// enable_height defines at which block height the base fee calculation is enabled.
	EnableHeight int64 `protobuf:"varint,5,opt,name=enable_height,json=enableHeight,proto3" json:"enable_height,omitempty"`
	// base_fee for EIP-1559 blocks.
	BaseFee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=base_fee,json=baseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"base_fee"`
	// min_gas_price defines the minimum gas price value for cosmos and eth transactions
	MinGasPrice github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=min_gas_price,json=minGasPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_gas_price"`
	// min_gas_multiplier bounds the minimum gas used to be charged
	// to senders based on gas limit
	MinGasMultiplier github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=min_gas_multiplier,json=minGasMultiplier,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_gas_multiplier"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b545c073c30863c, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetNoBaseFee() bool {
	if m != nil {
		return m.NoBaseFee
	}
	return false
}

func (m *Params) GetBaseFeeChangeDenominator() uint32 {
	if m != nil {
		return m.BaseFeeChangeDenominator
	}
	return 0
}

func (m *Params) GetElasticityMultiplier() uint32 {
	if m != nil {
		return m.ElasticityMultiplier
	}
	return 0
}

func (m *Params) GetEnableHeight() int64 {
	if m != nil {
		return m.EnableHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "artela.fee.v1.Params")
}

func init() { proto.RegisterFile("artela/fee/v1/fee.proto", fileDescriptor_5b545c073c30863c) }

var fileDescriptor_5b545c073c30863c = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x5f, 0xcb, 0xd3, 0x30,
	0x14, 0xc6, 0x1b, 0xf7, 0xbe, 0x7b, 0xfb, 0x66, 0x16, 0x46, 0x98, 0x58, 0x14, 0xba, 0xa2, 0x20,
	0xbd, 0x59, 0xcb, 0xd8, 0xb5, 0x37, 0x73, 0xa8, 0x13, 0x84, 0xd1, 0x4b, 0x11, 0x42, 0xda, 0x9d,
	0xb5, 0x61, 0x6d, 0x52, 0x9a, 0x6c, 0xba, 0x6f, 0xe1, 0xc7, 0xda, 0xe5, 0xae, 0x44, 0xbc, 0x18,
	0xb2, 0x7d, 0x11, 0x69, 0xbb, 0x7f, 0xb7, 0x7a, 0x75, 0x92, 0xf3, 0x3b, 0x79, 0xce, 0x49, 0xf2,
	0xe0, 0xe7, 0xac, 0xd4, 0x90, 0xb1, 0x60, 0x01, 0x10, 0xac, 0x87, 0x55, 0xf0, 0x8b, 0x52, 0x6a,
	0x49, 0xac, 0x06, 0xf8, 0x55, 0x66, 0x3d, 0x7c, 0xd1, 0x4b, 0x64, 0x22, 0x6b, 0x12, 0x54, 0xab,
	0xa6, 0xe8, 0xd5, 0xcf, 0x16, 0x6e, 0xcf, 0x58, 0xc9, 0x72, 0x45, 0x1c, 0xdc, 0x11, 0x92, 0x46,
	0x4c, 0x01, 0x5d, 0x00, 0xd8, 0xc8, 0x45, 0x9e, 0x19, 0x3e, 0x0a, 0x39, 0x66, 0x0a, 0xde, 0x03,
	0x90, 0xb7, 0xf8, 0xe5, 0x19, 0xd2, 0x38, 0x65, 0x22, 0x01, 0x3a, 0x07, 0x21, 0x73, 0x2e, 0x98,
	0x96, 0xa5, 0xfd, 0xc4, 0x45, 0x9e, 0x15, 0xda, 0x51, 0x53, 0xfd, 0xae, 0x2e, 0x98, 0x5c, 0x39,
	0x19, 0xe1, 0x67, 0x90, 0x31, 0xa5, 0x79, 0xcc, 0xf5, 0x86, 0xe6, 0xab, 0x4c, 0xf3, 0x22, 0xe3,
	0x50, 0xda, 0xad, 0xfa, 0x60, 0xef, 0x0a, 0x3f, 0x5f, 0x18, 0x79, 0x8d, 0x2d, 0x10, 0x2c, 0xca,
	0x80, 0xa6, 0xc0, 0x93, 0x54, 0xdb, 0xf7, 0x2e, 0xf2, 0x5a, 0xe1, 0xd3, 0x26, 0xf9, 0xb1, 0xce,
	0x91, 0x29, 0x36, 0x2f, 0x53, 0xb7, 0x5d, 0xe4, 0x3d, 0x8e, 0xfd, 0xed, 0xbe, 0x6f, 0xfc, 0xde,
	0xf7, 0xdf, 0x24, 0x5c, 0xa7, 0xab, 0xc8, 0x8f, 0x65, 0x1e, 0xc4, 0x52, 0xe5, 0x52, 0x9d, 0xc2,
	0x40, 0xcd, 0x97, 0x81, 0xde, 0x14, 0xa0, 0xfc, 0xa9, 0xd0, 0xe1, 0xc3, 0x69, 0x6a, 0x12, 0x62,
	0x2b, 0xe7, 0x82, 0x26, 0x4c, 0xd1, 0xa2, 0xe4, 0x31, 0xd8, 0x0f, 0xff, 0xac, 0x37, 0x81, 0x38,
	0xec, 0xe4, 0x5c, 0x7c, 0x60, 0x6a, 0x56, 0x49, 0x90, 0xaf, 0x98, 0x9c, 0x35, 0x6f, 0x6e, 0x6d,
	0xfe, 0x97, 0x70, 0xb7, 0x11, 0xbe, 0xbe, 0xd0, 0xa7, 0x3b, 0xf3, 0xae, 0x7b, 0x1f, 0x76, 0xb9,
	0xe0, 0x9a, 0xb3, 0xec, 0xf2, 0x7d, 0xe3, 0xe9, 0xf6, 0xe0, 0xa0, 0xdd, 0xc1, 0x41, 0x7f, 0x0e,
	0x0e, 0xfa, 0x71, 0x74, 0x8c, 0xdd, 0xd1, 0x31, 0x7e, 0x1d, 0x1d, 0xe3, 0x4b, 0x70, 0xd3, 0xab,
	0xb1, 0xc8, 0x40, 0x80, 0xfe, 0x26, 0xcb, 0xe5, 0x69, 0x5b, 0xd9, 0xe8, 0x7b, 0xed, 0xa7, 0xba,
	0x71, 0xd4, 0xae, 0xad, 0x32, 0xfa, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x34, 0x50, 0xe9, 0x1f, 0x6a,
	0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MinGasMultiplier.Size()
		i -= size
		if _, err := m.MinGasMultiplier.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.MinGasPrice.Size()
		i -= size
		if _, err := m.MinGasPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.BaseFee.Size()
		i -= size
		if _, err := m.BaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFee(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.EnableHeight != 0 {
		i = encodeVarintFee(dAtA, i, uint64(m.EnableHeight))
		i--
		dAtA[i] = 0x28
	}
	if m.ElasticityMultiplier != 0 {
		i = encodeVarintFee(dAtA, i, uint64(m.ElasticityMultiplier))
		i--
		dAtA[i] = 0x18
	}
	if m.BaseFeeChangeDenominator != 0 {
		i = encodeVarintFee(dAtA, i, uint64(m.BaseFeeChangeDenominator))
		i--
		dAtA[i] = 0x10
	}
	if m.NoBaseFee {
		i--
		if m.NoBaseFee {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintFee(dAtA []byte, offset int, v uint64) int {
	offset -= sovFee(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NoBaseFee {
		n += 2
	}
	if m.BaseFeeChangeDenominator != 0 {
		n += 1 + sovFee(uint64(m.BaseFeeChangeDenominator))
	}
	if m.ElasticityMultiplier != 0 {
		n += 1 + sovFee(uint64(m.ElasticityMultiplier))
	}
	if m.EnableHeight != 0 {
		n += 1 + sovFee(uint64(m.EnableHeight))
	}
	l = m.BaseFee.Size()
	n += 1 + l + sovFee(uint64(l))
	l = m.MinGasPrice.Size()
	n += 1 + l + sovFee(uint64(l))
	l = m.MinGasMultiplier.Size()
	n += 1 + l + sovFee(uint64(l))
	return n
}

func sovFee(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFee(x uint64) (n int) {
	return sovFee(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFee
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NoBaseFee", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.NoBaseFee = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFeeChangeDenominator", wireType)
			}
			m.BaseFeeChangeDenominator = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaseFeeChangeDenominator |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ElasticityMultiplier", wireType)
			}
			m.ElasticityMultiplier = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ElasticityMultiplier |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableHeight", wireType)
			}
			m.EnableHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EnableHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
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
				return ErrInvalidLengthFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinGasPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
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
				return ErrInvalidLengthFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinGasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinGasMultiplier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFee
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
				return ErrInvalidLengthFee
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFee
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinGasMultiplier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFee(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFee
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
func skipFee(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFee
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
					return 0, ErrIntOverflowFee
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
					return 0, ErrIntOverflowFee
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
				return 0, ErrInvalidLengthFee
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFee
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFee
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFee        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFee          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFee = fmt.Errorf("proto: unexpected end of group")
)
