// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/lockup/lock.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// PeriodLock is a single lock unit by period defined by the x/lockup module.
// It's a record of a locked coin at a specific time. It stores owner, duration,
// unlock time and the number of coins locked. A state of a period lock is
// created upon lock creation, and deleted once the lock has been matured after
// the `duration` has passed since unbonding started.
type PeriodLock struct {
	// ID is the unique id of the lock.
	// The ID of the lock is decided upon lock creation, incrementing by 1 for
	// every lock.
	ID uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Owner is the account address of the lock owner.
	// Only the owner can modify the state of the lock.
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty" yaml:"owner"`
	// Duration is the time needed for a lock to mature after unlocking has
	// started.
	Duration time.Duration `protobuf:"bytes,3,opt,name=duration,proto3,stdduration" json:"duration,omitempty" yaml:"duration"`
	// EndTime refers to the time at which the lock would mature and get deleted.
	// This value is first initialized when an unlock has started for the lock,
	// end time being block time + duration.
	EndTime time.Time `protobuf:"bytes,4,opt,name=end_time,json=endTime,proto3,stdtime" json:"end_time" yaml:"end_time"`
	// Coins are the tokens locked within the lock, kept in the module account.
	Coins github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,5,rep,name=coins,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"coins"`
	// updated_at is the hub timestamp when the lock was created or updated.
	// adding coins to existing lock will update the timestamp of the lock
	UpdatedAt time.Time `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3,stdtime" json:"updated_at" yaml:"creation_timestamp"`
}

func (m *PeriodLock) Reset()         { *m = PeriodLock{} }
func (m *PeriodLock) String() string { return proto.CompactTextString(m) }
func (*PeriodLock) ProtoMessage()    {}
func (*PeriodLock) Descriptor() ([]byte, []int) {
	return fileDescriptor_933c4724bc61cc7c, []int{0}
}
func (m *PeriodLock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PeriodLock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PeriodLock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PeriodLock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeriodLock.Merge(m, src)
}
func (m *PeriodLock) XXX_Size() int {
	return m.Size()
}
func (m *PeriodLock) XXX_DiscardUnknown() {
	xxx_messageInfo_PeriodLock.DiscardUnknown(m)
}

var xxx_messageInfo_PeriodLock proto.InternalMessageInfo

func (m *PeriodLock) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *PeriodLock) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *PeriodLock) GetDuration() time.Duration {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *PeriodLock) GetEndTime() time.Time {
	if m != nil {
		return m.EndTime
	}
	return time.Time{}
}

func (m *PeriodLock) GetCoins() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Coins
	}
	return nil
}

func (m *PeriodLock) GetUpdatedAt() time.Time {
	if m != nil {
		return m.UpdatedAt
	}
	return time.Time{}
}

// QueryCondition is a struct used for querying locks upon different conditions.
// Duration field and lock_age field could be empty, depending on the
// LockQueryType.
type QueryCondition struct {
	// Denom represents the token denomination we are looking to lock up
	Denom string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	// Duration is used to query locks with longer duration than the specified
	// duration. Duration field must not be nil when the lock query type is
	// `ByLockDuration`.
	Duration time.Duration `protobuf:"bytes,3,opt,name=duration,proto3,stdduration" json:"duration" yaml:"duration"`
	// LockAge is used to query locks with longer lock age than the specified
	// lock age. LockAge field must not be nil when the lock query type is
	// `ByLockAge`.
	LockAge time.Duration `protobuf:"bytes,5,opt,name=lock_age,json=lockAge,proto3,stdduration" json:"lock_age" yaml:"lock_age"`
}

func (m *QueryCondition) Reset()         { *m = QueryCondition{} }
func (m *QueryCondition) String() string { return proto.CompactTextString(m) }
func (*QueryCondition) ProtoMessage()    {}
func (*QueryCondition) Descriptor() ([]byte, []int) {
	return fileDescriptor_933c4724bc61cc7c, []int{1}
}
func (m *QueryCondition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCondition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCondition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCondition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCondition.Merge(m, src)
}
func (m *QueryCondition) XXX_Size() int {
	return m.Size()
}
func (m *QueryCondition) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCondition.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCondition proto.InternalMessageInfo

func (m *QueryCondition) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *QueryCondition) GetDuration() time.Duration {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *QueryCondition) GetLockAge() time.Duration {
	if m != nil {
		return m.LockAge
	}
	return 0
}

func init() {
	proto.RegisterType((*PeriodLock)(nil), "dymensionxyz.dymension.lockup.PeriodLock")
	proto.RegisterType((*QueryCondition)(nil), "dymensionxyz.dymension.lockup.QueryCondition")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/lockup/lock.proto", fileDescriptor_933c4724bc61cc7c)
}

