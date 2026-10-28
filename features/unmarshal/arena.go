package unmarshal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aperturerobotics/protobuf-go-lite/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// arenaComment marks a message that gets an UnmarshalVTArena method.
const arenaComment = "protobuf-go-lite:arena"

// hasArenaComment checks if a comments section has the arena comment.
func hasArenaComment(comments protogen.Comments) bool {
	for _, line := range strings.Split(strings.TrimSuffix(string(comments), "\n"), "\n") {
		if strings.TrimSpace(line) == arenaComment {
			return true
		}
	}
	return false
}

// arena generates UnmarshalVTArena for one marked message.
//
// The generated decoder counts the messages the input holds, copies the input
// once, and carves each message from a slab of its type. Bytes fields share
// the copy. Input it does not decode exactly as UnmarshalVT would goes to
// UnmarshalVT.
type arena struct {
	*unmarshal
	// root is the marked message.
	root *protogen.Message
	// name is the arena type name.
	name string
	// types are the messages the root reaches, root first.
	types []*protogen.Message
	// index maps a message to its position in types.
	index map[protoreflect.FullName]int
	// slab marks the types a field holds, which get a slab.
	slab []bool
	// ptrs marks the types a repeated field holds, which get a pointer slab.
	ptrs []bool
}

// isArenaMessage reports whether field holds a message carved from a slab.
func isArenaMessage(field *protogen.Field) bool {
	if field.Desc.Kind() != protoreflect.MessageKind || field.Desc.IsMap() {
		return false
	}
	return field.Oneof == nil || field.Oneof.Desc.IsSynthetic()
}

// isArenaBytes reports whether field holds bytes that share the copy.
func isArenaBytes(field *protogen.Field) bool {
	if field.Desc.Kind() != protoreflect.BytesKind {
		return false
	}
	return field.Oneof == nil || field.Oneof.Desc.IsSynthetic()
}

// arenaMessage generates UnmarshalVTArena for root.
func (p *unmarshal) arenaMessage(root *protogen.Message) {
	// Collect the messages root reaches through its message fields.
	a := &arena{
		unmarshal: p,
		root:      root,
		name:      "arena" + root.GoIdent.GoName,
		index:     make(map[protoreflect.FullName]int),
	}
	a.add(root)
	for i := 0; i < len(a.types); i++ {
		for _, field := range a.types[i].Fields {
			if field.Desc.Kind() == protoreflect.GroupKind || field.Desc.Cardinality() == protoreflect.Required {
				panic(fmt.Sprintf("%s: %s reaches %s, a group or required field, which UnmarshalVTArena does not decode", arenaComment, root.Desc.FullName(), field.Desc.FullName()))
			}
			if !isArenaMessage(field) {
				continue
			}
			k := a.add(field.Message)
			a.slab[k] = true
			if field.Desc.IsList() {
				a.ptrs[k] = true
			}
		}
	}

	// Generate the arena type, the method, and a count and decode per type.
	a.arenaType()
	a.method()
	for i := range a.types {
		a.count(i)
		a.decode(i)
	}
}

// add returns the index of message, adding it when new.
func (a *arena) add(message *protogen.Message) int {
	if k, ok := a.index[message.Desc.FullName()]; ok {
		return k
	}
	k := len(a.types)
	a.index[message.Desc.FullName()] = k
	a.types = append(a.types, message)
	a.slab = append(a.slab, false)
	a.ptrs = append(a.ptrs, false)
	return k
}

// typeName returns the Go type of the message at k.
func (a *arena) typeName(k int) string {
	return a.QualifiedGoIdent(a.types[k].GoIdent)
}

// arenaType generates the type holding the slabs.
func (a *arena) arenaType() {
	n := strconv.Itoa(len(a.types))
	a.P(`// `, a.name, ` holds the slabs `, a.root.GoIdent.GoName, `.UnmarshalVTArena decodes into.`)
	a.P(`type `, a.name, ` struct {`)
	a.P(`// buf is the copy of the input that bytes fields share.`)
	a.P(`buf []byte`)
	a.P(`// n and np count the messages and repeated field elements of each type.`)
	a.P(`n, np [`, n, `]int`)
	for k := range a.types {
		if a.slab[k] {
			a.P(`t`, k, ` []`, a.typeName(k))
		}
		if a.ptrs[k] {
			a.P(`p`, k, ` []*`, a.typeName(k))
		}
	}
	a.P(`}`)
	a.P()
}

