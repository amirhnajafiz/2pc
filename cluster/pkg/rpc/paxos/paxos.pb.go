// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: paxos.proto

package paxos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ballot-number is used for ordering the input transactions.
type BallotNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence int64  `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	NodeId   string `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *BallotNumber) Reset() {
	*x = BallotNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BallotNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BallotNumber) ProtoMessage() {}

func (x *BallotNumber) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BallotNumber.ProtoReflect.Descriptor instead.
func (*BallotNumber) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{0}
}

func (x *BallotNumber) GetSequence() int64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *BallotNumber) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// ping message is sent by the leader to all other nodes.
type PingMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastCommitted *BallotNumber `protobuf:"bytes,1,opt,name=last_committed,json=lastCommitted,proto3" json:"last_committed,omitempty"`
	NodeId        string        `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *PingMsg) Reset() {
	*x = PingMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingMsg) ProtoMessage() {}

func (x *PingMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingMsg.ProtoReflect.Descriptor instead.
func (*PingMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{1}
}

func (x *PingMsg) GetLastCommitted() *BallotNumber {
	if x != nil {
		return x.LastCommitted
	}
	return nil
}

func (x *PingMsg) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// pong message is sent by the followers to
// the leader when it wants to be synced.
type PongMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastCommitted *BallotNumber `protobuf:"bytes,1,opt,name=last_committed,json=lastCommitted,proto3" json:"last_committed,omitempty"`
	NodeId        string        `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *PongMsg) Reset() {
	*x = PongMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PongMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PongMsg) ProtoMessage() {}

func (x *PongMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PongMsg.ProtoReflect.Descriptor instead.
func (*PongMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{2}
}

func (x *PongMsg) GetLastCommitted() *BallotNumber {
	if x != nil {
		return x.LastCommitted
	}
	return nil
}

func (x *PongMsg) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// sync message and sync item is sent from the sender to receiver
type SyncItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record string `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
	Value  int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *SyncItem) Reset() {
	*x = SyncItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncItem) ProtoMessage() {}

func (x *SyncItem) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncItem.ProtoReflect.Descriptor instead.
func (*SyncItem) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{3}
}

func (x *SyncItem) GetRecord() string {
	if x != nil {
		return x.Record
	}
	return ""
}

func (x *SyncItem) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type SyncMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastCommitted *BallotNumber `protobuf:"bytes,1,opt,name=last_committed,json=lastCommitted,proto3" json:"last_committed,omitempty"`
	Items         []*SyncItem   `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SyncMsg) Reset() {
	*x = SyncMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncMsg) ProtoMessage() {}