var fileDescriptor_933c4724bc61cc7c = []byte{
	// 518 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x31, 0x6f, 0xd3, 0x40,
	0x18, 0xcd, 0xa5, 0x71, 0xeb, 0x1e, 0xa8, 0x44, 0x56, 0x07, 0x37, 0x08, 0x3b, 0xb2, 0x04, 0xf2,
	0x00, 0x77, 0x84, 0x6c, 0x6c, 0x75, 0xb3, 0xb4, 0x02, 0x89, 0x5a, 0x4c, 0x2c, 0xc1, 0xf6, 0x1d,
	0xae, 0x95, 0xd8, 0x67, 0xd9, 0xe7, 0x52, 0xf3, 0x2b, 0x3a, 0xf2, 0x1b, 0xf8, 0x23, 0x74, 0xec,
	0x08, 0x4b, 0x8a, 0x92, 0x8d, 0xb1, 0xbf, 0x00, 0xdd, 0x9d, 0x1d, 0xaa, 0xa2, 0x0a, 0x98, 0x2e,
	0x5f, 0xbe, 0xef, 0xbd, 0x7b, 0xdf, 0x7b, 0x3e, 0xe8, 0x92, 0x3a, 0xa5, 0x59, 0x99, 0xb0, 0xec,
	0xac, 0xfe, 0x84, 0xd7, 0x05, 0x9e, 0xb3, 0x68, 0x56, 0xe5, 0xf2, 0x40, 0x79, 0xc1, 0x38, 0x33,
	0x1e, 0xdd, 0x9c, 0x44, 0xeb, 0x02, 0xa9, 0xc9, 0xc1, 0x6e, 0xcc, 0x62, 0x26, 0x27, 0xb1, 0xf8,
	0xa5, 0x40, 0x03, 0x2b, 0x66, 0x2c, 0x9e, 0x53, 0x2c, 0xab, 0xb0, 0xfa, 0x80, 0x49, 0x55, 0x04,
	0x5c, 0xc0, 0x54, 0xdf, 0xbe, 0xdd, 0xe7, 0x49, 0x4a, 0x4b, 0x1e, 0xa4, 0x79, 0x4b, 0x10, 0xb1,
	0x32, 0x65, 0x25, 0x0e, 0x83, 0x92, 0xe2, 0xd3, 0x51, 0x48, 0x79, 0x30, 0xc2, 0x11, 0x4b, 0x1a,
	0x02, 0xe7, 0xeb, 0x06, 0x84, 0x6f, 0x68, 0x91, 0x30, 0xf2, 0x8a, 0x45, 0x33, 0x63, 0x07, 0x76,
	0x0f, 0x27, 0x26, 0x18, 0x02, 0xb7, 0xe7, 0x77, 0x0f, 0x27, 0xc6, 0x13, 0xa8, 0xb1, 0x8f, 0x19,
	0x2d, 0xcc, 0xee, 0x10, 0xb8, 0xdb, 0x5e, 0xff, 0x7a, 0x61, 0xdf, 0xaf, 0x83, 0x74, 0xfe, 0xd2,
	0x91, 0x7f, 0x3b, 0xbe, 0x6a, 0x1b, 0x27, 0x50, 0x6f, 0x95, 0x99, 0x1b, 0x43, 0xe0, 0xde, 0x7b,
	0xb1, 0x87, 0x94, 0x34, 0xd4, 0x4a, 0x43, 0x93, 0x66, 0xc0, 0x1b, 0x5d, 0x2c, 0xec, 0xce, 0xcf,
	0x85, 0x6d, 0xb4, 0x90, 0xa7, 0x2c, 0x4d, 0x38, 0x4d, 0x73, 0x5e, 0x5f, 0x2f, 0xec, 0x07, 0x8a,
	0xbf, 0xed, 0x39, 0x9f, 0xaf, 0x6c, 0xe0, 0xaf, 0xd9, 0x0d, 0x1f, 0xea, 0x34, 0x23, 0x53, 0xb1,
	0xa7, 0xd9, 0x93, 0x37, 0x0d, 0xfe, 0xb8, 0xe9, 0x6d, 0x6b, 0x82, 0xf7, 0x50, 0x5c, 0xf5, 0x9b,
	0xb4, 0x45, 0x3a, 0xe7, 0x82, 0x74, 0x8b, 0x66, 0x44, 0x8c, 0x1a, 0x01, 0xd4, 0x84, 0x25, 0xa5,
	0xa9, 0x0d, 0x37, 0xa4, 0x74, 0x65, 0x1a, 0x12, 0xa6, 0xa1, 0xc6, 0x34, 0x74, 0xc0, 0x92, 0xcc,
	0x7b, 0x2e, 0xf8, 0xbe, 0x5c, 0xd9, 0x6e, 0x9c, 0xf0, 0x93, 0x2a, 0x44, 0x11, 0x4b, 0x71, 0xe3,
	0xb0, 0x3a, 0x9e, 0x95, 0x64, 0x86, 0x79, 0x9d, 0xd3, 0x52, 0x02, 0x4a, 0x5f, 0x31, 0x1b, 0xef,
	0x21, 0xac, 0x72, 0x12, 0x70, 0x4a, 0xa6, 0x01, 0x37, 0x37, 0xff, 0x2a, 0xfc, 0x71, 0x23, 0x7c,
	0x4f, 0x09, 0x8f, 0x0a, 0x2a, 0xd7, 0x9f, 0xae, 0xf3, 0x55, 0x2b, 0x6c, 0x37, 0xa4, 0xfb, 0xdc,
	0xf9, 0x0e, 0xe0, 0xce, 0x71, 0x45, 0x8b, 0xfa, 0x80, 0x65, 0x24, 0x91, 0x5e, 0xed, 0x42, 0x8d,
	0xd0, 0x8c, 0xa5, 0x2a, 0x3d, 0x5f, 0x15, 0xc2, 0xc1, 0x7f, 0xcf, 0xea, 0x96, 0x81, 0x77, 0xa5,
	0x72, 0x0c, 0x75, 0xf1, 0x1d, 0x4f, 0x83, 0x98, 0x9a, 0xda, 0x7f, 0x72, 0xb6, 0x40, 0xc5, 0xb9,
	0x25, 0xca, 0xfd, 0x98, 0x1e, 0xf5, 0x74, 0xd0, 0xef, 0x1e, 0xf5, 0xf4, 0x5e, 0x5f, 0xf3, 0x5e,
	0x5f, 0x2c, 0x2d, 0x70, 0xb9, 0xb4, 0xc0, 0x8f, 0xa5, 0x05, 0xce, 0x57, 0x56, 0xe7, 0x72, 0x65,
	0x75, 0xbe, 0xad, 0xac, 0xce, 0xbb, 0xf1, 0x8d, 0x20, 0xee, 0x78, 0x8a, 0xa7, 0x63, 0x7c, 0xd6,
	0xbe, 0x47, 0x99, 0x4c, 0xb8, 0x29, 0x35, 0x8d, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x91,
	0x65, 0xf1, 0xbd, 0x03, 0x00, 0x00,
}

