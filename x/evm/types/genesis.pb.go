// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stargazer/evm/v1alpha1/genesis.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	proto "github.com/cosmos/gogoproto/proto"

	_ "github.com/cosmos/gogoproto/gogoproto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the evm module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// `address_to_contract` is a map of address to contract.
	AddressToContract map[string]*Contract `protobuf:"bytes,2,rep,name=address_to_contract,json=addressToContract,proto3" json:"address_to_contract,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// `hash_to_code` is a map of code hash to code.
	HashToCode map[string]string `protobuf:"bytes,3,rep,name=hash_to_code,json=hashToCode,proto3" json:"hash_to_code,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_134c80690d273a37, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetAddressToContract() map[string]*Contract {
	if m != nil {
		return m.AddressToContract
	}
	return nil
}

func (m *GenesisState) GetHashToCode() map[string]string {
	if m != nil {
		return m.HashToCode
	}
	return nil
}

// `Contract` defines the contract state.
type Contract struct {
	// `code_hash` is the hash of the contract code.
	CodeHash string `protobuf:"bytes,1,opt,name=code_hash,json=codeHash,proto3" json:"code_hash,omitempty"`
	// `slot_to_value` is a map of slot to value.
	SlotToValue map[string]string `protobuf:"bytes,2,rep,name=slot_to_value,json=slotToValue,proto3" json:"slot_to_value,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Contract) Reset()         { *m = Contract{} }
func (m *Contract) String() string { return proto.CompactTextString(m) }
func (*Contract) ProtoMessage()    {}
func (*Contract) Descriptor() ([]byte, []int) {
	return fileDescriptor_134c80690d273a37, []int{1}
}
func (m *Contract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Contract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Contract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Contract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Contract.Merge(m, src)
}
func (m *Contract) XXX_Size() int {
	return m.Size()
}
func (m *Contract) XXX_DiscardUnknown() {
	xxx_messageInfo_Contract.DiscardUnknown(m)
}

var xxx_messageInfo_Contract proto.InternalMessageInfo

func (m *Contract) GetCodeHash() string {
	if m != nil {
		return m.CodeHash
	}
	return ""
}

func (m *Contract) GetSlotToValue() map[string]string {
	if m != nil {
		return m.SlotToValue
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "stargazer.evm.v1alpha1.GenesisState")
	proto.RegisterMapType((map[string]*Contract)(nil), "stargazer.evm.v1alpha1.GenesisState.AddressToContractEntry")
	proto.RegisterMapType((map[string]string)(nil), "stargazer.evm.v1alpha1.GenesisState.HashToCodeEntry")
	proto.RegisterType((*Contract)(nil), "stargazer.evm.v1alpha1.Contract")
	proto.RegisterMapType((map[string]string)(nil), "stargazer.evm.v1alpha1.Contract.SlotToValueEntry")
}

func init() {
	proto.RegisterFile("stargazer/evm/v1alpha1/genesis.proto", fileDescriptor_134c80690d273a37)
}

var fileDescriptor_134c80690d273a37 = []byte{
	// 416 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4f, 0x6b, 0x1a, 0x41,
	0x18, 0xc6, 0x77, 0xdd, 0x56, 0x74, 0xb4, 0xd4, 0x6e, 0x45, 0x96, 0x2d, 0x6c, 0xc5, 0x16, 0xea,
	0x69, 0x16, 0x6d, 0x29, 0xa5, 0xff, 0x40, 0x4b, 0x49, 0x8e, 0x61, 0x35, 0x1e, 0x72, 0x91, 0xd1,
	0x9d, 0xec, 0x8a, 0xeb, 0xce, 0x32, 0x33, 0x59, 0x62, 0x3e, 0x45, 0x3e, 0x4e, 0x4e, 0x39, 0x7b,
	0xf4, 0x98, 0x53, 0x08, 0xfa, 0x45, 0xc2, 0xec, 0xac, 0x7f, 0x10, 0x25, 0xc9, 0x6d, 0x78, 0x79,
	0x9e, 0xe7, 0xf7, 0xce, 0xc3, 0x0b, 0x3e, 0x33, 0x8e, 0xa8, 0x87, 0xae, 0x30, 0xb5, 0x71, 0x3c,
	0xb1, 0xe3, 0x06, 0x0a, 0x22, 0x1f, 0x35, 0x6c, 0x0f, 0x87, 0x98, 0x8d, 0x18, 0x8c, 0x28, 0xe1,
	0x44, 0xaf, 0xac, 0x55, 0x10, 0xc7, 0x13, 0xb8, 0x52, 0x99, 0x65, 0x8f, 0x78, 0x24, 0x91, 0xd8,
	0xe2, 0x25, 0xd5, 0xe6, 0xa7, 0x03, 0x99, 0x11, 0xa2, 0x68, 0x92, 0x46, 0xd6, 0x6e, 0x34, 0x50,
	0x3c, 0x92, 0x90, 0x0e, 0x47, 0x1c, 0xeb, 0xbf, 0x41, 0x56, 0x0a, 0x0c, 0xb5, 0xaa, 0xd6, 0x0b,
	0x4d, 0x0b, 0xee, 0x87, 0xc2, 0x93, 0x44, 0xd5, 0x7e, 0x35, 0xbb, 0xff, 0xa8, 0x38, 0xa9, 0x47,
	0x1f, 0x83, 0xf7, 0xc8, 0x75, 0x29, 0x66, 0xac, 0xcf, 0x49, 0x7f, 0x48, 0x42, 0x4e, 0xd1, 0x90,
	0x1b, 0x99, 0xaa, 0x56, 0x2f, 0x34, 0x7f, 0x1d, 0x8a, 0xda, 0x5e, 0x00, 0xb6, 0xa4, 0xbf, 0x4b,
	0xfe, 0xa5, 0xee, 0xff, 0x21, 0xa7, 0x53, 0xe7, 0x1d, 0xda, 0x9d, 0xeb, 0x3d, 0x50, 0xf4, 0x11,
	0xf3, 0x25, 0xc9, 0xc5, 0x86, 0x96, 0x50, 0xbe, 0x3d, 0x8b, 0x72, 0x8c, 0x98, 0x2f, 0xa2, 0x5c,
	0x2c, 0xe3, 0x81, 0xbf, 0x1e, 0x98, 0xe7, 0xa0, 0xb2, 0x7f, 0x09, 0xbd, 0x04, 0xb4, 0x31, 0x9e,
	0x26, 0xcd, 0xe4, 0x1d, 0xf1, 0xd4, 0xbf, 0x83, 0xd7, 0x31, 0x0a, 0x2e, 0xb0, 0x91, 0x49, 0xda,
	0xaa, 0x1e, 0x82, 0xaf, 0x72, 0x1c, 0x29, 0xff, 0x99, 0xf9, 0xa1, 0x9a, 0x7f, 0xc0, 0xdb, 0x9d,
	0x35, 0xf6, 0x00, 0xca, 0xdb, 0x80, 0xfc, 0x96, 0xbd, 0x76, 0xab, 0x82, 0xdc, 0xba, 0x8b, 0x0f,
	0x20, 0x2f, 0x3a, 0xe8, 0x8b, 0x6f, 0xa4, 0xf6, 0x9c, 0x18, 0x08, 0x80, 0x7e, 0x0a, 0xde, 0xb0,
	0x80, 0x70, 0x51, 0xd4, 0x2a, 0x4b, 0x34, 0xd5, 0x78, 0x6a, 0x59, 0xd8, 0x09, 0x08, 0xef, 0x92,
	0x9e, 0xf0, 0xc8, 0x9a, 0x0a, 0x6c, 0x33, 0x31, 0xff, 0x82, 0xd2, 0xae, 0xe0, 0x25, 0x1f, 0x68,
	0xb7, 0x66, 0x0b, 0x4b, 0x9d, 0x2f, 0x2c, 0xf5, 0x61, 0x61, 0xa9, 0xd7, 0x4b, 0x4b, 0x99, 0x2f,
	0x2d, 0xe5, 0x6e, 0x69, 0x29, 0x67, 0x5f, 0xa2, 0xb1, 0x07, 0x07, 0x98, 0xa2, 0xa1, 0x8f, 0x46,
	0x21, 0x74, 0x71, 0x6c, 0x6f, 0x8e, 0xf9, 0x32, 0x39, 0x67, 0x3e, 0x8d, 0x30, 0x1b, 0x64, 0x93,
	0x2b, 0xfe, 0xfa, 0x18, 0x00, 0x00, 0xff, 0xff, 0x2a, 0xc1, 0xcb, 0xc6, 0x40, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.HashToCode) > 0 {
		for k := range m.HashToCode {
			v := m.HashToCode[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintGenesis(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintGenesis(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AddressToContract) > 0 {
		for k := range m.AddressToContract {
			v := m.AddressToContract[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintGenesis(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintGenesis(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Contract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Contract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Contract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SlotToValue) > 0 {
		for k := range m.SlotToValue {
			v := m.SlotToValue[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintGenesis(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintGenesis(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.CodeHash) > 0 {
		i -= len(m.CodeHash)
		copy(dAtA[i:], m.CodeHash)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.CodeHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.AddressToContract) > 0 {
		for k, v := range m.AddressToContract {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovGenesis(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovGenesis(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	if len(m.HashToCode) > 0 {
		for k, v := range m.HashToCode {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovGenesis(uint64(len(k))) + 1 + len(v) + sovGenesis(uint64(len(v)))
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *Contract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CodeHash)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.SlotToValue) > 0 {
		for k, v := range m.SlotToValue {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovGenesis(uint64(len(k))) + 1 + len(v) + sovGenesis(uint64(len(v)))
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressToContract", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.AddressToContract == nil {
				m.AddressToContract = make(map[string]*Contract)
			}
			var mapkey string
			var mapvalue *Contract
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthGenesis
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthGenesis
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &Contract{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.AddressToContract[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HashToCode", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HashToCode == nil {
				m.HashToCode = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.HashToCode[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *Contract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: Contract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Contract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CodeHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlotToValue", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SlotToValue == nil {
				m.SlotToValue = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenesis
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthGenesis
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthGenesis
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.SlotToValue[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
