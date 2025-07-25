// Code generated by protoc-gen-go-lite. DO NOT EDIT.
// protoc-gen-go-lite version: v0.8.1
// source: github.com/aperturerobotics/protobuf-go-lite/testproto/maps/maps.proto

package testproto_maps

import (
	fmt "fmt"
	io "io"
	strconv "strconv"
	strings "strings"
	unsafe "unsafe"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	json "github.com/aperturerobotics/protobuf-go-lite/json"
	timestamppb "github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
)

type MsgWithMaps struct {
	unknownFields []byte
	StringKeys    map[string]*timestamppb.Timestamp `protobuf:"bytes,1,rep,name=stringKeys,proto3" json:"stringKeys,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	IntKeys       map[uint32]*timestamppb.Timestamp `protobuf:"bytes,2,rep,name=intKeys,proto3" json:"intKeys,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MsgWithMaps) Reset() {
	*x = MsgWithMaps{}
}

func (*MsgWithMaps) ProtoMessage() {}

func (x *MsgWithMaps) GetStringKeys() map[string]*timestamppb.Timestamp {
	if x != nil {
		return x.StringKeys
	}
	return nil
}

func (x *MsgWithMaps) GetIntKeys() map[uint32]*timestamppb.Timestamp {
	if x != nil {
		return x.IntKeys
	}
	return nil
}

type MsgWithMaps_StringKeysEntry struct {
	unknownFields []byte
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MsgWithMaps_StringKeysEntry) Reset() {
	*x = MsgWithMaps_StringKeysEntry{}
}

func (*MsgWithMaps_StringKeysEntry) ProtoMessage() {}

func (x *MsgWithMaps_StringKeysEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *MsgWithMaps_StringKeysEntry) GetValue() *timestamppb.Timestamp {
	if x != nil {
		return x.Value
	}
	return nil
}

type MsgWithMaps_IntKeysEntry struct {
	unknownFields []byte
	Key           uint32                 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MsgWithMaps_IntKeysEntry) Reset() {
	*x = MsgWithMaps_IntKeysEntry{}
}

func (*MsgWithMaps_IntKeysEntry) ProtoMessage() {}

func (x *MsgWithMaps_IntKeysEntry) GetKey() uint32 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *MsgWithMaps_IntKeysEntry) GetValue() *timestamppb.Timestamp {
	if x != nil {
		return x.Value
	}
	return nil
}

func (m *MsgWithMaps) CloneVT() *MsgWithMaps {
	if m == nil {
		return (*MsgWithMaps)(nil)
	}
	r := new(MsgWithMaps)
	if rhs := m.StringKeys; rhs != nil {
		tmpContainer := make(map[string]*timestamppb.Timestamp, len(rhs))
		for k, v := range rhs {
			tmpContainer[k] = v.CloneVT()
		}
		r.StringKeys = tmpContainer
	}
	if rhs := m.IntKeys; rhs != nil {
		tmpContainer := make(map[uint32]*timestamppb.Timestamp, len(rhs))
		for k, v := range rhs {
			tmpContainer[k] = v.CloneVT()
		}
		r.IntKeys = tmpContainer
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *MsgWithMaps) CloneMessageVT() protobuf_go_lite.CloneMessage {
	return m.CloneVT()
}

func (this *MsgWithMaps) EqualVT(that *MsgWithMaps) bool {
	if this == that {
		return true
	} else if this == nil || that == nil {
		return false
	}
	if len(this.StringKeys) != len(that.StringKeys) {
		return false
	}
	for i, vx := range this.StringKeys {
		vy, ok := that.StringKeys[i]
		if !ok {
			return false
		}
		if p, q := vx, vy; p != q {
			if p == nil {
				p = &timestamppb.Timestamp{}
			}
			if q == nil {
				q = &timestamppb.Timestamp{}
			}
			if !p.EqualVT(q) {
				return false
			}
		}
	}
	if len(this.IntKeys) != len(that.IntKeys) {
		return false
	}
	for i, vx := range this.IntKeys {
		vy, ok := that.IntKeys[i]
		if !ok {
			return false
		}
		if p, q := vx, vy; p != q {
			if p == nil {
				p = &timestamppb.Timestamp{}
			}
			if q == nil {
				q = &timestamppb.Timestamp{}
			}
			if !p.EqualVT(q) {
				return false
			}
		}
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *MsgWithMaps) EqualMessageVT(thatMsg any) bool {
	that, ok := thatMsg.(*MsgWithMaps)
	if !ok {
		return false
	}
	return this.EqualVT(that)
}

// MarshalProtoJSON marshals the MsgWithMaps_StringKeysEntry message to JSON.
func (x *MsgWithMaps_StringKeysEntry) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.Key != "" || s.HasField("key") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("key")
		s.WriteString(x.Key)
	}
	if x.Value != nil || s.HasField("value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("value")
		x.Value.MarshalProtoJSON(s.WithField("value"))
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MsgWithMaps_StringKeysEntry to JSON.
func (x *MsgWithMaps_StringKeysEntry) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalProtoJSON unmarshals the MsgWithMaps_StringKeysEntry message from JSON.
func (x *MsgWithMaps_StringKeysEntry) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.Skip() // ignore unknown field
		case "key":
			s.AddField("key")
			x.Key = s.ReadString()
		case "value":
			if s.ReadNil() {
				x.Value = nil
				return
			}
			x.Value = &timestamppb.Timestamp{}
			x.Value.UnmarshalProtoJSON(s.WithField("value", true))
		}
	})
}