func (m *PeriodLock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PeriodLock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PeriodLock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.UpdatedAt):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintLock(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x32
	if len(m.Coins) > 0 {
		for iNdEx := len(m.Coins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Coins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.EndTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintLock(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x22
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.Duration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintLock(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x1a
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintLock(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if m.ID != 0 {
		i = encodeVarintLock(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryCondition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCondition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCondition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n4, err4 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.LockAge, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.LockAge):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintLock(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x2a
	n5, err5 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.Duration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration):])
	if err5 != nil {
		return 0, err5
	}
	i -= n5
	i = encodeVarintLock(dAtA, i, uint64(n5))
	i--
	dAtA[i] = 0x1a
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintLock(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}

func encodeVarintLock(dAtA []byte, offset int, v uint64) int {
	offset -= sovLock(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PeriodLock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovLock(uint64(m.ID))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovLock(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration)
	n += 1 + l + sovLock(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.EndTime)
	n += 1 + l + sovLock(uint64(l))
	if len(m.Coins) > 0 {
		for _, e := range m.Coins {
			l = e.Size()
			n += 1 + l + sovLock(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovLock(uint64(l))
	return n
}

func (m *QueryCondition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovLock(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.Duration)
	n += 1 + l + sovLock(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.LockAge)
	n += 1 + l + sovLock(uint64(l))
	return n
}

func sovLock(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLock(x uint64) (n int) {
	return sovLock(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PeriodLock) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLock
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
			return fmt.Errorf("proto: PeriodLock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PeriodLock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.Duration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.EndTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Coins = append(m.Coins, types.Coin{})
			if err := m.Coins[len(m.Coins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLock
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
func (m *QueryCondition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLock
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
			return fmt.Errorf("proto: QueryCondition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCondition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.Duration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockAge", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLock
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
				return ErrInvalidLengthLock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.LockAge, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLock
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
func skipLock(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLock
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
					return 0, ErrIntOverflowLock
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
					return 0, ErrIntOverflowLock
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
				return 0, ErrInvalidLengthLock
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLock
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLock
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLock        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLock          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLock = fmt.Errorf("proto: unexpected end of group")
)