// method generates UnmarshalVTArena on the root.
func (a *arena) method() {
	name := a.root.GoIdent.GoName
	a.P(`// UnmarshalVTArena unmarshals dAtA as UnmarshalVT does, in a few allocations:`)
	a.P(`// it copies dAtA once, shares the copy with bytes fields, and carves the`)
	a.P(`// messages it holds from slabs of their types. A message that is not empty,`)
	a.P(`// an unknown or repeated singular message field, and every error go to`)
	a.P(`// UnmarshalVT. A message taken from the result keeps the slabs alive.`)
	a.P(`func (m *`, name, `) UnmarshalVTArena(dAtA []byte) error {`)
	a.P(`if `, a.notEmpty(), ` {`)
	a.P(`return m.UnmarshalVT(dAtA)`)
	a.P(`}`)
	a.P(`var a `, a.name)
	a.P(`if !a.count0(dAtA) {`)
	a.P(`return m.UnmarshalVT(dAtA)`)
	a.P(`}`)
	a.P(`a.buf = `, a.Ident("bytes", "Clone"), `(dAtA)`)
	for k := range a.types {
		if a.slab[k] {
			a.P(`if a.n[`, k, `] != 0 {`)
			a.P(`a.t`, k, ` = make([]`, a.typeName(k), `, a.n[`, k, `])`)
			a.P(`}`)
		}
		if a.ptrs[k] {
			a.P(`if a.np[`, k, `] != 0 {`)
			a.P(`a.p`, k, ` = make([]*`, a.typeName(k), `, a.np[`, k, `])`)
			a.P(`}`)
		}
	}
	a.P(`var next `, name)
	a.P(`if err := a.decode0(&next, a.buf); err != nil {`)
	a.P(`return m.UnmarshalVT(dAtA)`)
	a.P(`}`)
	a.P(`*m = next`)
	a.P(`return nil`)
	a.P(`}`)
	a.P()
}

// notEmpty returns an expression true when m holds a decoded field.
func (a *arena) notEmpty() string {
	conds := []string{"m.unknownFields != nil"}
	seenOneof := make(map[string]bool)
	for _, field := range a.root.Fields {
		if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
			if !seenOneof[field.Oneof.GoName] {
				seenOneof[field.Oneof.GoName] = true
				conds = append(conds, "m."+field.Oneof.GoName+" != nil")
			}
			continue
		}
		f := "m." + field.GoName
		switch {
		case field.Desc.IsList() || field.Desc.IsMap() || a.FieldSemantics(field).Pointer:
			conds = append(conds, f+" != nil")
		case field.Desc.Kind() == protoreflect.MessageKind || field.Desc.Kind() == protoreflect.BytesKind:
			conds = append(conds, f+" != nil")
		case field.Desc.Kind() == protoreflect.StringKind:
			conds = append(conds, f+` != ""`)
		case field.Desc.Kind() == protoreflect.BoolKind:
			conds = append(conds, f)
		default:
			conds = append(conds, f+" != 0")
		}
	}
	return strings.Join(conds, " || ")
}

// count generates the counting pass of the message at i.
func (a *arena) count(i int) {
	message := a.types[i]
	a.P(`// count`, i, ` counts the messages the `, message.GoIdent.GoName, ` in dAtA holds. It returns`)
	a.P(`// false for input UnmarshalVTArena leaves to UnmarshalVT.`)
	a.P(`func (a *`, a.name, `) count`, i, `(dAtA []byte) bool {`)

	// Number the singular message fields for the repeat check.
	singular := make(map[*protogen.Field]int)
	var others []string
	for _, field := range message.Fields {
		if isArenaMessage(field) && !field.Desc.IsList() {
			singular[field] = len(singular)
		}
		if !isArenaMessage(field) {
			others = append(others, strconv.Itoa(int(field.Desc.Number())))
		}
	}
	if len(singular) != 0 {
		a.P(`var seen [`, (len(singular)+63)/64, `]uint64`)
	}

	a.P(`for iNdEx := 0; iNdEx < len(dAtA); {`)
	if len(others) == len(message.Fields) {
		a.P(`wire, _, err := `, a.Helper("DecodeVarint"), `(dAtA, iNdEx)`)
	} else {
		a.P(`wire, next, err := `, a.Helper("DecodeVarint"), `(dAtA, iNdEx)`)
	}
	a.P(`if err != nil {`)
	a.P(`return false`)
	a.P(`}`)
	a.P(`switch wire >> 3 {`)
	for _, field := range message.Fields {
		if !isArenaMessage(field) {
			continue
		}
		k := a.index[field.Message.Desc.FullName()]
		a.P(`case `, field.Desc.Number(), `:`)
		a.P(`if wire&7 != 2 {`)
		a.P(`return false`)
		a.P(`}`)
		if bit, ok := singular[field]; ok {
			word, mask := strconv.Itoa(bit/64), fmt.Sprintf("0x%x", uint64(1)<<(bit%64))
			a.P(`if seen[`, word, `]&`, mask, ` != 0 {`)
			a.P(`return false`)
			a.P(`}`)
			a.P(`seen[`, word, `] |= `, mask)
		}
		a.P(`start, end, err := `, a.Helper("DecodeLengthDelimited"), `(dAtA, next)`)
		a.P(`if err != nil || !a.count`, k, `(dAtA[start:end]) {`)
		a.P(`return false`)
		a.P(`}`)
		a.P(`a.n[`, k, `]++`)
		if field.Desc.IsList() {
			a.P(`a.np[`, k, `]++`)
		}
		a.P(`iNdEx = end`)
	}
	if len(others) != 0 {
		a.P(`case `, strings.Join(others, ", "), `:`)
		a.P(`skippy, err := `, a.Helper("Skip"), `(dAtA[iNdEx:])`)
		a.P(`if err != nil {`)
		a.P(`return false`)
		a.P(`}`)
		a.P(`iNdEx += skippy`)
	}
	a.P(`default:`)
	a.P(`return false`)
	a.P(`}`)
	a.P(`}`)
	a.P(`return true`)
	a.P(`}`)
	a.P()
}