// UnmarshalJSON unmarshals the MsgWithMaps_StringKeysEntry from JSON.
func (x *MsgWithMaps_StringKeysEntry) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals the MsgWithMaps_IntKeysEntry message to JSON.
func (x *MsgWithMaps_IntKeysEntry) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.Key != 0 || s.HasField("key") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("key")
		s.WriteUint32(x.Key)
	}
	if x.Value != nil || s.HasField("value") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("value")
		x.Value.MarshalProtoJSON(s.WithField("value"))
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MsgWithMaps_IntKeysEntry to JSON.
func (x *MsgWithMaps_IntKeysEntry) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalProtoJSON unmarshals the MsgWithMaps_IntKeysEntry message from JSON.
func (x *MsgWithMaps_IntKeysEntry) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.Skip() // ignore unknown field
		case "key":
			s.AddField("key")
			x.Key = s.ReadUint32()
		case "value":
			if s.ReadNil() {
				x.Value = nil
				return
			}
			x.Value = &timestamppb.Timestamp{}
			x.Value.UnmarshalProtoJSON(s.WithField("value", true))
		}
	})
}

// UnmarshalJSON unmarshals the MsgWithMaps_IntKeysEntry from JSON.
func (x *MsgWithMaps_IntKeysEntry) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

// MarshalProtoJSON marshals the MsgWithMaps message to JSON.
func (x *MsgWithMaps) MarshalProtoJSON(s *json.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.StringKeys != nil || s.HasField("stringKeys") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("stringKeys")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.StringKeys {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectStringField(k)
			v.MarshalProtoJSON(s.WithField("stringKeys"))
		}
		s.WriteObjectEnd()
	}
	if x.IntKeys != nil || s.HasField("intKeys") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("intKeys")
		s.WriteObjectStart()
		var wroteElement bool
		for k, v := range x.IntKeys {
			s.WriteMoreIf(&wroteElement)
			s.WriteObjectUint32Field(k)
			v.MarshalProtoJSON(s.WithField("intKeys"))
		}
		s.WriteObjectEnd()
	}
	s.WriteObjectEnd()
}

// MarshalJSON marshals the MsgWithMaps to JSON.
func (x *MsgWithMaps) MarshalJSON() ([]byte, error) {
	return json.DefaultMarshalerConfig.Marshal(x)
}

