// Code support by protoc-gen-gogo. DO NOT EDIT.
// source: artela/evm/v1/query.proto

package process

import (
	context "context"
	fmt "fmt"
	"github.com/artela-network/artela/x/evm/process/support"

	//"github.com/artela-network/artela/x/evm/process"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/types"
	github_com_gogo_protobuf_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this support file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryAccountRequest is the request type for the Query/Account RPC method.
type QueryAccountRequest struct {
	// address is the ethereum hex address to query the account for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryAccountRequest) Reset()         { *m = QueryAccountRequest{} }
func (m *QueryAccountRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAccountRequest) ProtoMessage()    {}
func (*QueryAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{0}
}
func (m *QueryAccountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAccountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAccountRequest.Merge(m, src)
}
func (m *QueryAccountRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAccountRequest proto.InternalMessageInfo

// QueryAccountResponse is the response type for the Query/Account RPC method.
type QueryAccountResponse struct {
	// balance is the balance of the EVM denomination.
	Balance string `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
	// code_hash is the hex-formatted code bytes from the EOA.
	CodeHash string `protobuf:"bytes,2,opt,name=code_hash,json=codeHash,proto3" json:"code_hash,omitempty"`
	// nonce is the account's sequence number.
	Nonce uint64 `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *QueryAccountResponse) Reset()         { *m = QueryAccountResponse{} }
func (m *QueryAccountResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAccountResponse) ProtoMessage()    {}
func (*QueryAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{1}
}
func (m *QueryAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAccountResponse.Merge(m, src)
}
func (m *QueryAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAccountResponse proto.InternalMessageInfo

func (m *QueryAccountResponse) GetBalance() string {
	if m != nil {
		return m.Balance
	}
	return ""
}

func (m *QueryAccountResponse) GetCodeHash() string {
	if m != nil {
		return m.CodeHash
	}
	return ""
}

func (m *QueryAccountResponse) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

// QueryCosmosAccountRequest is the request type for the Query/CosmosAccount RPC
// method.
type QueryCosmosAccountRequest struct {
	// address is the ethereum hex address to query the account for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryCosmosAccountRequest) Reset()         { *m = QueryCosmosAccountRequest{} }
func (m *QueryCosmosAccountRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCosmosAccountRequest) ProtoMessage()    {}
func (*QueryCosmosAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{2}
}
func (m *QueryCosmosAccountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCosmosAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCosmosAccountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCosmosAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCosmosAccountRequest.Merge(m, src)
}
func (m *QueryCosmosAccountRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCosmosAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCosmosAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCosmosAccountRequest proto.InternalMessageInfo

// QueryCosmosAccountResponse is the response type for the Query/CosmosAccount
// RPC method.
type QueryCosmosAccountResponse struct {
	// cosmos_address is the cosmos address of the account.
	CosmosAddress string `protobuf:"bytes,1,opt,name=cosmos_address,json=cosmosAddress,proto3" json:"cosmos_address,omitempty"`
	// sequence is the account's sequence number.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// account_number is the account number
	AccountNumber uint64 `protobuf:"varint,3,opt,name=account_number,json=accountNumber,proto3" json:"account_number,omitempty"`
}

func (m *QueryCosmosAccountResponse) Reset()         { *m = QueryCosmosAccountResponse{} }
func (m *QueryCosmosAccountResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCosmosAccountResponse) ProtoMessage()    {}
func (*QueryCosmosAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{3}
}
func (m *QueryCosmosAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCosmosAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCosmosAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCosmosAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCosmosAccountResponse.Merge(m, src)
}
func (m *QueryCosmosAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCosmosAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCosmosAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCosmosAccountResponse proto.InternalMessageInfo

func (m *QueryCosmosAccountResponse) GetCosmosAddress() string {
	if m != nil {
		return m.CosmosAddress
	}
	return ""
}

func (m *QueryCosmosAccountResponse) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *QueryCosmosAccountResponse) GetAccountNumber() uint64 {
	if m != nil {
		return m.AccountNumber
	}
	return 0
}

// QueryValidatorAccountRequest is the request type for the
// Query/ValidatorAccount RPC method.
type QueryValidatorAccountRequest struct {
	// cons_address is the validator cons address to query the account for.
	ConsAddress string `protobuf:"bytes,1,opt,name=cons_address,json=consAddress,proto3" json:"cons_address,omitempty"`
}

func (m *QueryValidatorAccountRequest) Reset()         { *m = QueryValidatorAccountRequest{} }
func (m *QueryValidatorAccountRequest) String() string { return proto.CompactTextString(m) }
func (*QueryValidatorAccountRequest) ProtoMessage()    {}
func (*QueryValidatorAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{4}
}
func (m *QueryValidatorAccountRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryValidatorAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryValidatorAccountRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryValidatorAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryValidatorAccountRequest.Merge(m, src)
}
func (m *QueryValidatorAccountRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryValidatorAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryValidatorAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryValidatorAccountRequest proto.InternalMessageInfo

// QueryValidatorAccountResponse is the response type for the
// Query/ValidatorAccount RPC method.
type QueryValidatorAccountResponse struct {
	// account_address is the cosmos address of the account in bech32 format.
	AccountAddress string `protobuf:"bytes,1,opt,name=account_address,json=accountAddress,proto3" json:"account_address,omitempty"`
	// sequence is the account's sequence number.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// account_number is the account number
	AccountNumber uint64 `protobuf:"varint,3,opt,name=account_number,json=accountNumber,proto3" json:"account_number,omitempty"`
}

func (m *QueryValidatorAccountResponse) Reset()         { *m = QueryValidatorAccountResponse{} }
func (m *QueryValidatorAccountResponse) String() string { return proto.CompactTextString(m) }
func (*QueryValidatorAccountResponse) ProtoMessage()    {}
func (*QueryValidatorAccountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{5}
}
func (m *QueryValidatorAccountResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryValidatorAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryValidatorAccountResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryValidatorAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryValidatorAccountResponse.Merge(m, src)
}
func (m *QueryValidatorAccountResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryValidatorAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryValidatorAccountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryValidatorAccountResponse proto.InternalMessageInfo

func (m *QueryValidatorAccountResponse) GetAccountAddress() string {
	if m != nil {
		return m.AccountAddress
	}
	return ""
}

func (m *QueryValidatorAccountResponse) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *QueryValidatorAccountResponse) GetAccountNumber() uint64 {
	if m != nil {
		return m.AccountNumber
	}
	return 0
}

// QueryBalanceRequest is the request type for the Query/Balance RPC method.
type QueryBalanceRequest struct {
	// address is the ethereum hex address to query the balance for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryBalanceRequest) Reset()         { *m = QueryBalanceRequest{} }
func (m *QueryBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBalanceRequest) ProtoMessage()    {}
func (*QueryBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{6}
}
func (m *QueryBalanceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBalanceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBalanceRequest.Merge(m, src)
}
func (m *QueryBalanceRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBalanceRequest proto.InternalMessageInfo

// QueryBalanceResponse is the response type for the Query/Balance RPC method.
type QueryBalanceResponse struct {
	// balance is the balance of the EVM denomination.
	Balance string `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (m *QueryBalanceResponse) Reset()         { *m = QueryBalanceResponse{} }
func (m *QueryBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBalanceResponse) ProtoMessage()    {}
func (*QueryBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{7}
}
func (m *QueryBalanceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBalanceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBalanceResponse.Merge(m, src)
}
func (m *QueryBalanceResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBalanceResponse proto.InternalMessageInfo

func (m *QueryBalanceResponse) GetBalance() string {
	if m != nil {
		return m.Balance
	}
	return ""
}

// QueryStorageRequest is the request type for the Query/Storage RPC method.
type QueryStorageRequest struct {
	// address is the ethereum hex address to query the storage state for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// key defines the key of the storage state
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *QueryStorageRequest) Reset()         { *m = QueryStorageRequest{} }
func (m *QueryStorageRequest) String() string { return proto.CompactTextString(m) }
func (*QueryStorageRequest) ProtoMessage()    {}
func (*QueryStorageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{8}
}
func (m *QueryStorageRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryStorageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryStorageRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryStorageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStorageRequest.Merge(m, src)
}
func (m *QueryStorageRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryStorageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStorageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStorageRequest proto.InternalMessageInfo

// QueryStorageResponse is the response type for the Query/Storage RPC
// method.
type QueryStorageResponse struct {
	// value defines the storage state value hash associated with the given key.
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *QueryStorageResponse) Reset()         { *m = QueryStorageResponse{} }
func (m *QueryStorageResponse) String() string { return proto.CompactTextString(m) }
func (*QueryStorageResponse) ProtoMessage()    {}
func (*QueryStorageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{9}
}
func (m *QueryStorageResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryStorageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryStorageResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryStorageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStorageResponse.Merge(m, src)
}
func (m *QueryStorageResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryStorageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStorageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStorageResponse proto.InternalMessageInfo

func (m *QueryStorageResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// QueryCodeRequest is the request type for the Query/Code RPC method.
type QueryCodeRequest struct {
	// address is the ethereum hex address to query the code for.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryCodeRequest) Reset()         { *m = QueryCodeRequest{} }
func (m *QueryCodeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCodeRequest) ProtoMessage()    {}
func (*QueryCodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{10}
}
func (m *QueryCodeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCodeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCodeRequest.Merge(m, src)
}
func (m *QueryCodeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCodeRequest proto.InternalMessageInfo

// QueryCodeResponse is the response type for the Query/Code RPC
// method.
type QueryCodeResponse struct {
	// code represents the code bytes from an ethereum address.
	Code []byte `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (m *QueryCodeResponse) Reset()         { *m = QueryCodeResponse{} }
func (m *QueryCodeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCodeResponse) ProtoMessage()    {}
func (*QueryCodeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{11}
}
func (m *QueryCodeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCodeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCodeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCodeResponse.Merge(m, src)
}
func (m *QueryCodeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCodeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCodeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCodeResponse proto.InternalMessageInfo

func (m *QueryCodeResponse) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

// QueryTxLogsRequest is the request type for the Query/TxLogs RPC method.
type QueryTxLogsRequest struct {
	// hash is the ethereum process hex hash to query the logs for.
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryTxLogsRequest) Reset()         { *m = QueryTxLogsRequest{} }
func (m *QueryTxLogsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTxLogsRequest) ProtoMessage()    {}
func (*QueryTxLogsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{12}
}
func (m *QueryTxLogsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTxLogsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTxLogsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTxLogsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTxLogsRequest.Merge(m, src)
}
func (m *QueryTxLogsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTxLogsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTxLogsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTxLogsRequest proto.InternalMessageInfo

// QueryTxLogsResponse is the response type for the Query/TxLogs RPC method.
type QueryTxLogsResponse struct {
	// logs represents the ethereum logs support from the given process.
	Logs []*support.Log `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryTxLogsResponse) Reset()         { *m = QueryTxLogsResponse{} }
func (m *QueryTxLogsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTxLogsResponse) ProtoMessage()    {}
func (*QueryTxLogsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{13}
}
func (m *QueryTxLogsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTxLogsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTxLogsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTxLogsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTxLogsResponse.Merge(m, src)
}
func (m *QueryTxLogsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTxLogsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTxLogsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTxLogsResponse proto.InternalMessageInfo

func (m *QueryTxLogsResponse) GetLogs() []*support.Log {
	if m != nil {
		return m.Logs
	}
	return nil
}

func (m *QueryTxLogsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryParamsRequest defines the request type for querying x/evm parameters.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{14}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse defines the response type for querying x/evm parameters.
type QueryParamsResponse struct {
	// params define the evm module parameters.
	Params support.Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{15}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() support.Params {
	if m != nil {
		return m.Params
	}
	return support.Params{}
}

// EthCallRequest defines EthCall request
type EthCallRequest struct {
	// args uses the same json format as the json rpc api.
	Args []byte `protobuf:"bytes,1,opt,name=args,proto3" json:"args,omitempty"`
	// gas_cap defines the default gas cap to be used
	GasCap uint64 `protobuf:"varint,2,opt,name=gas_cap,json=gasCap,proto3" json:"gas_cap,omitempty"`
	// proposer_address of the requested block in hex format
	ProposerAddress github_com_cosmos_cosmos_sdk_types.ConsAddress `protobuf:"bytes,3,opt,name=proposer_address,json=proposerAddress,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ConsAddress" json:"proposer_address,omitempty"`
	// chain_id is the eip155 chain id parsed from the requested block header
	ChainId int64 `protobuf:"varint,4,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (m *EthCallRequest) Reset()         { *m = EthCallRequest{} }
func (m *EthCallRequest) String() string { return proto.CompactTextString(m) }
func (*EthCallRequest) ProtoMessage()    {}
func (*EthCallRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{16}
}
func (m *EthCallRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EthCallRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EthCallRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EthCallRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EthCallRequest.Merge(m, src)
}
func (m *EthCallRequest) XXX_Size() int {
	return m.Size()
}
func (m *EthCallRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EthCallRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EthCallRequest proto.InternalMessageInfo

func (m *EthCallRequest) GetArgs() []byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *EthCallRequest) GetGasCap() uint64 {
	if m != nil {
		return m.GasCap
	}
	return 0
}

func (m *EthCallRequest) GetProposerAddress() github_com_cosmos_cosmos_sdk_types.ConsAddress {
	if m != nil {
		return m.ProposerAddress
	}
	return nil
}

func (m *EthCallRequest) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

// EstimateGasResponse defines EstimateGas response
type EstimateGasResponse struct {
	// gas returns the estimated gas
	Gas uint64 `protobuf:"varint,1,opt,name=gas,proto3" json:"gas,omitempty"`
}

func (m *EstimateGasResponse) Reset()         { *m = EstimateGasResponse{} }
func (m *EstimateGasResponse) String() string { return proto.CompactTextString(m) }
func (*EstimateGasResponse) ProtoMessage()    {}
func (*EstimateGasResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{17}
}
func (m *EstimateGasResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EstimateGasResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EstimateGasResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EstimateGasResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EstimateGasResponse.Merge(m, src)
}
func (m *EstimateGasResponse) XXX_Size() int {
	return m.Size()
}
func (m *EstimateGasResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EstimateGasResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EstimateGasResponse proto.InternalMessageInfo

func (m *EstimateGasResponse) GetGas() uint64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

// QueryTraceTxRequest defines TraceTx request
type QueryTraceTxRequest struct {
	// msg is the MsgEthereumTx for the requested process
	Msg *MsgEthereumTx `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	// trace_config holds extra parameters to trace functions.
	TraceConfig *support.TraceConfig `protobuf:"bytes,3,opt,name=trace_config,json=traceConfig,proto3" json:"trace_config,omitempty"`
	// predecessors is an array of transactions included in the same block
	// need to be replayed first to get correct context for tracing.
	Predecessors []*MsgEthereumTx `protobuf:"bytes,4,rep,name=predecessors,proto3" json:"predecessors,omitempty"`
	// block_number of requested process
	BlockNumber int64 `protobuf:"varint,5,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// block_hash of requested process
	BlockHash string `protobuf:"bytes,6,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// block_time of requested process
	BlockTime time.Time `protobuf:"bytes,7,opt,name=block_time,json=blockTime,proto3,stdtime" json:"block_time"`
	// proposer_address is the proposer of the requested block
	ProposerAddress github_com_cosmos_cosmos_sdk_types.ConsAddress `protobuf:"bytes,8,opt,name=proposer_address,json=proposerAddress,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ConsAddress" json:"proposer_address,omitempty"`
	// chain_id is the the eip155 chain id parsed from the requested block header
	ChainId int64 `protobuf:"varint,9,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (m *QueryTraceTxRequest) Reset()         { *m = QueryTraceTxRequest{} }
func (m *QueryTraceTxRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTraceTxRequest) ProtoMessage()    {}
func (*QueryTraceTxRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{18}
}
func (m *QueryTraceTxRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTraceTxRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTraceTxRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTraceTxRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTraceTxRequest.Merge(m, src)
}
func (m *QueryTraceTxRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTraceTxRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTraceTxRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTraceTxRequest proto.InternalMessageInfo

func (m *QueryTraceTxRequest) GetMsg() *MsgEthereumTx {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *QueryTraceTxRequest) GetTraceConfig() *support.TraceConfig {
	if m != nil {
		return m.TraceConfig
	}
	return nil
}

func (m *QueryTraceTxRequest) GetPredecessors() []*MsgEthereumTx {
	if m != nil {
		return m.Predecessors
	}
	return nil
}

func (m *QueryTraceTxRequest) GetBlockNumber() int64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *QueryTraceTxRequest) GetBlockHash() string {
	if m != nil {
		return m.BlockHash
	}
	return ""
}

func (m *QueryTraceTxRequest) GetBlockTime() time.Time {
	if m != nil {
		return m.BlockTime
	}
	return time.Time{}
}

func (m *QueryTraceTxRequest) GetProposerAddress() github_com_cosmos_cosmos_sdk_types.ConsAddress {
	if m != nil {
		return m.ProposerAddress
	}
	return nil
}

func (m *QueryTraceTxRequest) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

// QueryTraceTxResponse defines TraceTx response
type QueryTraceTxResponse struct {
	// data is the response serialized in bytes
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *QueryTraceTxResponse) Reset()         { *m = QueryTraceTxResponse{} }
func (m *QueryTraceTxResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTraceTxResponse) ProtoMessage()    {}
func (*QueryTraceTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{19}
}
func (m *QueryTraceTxResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTraceTxResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTraceTxResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTraceTxResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTraceTxResponse.Merge(m, src)
}
func (m *QueryTraceTxResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTraceTxResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTraceTxResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTraceTxResponse proto.InternalMessageInfo

func (m *QueryTraceTxResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// QueryTraceBlockRequest defines TraceTx request
type QueryTraceBlockRequest struct {
	// txs is an array of messages in the block
	Txs []*MsgEthereumTx `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	// trace_config holds extra parameters to trace functions.
	TraceConfig *support.TraceConfig `protobuf:"bytes,3,opt,name=trace_config,json=traceConfig,proto3" json:"trace_config,omitempty"`
	// block_number of the traced block
	BlockNumber int64 `protobuf:"varint,5,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// block_hash (hex) of the traced block
	BlockHash string `protobuf:"bytes,6,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// block_time of the traced block
	BlockTime time.Time `protobuf:"bytes,7,opt,name=block_time,json=blockTime,proto3,stdtime" json:"block_time"`
	// proposer_address is the address of the requested block
	ProposerAddress github_com_cosmos_cosmos_sdk_types.ConsAddress `protobuf:"bytes,8,opt,name=proposer_address,json=proposerAddress,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ConsAddress" json:"proposer_address,omitempty"`
	// chain_id is the eip155 chain id parsed from the requested block header
	ChainId int64 `protobuf:"varint,9,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (m *QueryTraceBlockRequest) Reset()         { *m = QueryTraceBlockRequest{} }
func (m *QueryTraceBlockRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTraceBlockRequest) ProtoMessage()    {}
func (*QueryTraceBlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{20}
}
func (m *QueryTraceBlockRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTraceBlockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTraceBlockRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTraceBlockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTraceBlockRequest.Merge(m, src)
}
func (m *QueryTraceBlockRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTraceBlockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTraceBlockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTraceBlockRequest proto.InternalMessageInfo

func (m *QueryTraceBlockRequest) GetTxs() []*MsgEthereumTx {
	if m != nil {
		return m.Txs
	}
	return nil
}

func (m *QueryTraceBlockRequest) GetTraceConfig() *support.TraceConfig {
	if m != nil {
		return m.TraceConfig
	}
	return nil
}

func (m *QueryTraceBlockRequest) GetBlockNumber() int64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *QueryTraceBlockRequest) GetBlockHash() string {
	if m != nil {
		return m.BlockHash
	}
	return ""
}

func (m *QueryTraceBlockRequest) GetBlockTime() time.Time {
	if m != nil {
		return m.BlockTime
	}
	return time.Time{}
}

func (m *QueryTraceBlockRequest) GetProposerAddress() github_com_cosmos_cosmos_sdk_types.ConsAddress {
	if m != nil {
		return m.ProposerAddress
	}
	return nil
}

func (m *QueryTraceBlockRequest) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

// QueryTraceBlockResponse defines TraceBlock response
type QueryTraceBlockResponse struct {
	// data is the response serialized in bytes
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *QueryTraceBlockResponse) Reset()         { *m = QueryTraceBlockResponse{} }
func (m *QueryTraceBlockResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTraceBlockResponse) ProtoMessage()    {}
func (*QueryTraceBlockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{21}
}
func (m *QueryTraceBlockResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTraceBlockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTraceBlockResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTraceBlockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTraceBlockResponse.Merge(m, src)
}
func (m *QueryTraceBlockResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTraceBlockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTraceBlockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTraceBlockResponse proto.InternalMessageInfo

func (m *QueryTraceBlockResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// QueryBaseFeeRequest defines the request type for querying the EIP1559 base
// fee.
type QueryBaseFeeRequest struct {
}

func (m *QueryBaseFeeRequest) Reset()         { *m = QueryBaseFeeRequest{} }
func (m *QueryBaseFeeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBaseFeeRequest) ProtoMessage()    {}
func (*QueryBaseFeeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{22}
}
func (m *QueryBaseFeeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBaseFeeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBaseFeeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBaseFeeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBaseFeeRequest.Merge(m, src)
}
func (m *QueryBaseFeeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryBaseFeeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBaseFeeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBaseFeeRequest proto.InternalMessageInfo

// QueryBaseFeeResponse returns the EIP1559 base fee.
type QueryBaseFeeResponse struct {
	// base_fee is the EIP1559 base fee
	BaseFee *github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=base_fee,json=baseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"base_fee,omitempty"`
}

func (m *QueryBaseFeeResponse) Reset()         { *m = QueryBaseFeeResponse{} }
func (m *QueryBaseFeeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBaseFeeResponse) ProtoMessage()    {}
func (*QueryBaseFeeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d7bc138cc47c0d0, []int{23}
}
func (m *QueryBaseFeeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBaseFeeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBaseFeeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBaseFeeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBaseFeeResponse.Merge(m, src)
}
func (m *QueryBaseFeeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryBaseFeeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBaseFeeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBaseFeeResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryAccountRequest)(nil), "artela.evm.v1.QueryAccountRequest")
	proto.RegisterType((*QueryAccountResponse)(nil), "artela.evm.v1.QueryAccountResponse")
	proto.RegisterType((*QueryCosmosAccountRequest)(nil), "artela.evm.v1.QueryCosmosAccountRequest")
	proto.RegisterType((*QueryCosmosAccountResponse)(nil), "artela.evm.v1.QueryCosmosAccountResponse")
	proto.RegisterType((*QueryValidatorAccountRequest)(nil), "artela.evm.v1.QueryValidatorAccountRequest")
	proto.RegisterType((*QueryValidatorAccountResponse)(nil), "artela.evm.v1.QueryValidatorAccountResponse")
	proto.RegisterType((*QueryBalanceRequest)(nil), "artela.evm.v1.QueryBalanceRequest")
	proto.RegisterType((*QueryBalanceResponse)(nil), "artela.evm.v1.QueryBalanceResponse")
	proto.RegisterType((*QueryStorageRequest)(nil), "artela.evm.v1.QueryStorageRequest")
	proto.RegisterType((*QueryStorageResponse)(nil), "artela.evm.v1.QueryStorageResponse")
	proto.RegisterType((*QueryCodeRequest)(nil), "artela.evm.v1.QueryCodeRequest")
	proto.RegisterType((*QueryCodeResponse)(nil), "artela.evm.v1.QueryCodeResponse")
	proto.RegisterType((*QueryTxLogsRequest)(nil), "artela.evm.v1.QueryTxLogsRequest")
	proto.RegisterType((*QueryTxLogsResponse)(nil), "artela.evm.v1.QueryTxLogsResponse")
	proto.RegisterType((*QueryParamsRequest)(nil), "artela.evm.v1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "artela.evm.v1.QueryParamsResponse")
	proto.RegisterType((*EthCallRequest)(nil), "artela.evm.v1.EthCallRequest")
	proto.RegisterType((*EstimateGasResponse)(nil), "artela.evm.v1.EstimateGasResponse")
	proto.RegisterType((*QueryTraceTxRequest)(nil), "artela.evm.v1.QueryTraceTxRequest")
	proto.RegisterType((*QueryTraceTxResponse)(nil), "artela.evm.v1.QueryTraceTxResponse")
	proto.RegisterType((*QueryTraceBlockRequest)(nil), "artela.evm.v1.QueryTraceBlockRequest")
	proto.RegisterType((*QueryTraceBlockResponse)(nil), "artela.evm.v1.QueryTraceBlockResponse")
	proto.RegisterType((*QueryBaseFeeRequest)(nil), "artela.evm.v1.QueryBaseFeeRequest")
	proto.RegisterType((*QueryBaseFeeResponse)(nil), "artela.evm.v1.QueryBaseFeeResponse")
}

func init() { proto.RegisterFile("artela/evm/v1/query.proto", fileDescriptor_8d7bc138cc47c0d0) }

var fileDescriptor_8d7bc138cc47c0d0 = []byte{
	// 1429 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe4, 0x56, 0x4f, 0x6f, 0x1b, 0x45,
	0x14, 0xcf, 0xc6, 0x4e, 0xec, 0x3c, 0x27, 0x6d, 0x98, 0xa6, 0x4d, 0xe2, 0x26, 0x76, 0xba, 0xa5,
	0x89, 0x5b, 0xda, 0x5d, 0x92, 0x4a, 0x20, 0x90, 0x10, 0xd4, 0x51, 0x5a, 0xfa, 0x07, 0x54, 0x96,
	0x88, 0x03, 0x52, 0x65, 0x8d, 0xd7, 0xd3, 0xb5, 0x15, 0x7b, 0xc7, 0xdd, 0x19, 0x1b, 0x87, 0x12,
	0x81, 0x7a, 0x40, 0x95, 0xb8, 0x54, 0x42, 0xdc, 0x7b, 0xe2, 0x2b, 0xf0, 0x15, 0x7a, 0xac, 0xc4,
	0x01, 0xc4, 0xa1, 0xa0, 0x96, 0x03, 0x9f, 0x81, 0x13, 0x9a, 0x3f, 0x6b, 0x7b, 0x37, 0xeb, 0x24,
	0xe5, 0xcf, 0x89, 0xd3, 0xee, 0xcc, 0xbc, 0xf7, 0x7e, 0xbf, 0x37, 0xef, 0xcd, 0x7b, 0x0f, 0x16,
	0x71, 0xc0, 0x49, 0x13, 0xdb, 0xa4, 0xdb, 0xb2, 0xbb, 0xeb, 0xf6, 0xbd, 0x0e, 0x09, 0x76, 0xad,
	0x76, 0x40, 0x39, 0x45, 0x33, 0xea, 0xc8, 0x22, 0xdd, 0x96, 0xd5, 0x5d, 0xcf, 0x5f, 0x70, 0x29,
	0x6b, 0x51, 0x66, 0x57, 0x31, 0x23, 0x4a, 0xce, 0xee, 0xae, 0x57, 0x09, 0xc7, 0xeb, 0x76, 0x1b,
	0x7b, 0x0d, 0x1f, 0xf3, 0x06, 0xf5, 0x95, 0x6a, 0x7e, 0x3e, 0x6a, 0x55, 0x58, 0x50, 0x07, 0xa7,
	0xa2, 0x07, 0xbc, 0xa7, 0xf7, 0xe7, 0x3c, 0xea, 0x51, 0xf9, 0x6b, 0x8b, 0x3f, 0xbd, 0xbb, 0xe4,
	0x51, 0xea, 0x35, 0x89, 0x8d, 0xdb, 0x0d, 0x1b, 0xfb, 0x3e, 0xe5, 0x12, 0x83, 0xe9, 0xd3, 0xa2,
	0x3e, 0x95, 0xab, 0x6a, 0xe7, 0xae, 0xcd, 0x1b, 0x2d, 0xc2, 0x38, 0x6e, 0xb5, 0x95, 0x80, 0xf9,
	0x16, 0x9c, 0xf8, 0x48, 0xf0, 0xbc, 0xe2, 0xba, 0xb4, 0xe3, 0x73, 0x87, 0xdc, 0xeb, 0x10, 0xc6,
	0xd1, 0x02, 0x64, 0x70, 0xad, 0x16, 0x10, 0xc6, 0x16, 0x8c, 0x15, 0xa3, 0x34, 0xe5, 0x84, 0xcb,
	0xb7, 0xb3, 0x0f, 0x1f, 0x17, 0xc7, 0xfe, 0x78, 0x5c, 0x1c, 0x33, 0x5d, 0x98, 0x8b, 0xaa, 0xb2,
	0x36, 0xf5, 0x19, 0x11, 0xba, 0x55, 0xdc, 0xc4, 0xbe, 0x4b, 0x42, 0x5d, 0xbd, 0x44, 0xa7, 0x61,
	0xca, 0xa5, 0x35, 0x52, 0xa9, 0x63, 0x56, 0x5f, 0x18, 0x97, 0x67, 0x59, 0xb1, 0xf1, 0x3e, 0x66,
	0x75, 0x34, 0x07, 0x13, 0x3e, 0x15, 0x4a, 0xa9, 0x15, 0xa3, 0x94, 0x76, 0xd4, 0xc2, 0x7c, 0x17,
	0x16, 0x25, 0xc8, 0xa6, 0xbc, 0xd8, 0xbf, 0xc1, 0xf2, 0x6b, 0x03, 0xf2, 0x49, 0x16, 0x34, 0xd9,
	0x73, 0x70, 0x4c, 0xc5, 0xac, 0x12, 0xb5, 0x34, 0xa3, 0x76, 0xaf, 0xa8, 0x4d, 0x94, 0x87, 0x2c,
	0x13, 0xa0, 0x82, 0xdf, 0xb8, 0xe4, 0xd7, 0x5f, 0x0b, 0x13, 0x58, 0x59, 0xad, 0xf8, 0x9d, 0x56,
	0x95, 0x04, 0xda, 0x83, 0x19, 0xbd, 0xfb, 0xa1, 0xdc, 0x34, 0x6f, 0xc2, 0x92, 0xe4, 0xf1, 0x09,
	0x6e, 0x36, 0x6a, 0x98, 0xd3, 0x20, 0xe6, 0xcc, 0x19, 0x98, 0x76, 0xa9, 0x1f, 0xe7, 0x91, 0x13,
	0x7b, 0x57, 0xf6, 0x79, 0xf5, 0x8d, 0x01, 0xcb, 0x23, 0xac, 0x69, 0xc7, 0xd6, 0xe0, 0x78, 0xc8,
	0x2a, 0x6a, 0x31, 0x24, 0xfb, 0x2f, 0xba, 0x16, 0x26, 0x51, 0x59, 0xc5, 0xf9, 0x65, 0xc2, 0xf3,
	0xba, 0x4e, 0xa2, 0xbe, 0xea, 0x61, 0x49, 0x64, 0xde, 0xd4, 0x60, 0x1f, 0x73, 0x1a, 0x60, 0xef,
	0x70, 0x30, 0x34, 0x0b, 0xa9, 0x1d, 0xb2, 0xab, 0xf3, 0x4d, 0xfc, 0x0e, 0xc1, 0x5f, 0xd4, 0xf0,
	0x7d, 0x63, 0x1a, 0x7e, 0x0e, 0x26, 0xba, 0xb8, 0xd9, 0x09, 0xc1, 0xd5, 0xc2, 0x7c, 0x03, 0x66,
	0x75, 0x2a, 0xd5, 0x5e, 0xca, 0xc9, 0x35, 0x78, 0x65, 0x48, 0x4f, 0x43, 0x20, 0x48, 0x8b, 0xdc,
	0x97, 0x5a, 0xd3, 0x8e, 0xfc, 0x37, 0x3f, 0x07, 0x24, 0x05, 0xb7, 0x7b, 0xb7, 0xa8, 0xc7, 0x42,
	0x08, 0x04, 0x69, 0xf9, 0x62, 0x94, 0x7d, 0xf9, 0x8f, 0xae, 0x02, 0x0c, 0x2a, 0x8a, 0xf4, 0x2d,
	0xb7, 0xb1, 0x6a, 0xa9, 0xa4, 0xb5, 0x44, 0xf9, 0xb1, 0x54, 0x99, 0xd2, 0xe5, 0xc7, 0xba, 0x3d,
	0xb8, 0x2a, 0x67, 0x48, 0x33, 0xfa, 0x50, 0x4e, 0x44, 0xc0, 0x35, 0xcf, 0x55, 0x48, 0x37, 0xa9,
	0x27, 0xbc, 0x4b, 0x95, 0x72, 0x1b, 0xc8, 0x8a, 0x54, 0x3c, 0xeb, 0x16, 0xf5, 0x1c, 0x79, 0x8e,
	0xae, 0x25, 0x30, 0x5a, 0x3b, 0x94, 0x91, 0x02, 0x19, 0xa6, 0x64, 0xce, 0xe9, 0x4b, 0xb8, 0x8d,
	0x03, 0xdc, 0x0a, 0x2f, 0xc1, 0xbc, 0xa1, 0xd9, 0x85, 0xbb, 0x9a, 0xdd, 0x65, 0x98, 0x6c, 0xcb,
	0x1d, 0x79, 0x3b, 0xb9, 0x8d, 0x93, 0x31, 0x7e, 0x4a, 0xbc, 0x9c, 0x7e, 0xf2, 0xac, 0x38, 0xe6,
	0x68, 0x51, 0xf3, 0x07, 0x03, 0x8e, 0x6d, 0xf1, 0xfa, 0x26, 0x6e, 0x36, 0x87, 0xee, 0x18, 0x07,
	0x1e, 0x0b, 0xa3, 0x21, 0xfe, 0xd1, 0x3c, 0x64, 0x3c, 0xcc, 0x2a, 0x2e, 0x6e, 0xeb, 0x87, 0x31,
	0xe9, 0x61, 0xb6, 0x89, 0xdb, 0xe8, 0x0e, 0xcc, 0xb6, 0x03, 0xda, 0xa6, 0x8c, 0x04, 0xfd, 0xc7,
	0x25, 0x1e, 0xc6, 0x74, 0x79, 0xe3, 0xcf, 0x67, 0x45, 0xcb, 0x6b, 0xf0, 0x7a, 0xa7, 0x6a, 0xb9,
	0xb4, 0x65, 0xeb, 0x7e, 0xa0, 0x3e, 0x97, 0x58, 0x6d, 0xc7, 0xe6, 0xbb, 0x6d, 0xc2, 0xac, 0xcd,
	0xc1, 0xab, 0x76, 0x8e, 0x87, 0xb6, 0xc2, 0x17, 0xb9, 0x08, 0x59, 0xb7, 0x8e, 0x1b, 0x7e, 0xa5,
	0x51, 0x5b, 0x48, 0xaf, 0x18, 0xa5, 0x94, 0x93, 0x91, 0xeb, 0xeb, 0x35, 0x73, 0x0d, 0x4e, 0x6c,
	0x31, 0xde, 0x68, 0x61, 0x4e, 0xae, 0xe1, 0xc1, 0x2d, 0xcc, 0x42, 0xca, 0xc3, 0x8a, 0x7c, 0xda,
	0x11, 0xbf, 0xe6, 0x4f, 0xa9, 0x30, 0x9a, 0x01, 0x76, 0xc9, 0x76, 0x2f, 0xf4, 0xd3, 0x82, 0x54,
	0x8b, 0x79, 0xfa, 0xb2, 0x96, 0x62, 0x97, 0xf5, 0x01, 0xf3, 0xb6, 0x78, 0x9d, 0x04, 0xa4, 0xd3,
	0xda, 0xee, 0x39, 0x42, 0x10, 0xbd, 0x03, 0xd3, 0x5c, 0x58, 0xa8, 0xb8, 0xd4, 0xbf, 0xdb, 0xf0,
	0xa4, 0x9b, 0xb9, 0x8d, 0x7c, 0x4c, 0x51, 0x82, 0x6c, 0x4a, 0x09, 0x27, 0xc7, 0x07, 0x0b, 0xf4,
	0x1e, 0x4c, 0xb7, 0x03, 0x52, 0x23, 0x2e, 0x61, 0x8c, 0x06, 0x6c, 0x21, 0x2d, 0x93, 0xe8, 0x60,
	0xdc, 0x88, 0x86, 0x28, 0x8b, 0xd5, 0x26, 0x75, 0x77, 0xc2, 0x02, 0x34, 0x21, 0x2f, 0x24, 0x27,
	0xf7, 0x54, 0xf9, 0x41, 0xcb, 0x00, 0x4a, 0x44, 0xbe, 0x92, 0x49, 0xf9, 0x4a, 0xa6, 0xe4, 0x8e,
	0x6c, 0x2c, 0x9b, 0xe1, 0xb1, 0xe8, 0x7d, 0x0b, 0x19, 0xed, 0x80, 0x6a, 0x8c, 0x56, 0xd8, 0x18,
	0xad, 0xed, 0xb0, 0x31, 0x96, 0xb3, 0x22, 0x57, 0x1e, 0xfd, 0x5a, 0x34, 0xb4, 0x11, 0x71, 0x92,
	0x18, 0xf2, 0xec, 0x7f, 0x13, 0xf2, 0xa9, 0x48, 0xc8, 0x6f, 0xa4, 0xb3, 0xe3, 0xb3, 0x29, 0x27,
	0xcb, 0x7b, 0x95, 0x86, 0x5f, 0x23, 0x3d, 0xf3, 0x82, 0x2e, 0x59, 0xfd, 0xc0, 0x0e, 0xea, 0x49,
	0x0d, 0x73, 0x1c, 0x66, 0xb0, 0xf8, 0x37, 0x1f, 0xa6, 0xe0, 0xd4, 0x40, 0xb8, 0x2c, 0xbc, 0x19,
	0x4a, 0x04, 0xde, 0x0b, 0x5f, 0xf5, 0x21, 0x89, 0xc0, 0x7b, 0xec, 0x9f, 0x26, 0xc2, 0xff, 0x3d,
	0x8c, 0xe6, 0x25, 0x98, 0xdf, 0x17, 0x89, 0x03, 0x22, 0x77, 0xb2, 0xdf, 0x52, 0x19, 0xb9, 0x4a,
	0xc2, 0xd2, 0x6d, 0xde, 0xe9, 0xb7, 0x4b, 0xbd, 0xad, 0x4d, 0x6c, 0x41, 0x56, 0x94, 0xd8, 0xca,
	0x5d, 0xa2, 0x5b, 0x56, 0xf9, 0xc2, 0x2f, 0xcf, 0x8a, 0xab, 0x47, 0xf0, 0xe7, 0xba, 0xcf, 0x45,
	0x6f, 0x95, 0xe6, 0x36, 0xbe, 0x9a, 0x86, 0x09, 0x69, 0x1f, 0x7d, 0x01, 0x19, 0x3d, 0x51, 0x20,
	0x33, 0x16, 0xe3, 0x84, 0x79, 0x31, 0x7f, 0xf6, 0x40, 0x19, 0x45, 0xd2, 0x2c, 0x3d, 0xf8, 0xf1,
	0xf7, 0x6f, 0xc7, 0x4d, 0xb4, 0x62, 0x47, 0x27, 0x5c, 0x3d, 0x4c, 0xd8, 0xf7, 0x75, 0x44, 0xf6,
	0xd0, 0x77, 0x06, 0xcc, 0x44, 0xe6, 0x35, 0x54, 0x4a, 0x02, 0x48, 0x1a, 0x0a, 0xf3, 0xe7, 0x8f,
	0x20, 0xa9, 0x09, 0xd9, 0x92, 0xd0, 0x79, 0xb4, 0x16, 0x23, 0x14, 0x4e, 0x84, 0xfb, 0x78, 0x7d,
	0x6f, 0xc0, 0x6c, 0x7c, 0xe2, 0x42, 0xaf, 0x25, 0x01, 0x8e, 0x98, 0xf2, 0xf2, 0x17, 0x8f, 0x26,
	0xac, 0x09, 0xbe, 0x29, 0x09, 0xae, 0x23, 0x3b, 0x46, 0xb0, 0x1b, 0x2a, 0x0c, 0x38, 0x0e, 0xcf,
	0x8e, 0x7b, 0x68, 0x0f, 0x32, 0x7a, 0xa2, 0x4a, 0x0e, 0x5f, 0x74, 0x52, 0x4b, 0x0e, 0x5f, 0x6c,
	0x24, 0x33, 0xcf, 0x4b, 0x32, 0x67, 0xd1, 0x99, 0x18, 0x19, 0x3d, 0x98, 0xb1, 0xa1, 0x7b, 0x7a,
	0x60, 0x40, 0x46, 0x8f, 0x54, 0xc9, 0xf8, 0xd1, 0xe1, 0x2d, 0x19, 0x3f, 0x36, 0x93, 0x99, 0x96,
	0xc4, 0x2f, 0xa1, 0xd5, 0x18, 0x3e, 0x53, 0x72, 0x03, 0x78, 0xfb, 0xfe, 0x0e, 0xd9, 0xdd, 0x43,
	0xf7, 0x20, 0x2d, 0x06, 0x2e, 0x54, 0x4c, 0x4e, 0x88, 0xfe, 0x08, 0x97, 0x5f, 0x19, 0x2d, 0xa0,
	0xa1, 0x57, 0x25, 0xf4, 0x0a, 0x2a, 0xec, 0x4b, 0x94, 0x5a, 0xc4, 0x6f, 0x1f, 0x26, 0xd5, 0xc0,
	0x81, 0xce, 0x24, 0xd9, 0x8c, 0x4c, 0x34, 0x79, 0xf3, 0x20, 0x11, 0x0d, 0xbc, 0x2c, 0x81, 0xe7,
	0xd1, 0xc9, 0x18, 0xb0, 0x1a, 0x64, 0x10, 0x85, 0x8c, 0x9e, 0x63, 0xd0, 0x72, 0xcc, 0x5a, 0x74,
	0xbe, 0xc9, 0xbf, 0x7a, 0x60, 0x85, 0x0f, 0xe1, 0x8a, 0x12, 0x6e, 0x11, 0xcd, 0xc7, 0xe0, 0x08,
	0xaf, 0x57, 0x5c, 0x81, 0xd2, 0x81, 0xdc, 0xd0, 0xfc, 0x71, 0x18, 0x68, 0xdc, 0xc3, 0x84, 0xd1,
	0xc5, 0x3c, 0x2b, 0x21, 0x97, 0xd1, 0xe9, 0x38, 0xa4, 0x96, 0xad, 0x78, 0x98, 0x21, 0x06, 0x19,
	0xdd, 0xee, 0x92, 0xd3, 0x29, 0x3a, 0xe4, 0x24, 0xa7, 0x53, 0xac, 0x5f, 0x8e, 0xf4, 0x55, 0x75,
	0x39, 0xde, 0x43, 0x5f, 0x02, 0x0c, 0x8a, 0x35, 0x3a, 0x37, 0xd2, 0xe6, 0x70, 0x5b, 0xcd, 0xaf,
	0x1e, 0x26, 0xa6, 0xd1, 0x4d, 0x89, 0xbe, 0x84, 0xf2, 0x89, 0xe8, 0xb2, 0x61, 0x09, 0xaf, 0x75,
	0x9d, 0x1f, 0xf5, 0x88, 0x87, 0x7b, 0xc3, 0xa8, 0x47, 0x1c, 0x69, 0x14, 0x23, 0xbd, 0x0e, 0xbb,
	0x47, 0xf9, 0xfa, 0x93, 0xe7, 0x05, 0xe3, 0xe9, 0xf3, 0x82, 0xf1, 0xdb, 0xf3, 0x82, 0xf1, 0xe8,
	0x45, 0x61, 0xec, 0xe9, 0x8b, 0xc2, 0xd8, 0xcf, 0x2f, 0x0a, 0x63, 0x9f, 0xda, 0x43, 0xdd, 0x44,
	0x29, 0x5f, 0xf2, 0x09, 0xff, 0x8c, 0x06, 0x3b, 0xa1, 0xad, 0xee, 0xba, 0xdd, 0x93, 0x06, 0x65,
	0x6b, 0xa9, 0x4e, 0xca, 0xae, 0x7c, 0xf9, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x30, 0x1a,
	0x43, 0x40, 0x11, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this support file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Account queries an Ethereum account.
	Account(ctx context.Context, in *QueryAccountRequest, opts ...grpc.CallOption) (*QueryAccountResponse, error)
	// CosmosAccount queries an Ethereum account's Cosmos Address.
	CosmosAccount(ctx context.Context, in *QueryCosmosAccountRequest, opts ...grpc.CallOption) (*QueryCosmosAccountResponse, error)
	// ValidatorAccount queries an Ethereum account's from a validator consensus
	// Address.
	ValidatorAccount(ctx context.Context, in *QueryValidatorAccountRequest, opts ...grpc.CallOption) (*QueryValidatorAccountResponse, error)
	// Balance queries the balance of a the EVM denomination for a single
	// EthAccount.
	Balance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error)
	// Storage queries the balance of all coins for a single account.
	Storage(ctx context.Context, in *QueryStorageRequest, opts ...grpc.CallOption) (*QueryStorageResponse, error)
	// Code queries the balance of all coins for a single account.
	Code(ctx context.Context, in *QueryCodeRequest, opts ...grpc.CallOption) (*QueryCodeResponse, error)
	// Params queries the parameters of x/evm module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// EthCall implements the `eth_call` rpc api
	EthCall(ctx context.Context, in *EthCallRequest, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error)
	// EstimateGas implements the `eth_estimateGas` rpc api
	EstimateGas(ctx context.Context, in *EthCallRequest, opts ...grpc.CallOption) (*EstimateGasResponse, error)
	// TraceTx implements the `debug_traceTransaction` rpc api
	TraceTx(ctx context.Context, in *QueryTraceTxRequest, opts ...grpc.CallOption) (*QueryTraceTxResponse, error)
	// TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api
	TraceBlock(ctx context.Context, in *QueryTraceBlockRequest, opts ...grpc.CallOption) (*QueryTraceBlockResponse, error)
	// BaseFee queries the base fee of the parent block of the current block,
	// it's similar to feemarket module's method, but also checks london hardfork status.
	BaseFee(ctx context.Context, in *QueryBaseFeeRequest, opts ...grpc.CallOption) (*QueryBaseFeeResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Account(ctx context.Context, in *QueryAccountRequest, opts ...grpc.CallOption) (*QueryAccountResponse, error) {
	out := new(QueryAccountResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/Account", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CosmosAccount(ctx context.Context, in *QueryCosmosAccountRequest, opts ...grpc.CallOption) (*QueryCosmosAccountResponse, error) {
	out := new(QueryCosmosAccountResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/CosmosAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ValidatorAccount(ctx context.Context, in *QueryValidatorAccountRequest, opts ...grpc.CallOption) (*QueryValidatorAccountResponse, error) {
	out := new(QueryValidatorAccountResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/ValidatorAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Balance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error) {
	out := new(QueryBalanceResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/Balance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Storage(ctx context.Context, in *QueryStorageRequest, opts ...grpc.CallOption) (*QueryStorageResponse, error) {
	out := new(QueryStorageResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/Storage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Code(ctx context.Context, in *QueryCodeRequest, opts ...grpc.CallOption) (*QueryCodeResponse, error) {
	out := new(QueryCodeResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/Code", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EthCall(ctx context.Context, in *EthCallRequest, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error) {
	out := new(MsgEthereumTxResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/EthCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EstimateGas(ctx context.Context, in *EthCallRequest, opts ...grpc.CallOption) (*EstimateGasResponse, error) {
	out := new(EstimateGasResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/EstimateGas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TraceTx(ctx context.Context, in *QueryTraceTxRequest, opts ...grpc.CallOption) (*QueryTraceTxResponse, error) {
	out := new(QueryTraceTxResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/TraceTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TraceBlock(ctx context.Context, in *QueryTraceBlockRequest, opts ...grpc.CallOption) (*QueryTraceBlockResponse, error) {
	out := new(QueryTraceBlockResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/TraceBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BaseFee(ctx context.Context, in *QueryBaseFeeRequest, opts ...grpc.CallOption) (*QueryBaseFeeResponse, error) {
	out := new(QueryBaseFeeResponse)
	err := c.cc.Invoke(ctx, "/artela.evm.v1.Query/BaseFee", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Account queries an Ethereum account.
	Account(context.Context, *QueryAccountRequest) (*QueryAccountResponse, error)
	// CosmosAccount queries an Ethereum account's Cosmos Address.
	CosmosAccount(context.Context, *QueryCosmosAccountRequest) (*QueryCosmosAccountResponse, error)
	// ValidatorAccount queries an Ethereum account's from a validator consensus
	// Address.
	ValidatorAccount(context.Context, *QueryValidatorAccountRequest) (*QueryValidatorAccountResponse, error)
	// Balance queries the balance of a the EVM denomination for a single
	// EthAccount.
	Balance(context.Context, *QueryBalanceRequest) (*QueryBalanceResponse, error)
	// Storage queries the balance of all coins for a single account.
	Storage(context.Context, *QueryStorageRequest) (*QueryStorageResponse, error)
	// Code queries the balance of all coins for a single account.
	Code(context.Context, *QueryCodeRequest) (*QueryCodeResponse, error)
	// Params queries the parameters of x/evm module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// EthCall implements the `eth_call` rpc api
	EthCall(context.Context, *EthCallRequest) (*MsgEthereumTxResponse, error)
	// EstimateGas implements the `eth_estimateGas` rpc api
	EstimateGas(context.Context, *EthCallRequest) (*EstimateGasResponse, error)
	// TraceTx implements the `debug_traceTransaction` rpc api
	TraceTx(context.Context, *QueryTraceTxRequest) (*QueryTraceTxResponse, error)
	// TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api
	TraceBlock(context.Context, *QueryTraceBlockRequest) (*QueryTraceBlockResponse, error)
	// BaseFee queries the base fee of the parent block of the current block,
	// it's similar to feemarket module's method, but also checks london hardfork status.
	BaseFee(context.Context, *QueryBaseFeeRequest) (*QueryBaseFeeResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Account(ctx context.Context, req *QueryAccountRequest) (*QueryAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Account not implemented")
}
func (*UnimplementedQueryServer) CosmosAccount(ctx context.Context, req *QueryCosmosAccountRequest) (*QueryCosmosAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CosmosAccount not implemented")
}
func (*UnimplementedQueryServer) ValidatorAccount(ctx context.Context, req *QueryValidatorAccountRequest) (*QueryValidatorAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidatorAccount not implemented")
}
func (*UnimplementedQueryServer) Balance(ctx context.Context, req *QueryBalanceRequest) (*QueryBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
}
func (*UnimplementedQueryServer) Storage(ctx context.Context, req *QueryStorageRequest) (*QueryStorageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Storage not implemented")
}
func (*UnimplementedQueryServer) Code(ctx context.Context, req *QueryCodeRequest) (*QueryCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Code not implemented")
}
func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) EthCall(ctx context.Context, req *EthCallRequest) (*MsgEthereumTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EthCall not implemented")
}
func (*UnimplementedQueryServer) EstimateGas(ctx context.Context, req *EthCallRequest) (*EstimateGasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimateGas not implemented")
}
func (*UnimplementedQueryServer) TraceTx(ctx context.Context, req *QueryTraceTxRequest) (*QueryTraceTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TraceTx not implemented")
}
func (*UnimplementedQueryServer) TraceBlock(ctx context.Context, req *QueryTraceBlockRequest) (*QueryTraceBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TraceBlock not implemented")
}
func (*UnimplementedQueryServer) BaseFee(ctx context.Context, req *QueryBaseFeeRequest) (*QueryBaseFeeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BaseFee not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Account_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Account(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/Account",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Account(ctx, req.(*QueryAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CosmosAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCosmosAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CosmosAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/CosmosAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CosmosAccount(ctx, req.(*QueryCosmosAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ValidatorAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryValidatorAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ValidatorAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/ValidatorAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ValidatorAccount(ctx, req.(*QueryValidatorAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Balance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Balance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/Balance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Balance(ctx, req.(*QueryBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Storage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Storage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/Storage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Storage(ctx, req.(*QueryStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Code_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Code(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/Code",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Code(ctx, req.(*QueryCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EthCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EthCallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EthCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/EthCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EthCall(ctx, req.(*EthCallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EstimateGas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EthCallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EstimateGas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/EstimateGas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EstimateGas(ctx, req.(*EthCallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraceTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraceTxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraceTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/TraceTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraceTx(ctx, req.(*QueryTraceTxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TraceBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTraceBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TraceBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/TraceBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TraceBlock(ctx, req.(*QueryTraceBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BaseFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBaseFeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BaseFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/artela.evm.v1.Query/BaseFee",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BaseFee(ctx, req.(*QueryBaseFeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "artela.evm.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Account",
			Handler:    _Query_Account_Handler,
		},
		{
			MethodName: "CosmosAccount",
			Handler:    _Query_CosmosAccount_Handler,
		},
		{
			MethodName: "ValidatorAccount",
			Handler:    _Query_ValidatorAccount_Handler,
		},
		{
			MethodName: "Balance",
			Handler:    _Query_Balance_Handler,
		},
		{
			MethodName: "Storage",
			Handler:    _Query_Storage_Handler,
		},
		{
			MethodName: "Code",
			Handler:    _Query_Code_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "EthCall",
			Handler:    _Query_EthCall_Handler,
		},
		{
			MethodName: "EstimateGas",
			Handler:    _Query_EstimateGas_Handler,
		},
		{
			MethodName: "TraceTx",
			Handler:    _Query_TraceTx_Handler,
		},
		{
			MethodName: "TraceBlock",
			Handler:    _Query_TraceBlock_Handler,
		},
		{
			MethodName: "BaseFee",
			Handler:    _Query_BaseFee_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "artela/evm/v1/query.proto",
}

func (m *QueryAccountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAccountRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAccountRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Nonce != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x18
	}
	if len(m.CodeHash) > 0 {
		i -= len(m.CodeHash)
		copy(dAtA[i:], m.CodeHash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.CodeHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Balance) > 0 {
		i -= len(m.Balance)
		copy(dAtA[i:], m.Balance)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Balance)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryCosmosAccountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCosmosAccountRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCosmosAccountRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryCosmosAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCosmosAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCosmosAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AccountNumber != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.AccountNumber))
		i--
		dAtA[i] = 0x18
	}
	if m.Sequence != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if len(m.CosmosAddress) > 0 {
		i -= len(m.CosmosAddress)
		copy(dAtA[i:], m.CosmosAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.CosmosAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryValidatorAccountRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryValidatorAccountRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryValidatorAccountRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ConsAddress) > 0 {
		i -= len(m.ConsAddress)
		copy(dAtA[i:], m.ConsAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ConsAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryValidatorAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryValidatorAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryValidatorAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AccountNumber != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.AccountNumber))
		i--
		dAtA[i] = 0x18
	}
	if m.Sequence != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if len(m.AccountAddress) > 0 {
		i -= len(m.AccountAddress)
		copy(dAtA[i:], m.AccountAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.AccountAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryBalanceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBalanceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBalanceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryBalanceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBalanceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBalanceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Balance) > 0 {
		i -= len(m.Balance)
		copy(dAtA[i:], m.Balance)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Balance)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryStorageRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryStorageRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryStorageRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryStorageResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryStorageResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryStorageResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryCodeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCodeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCodeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryCodeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCodeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCodeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Code) > 0 {
		i -= len(m.Code)
		copy(dAtA[i:], m.Code)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Code)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTxLogsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTxLogsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTxLogsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTxLogsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTxLogsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTxLogsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Logs) > 0 {
		for iNdEx := len(m.Logs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Logs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *EthCallRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EthCallRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EthCallRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ChainId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ProposerAddress) > 0 {
		i -= len(m.ProposerAddress)
		copy(dAtA[i:], m.ProposerAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ProposerAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if m.GasCap != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.GasCap))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Args) > 0 {
		i -= len(m.Args)
		copy(dAtA[i:], m.Args)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Args)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EstimateGasResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EstimateGasResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EstimateGasResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Gas != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Gas))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryTraceTxRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTraceTxRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTraceTxRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ChainId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x48
	}
	if len(m.ProposerAddress) > 0 {
		i -= len(m.ProposerAddress)
		copy(dAtA[i:], m.ProposerAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ProposerAddress)))
		i--
		dAtA[i] = 0x42
	}
	n4, err4 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.BlockTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTime):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintQuery(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x3a
	if len(m.BlockHash) > 0 {
		i -= len(m.BlockHash)
		copy(dAtA[i:], m.BlockHash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.BlockHash)))
		i--
		dAtA[i] = 0x32
	}
	if m.BlockNumber != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.BlockNumber))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Predecessors) > 0 {
		for iNdEx := len(m.Predecessors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Predecessors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.TraceConfig != nil {
		{
			size, err := m.TraceConfig.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Msg != nil {
		{
			size, err := m.Msg.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTraceTxResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTraceTxResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTraceTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTraceBlockRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTraceBlockRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTraceBlockRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ChainId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x48
	}
	if len(m.ProposerAddress) > 0 {
		i -= len(m.ProposerAddress)
		copy(dAtA[i:], m.ProposerAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ProposerAddress)))
		i--
		dAtA[i] = 0x42
	}
	n7, err7 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.BlockTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTime):])
	if err7 != nil {
		return 0, err7
	}
	i -= n7
	i = encodeVarintQuery(dAtA, i, uint64(n7))
	i--
	dAtA[i] = 0x3a
	if len(m.BlockHash) > 0 {
		i -= len(m.BlockHash)
		copy(dAtA[i:], m.BlockHash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.BlockHash)))
		i--
		dAtA[i] = 0x32
	}
	if m.BlockNumber != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.BlockNumber))
		i--
		dAtA[i] = 0x28
	}
	if m.TraceConfig != nil {
		{
			size, err := m.TraceConfig.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Txs) > 0 {
		for iNdEx := len(m.Txs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Txs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryTraceBlockResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTraceBlockResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTraceBlockResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryBaseFeeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBaseFeeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBaseFeeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryBaseFeeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBaseFeeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBaseFeeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BaseFee != nil {
		{
			size := m.BaseFee.Size()
			i -= size
			if _, err := m.BaseFee.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryAccountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Balance)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.CodeHash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Nonce != 0 {
		n += 1 + sovQuery(uint64(m.Nonce))
	}
	return n
}

func (m *QueryCosmosAccountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryCosmosAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CosmosAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovQuery(uint64(m.Sequence))
	}
	if m.AccountNumber != 0 {
		n += 1 + sovQuery(uint64(m.AccountNumber))
	}
	return n
}

func (m *QueryValidatorAccountRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ConsAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryValidatorAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AccountAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovQuery(uint64(m.Sequence))
	}
	if m.AccountNumber != 0 {
		n += 1 + sovQuery(uint64(m.AccountNumber))
	}
	return n
}

func (m *QueryBalanceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryBalanceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Balance)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryStorageRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryStorageResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryCodeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryCodeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTxLogsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTxLogsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Logs) > 0 {
		for _, e := range m.Logs {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *EthCallRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Args)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.GasCap != 0 {
		n += 1 + sovQuery(uint64(m.GasCap))
	}
	l = len(m.ProposerAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovQuery(uint64(m.ChainId))
	}
	return n
}

func (m *EstimateGasResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Gas != 0 {
		n += 1 + sovQuery(uint64(m.Gas))
	}
	return n
}

func (m *QueryTraceTxRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.TraceConfig != nil {
		l = m.TraceConfig.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	if len(m.Predecessors) > 0 {
		for _, e := range m.Predecessors {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.BlockNumber != 0 {
		n += 1 + sovQuery(uint64(m.BlockNumber))
	}
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTime)
	n += 1 + l + sovQuery(uint64(l))
	l = len(m.ProposerAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovQuery(uint64(m.ChainId))
	}
	return n
}

func (m *QueryTraceTxResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTraceBlockRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Txs) > 0 {
		for _, e := range m.Txs {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.TraceConfig != nil {
		l = m.TraceConfig.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.BlockNumber != 0 {
		n += 1 + sovQuery(uint64(m.BlockNumber))
	}
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.BlockTime)
	n += 1 + l + sovQuery(uint64(l))
	l = len(m.ProposerAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovQuery(uint64(m.ChainId))
	}
	return n
}

func (m *QueryTraceBlockResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryBaseFeeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryBaseFeeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BaseFee != nil {
		l = m.BaseFee.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryAccountRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAccountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAccountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryAccountResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Balance = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CodeHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCosmosAccountRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCosmosAccountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCosmosAccountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCosmosAccountResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCosmosAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCosmosAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CosmosAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountNumber", wireType)
			}
			m.AccountNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AccountNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryValidatorAccountRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryValidatorAccountRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryValidatorAccountRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConsAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryValidatorAccountResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryValidatorAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryValidatorAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountNumber", wireType)
			}
			m.AccountNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AccountNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryBalanceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBalanceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBalanceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryBalanceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBalanceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBalanceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Balance = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryStorageRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryStorageRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryStorageRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryStorageResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryStorageResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryStorageResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCodeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCodeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCodeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryCodeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryCodeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCodeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = append(m.Code[:0], dAtA[iNdEx:postIndex]...)
			if m.Code == nil {
				m.Code = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTxLogsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTxLogsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTxLogsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTxLogsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTxLogsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTxLogsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Logs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Logs = append(m.Logs, &support.Log{})
			if err := m.Logs[len(m.Logs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
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
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *EthCallRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: EthCallRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EthCallRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Args = append(m.Args[:0], dAtA[iNdEx:postIndex]...)
			if m.Args == nil {
				m.Args = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasCap", wireType)
			}
			m.GasCap = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasCap |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposerAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProposerAddress = append(m.ProposerAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ProposerAddress == nil {
				m.ProposerAddress = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *EstimateGasResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: EstimateGasResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EstimateGasResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			m.Gas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Gas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTraceTxRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTraceTxRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTraceTxRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Msg == nil {
				m.Msg = &MsgEthereumTx{}
			}
			if err := m.Msg.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TraceConfig == nil {
				m.TraceConfig = &support.TraceConfig{}
			}
			if err := m.TraceConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Predecessors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Predecessors = append(m.Predecessors, &MsgEthereumTx{})
			if err := m.Predecessors[len(m.Predecessors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockNumber", wireType)
			}
			m.BlockNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockNumber |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlockHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.BlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposerAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProposerAddress = append(m.ProposerAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ProposerAddress == nil {
				m.ProposerAddress = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTraceTxResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTraceTxResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTraceTxResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTraceBlockRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTraceBlockRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTraceBlockRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Txs = append(m.Txs, &MsgEthereumTx{})
			if err := m.Txs[len(m.Txs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TraceConfig == nil {
				m.TraceConfig = &support.TraceConfig{}
			}
			if err := m.TraceConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockNumber", wireType)
			}
			m.BlockNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockNumber |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlockHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.BlockTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposerAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProposerAddress = append(m.ProposerAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ProposerAddress == nil {
				m.ProposerAddress = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTraceBlockResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTraceBlockResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTraceBlockResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryBaseFeeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBaseFeeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBaseFeeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryBaseFeeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryBaseFeeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBaseFeeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.BaseFee = &v
			if err := m.BaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
