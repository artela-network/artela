// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/crypto/v1/ethsecp256k1/keys.proto

package ethsecp256k1

import (
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	_ "github.com/cosmos/gogoproto/gogoproto"
	"github.com/cosmos/gogoproto/proto"
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

// PubKey defines a type alias for an ecdsa.PublicKey that implements
// Tendermint's PubKey interface. It represents the 33-byte compressed public
// key format.
type PubKey struct {
	// key is the public key in byte form
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (pubKey *PubKey) Reset() { *pubKey = PubKey{} }
func (*PubKey) ProtoMessage() {}
func (*PubKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_56b55ddbae0a9542, []int{0}
}
func (pubKey *PubKey) XXX_Unmarshal(b []byte) error {
	return pubKey.Unmarshal(b)
}
func (pubKey *PubKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PubKey.Marshal(b, pubKey, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := pubKey.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (pubKey *PubKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PubKey.Merge(pubKey, src)
}
func (pubKey *PubKey) XXX_Size() int {
	return pubKey.Size()
}
func (pubKey *PubKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PubKey.DiscardUnknown(pubKey)
}

var xxx_messageInfo_PubKey proto.InternalMessageInfo

func (pubKey *PubKey) GetKey() []byte {
	if pubKey != nil {
		return pubKey.Key
	}
	return nil
}

// PrivKey defines a type alias for an ecdsa.PrivateKey that implements
// Tendermint's PrivateKey interface.
type PrivKey struct {
	// key is the private key in byte form
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (privKey *PrivKey) Reset()         { *privKey = PrivKey{} }
func (privKey *PrivKey) String() string { return proto.CompactTextString(privKey) }
func (*PrivKey) ProtoMessage()          {}
func (*PrivKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_56b55ddbae0a9542, []int{1}
}
func (privKey *PrivKey) XXX_Unmarshal(b []byte) error {
	return privKey.Unmarshal(b)
}
func (privKey *PrivKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PrivKey.Marshal(b, privKey, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := privKey.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (privKey *PrivKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivKey.Merge(privKey, src)
}
func (privKey *PrivKey) XXX_Size() int {
	return privKey.Size()
}
func (privKey *PrivKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivKey.DiscardUnknown(privKey)
}

var xxx_messageInfo_PrivKey proto.InternalMessageInfo

func (privKey *PrivKey) GetKey() []byte {
	if privKey != nil {
		return privKey.Key
	}
	return nil
}

func init() {
	proto.RegisterType((*PubKey)(nil), "artela.crypto.v1.ethsecp256k1.PubKey")
	proto.RegisterType((*PrivKey)(nil), "artela.crypto.v1.ethsecp256k1.PrivKey")
}

func init() {
	proto.RegisterFile("artela/crypto/v1/ethsecp256k1/keys.proto", fileDescriptor_56b55ddbae0a9542)
}

var fileDescriptor_56b55ddbae0a9542 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x48, 0x2c, 0x2a, 0x49,
	0xcd, 0x49, 0xd4, 0x4f, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0x2d, 0xc9,
	0x28, 0x4e, 0x4d, 0x2e, 0x30, 0x32, 0x35, 0xcb, 0x36, 0xd4, 0xcf, 0x4e, 0xad, 0x2c, 0xd6, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x85, 0xa8, 0xd4, 0x83, 0xa8, 0xd4, 0x2b, 0x33, 0xd4, 0x43,
	0x56, 0x29, 0x25, 0x92, 0x9e, 0x9f, 0x9e, 0x0f, 0x56, 0xa9, 0x0f, 0x62, 0x41, 0x34, 0x29, 0x29,
	0x70, 0xb1, 0x05, 0x94, 0x26, 0x79, 0xa7, 0x56, 0x0a, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a,
	0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0x98, 0x56, 0x2c, 0x33, 0x16, 0xc8, 0x33, 0x28, 0x49,
	0x73, 0xb1, 0x07, 0x14, 0x65, 0x96, 0x61, 0x55, 0xe2, 0x14, 0x75, 0xe2, 0x91, 0x1c, 0xe3, 0x85,
	0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3,
	0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x0e, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9,
	0xfa, 0x10, 0x87, 0xe9, 0xe6, 0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0x43, 0xb9, 0x50, 0xaf, 0xa4,
	0x16, 0xa5, 0x96, 0xe6, 0xc2, 0x7c, 0x87, 0xec, 0xe0, 0x24, 0x36, 0xb0, 0x0b, 0x8d, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x8e, 0xb4, 0xc3, 0x6c, 0x02, 0x01, 0x00, 0x00,
}

func (pubKey *PubKey) Marshal() (dAtA []byte, err error) {
	size := pubKey.Size()
	dAtA = make([]byte, size)
	n, err := pubKey.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (pubKey *PubKey) MarshalTo(dAtA []byte) (int, error) {
	size := pubKey.Size()
	return pubKey.MarshalToSizedBuffer(dAtA[:size])
}

func (pubKey *PubKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(pubKey.Key) > 0 {
		i -= len(pubKey.Key)
		copy(dAtA[i:], pubKey.Key)
		i = encodeVarintKeys(dAtA, i, uint64(len(pubKey.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (privKey *PrivKey) Marshal() (dAtA []byte, err error) {
	size := privKey.Size()
	dAtA = make([]byte, size)
	n, err := privKey.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (privKey *PrivKey) MarshalTo(dAtA []byte) (int, error) {
	size := privKey.Size()
	return privKey.MarshalToSizedBuffer(dAtA[:size])
}

func (privKey *PrivKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(privKey.Key) > 0 {
		i -= len(privKey.Key)
		copy(dAtA[i:], privKey.Key)
		i = encodeVarintKeys(dAtA, i, uint64(len(privKey.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintKeys(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeys(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (pubKey *PubKey) Size() (n int) {
	if pubKey == nil {
		return 0
	}
	var l int
	_ = l
	l = len(pubKey.Key)
	if l > 0 {
		n += 1 + l + sovKeys(uint64(l))
	}
	return n
}

func (privKey *PrivKey) Size() (n int) {
	if privKey == nil {
		return 0
	}
	var l int
	_ = l
	l = len(privKey.Key)
	if l > 0 {
		n += 1 + l + sovKeys(uint64(l))
	}
	return n
}

func sovKeys(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeys(x uint64) (n int) {
	return sovKeys(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (pubKey *PubKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeys
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
			return fmt.Errorf("proto: PubKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PubKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthKeys
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeys
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			pubKey.Key = append(pubKey.Key[:0], dAtA[iNdEx:postIndex]...)
			if pubKey.Key == nil {
				pubKey.Key = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeys(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthKeys
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
func (privKey *PrivKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeys
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
			return fmt.Errorf("proto: PrivKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PrivKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthKeys
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeys
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			privKey.Key = append(privKey.Key[:0], dAtA[iNdEx:postIndex]...)
			if privKey.Key == nil {
				privKey.Key = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeys(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthKeys
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
func skipKeys(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeys
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
					return 0, ErrIntOverflowKeys
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
					return 0, ErrIntOverflowKeys
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
				return 0, ErrInvalidLengthKeys
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeys
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeys
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthKeys        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowKeys          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupKeys = fmt.Errorf("proto: unexpected end of group")
)