// UnmarshalProtoJSON unmarshals the MsgWithMaps message from JSON.
func (x *MsgWithMaps) UnmarshalProtoJSON(s *json.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.Skip() // ignore unknown field
		case "stringKeys":
			s.AddField("stringKeys")
			if s.ReadNil() {
				x.StringKeys = nil
				return
			}
			x.StringKeys = make(map[string]*timestamppb.Timestamp)
			s.ReadStringMap(func(key string) {
				var v timestamppb.Timestamp
				v.UnmarshalProtoJSON(s)
				x.StringKeys[key] = &v
			})
		case "intKeys":
			s.AddField("intKeys")
			if s.ReadNil() {
				x.IntKeys = nil
				return
			}
			x.IntKeys = make(map[uint32]*timestamppb.Timestamp)
			s.ReadUint32Map(func(key uint32) {
				var v timestamppb.Timestamp
				v.UnmarshalProtoJSON(s)
				x.IntKeys[key] = &v
			})
		}
	})
}

// UnmarshalJSON unmarshals the MsgWithMaps from JSON.
func (x *MsgWithMaps) UnmarshalJSON(b []byte) error {
	return json.DefaultUnmarshalerConfig.Unmarshal(b, x)
}

func (m *MsgWithMaps) MarshalVT() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVT(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithMaps) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *MsgWithMaps) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if len(m.IntKeys) > 0 {
		for k := range m.IntKeys {
			v := m.IntKeys[k]
			baseI := i
			size, err := v.MarshalToSizedBufferVT(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(size))
			i--
			dAtA[i] = 0x12
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.StringKeys) > 0 {
		for k := range m.StringKeys {
			v := m.StringKeys[k]
			baseI := i
			size, err := v.MarshalToSizedBufferVT(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(size))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithMaps) MarshalVTStrict() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVTStrict(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithMaps) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *MsgWithMaps) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if len(m.IntKeys) > 0 {
		for k := range m.IntKeys {
			v := m.IntKeys[k]
			baseI := i
			size, err := v.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(size))
			i--
			dAtA[i] = 0x12
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.StringKeys) > 0 {
		for k := range m.StringKeys {
			v := m.StringKeys[k]
			baseI := i
			size, err := v.MarshalToSizedBufferVTStrict(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(size))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = protobuf_go_lite.EncodeVarint(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithMaps) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.StringKeys) > 0 {
		for k, v := range m.StringKeys {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.SizeVT()
			}
			l += 1 + protobuf_go_lite.SizeOfVarint(uint64(l))
			mapEntrySize := 1 + len(k) + protobuf_go_lite.SizeOfVarint(uint64(len(k))) + l
			n += mapEntrySize + 1 + protobuf_go_lite.SizeOfVarint(uint64(mapEntrySize))
		}
	}
	if len(m.IntKeys) > 0 {
		for k, v := range m.IntKeys {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.SizeVT()
			}
			l += 1 + protobuf_go_lite.SizeOfVarint(uint64(l))
			mapEntrySize := 1 + protobuf_go_lite.SizeOfVarint(uint64(k)) + l
			n += mapEntrySize + 1 + protobuf_go_lite.SizeOfVarint(uint64(mapEntrySize))
		}
	}
	n += len(m.unknownFields)
	return n
}

func (x *MsgWithMaps_StringKeysEntry) MarshalProtoText() string {
	var sb strings.Builder
	sb.WriteString("StringKeysEntry {")
	if x.Key != "" {
		if sb.Len() > 17 {
			sb.WriteString(" ")
		}
		sb.WriteString("key: ")
		sb.WriteString(strconv.Quote(x.Key))
	}
	if x.Value != nil {
		if sb.Len() > 17 {
			sb.WriteString(" ")
		}
		sb.WriteString("value: ")
		sb.WriteString(x.Value.MarshalProtoText())
	}
	sb.WriteString("}")
	return sb.String()
}

func (x *MsgWithMaps_StringKeysEntry) String() string {
	return x.MarshalProtoText()
}
func (x *MsgWithMaps_IntKeysEntry) MarshalProtoText() string {
	var sb strings.Builder
	sb.WriteString("IntKeysEntry {")
	if x.Key != 0 {
		if sb.Len() > 14 {
			sb.WriteString(" ")
		}
		sb.WriteString("key: ")
		sb.WriteString(strconv.FormatUint(uint64(x.Key), 10))
	}
	if x.Value != nil {
		if sb.Len() > 14 {
			sb.WriteString(" ")
		}
		sb.WriteString("value: ")
		sb.WriteString(x.Value.MarshalProtoText())
	}
	sb.WriteString("}")
	return sb.String()
}

