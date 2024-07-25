// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/evm/v1/txs.proto

package txs

import (
	context "context"
	encoding_binary "encoding/binary"
	fmt "fmt"
	"github.com/artela-network/artela/x/evm/txs/support"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// MsgEthereumTx encapsulates an Ethereum txs as an SDK message.
type MsgEthereumTx struct {
	// data is inner txs data of the Ethereum txs
	Data *types.Any `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	// size is the encoded storage size of the txs (DEPRECATED)
	Size_ float64 `protobuf:"fixed64,2,opt,name=size,proto3" json:"-"`
	// hash of the txs in hex format
	Hash string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty" rlp:"-"`
	// from is the ethereum signer address in hex format. This address value is checked
	// against the address derived from the signature (V, R, S) using the
	// secp256k1 elliptic curve
	From string `protobuf:"bytes,4,opt,name=from,proto3" json:"from,omitempty"`
}

func (msg *MsgEthereumTx) Reset()         { *msg = MsgEthereumTx{} }
func (msg *MsgEthereumTx) String() string { return proto.CompactTextString(msg) }
func (*MsgEthereumTx) ProtoMessage()      {}
func (*MsgEthereumTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{0}
}
func (msg *MsgEthereumTx) XXX_Unmarshal(b []byte) error {
	return msg.Unmarshal(b)
}
func (msg *MsgEthereumTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEthereumTx.Marshal(b, msg, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := msg.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (msg *MsgEthereumTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEthereumTx.Merge(msg, src)
}
func (msg *MsgEthereumTx) XXX_Size() int {
	return msg.Size()
}
func (msg *MsgEthereumTx) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEthereumTx.DiscardUnknown(msg)
}

var xxx_messageInfo_MsgEthereumTx proto.InternalMessageInfo

// LegacyTx is the txs data of regular Ethereum transactions.
// NOTE: All non-protected transactions (i.e non EIP155 signed) will fail if the
// AllowUnprotectedTxs parameter is disabled.
type LegacyTx struct {
	// nonce corresponds to the account nonce (txs sequence).
	Nonce uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// gas_price defines the value for each gas unit
	GasPrice *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=gas_price,json=gasPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"gas_price,omitempty"`
	// gas defines the gas limit defined for the txs.
	GasLimit uint64 `protobuf:"varint,3,opt,name=gas,proto3" json:"gas,omitempty"`
	// to is the hex formatted address of the recipient
	To string `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	// value defines the unsigned integer value of the txs amount.
	Amount *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=value,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"value,omitempty"`
	// data is the data payload bytes of the txs.
	Data []byte `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	// v defines the signature value
	V []byte `protobuf:"bytes,7,opt,name=v,proto3" json:"v,omitempty"`
	// r defines the signature value
	R []byte `protobuf:"bytes,8,opt,name=r,proto3" json:"r,omitempty"`
	// s define the signature value
	S []byte `protobuf:"bytes,9,opt,name=s,proto3" json:"s,omitempty"`
}

func (tx *LegacyTx) Reset()         { *tx = LegacyTx{} }
func (tx *LegacyTx) String() string { return proto.CompactTextString(tx) }
func (*LegacyTx) ProtoMessage()     {}
func (*LegacyTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{1}
}
func (tx *LegacyTx) XXX_Unmarshal(b []byte) error {
	return tx.Unmarshal(b)
}
func (tx *LegacyTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyTx.Marshal(b, tx, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := tx.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (tx *LegacyTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyTx.Merge(tx, src)
}
func (tx *LegacyTx) XXX_Size() int {
	return tx.Size()
}
func (tx *LegacyTx) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyTx.DiscardUnknown(tx)
}

var xxx_messageInfo_LegacyTx proto.InternalMessageInfo

// AccessListTx is the data of EIP-2930 access list transactions.
type AccessListTx struct {
	// chain_id of the destination EVM chain
	ChainID *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"chainID"`
	// nonce corresponds to the account nonce (txs sequence).
	Nonce uint64 `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// gas_price defines the value for each gas unit
	GasPrice *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=gas_price,json=gasPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"gas_price,omitempty"`
	// gas defines the gas limit defined for the txs.
	GasLimit uint64 `protobuf:"varint,4,opt,name=gas,proto3" json:"gas,omitempty"`
	// to is the recipient address in hex format
	To string `protobuf:"bytes,5,opt,name=to,proto3" json:"to,omitempty"`
	// value defines the unsigned integer value of the txs amount.
	Amount *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=value,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"value,omitempty"`
	// data is the data payload bytes of the txs.
	Data []byte `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	// accesses is an array of access tuples
	Accesses AccessList `protobuf:"bytes,8,rep,name=accesses,proto3,castrepeated=AccessList" json:"accessList"`
	// v defines the signature value
	V []byte `protobuf:"bytes,9,opt,name=v,proto3" json:"v,omitempty"`
	// r defines the signature value
	R []byte `protobuf:"bytes,10,opt,name=r,proto3" json:"r,omitempty"`
	// s define the signature value
	S []byte `protobuf:"bytes,11,opt,name=s,proto3" json:"s,omitempty"`
}

func (tx *AccessListTx) Reset()         { *tx = AccessListTx{} }
func (tx *AccessListTx) String() string { return proto.CompactTextString(tx) }
func (*AccessListTx) ProtoMessage()     {}
func (*AccessListTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{2}
}
func (tx *AccessListTx) XXX_Unmarshal(b []byte) error {
	return tx.Unmarshal(b)
}
func (tx *AccessListTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccessListTx.Marshal(b, tx, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := tx.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (tx *AccessListTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessListTx.Merge(tx, src)
}
func (tx *AccessListTx) XXX_Size() int {
	return tx.Size()
}
func (tx *AccessListTx) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessListTx.DiscardUnknown(tx)
}

var xxx_messageInfo_AccessListTx proto.InternalMessageInfo

// DynamicFeeTx is the data of EIP-1559 dinamic fee transactions.
type DynamicFeeTx struct {
	// chain_id of the destination EVM chain
	ChainID *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"chainID"`
	// nonce corresponds to the account nonce (txs sequence).
	Nonce uint64 `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// gas_tip_cap defines the max value for the gas tip
	GasTipCap *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=gas_tip_cap,json=gasTipCap,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"gas_tip_cap,omitempty"`
	// gas_fee_cap defines the max value for the gas fee
	GasFeeCap *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=gas_fee_cap,json=gasFeeCap,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"gas_fee_cap,omitempty"`
	// gas defines the gas limit defined for the txs.
	GasLimit uint64 `protobuf:"varint,5,opt,name=gas,proto3" json:"gas,omitempty"`
	// to is the hex formatted address of the recipient
	To string `protobuf:"bytes,6,opt,name=to,proto3" json:"to,omitempty"`
	// value defines the the txs amount.
	Amount *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=value,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"value,omitempty"`
	// data is the data payload bytes of the txs.
	Data []byte `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
	// accesses is an array of access tuples
	Accesses AccessList `protobuf:"bytes,9,rep,name=accesses,proto3,castrepeated=AccessList" json:"accessList"`
	// v defines the signature value
	V []byte `protobuf:"bytes,10,opt,name=v,proto3" json:"v,omitempty"`
	// r defines the signature value
	R []byte `protobuf:"bytes,11,opt,name=r,proto3" json:"r,omitempty"`
	// s define the signature value
	S []byte `protobuf:"bytes,12,opt,name=s,proto3" json:"s,omitempty"`
}

func (tx *DynamicFeeTx) Reset()         { *tx = DynamicFeeTx{} }
func (tx *DynamicFeeTx) String() string { return proto.CompactTextString(tx) }
func (*DynamicFeeTx) ProtoMessage()     {}
func (*DynamicFeeTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{3}
}
func (tx *DynamicFeeTx) XXX_Unmarshal(b []byte) error {
	return tx.Unmarshal(b)
}
func (tx *DynamicFeeTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DynamicFeeTx.Marshal(b, tx, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := tx.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (tx *DynamicFeeTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynamicFeeTx.Merge(tx, src)
}
func (tx *DynamicFeeTx) XXX_Size() int {
	return tx.Size()
}
func (tx *DynamicFeeTx) XXX_DiscardUnknown() {
	xxx_messageInfo_DynamicFeeTx.DiscardUnknown(tx)
}

var xxx_messageInfo_DynamicFeeTx proto.InternalMessageInfo

// ExtensionOptionsEthereumTx is an extension option for ethereum transactions
type ExtensionOptionsEthereumTx struct {
}

func (m *ExtensionOptionsEthereumTx) Reset()         { *m = ExtensionOptionsEthereumTx{} }
func (m *ExtensionOptionsEthereumTx) String() string { return proto.CompactTextString(m) }
func (*ExtensionOptionsEthereumTx) ProtoMessage()    {}
func (*ExtensionOptionsEthereumTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{4}
}
func (m *ExtensionOptionsEthereumTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExtensionOptionsEthereumTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExtensionOptionsEthereumTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExtensionOptionsEthereumTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtensionOptionsEthereumTx.Merge(m, src)
}
func (m *ExtensionOptionsEthereumTx) XXX_Size() int {
	return m.Size()
}
func (m *ExtensionOptionsEthereumTx) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtensionOptionsEthereumTx.DiscardUnknown(m)
}

var xxx_messageInfo_ExtensionOptionsEthereumTx proto.InternalMessageInfo

// MsgEthereumTxResponse defines the Msg/EthereumTx response type.
type MsgEthereumTxResponse struct {
	// hash of the ethereum txs in hex format. This hash differs from the
	// Tendermint sha256 hash of the txs bytes. See
	// https://github.com/tendermint/tendermint/issues/6539 for reference
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	// logs contains the txs hash and the proto-compatible ethereum
	// logs.
	Logs []*support.Log `protobuf:"bytes,2,rep,name=logs,proto3" json:"logs,omitempty"`
	// ret is the returned data from evm function (result or data supplied with revert
	// opcode)
	Ret []byte `protobuf:"bytes,3,opt,name=ret,proto3" json:"ret,omitempty"`
	// vm_error is the error returned by vm execution
	VmError string `protobuf:"bytes,4,opt,name=vm_error,json=vmError,proto3" json:"vm_error,omitempty"`
	// gas_used specifies how much gas was consumed by the txs
	GasUsed uint64 `protobuf:"varint,5,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// cumulative gas used
	CumulativeGasUsed uint64 `protobuf:"varint,6,opt,name=cumulative_gas_used,json=cumulativeGasUsed,proto3" json:"cumulative_gas_used,omitempty"`
}

func (m *MsgEthereumTxResponse) Reset()         { *m = MsgEthereumTxResponse{} }
func (m *MsgEthereumTxResponse) String() string { return proto.CompactTextString(m) }
func (*MsgEthereumTxResponse) ProtoMessage()    {}
func (*MsgEthereumTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{5}
}
func (m *MsgEthereumTxResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgEthereumTxResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEthereumTxResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgEthereumTxResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEthereumTxResponse.Merge(m, src)
}
func (m *MsgEthereumTxResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgEthereumTxResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEthereumTxResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEthereumTxResponse proto.InternalMessageInfo

// MsgUpdateParams defines a Msg for updating the x/evm module parameters.
type MsgUpdateParams struct {
	// authority is the address of the governance account.
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// params defines the x/evm parameters to update.
	// NOTE: All parameters must be supplied.
	Params support.Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdateParams) Reset()         { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParams) ProtoMessage()    {}
func (*MsgUpdateParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{6}
}
func (m *MsgUpdateParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParams.Merge(m, src)
}
func (m *MsgUpdateParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParams proto.InternalMessageInfo

func (m *MsgUpdateParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdateParams) GetParams() support.Params {
	if m != nil {
		return m.Params
	}
	return support.Params{}
}

// MsgUpdateParamsResponse defines the response structure for executing a
// MsgUpdateParams message.
type MsgUpdateParamsResponse struct {
}

func (m *MsgUpdateParamsResponse) Reset()         { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParamsResponse) ProtoMessage()    {}
func (*MsgUpdateParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c43c0836c37bbe6, []int{7}
}
func (m *MsgUpdateParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParamsResponse.Merge(m, src)
}
func (m *MsgUpdateParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgEthereumTx)(nil), "artela.evm.v1.MsgEthereumTx")
	proto.RegisterType((*LegacyTx)(nil), "artela.evm.v1.LegacyTx")
	proto.RegisterType((*AccessListTx)(nil), "artela.evm.v1.AccessListTx")
	proto.RegisterType((*DynamicFeeTx)(nil), "artela.evm.v1.DynamicFeeTx")
	proto.RegisterType((*ExtensionOptionsEthereumTx)(nil), "artela.evm.v1.ExtensionOptionsEthereumTx")
	proto.RegisterType((*MsgEthereumTxResponse)(nil), "artela.evm.v1.MsgEthereumTxResponse")
	proto.RegisterType((*MsgUpdateParams)(nil), "artela.evm.v1.MsgUpdateParams")
	proto.RegisterType((*MsgUpdateParamsResponse)(nil), "artela.evm.v1.MsgUpdateParamsResponse")
}

func init() { proto.RegisterFile("artela/evm/v1/txs.proto", fileDescriptor_3c43c0836c37bbe6) }

var fileDescriptor_3c43c0836c37bbe6 = []byte{
	// 997 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x56, 0xcf, 0x6f, 0xe3, 0x44,
	0x14, 0x8e, 0x13, 0xe7, 0xd7, 0x24, 0x5b, 0x60, 0xe8, 0x52, 0x27, 0xaa, 0xe2, 0xc8, 0x42, 0x55,
	0xb4, 0x52, 0x6c, 0xb5, 0x8b, 0x38, 0xf4, 0x44, 0xb3, 0xed, 0x56, 0x5d, 0xb5, 0x62, 0x65, 0xb2,
	0x08, 0xc1, 0x21, 0x9a, 0x3a, 0x53, 0xc7, 0xda, 0xd8, 0x63, 0x79, 0xc6, 0x26, 0xe1, 0xb8, 0x27,
	0x4e, 0x80, 0xc4, 0x3f, 0xc0, 0x99, 0x13, 0x12, 0x7b, 0xe5, 0xbe, 0xe2, 0xb4, 0xb0, 0x17, 0xc4,
	0x21, 0xa0, 0x14, 0x09, 0xa9, 0x07, 0x0e, 0xfc, 0x05, 0x68, 0x66, 0x9c, 0xa6, 0x49, 0xd5, 0x0a,
	0x96, 0x4a, 0x7b, 0xca, 0x3c, 0x7f, 0x6f, 0xbe, 0x79, 0xf3, 0xbe, 0x6f, 0x66, 0x02, 0xde, 0x42,
	0x11, 0xc3, 0x43, 0x64, 0xe1, 0xc4, 0xb7, 0x92, 0x4d, 0x8b, 0x8d, 0xcc, 0x30, 0x22, 0x8c, 0xc0,
	0x5b, 0xf2, 0xbb, 0x89, 0x13, 0xdf, 0x4c, 0x36, 0xeb, 0x6b, 0x0e, 0xa1, 0x3e, 0xa1, 0x96, 0x4f,
	0x5d, 0x9e, 0xe6, 0x53, 0x57, 0xe6, 0xd5, 0x6b, 0x12, 0xe8, 0x89, 0xc8, 0x92, 0x41, 0x0a, 0xad,
	0x2d, 0x52, 0x73, 0x26, 0x09, 0xac, 0xba, 0xc4, 0x25, 0x72, 0x02, 0x1f, 0xa5, 0x5f, 0xd7, 0x5d,
	0x42, 0xdc, 0x21, 0xb6, 0x50, 0xe8, 0x59, 0x28, 0x08, 0x08, 0x43, 0xcc, 0x23, 0xc1, 0x8c, 0xac,
	0x96, 0xa2, 0x22, 0x3a, 0x8e, 0x4f, 0x2c, 0x14, 0x8c, 0x25, 0x64, 0x7c, 0xa9, 0x80, 0x5b, 0x47,
	0xd4, 0xdd, 0x63, 0x03, 0x1c, 0xe1, 0xd8, 0xef, 0x8e, 0x60, 0x0b, 0xa8, 0x7d, 0xc4, 0x90, 0xa6,
	0x34, 0x95, 0x56, 0x65, 0x6b, 0xd5, 0x94, 0x73, 0xcd, 0xd9, 0x5c, 0x73, 0x27, 0x18, 0xdb, 0x22,
	0x03, 0xd6, 0x80, 0x4a, 0xbd, 0xcf, 0xb0, 0x96, 0x6d, 0x2a, 0x2d, 0xa5, 0x93, 0x3f, 0x9b, 0xe8,
	0x4a, 0xdb, 0x16, 0x9f, 0xa0, 0x0e, 0xd4, 0x01, 0xa2, 0x03, 0x2d, 0xd7, 0x54, 0x5a, 0xe5, 0x4e,
	0xe5, 0xef, 0x89, 0x5e, 0x8c, 0x86, 0xe1, 0xb6, 0xd1, 0x36, 0x6c, 0x01, 0x40, 0x08, 0xd4, 0x93,
	0x88, 0xf8, 0x9a, 0xca, 0x13, 0x6c, 0x31, 0xde, 0x56, 0x3f, 0xff, 0x46, 0xcf, 0x18, 0xdf, 0x67,
	0x41, 0xe9, 0x10, 0xbb, 0xc8, 0x19, 0x77, 0x47, 0x70, 0x15, 0xe4, 0x03, 0x12, 0x38, 0x58, 0x54,
	0xa3, 0xda, 0x32, 0x80, 0xfb, 0xa0, 0xec, 0x22, 0xde, 0x36, 0xcf, 0x91, 0xab, 0x97, 0x3b, 0x77,
	0x7e, 0x9d, 0xe8, 0x1b, 0xae, 0xc7, 0x06, 0xf1, 0xb1, 0xe9, 0x10, 0x3f, 0x6d, 0x66, 0xfa, 0xd3,
	0xa6, 0xfd, 0xc7, 0x16, 0x1b, 0x87, 0x98, 0x9a, 0x07, 0x01, 0xb3, 0x4b, 0x2e, 0xa2, 0x0f, 0xf9,
	0x5c, 0xd8, 0x00, 0x39, 0x17, 0x51, 0x51, 0xa5, 0xda, 0xa9, 0x4e, 0x27, 0x7a, 0x69, 0x1f, 0xd1,
	0x43, 0xcf, 0xf7, 0x98, 0xcd, 0x01, 0xb8, 0x02, 0xb2, 0x8c, 0xa4, 0x35, 0x66, 0x19, 0x81, 0x0f,
	0x40, 0x3e, 0x41, 0xc3, 0x18, 0x6b, 0x79, 0xb1, 0xe8, 0x3b, 0xff, 0x7e, 0xd1, 0xe9, 0x44, 0x2f,
	0xec, 0xf8, 0x24, 0x0e, 0x98, 0x2d, 0x29, 0x78, 0x07, 0x44, 0x9f, 0x0b, 0x4d, 0xa5, 0x55, 0x4d,
	0x3b, 0x5a, 0x05, 0x4a, 0xa2, 0x15, 0xc5, 0x07, 0x25, 0xe1, 0x51, 0xa4, 0x95, 0x64, 0x14, 0xf1,
	0x88, 0x6a, 0x65, 0x19, 0xd1, 0xed, 0x15, 0xde, 0xab, 0x1f, 0x9f, 0xb6, 0x0b, 0xdd, 0xd1, 0x2e,
	0x62, 0xc8, 0xf8, 0x2b, 0x07, 0xaa, 0x3b, 0x8e, 0x83, 0x29, 0x3d, 0xf4, 0x28, 0xeb, 0x8e, 0xe0,
	0x27, 0xa0, 0xe4, 0x0c, 0x90, 0x17, 0xf4, 0xbc, 0xbe, 0x68, 0x5e, 0xb9, 0xf3, 0xde, 0x7f, 0xaa,
	0xb6, 0x78, 0x8f, 0xcf, 0x3e, 0xd8, 0x3d, 0x9b, 0xe8, 0x45, 0x47, 0x0e, 0xed, 0x74, 0xd0, 0x9f,
	0xcb, 0x92, 0xbd, 0x52, 0x96, 0xdc, 0xff, 0x97, 0x45, 0xbd, 0x5e, 0x96, 0xfc, 0x65, 0x59, 0x0a,
	0x37, 0x27, 0x4b, 0xf1, 0x82, 0x2c, 0x1f, 0x81, 0x12, 0x12, 0xbd, 0xc5, 0x54, 0x2b, 0x35, 0x73,
	0xad, 0xca, 0x56, 0xdd, 0x5c, 0x38, 0xe2, 0xa6, 0x6c, 0x7d, 0x37, 0x0e, 0x87, 0xb8, 0xd3, 0x7c,
	0x36, 0xd1, 0x33, 0x67, 0x13, 0x1d, 0xa0, 0x73, 0x3d, 0xbe, 0xfd, 0x4d, 0x07, 0x73, 0x75, 0xec,
	0x73, 0x36, 0x29, 0x78, 0x79, 0x41, 0x70, 0xb0, 0x20, 0x78, 0xe5, 0x2a, 0xc1, 0x7f, 0x50, 0x41,
	0x75, 0x77, 0x1c, 0x20, 0xdf, 0x73, 0xee, 0x63, 0xfc, 0x6a, 0x04, 0x7f, 0x00, 0x2a, 0x5c, 0x70,
	0xe6, 0x85, 0x3d, 0x07, 0x85, 0x2f, 0x21, 0x39, 0xf7, 0x4b, 0xd7, 0x0b, 0xef, 0xa1, 0x70, 0xc6,
	0x75, 0x82, 0xb1, 0xe0, 0x52, 0x5f, 0x8a, 0xeb, 0x3e, 0xc6, 0x9c, 0x2b, 0xf5, 0x4f, 0xfe, 0x7a,
	0xff, 0x14, 0x2e, 0xfb, 0xa7, 0x78, 0x73, 0xfe, 0x29, 0x5d, 0xe1, 0x9f, 0xf2, 0xcd, 0xfb, 0x07,
	0x2c, 0xf8, 0xa7, 0xb2, 0xe0, 0x9f, 0xea, 0x55, 0xfe, 0x31, 0x40, 0x7d, 0x6f, 0xc4, 0x70, 0x40,
	0x3d, 0x12, 0xbc, 0x1f, 0x8a, 0xd7, 0x62, 0xfe, 0x08, 0xa4, 0x57, 0xf1, 0x4f, 0x0a, 0xb8, 0xbd,
	0xf0, 0x38, 0xd8, 0x98, 0x86, 0x24, 0xa0, 0x62, 0x97, 0xe2, 0x7e, 0x57, 0xe4, 0xf5, 0x2d, 0xae,
	0xf4, 0x0d, 0xa0, 0x0e, 0x89, 0x4b, 0xb5, 0xac, 0xd8, 0x21, 0x5c, 0xda, 0xe1, 0x21, 0x71, 0x6d,
	0x81, 0xc3, 0xd7, 0x41, 0x2e, 0xc2, 0x4c, 0xb8, 0xa5, 0x6a, 0xf3, 0x21, 0xac, 0x81, 0x52, 0xe2,
	0xf7, 0x70, 0x14, 0x91, 0x28, 0xbd, 0x6c, 0x8b, 0x89, 0xbf, 0xc7, 0x43, 0x0e, 0x71, 0x5b, 0xc4,
	0x14, 0xf7, 0xa5, 0x9e, 0x76, 0xd1, 0x45, 0xf4, 0x11, 0xc5, 0x7d, 0x68, 0x82, 0x37, 0x9d, 0xd8,
	0x8f, 0x87, 0x88, 0x79, 0x09, 0xee, 0x9d, 0x67, 0x15, 0x44, 0xd6, 0x1b, 0x73, 0x68, 0x5f, 0xe6,
	0xa7, 0x7b, 0xfa, 0x42, 0x01, 0xaf, 0x1d, 0x51, 0xf7, 0x51, 0xd8, 0x47, 0x0c, 0x3f, 0x44, 0x11,
	0xf2, 0x29, 0x7c, 0x17, 0x94, 0x51, 0xcc, 0x06, 0x24, 0xf2, 0xd8, 0x38, 0x3d, 0x3b, 0xda, 0xcf,
	0x4f, 0xdb, 0xab, 0xe9, 0x8b, 0xbc, 0xd3, 0xef, 0x47, 0x98, 0xd2, 0x0f, 0x58, 0xe4, 0x05, 0xae,
	0x3d, 0x4f, 0x85, 0x77, 0x41, 0x21, 0x14, 0x0c, 0xe2, 0x58, 0x54, 0xb6, 0x6e, 0x2f, 0xed, 0x59,
	0xd2, 0x77, 0x54, 0x2e, 0xa8, 0x9d, 0xa6, 0x6e, 0xaf, 0x3c, 0xf9, 0xf3, 0xbb, 0x3b, 0x73, 0x12,
	0xa3, 0x06, 0xd6, 0x96, 0xea, 0x99, 0x75, 0x79, 0xeb, 0x85, 0x02, 0x72, 0x47, 0xd4, 0x85, 0x0c,
	0x80, 0x0b, 0x0f, 0xf4, 0xfa, 0xd2, 0x2a, 0x0b, 0x0a, 0xd5, 0xdf, 0xbe, 0x0e, 0x9d, 0x31, 0x1b,
	0xc6, 0x93, 0x17, 0x7f, 0x7c, 0x9d, 0x5d, 0x37, 0xea, 0xd6, 0xd2, 0xff, 0x8c, 0x34, 0xb5, 0xc7,
	0x46, 0xf0, 0x43, 0x50, 0x5d, 0xe8, 0x52, 0xe3, 0x32, 0xf3, 0x45, 0xbc, 0xbe, 0x71, 0x3d, 0x3e,
	0x5b, 0xbb, 0x73, 0xf0, 0x6c, 0xda, 0x50, 0x9e, 0x4f, 0x1b, 0xca, 0xef, 0xd3, 0x86, 0xf2, 0xd5,
	0x69, 0x23, 0xf3, 0xfc, 0xb4, 0x91, 0xf9, 0xe5, 0xb4, 0x91, 0xf9, 0xd8, 0xba, 0x70, 0xe8, 0x24,
	0x57, 0x3b, 0xc0, 0xec, 0x53, 0x12, 0x3d, 0x9e, 0x95, 0x99, 0x6c, 0x5a, 0x23, 0x51, 0xab, 0x38,
	0x81, 0xc7, 0x05, 0xf1, 0xaf, 0xe4, 0xee, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9b, 0xec, 0xd3,
	0xde, 0x89, 0x09, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this support file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// EthereumTx defines a method submitting Ethereum transactions.
	EthereumTx(ctx context.Context, in *MsgEthereumTx, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error)
	// UpdateParams defined a governance operation for updating the x/evm module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) EthereumTx(ctx context.Context, in *MsgEthereumTx, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error) {
	out := new(MsgEthereumTxResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Msg/EthereumTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// EthereumTx defines a method submitting Ethereum transactions.
	EthereumTx(context.Context, *MsgEthereumTx) (*MsgEthereumTxResponse, error)
	// UpdateParams defined a governance operation for updating the x/evm module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) EthereumTx(ctx context.Context, req *MsgEthereumTx) (*MsgEthereumTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EthereumTx not implemented")
}
func (*UnimplementedMsgServer) UpdateParams(ctx context.Context, req *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_EthereumTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEthereumTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EthereumTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Msg/EthereumTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EthereumTx(ctx, req.(*MsgEthereumTx))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "artela.evm.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EthereumTx",
			Handler:    _Msg_EthereumTx_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "artela/evm/v1/txs.proto",
}

func (msg *MsgEthereumTx) Marshal() (dAtA []byte, err error) {
	size := msg.Size()
	dAtA = make([]byte, size)
	n, err := msg.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (msg *MsgEthereumTx) MarshalTo(dAtA []byte) (int, error) {
	size := msg.Size()
	return msg.MarshalToSizedBuffer(dAtA[:size])
}

func (msg *MsgEthereumTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(msg.From) > 0 {
		i -= len(msg.From)
		copy(dAtA[i:], msg.From)
		i = encodeVarintTx(dAtA, i, uint64(len(msg.From)))
		i--
		dAtA[i] = 0x22
	}
	if len(msg.Hash) > 0 {
		i -= len(msg.Hash)
		copy(dAtA[i:], msg.Hash)
		i = encodeVarintTx(dAtA, i, uint64(len(msg.Hash)))
		i--
		dAtA[i] = 0x1a
	}
	if msg.Size_ != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(msg.Size_))))
		i--
		dAtA[i] = 0x11
	}
	if msg.Data != nil {
		{
			size, err := msg.Data.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (tx *LegacyTx) Marshal() (dAtA []byte, err error) {
	size := tx.Size()
	dAtA = make([]byte, size)
	n, err := tx.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (tx *LegacyTx) MarshalTo(dAtA []byte) (int, error) {
	size := tx.Size()
	return tx.MarshalToSizedBuffer(dAtA[:size])
}

func (tx *LegacyTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(tx.S) > 0 {
		i -= len(tx.S)
		copy(dAtA[i:], tx.S)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.S)))
		i--
		dAtA[i] = 0x4a
	}
	if len(tx.R) > 0 {
		i -= len(tx.R)
		copy(dAtA[i:], tx.R)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.R)))
		i--
		dAtA[i] = 0x42
	}
	if len(tx.V) > 0 {
		i -= len(tx.V)
		copy(dAtA[i:], tx.V)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.V)))
		i--
		dAtA[i] = 0x3a
	}
	if len(tx.Data) > 0 {
		i -= len(tx.Data)
		copy(dAtA[i:], tx.Data)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.Data)))
		i--
		dAtA[i] = 0x32
	}
	if tx.Amount != nil {
		{
			size := tx.Amount.Size()
			i -= size
			if _, err := tx.Amount.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if len(tx.To) > 0 {
		i -= len(tx.To)
		copy(dAtA[i:], tx.To)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.To)))
		i--
		dAtA[i] = 0x22
	}
	if tx.GasLimit != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.GasLimit))
		i--
		dAtA[i] = 0x18
	}
	if tx.GasPrice != nil {
		{
			size := tx.GasPrice.Size()
			i -= size
			if _, err := tx.GasPrice.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if tx.Nonce != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.Nonce))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (tx *AccessListTx) Marshal() (dAtA []byte, err error) {
	size := tx.Size()
	dAtA = make([]byte, size)
	n, err := tx.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (tx *AccessListTx) MarshalTo(dAtA []byte) (int, error) {
	size := tx.Size()
	return tx.MarshalToSizedBuffer(dAtA[:size])
}

func (tx *AccessListTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(tx.S) > 0 {
		i -= len(tx.S)
		copy(dAtA[i:], tx.S)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.S)))
		i--
		dAtA[i] = 0x5a
	}
	if len(tx.R) > 0 {
		i -= len(tx.R)
		copy(dAtA[i:], tx.R)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.R)))
		i--
		dAtA[i] = 0x52
	}
	if len(tx.V) > 0 {
		i -= len(tx.V)
		copy(dAtA[i:], tx.V)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.V)))
		i--
		dAtA[i] = 0x4a
	}
	if len(tx.Accesses) > 0 {
		for iNdEx := len(tx.Accesses) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := tx.Accesses[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(tx.Data) > 0 {
		i -= len(tx.Data)
		copy(dAtA[i:], tx.Data)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.Data)))
		i--
		dAtA[i] = 0x3a
	}
	if tx.Amount != nil {
		{
			size := tx.Amount.Size()
			i -= size
			if _, err := tx.Amount.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if len(tx.To) > 0 {
		i -= len(tx.To)
		copy(dAtA[i:], tx.To)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.To)))
		i--
		dAtA[i] = 0x2a
	}
	if tx.GasLimit != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.GasLimit))
		i--
		dAtA[i] = 0x20
	}
	if tx.GasPrice != nil {
		{
			size := tx.GasPrice.Size()
			i -= size
			if _, err := tx.GasPrice.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if tx.Nonce != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.Nonce))
		i--
		dAtA[i] = 0x10
	}
	if tx.ChainID != nil {
		{
			size := tx.ChainID.Size()
			i -= size
			if _, err := tx.ChainID.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (tx *DynamicFeeTx) Marshal() (dAtA []byte, err error) {
	size := tx.Size()
	dAtA = make([]byte, size)
	n, err := tx.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (tx *DynamicFeeTx) MarshalTo(dAtA []byte) (int, error) {
	size := tx.Size()
	return tx.MarshalToSizedBuffer(dAtA[:size])
}

func (tx *DynamicFeeTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(tx.S) > 0 {
		i -= len(tx.S)
		copy(dAtA[i:], tx.S)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.S)))
		i--
		dAtA[i] = 0x62
	}
	if len(tx.R) > 0 {
		i -= len(tx.R)
		copy(dAtA[i:], tx.R)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.R)))
		i--
		dAtA[i] = 0x5a
	}
	if len(tx.V) > 0 {
		i -= len(tx.V)
		copy(dAtA[i:], tx.V)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.V)))
		i--
		dAtA[i] = 0x52
	}
	if len(tx.Accesses) > 0 {
		for iNdEx := len(tx.Accesses) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := tx.Accesses[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(tx.Data) > 0 {
		i -= len(tx.Data)
		copy(dAtA[i:], tx.Data)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.Data)))
		i--
		dAtA[i] = 0x42
	}
	if tx.Amount != nil {
		{
			size := tx.Amount.Size()
			i -= size
			if _, err := tx.Amount.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if len(tx.To) > 0 {
		i -= len(tx.To)
		copy(dAtA[i:], tx.To)
		i = encodeVarintTx(dAtA, i, uint64(len(tx.To)))
		i--
		dAtA[i] = 0x32
	}
	if tx.GasLimit != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.GasLimit))
		i--
		dAtA[i] = 0x28
	}
	if tx.GasFeeCap != nil {
		{
			size := tx.GasFeeCap.Size()
			i -= size
			if _, err := tx.GasFeeCap.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if tx.GasTipCap != nil {
		{
			size := tx.GasTipCap.Size()
			i -= size
			if _, err := tx.GasTipCap.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if tx.Nonce != 0 {
		i = encodeVarintTx(dAtA, i, uint64(tx.Nonce))
		i--
		dAtA[i] = 0x10
	}
	if tx.ChainID != nil {
		{
			size := tx.ChainID.Size()
			i -= size
			if _, err := tx.ChainID.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ExtensionOptionsEthereumTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtensionOptionsEthereumTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExtensionOptionsEthereumTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgEthereumTxResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgEthereumTxResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgEthereumTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CumulativeGasUsed != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.CumulativeGasUsed))
		i--
		dAtA[i] = 0x30
	}
	if m.GasUsed != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.GasUsed))
		i--
		dAtA[i] = 0x28
	}
	if len(m.VmError) > 0 {
		i -= len(m.VmError)
		copy(dAtA[i:], m.VmError)
		i = encodeVarintTx(dAtA, i, uint64(len(m.VmError)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Ret) > 0 {
		i -= len(m.Ret)
		copy(dAtA[i:], m.Ret)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Ret)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Logs) > 0 {
		for iNdEx := len(m.Logs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Logs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (msg *MsgEthereumTx) Size() (n int) {
	if msg == nil {
		return 0
	}
	var l int
	_ = l
	if msg.Data != nil {
		l = msg.Data.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if msg.Size_ != 0 {
		n += 9
	}
	l = len(msg.Hash)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(msg.From)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (tx *LegacyTx) Size() (n int) {
	if tx == nil {
		return 0
	}
	var l int
	_ = l
	if tx.Nonce != 0 {
		n += 1 + sovTx(uint64(tx.Nonce))
	}
	if tx.GasPrice != nil {
		l = tx.GasPrice.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.GasLimit != 0 {
		n += 1 + sovTx(uint64(tx.GasLimit))
	}
	l = len(tx.To)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.Amount != nil {
		l = tx.Amount.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.Data)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.V)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.R)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.S)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (tx *AccessListTx) Size() (n int) {
	if tx == nil {
		return 0
	}
	var l int
	_ = l
	if tx.ChainID != nil {
		l = tx.ChainID.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.Nonce != 0 {
		n += 1 + sovTx(uint64(tx.Nonce))
	}
	if tx.GasPrice != nil {
		l = tx.GasPrice.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.GasLimit != 0 {
		n += 1 + sovTx(uint64(tx.GasLimit))
	}
	l = len(tx.To)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.Amount != nil {
		l = tx.Amount.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.Data)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(tx.Accesses) > 0 {
		for _, e := range tx.Accesses {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(tx.V)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.R)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.S)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (tx *DynamicFeeTx) Size() (n int) {
	if tx == nil {
		return 0
	}
	var l int
	_ = l
	if tx.ChainID != nil {
		l = tx.ChainID.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.Nonce != 0 {
		n += 1 + sovTx(uint64(tx.Nonce))
	}
	if tx.GasTipCap != nil {
		l = tx.GasTipCap.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.GasFeeCap != nil {
		l = tx.GasFeeCap.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.GasLimit != 0 {
		n += 1 + sovTx(uint64(tx.GasLimit))
	}
	l = len(tx.To)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if tx.Amount != nil {
		l = tx.Amount.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.Data)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(tx.Accesses) > 0 {
		for _, e := range tx.Accesses {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(tx.V)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.R)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(tx.S)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *ExtensionOptionsEthereumTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgEthereumTxResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.Logs) > 0 {
		for _, e := range m.Logs {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Ret)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.VmError)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.GasUsed != 0 {
		n += 1 + sovTx(uint64(m.GasUsed))
	}
	if m.CumulativeGasUsed != 0 {
		n += 1 + sovTx(uint64(m.CumulativeGasUsed))
	}
	return n
}

func (m *MsgUpdateParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgUpdateParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (msg *MsgEthereumTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgEthereumTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEthereumTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if msg.Data == nil {
				msg.Data = &types.Any{}
			}
			if err := msg.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			msg.Size_ = float64(math.Float64frombits(v))
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			msg.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			msg.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (tx *LegacyTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: LegacyTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			tx.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.GasPrice = &v
			if err := tx.GasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasLimit", wireType)
			}
			tx.GasLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.GasLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.To = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.Amount = &v
			if err := tx.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.Data = append(tx.Data[:0], dAtA[iNdEx:postIndex]...)
			if tx.Data == nil {
				tx.Data = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.V = append(tx.V[:0], dAtA[iNdEx:postIndex]...)
			if tx.V == nil {
				tx.V = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field R", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.R = append(tx.R[:0], dAtA[iNdEx:postIndex]...)
			if tx.R == nil {
				tx.R = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field S", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.S = append(tx.S[:0], dAtA[iNdEx:postIndex]...)
			if tx.S == nil {
				tx.S = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (tx *AccessListTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: AccessListTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccessListTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.ChainID = &v
			if err := tx.ChainID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			tx.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.GasPrice = &v
			if err := tx.GasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasLimit", wireType)
			}
			tx.GasLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.GasLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.To = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.Amount = &v
			if err := tx.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.Data = append(tx.Data[:0], dAtA[iNdEx:postIndex]...)
			if tx.Data == nil {
				tx.Data = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Accesses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.Accesses = append(tx.Accesses, support.AccessTuple{})
			if err := tx.Accesses[len(tx.Accesses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.V = append(tx.V[:0], dAtA[iNdEx:postIndex]...)
			if tx.V == nil {
				tx.V = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field R", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.R = append(tx.R[:0], dAtA[iNdEx:postIndex]...)
			if tx.R == nil {
				tx.R = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field S", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.S = append(tx.S[:0], dAtA[iNdEx:postIndex]...)
			if tx.S == nil {
				tx.S = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (tx *DynamicFeeTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: DynamicFeeTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DynamicFeeTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.ChainID = &v
			if err := tx.ChainID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			tx.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasTipCap", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.GasTipCap = &v
			if err := tx.GasTipCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasFeeCap", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.GasFeeCap = &v
			if err := tx.GasFeeCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasLimit", wireType)
			}
			tx.GasLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				tx.GasLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.To = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			tx.Amount = &v
			if err := tx.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.Data = append(tx.Data[:0], dAtA[iNdEx:postIndex]...)
			if tx.Data == nil {
				tx.Data = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Accesses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.Accesses = append(tx.Accesses, support.AccessTuple{})
			if err := tx.Accesses[len(tx.Accesses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.V = append(tx.V[:0], dAtA[iNdEx:postIndex]...)
			if tx.V == nil {
				tx.V = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field R", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.R = append(tx.R[:0], dAtA[iNdEx:postIndex]...)
			if tx.R == nil {
				tx.R = []byte{}
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field S", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			tx.S = append(tx.S[:0], dAtA[iNdEx:postIndex]...)
			if tx.S == nil {
				tx.S = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *ExtensionOptionsEthereumTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: ExtensionOptionsEthereumTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtensionOptionsEthereumTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgEthereumTxResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgEthereumTxResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEthereumTxResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Logs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Logs = append(m.Logs, &support.Log{})
			if err := m.Logs[len(m.Logs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ret", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ret = append(m.Ret[:0], dAtA[iNdEx:postIndex]...)
			if m.Ret == nil {
				m.Ret = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VmError", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VmError = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasUsed", wireType)
			}
			m.GasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumulativeGasUsed", wireType)
			}
			m.CumulativeGasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CumulativeGasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgUpdateParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgUpdateParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