func (x *SyncMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncMsg.ProtoReflect.Descriptor instead.
func (*SyncMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{4}
}

func (x *SyncMsg) GetLastCommitted() *BallotNumber {
	if x != nil {
		return x.LastCommitted
	}
	return nil
}

func (x *SyncMsg) GetItems() []*SyncItem {
	if x != nil {
		return x.Items
	}
	return nil
}

// Request is used by the clients to forward a client's request to other nodes.
type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender        string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver      string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount        int64  `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	SessionId     int64  `protobuf:"varint,4,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Client        string `protobuf:"bytes,5,opt,name=client,proto3" json:"client,omitempty"`
	ReturnAddress string `protobuf:"bytes,6,opt,name=return_address,json=returnAddress,proto3" json:"return_address,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{5}
}

func (x *Request) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Request) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

func (x *Request) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Request) GetSessionId() int64 {
	if x != nil {
		return x.SessionId
	}
	return 0
}

func (x *Request) GetClient() string {
	if x != nil {
		return x.Client
	}
	return ""
}

func (x *Request) GetReturnAddress() string {
	if x != nil {
		return x.ReturnAddress
	}
	return ""
}

// accept message is sent by the leader to the followers on new request.
type AcceptMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BallotNumber *BallotNumber `protobuf:"bytes,1,opt,name=ballot_number,json=ballotNumber,proto3" json:"ballot_number,omitempty"`
	Request      *Request      `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	NodeId       string        `protobuf:"bytes,3,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	CrossShard   bool          `protobuf:"varint,4,opt,name=cross_shard,json=crossShard,proto3" json:"cross_shard,omitempty"`
}

func (x *AcceptMsg) Reset() {
	*x = AcceptMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptMsg) ProtoMessage() {}

func (x *AcceptMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptMsg.ProtoReflect.Descriptor instead.
func (*AcceptMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{6}
}

func (x *AcceptMsg) GetBallotNumber() *BallotNumber {
	if x != nil {
		return x.BallotNumber
	}
	return nil
}

func (x *AcceptMsg) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *AcceptMsg) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *AcceptMsg) GetCrossShard() bool {
	if x != nil {
		return x.CrossShard
	}
	return false
}

// accepted message is sent by the followers to the leader on new accept.
type AcceptedMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptedNumber *BallotNumber `protobuf:"bytes,1,opt,name=accepted_number,json=acceptedNumber,proto3" json:"accepted_number,omitempty"`
	AcceptedValue  *AcceptMsg    `protobuf:"bytes,2,opt,name=accepted_value,json=acceptedValue,proto3" json:"accepted_value,omitempty"`
}

func (x *AcceptedMsg) Reset() {
	*x = AcceptedMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptedMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptedMsg) ProtoMessage() {}

func (x *AcceptedMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptedMsg.ProtoReflect.Descriptor instead.
func (*AcceptedMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{7}
}

func (x *AcceptedMsg) GetAcceptedNumber() *BallotNumber {
	if x != nil {
		return x.AcceptedNumber
	}
	return nil
}

func (x *AcceptedMsg) GetAcceptedValue() *AcceptMsg {
	if x != nil {
		return x.AcceptedValue
	}
	return nil
}

// commit message is sent by the leader to followers to commit a procedure.
type CommitMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptedNumber *BallotNumber `protobuf:"bytes,1,opt,name=accepted_number,json=acceptedNumber,proto3" json:"accepted_number,omitempty"`
	AcceptedValue  *AcceptMsg    `protobuf:"bytes,2,opt,name=accepted_value,json=acceptedValue,proto3" json:"accepted_value,omitempty"`
}

func (x *CommitMsg) Reset() {
	*x = CommitMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_paxos_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitMsg) ProtoMessage() {}

func (x *CommitMsg) ProtoReflect() protoreflect.Message {
	mi := &file_paxos_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitMsg.ProtoReflect.Descriptor instead.
func (*CommitMsg) Descriptor() ([]byte, []int) {
	return file_paxos_proto_rawDescGZIP(), []int{8}
}

func (x *CommitMsg) GetAcceptedNumber() *BallotNumber {
	if x != nil {
		return x.AcceptedNumber
	}
	return nil
}

func (x *CommitMsg) GetAcceptedValue() *AcceptMsg {
	if x != nil {
		return x.AcceptedValue
	}
	return nil
}

var File_paxos_proto protoreflect.FileDescriptor

var file_paxos_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x61, 0x78, 0x6f, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x43, 0x0a, 0x0c, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x5e, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x67, 0x4d, 0x73,
	0x67, 0x12, 0x3a, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x78, 0x6f,
	0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0d,
	0x6c, 0x61, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x5e, 0x0a, 0x07, 0x50, 0x6f, 0x6e, 0x67, 0x4d, 0x73,
	0x67, 0x12, 0x3a, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x78, 0x6f,
	0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0d,
	0x6c, 0x61, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x38, 0x0a, 0x08, 0x53, 0x79, 0x6e, 0x63, 0x49, 0x74,
	0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x6c, 0x0a, 0x07, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x12, 0x3a, 0x0a, 0x0e, 0x6c,
	0x61, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c,
	0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x53,
	0x79, 0x6e, 0x63, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xb3,
	0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a,
	0x0e, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x22, 0xa9, 0x01, 0x0a, 0x09, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d,
	0x73, 0x67, 0x12, 0x38, 0x0a, 0x0d, 0x62, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x5f, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x78, 0x6f,
	0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0c,
	0x62, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x07,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x22, 0x84, 0x01, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x73, 0x67,
	0x12, 0x3c, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x78, 0x6f,
	0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0e,
	0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x37,
	0x0a, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x82, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x3c, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65,
	0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x42, 0x61, 0x6c, 0x6c, 0x6f, 0x74, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61,
	0x78, 0x6f, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x0d, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xc3, 0x02, 0x0a,
	0x05, 0x50, 0x61, 0x78, 0x6f, 0x73, 0x12, 0x30, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x0e,
	0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x4d, 0x73, 0x67, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67,
	0x12, 0x0e, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x4d, 0x73, 0x67,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04, 0x53, 0x79,
	0x6e, 0x63, 0x12, 0x0e, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x4d,
	0x73, 0x67, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x06,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x38, 0x0a, 0x08, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x12, 0x12,
	0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4d,
	0x73, 0x67, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x06,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x61, 0x78, 0x6f,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_paxos_proto_rawDescOnce sync.Once
	file_paxos_proto_rawDescData = file_paxos_proto_rawDesc
)

func file_paxos_proto_rawDescGZIP() []byte {
	file_paxos_proto_rawDescOnce.Do(func() {
		file_paxos_proto_rawDescData = protoimpl.X.CompressGZIP(file_paxos_proto_rawDescData)
	})
	return file_paxos_proto_rawDescData
}

var file_paxos_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_paxos_proto_goTypes = []any{
	(*BallotNumber)(nil),  // 0: paxos.BallotNumber
	(*PingMsg)(nil),       // 1: paxos.PingMsg
	(*PongMsg)(nil),       // 2: paxos.PongMsg
	(*SyncItem)(nil),      // 3: paxos.SyncItem
	(*SyncMsg)(nil),       // 4: paxos.SyncMsg
	(*Request)(nil),       // 5: paxos.Request
	(*AcceptMsg)(nil),     // 6: paxos.AcceptMsg
	(*AcceptedMsg)(nil),   // 7: paxos.AcceptedMsg
	(*CommitMsg)(nil),     // 8: paxos.CommitMsg
	(*emptypb.Empty)(nil), // 9: google.protobuf.Empty
}
var file_paxos_proto_depIdxs = []int32{
	0,  // 0: paxos.PingMsg.last_committed:type_name -> paxos.BallotNumber
	0,  // 1: paxos.PongMsg.last_committed:type_name -> paxos.BallotNumber
	0,  // 2: paxos.SyncMsg.last_committed:type_name -> paxos.BallotNumber
	3,  // 3: paxos.SyncMsg.items:type_name -> paxos.SyncItem
	0,  // 4: paxos.AcceptMsg.ballot_number:type_name -> paxos.BallotNumber
	5,  // 5: paxos.AcceptMsg.request:type_name -> paxos.Request
	0,  // 6: paxos.AcceptedMsg.accepted_number:type_name -> paxos.BallotNumber
	6,  // 7: paxos.AcceptedMsg.accepted_value:type_name -> paxos.AcceptMsg
	0,  // 8: paxos.CommitMsg.accepted_number:type_name -> paxos.BallotNumber
	6,  // 9: paxos.CommitMsg.accepted_value:type_name -> paxos.AcceptMsg
	1,  // 10: paxos.Paxos.Ping:input_type -> paxos.PingMsg
	2,  // 11: paxos.Paxos.Pong:input_type -> paxos.PongMsg
	4,  // 12: paxos.Paxos.Sync:input_type -> paxos.SyncMsg
	6,  // 13: paxos.Paxos.Accept:input_type -> paxos.AcceptMsg
	7,  // 14: paxos.Paxos.Accepted:input_type -> paxos.AcceptedMsg
	8,  // 15: paxos.Paxos.Commit:input_type -> paxos.CommitMsg
	9,  // 16: paxos.Paxos.Ping:output_type -> google.protobuf.Empty
	9,  // 17: paxos.Paxos.Pong:output_type -> google.protobuf.Empty
	9,  // 18: paxos.Paxos.Sync:output_type -> google.protobuf.Empty
	9,  // 19: paxos.Paxos.Accept:output_type -> google.protobuf.Empty
	9,  // 20: paxos.Paxos.Accepted:output_type -> google.protobuf.Empty
	9,  // 21: paxos.Paxos.Commit:output_type -> google.protobuf.Empty
	16, // [16:22] is the sub-list for method output_type
	10, // [10:16] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_paxos_proto_init() }
func file_paxos_proto_init() {
	if File_paxos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_paxos_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*BallotNumber); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PingMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PongMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SyncItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SyncMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*AcceptMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*AcceptedMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_paxos_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*CommitMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_paxos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_paxos_proto_goTypes,
		DependencyIndexes: file_paxos_proto_depIdxs,
		MessageInfos:      file_paxos_proto_msgTypes,
	}.Build()
	File_paxos_proto = out.File
	file_paxos_proto_rawDesc = nil
	file_paxos_proto_goTypes = nil
	file_paxos_proto_depIdxs = nil
}