func (x *MsgWithMaps_IntKeysEntry) String() string {
	return x.MarshalProtoText()
}
func (x *MsgWithMaps) MarshalProtoText() string {
	var sb strings.Builder
	sb.WriteString("MsgWithMaps {")
	if len(x.StringKeys) > 0 {
		if sb.Len() > 13 {
			sb.WriteString(" ")
		}
		sb.WriteString("stringKeys: {")
		for k, v := range x.StringKeys {
			sb.WriteString(" ")
			sb.WriteString(strconv.Quote(k))
			sb.WriteString(": ")
			sb.WriteString(v.MarshalProtoText())
		}
		sb.WriteString(" }")
	}
	if len(x.IntKeys) > 0 {
		if sb.Len() > 13 {
			sb.WriteString(" ")
		}
		sb.WriteString("intKeys: {")
		for k, v := range x.IntKeys {
			sb.WriteString(" ")
			sb.WriteString(strconv.FormatUint(uint64(k), 10))
			sb.WriteString(": ")
			sb.WriteString(v.MarshalProtoText())
		}
		sb.WriteString(" }")
	}
	sb.WriteString("}")
	return sb.String()
}

func (x *MsgWithMaps) String() string {
	return x.MarshalProtoText()
}
func (m *MsgWithMaps) UnmarshalVT(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return protobuf_go_lite.ErrIntOverflow
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
			return fmt.Errorf("proto: MsgWithMaps: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithMaps: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StringKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protobuf_go_lite.ErrIntOverflow
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
				return protobuf_go_lite.ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StringKeys == nil {
				m.StringKeys = make(map[string]*timestamppb.Timestamp)
			}
			var mapkey string
			var mapvalue *timestamppb.Timestamp
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protobuf_go_lite.ErrIntOverflow
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
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return protobuf_go_lite.ErrInvalidLength
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
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &timestamppb.Timestamp{}
					if err := mapvalue.UnmarshalVT(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.StringKeys[mapkey] = mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IntKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protobuf_go_lite.ErrIntOverflow
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
				return protobuf_go_lite.ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IntKeys == nil {
				m.IntKeys = make(map[uint32]*timestamppb.Timestamp)
			}
			var mapkey uint32
			var mapvalue *timestamppb.Timestamp
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protobuf_go_lite.ErrIntOverflow
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protobuf_go_lite.ErrIntOverflow
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &timestamppb.Timestamp{}
					if err := mapvalue.UnmarshalVT(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.IntKeys[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgWithMaps) UnmarshalVTUnsafe(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return protobuf_go_lite.ErrIntOverflow
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
			return fmt.Errorf("proto: MsgWithMaps: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithMaps: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StringKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protobuf_go_lite.ErrIntOverflow
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
				return protobuf_go_lite.ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StringKeys == nil {
				m.StringKeys = make(map[string]*timestamppb.Timestamp)
			}
			var mapkey string
			var mapvalue *timestamppb.Timestamp
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protobuf_go_lite.ErrIntOverflow
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
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					if intStringLenmapkey == 0 {
						mapkey = ""
					} else {
						mapkey = unsafe.String(&dAtA[iNdEx], intStringLenmapkey)
					}
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &timestamppb.Timestamp{}
					if err := mapvalue.UnmarshalVTUnsafe(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.StringKeys[mapkey] = mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IntKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protobuf_go_lite.ErrIntOverflow
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
				return protobuf_go_lite.ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IntKeys == nil {
				m.IntKeys = make(map[uint32]*timestamppb.Timestamp)
			}
			var mapkey uint32
			var mapvalue *timestamppb.Timestamp
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protobuf_go_lite.ErrIntOverflow
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protobuf_go_lite.ErrIntOverflow
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protobuf_go_lite.ErrIntOverflow
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
						return protobuf_go_lite.ErrInvalidLength
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &timestamppb.Timestamp{}
					if err := mapvalue.UnmarshalVTUnsafe(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return protobuf_go_lite.ErrInvalidLength
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.IntKeys[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := protobuf_go_lite.Skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return protobuf_go_lite.ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