// decode generates the decoding pass of the message at i.
func (a *arena) decode(i int) {
	message := a.types[i]
	ccTypeName := a.typeName(i)
	a.P(`// decode`, i, ` decodes the `, message.GoIdent.GoName, ` in dAtA, a part of a.buf, into m.`)
	a.P(`func (a *`, a.name, `) decode`, i, `(m *`, ccTypeName, `, dAtA []byte) error {`)

	// Decode each field, carving messages from the slabs.
	a.P(`l := len(dAtA)`)
	a.P(`iNdEx := 0`)
	a.P(`var err error`)
	a.P(`for iNdEx < l {`)
	a.P(`var wire uint64`)
	a.decodeVarint("wire", "uint64")
	a.P(`fieldNum := int32(wire >> 3)`)
	a.P(`wireType := int(wire & 0x7)`)
	a.P(`switch fieldNum {`)
	required := message.Desc.RequiredNumbers()
	for _, field := range message.Fields {
		switch {
		case isArenaMessage(field):
			a.messageField(field)
		case isArenaBytes(field):
			a.bytesField(field)
		default:
			a.field(field, message, required)
		}
	}
	a.P(`default:`)
	a.P(`return `, fmtPackage.Ident("Errorf"), `("proto: `, message.GoIdent.GoName, `: unknown field %d in arena decode", fieldNum)`)
	a.P(`}`)
	a.P(`}`)
	a.P(`if iNdEx > l {`)
	a.P(`return `, a.Ident("io", `ErrUnexpectedEOF`))
	a.P(`}`)
	a.P(`return nil`)
	a.P(`}`)
	a.P()
}

// wireTypeCheck generates the length-delimited wire type check of field.
func (a *arena) wireTypeCheck(field *protogen.Field) {
	a.P(`if wireType != 2 {`)
	a.P(`return `, fmtPackage.Ident("Errorf"), `("proto: wrong wireType = %d for field `, field.GoName, `", wireType)`)
	a.P(`}`)
}

// messageField generates the decode of a message field carved from a slab.
func (a *arena) messageField(field *protogen.Field) {
	k := a.index[field.Message.Desc.FullName()]
	a.P(`case `, field.Desc.Number(), `:`)
	a.wireTypeCheck(field)
	a.P(`msgStart, postIndex, err := `, a.Helper("DecodeLengthDelimited"), `(dAtA, iNdEx)`)
	a.P(`if err != nil {`)
	a.P(`return err`)
	a.P(`}`)
	a.P(`if len(a.t`, k, `) == 0 {`)
	a.P(`return `, a.Helper("ErrInvalidLength"))
	a.P(`}`)
	a.P(`v := &a.t`, k, `[0]`)
	a.P(`a.t`, k, ` = a.t`, k, `[1:]`)
	a.P(`if err := a.decode`, k, `(v, dAtA[msgStart:postIndex]); err != nil {`)
	a.P(`return err`)
	a.P(`}`)
	if field.Desc.IsList() {
		// Size the slice from the pointer slab on the field's first occurrence.
		a.P(`if m.`, field.GoName, ` == nil {`)
		a.P(`n, err := `, a.Helper("CountTag"), `(dAtA[postIndex:], `, uint64(field.Desc.Number())<<3|2, `)`)
		a.P(`if err != nil {`)
		a.P(`return err`)
		a.P(`}`)
		a.P(`n++`)
		a.P(`if len(a.p`, k, `) < n {`)
		a.P(`return `, a.Helper("ErrInvalidLength"))
		a.P(`}`)
		a.P(`m.`, field.GoName, ` = a.p`, k, `[:0:n]`)
		a.P(`a.p`, k, ` = a.p`, k, `[n:]`)
		a.P(`}`)
		a.P(`m.`, field.GoName, ` = append(m.`, field.GoName, `, v)`)
	} else {
		a.P(`m.`, field.GoName, ` = v`)
	}
	a.P(`iNdEx = postIndex`)
}

// bytesField generates the decode of a bytes field sharing the copy.
func (a *arena) bytesField(field *protogen.Field) {
	a.P(`case `, field.Desc.Number(), `:`)
	a.wireTypeCheck(field)
	if field.Desc.IsList() {
		a.P(`var v []byte`)
		a.P(`v, iNdEx, err = `, a.Helper("DecodeBytesCapped"), `(dAtA, iNdEx)`)
		a.P(`if err != nil {`)
		a.P(`return err`)
		a.P(`}`)
		a.P(`m.`, field.GoName, ` = append(m.`, field.GoName, `, v)`)
		return
	}
	a.P(`m.`, field.GoName, `, iNdEx, err = `, a.Helper("DecodeBytesCapped"), `(dAtA, iNdEx)`)
	a.P(`if err != nil {`)
	a.P(`return err`)
	a.P(`}`)
}
