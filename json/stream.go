package json

import (
	"io"
	"strconv"
	"unicode/utf8"
)

const jsonHex = "0123456789abcdef"

func appendJSONString(dst []byte, str string) []byte {
	dst = append(dst, '"')
	for len(str) != 0 {
		r, size := utf8.DecodeRuneInString(str)
		str = str[size:]
		switch r {
		case '"', '\\':
			dst = append(dst, '\\', byte(r))
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		case '\u2028', '\u2029':
			dst = append(dst, '\\', 'u', '2', '0', '2', jsonHex[r&0xf])
		default:
			if r < ' ' {
				dst = append(dst, '\\', 'u', '0', '0', jsonHex[r>>4], jsonHex[r&0xf])
			} else {
				dst = utf8.AppendRune(dst, r)
			}
		}
	}
	return append(dst, '"')
}

// JsonStream is an outgoing stream of json.
type JsonStream struct {
	wr  io.Writer
	err error
}

// NewJsonStream creates a new JsonStream that writes to wr.
func NewJsonStream(wr io.Writer) *JsonStream {
	return &JsonStream{wr: wr}
}

// Write writes the contents of p into the stream.
// It returns the number of bytes written and any write error encountered.
func (s *JsonStream) Write(p []byte) (n int, err error) {
	if s.err != nil {
		return 0, s.err
	}
	n, err = s.wr.Write(p)
	if err != nil {
		s.err = err
	}
	return n, err
}

// WriteString writes a quoted string into the stream.
func (s *JsonStream) WriteString(str string) {
	if s.err == nil {
		encoded := appendJSONString(make([]byte, 0, len(str)+2), str)
		_, s.err = s.wr.Write(encoded)
	}
}

// WriteFloat32 writes a float32 value into the stream.
func (s *JsonStream) WriteFloat32(f float32) {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, strconv.FormatFloat(float64(f), 'f', -1, 32))
	}
}

// WriteFloat64 writes a float64 value into the stream.
func (s *JsonStream) WriteFloat64(f float64) {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, strconv.FormatFloat(f, 'f', -1, 64))
	}
}

// WriteInt32 writes an int32 value into the stream.
func (s *JsonStream) WriteInt32(i int32) {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, strconv.FormatInt(int64(i), 10))
	}
}

// WriteUint32 writes a uint32 value into the stream.
func (s *JsonStream) WriteUint32(u uint32) {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, strconv.FormatUint(uint64(u), 10))
	}
}

// WriteBool writes a boolean value into the stream.
func (s *JsonStream) WriteBool(b bool) {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, strconv.FormatBool(b))
	}
}

// WriteNil writes a null value into the stream.
func (s *JsonStream) WriteNil() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, "null")
	}
}

// WriteObjectStart writes the start of a JSON object into the stream.
func (s *JsonStream) WriteObjectStart() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, "{")
	}
}

// WriteObjectField writes a field name into the stream.
func (s *JsonStream) WriteObjectField(field string) {
	if s.err != nil {
		return
	}
	s.WriteString(field)
	if s.err != nil {
		return
	}
	_, s.err = io.WriteString(s.wr, ":")
}

// WriteObjectEnd writes the end of a JSON object into the stream.
func (s *JsonStream) WriteObjectEnd() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, "}")
	}
}

// WriteArrayStart writes the start of a JSON array into the stream.
func (s *JsonStream) WriteArrayStart() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, "[")
	}
}

// WriteArrayEnd writes the end of a JSON array into the stream.
func (s *JsonStream) WriteArrayEnd() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, "]")
	}
}

// WriteMore writes a comma to separate elements in the stream.
func (s *JsonStream) WriteMore() {
	if s.err == nil {
		_, s.err = io.WriteString(s.wr, ",")
	}
}

// Error returns any error that has occurred on the stream.
func (s *JsonStream) Error() error {
	return s.err
}
