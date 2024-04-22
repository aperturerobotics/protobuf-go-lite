package protobuf_go_lite

import (
	"io"
	"math/bits"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidLength is returned when decoding a negative length.
	ErrInvalidLength = errors.New("proto: negative length found during unmarshaling")
	// ErrIntOverflow is returned when decoding a varint representation of an integer that overflows 64 bits.
	ErrIntOverflow = errors.New("proto: integer overflow")
	// ErrUnexpectedEndOfGroup is returned when decoding a group end without a corresponding group start.
	ErrUnexpectedEndOfGroup = errors.New("proto: unexpected end of group")
)

// Message is the base vtprotobuf message marshal/unmarshal interface.
type Message interface {
	// MarshalVT marshals the message with vtprotobuf.
	MarshalVT() ([]byte, error)
	// UnmarshalVT unmarshals the message object with vtprotobuf.
	UnmarshalVT(data []byte) error
	// Reset resets the message.
	Reset()
}

// EqualVT is a message with a EqualVT function (VTProtobuf).
type EqualVT[T comparable] interface {
	comparable
	// EqualVT compares against the other message for equality.
	EqualVT(other T) bool
}

// CompareEqualVT returns a compare function to compare two VTProtobuf messages.
func CompareEqualVT[T EqualVT[T]]() func(t1, t2 T) bool {
	return func(t1, t2 T) bool {
		return IsEqualVT(t1, t2)
	}
}

// CompareComparable returns a compare function to compare two comparable types.
func CompareComparable[T comparable]() func(t1, t2 T) bool {
	return func(t1, t2 T) bool {
		return t1 == t2
	}
}

// IsEqualVT checks if two EqualVT objects are equal.
func IsEqualVT[T EqualVT[T]](t1, t2 T) bool {
	var empty T
	t1Empty, t2Empty := t1 == empty, t2 == empty
	if t1Empty != t2Empty {
		return false
	}
	if t1Empty {
		return true
	}
	return t1.EqualVT(t2)
}

// EncodeVarint encodes a uint64 into a varint-encoded byte slice and returns the offset of the encoded value.
// The provided offset is the offset after the last byte of the encoded value.
func EncodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= SizeOfVarint(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

// AppendVarint appends v to b as a varint-encoded uint64.
func AppendVarint(b []byte, v uint64) []byte {
	switch {
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	}
	return b
}

// SizeOfVarint returns the size of the varint-encoded value.
func SizeOfVarint(x uint64) (n int) {
	return (bits.Len64(x|1) + 6) / 7
}

// SizeOfZigzag returns the size of the zigzag-encoded value.
func SizeOfZigzag(x uint64) (n int) {
	return SizeOfVarint(uint64((x << 1) ^ uint64((int64(x) >> 63)))) //nolint
}

// Skip the first record of the byte slice and return the offset of the next record.
func Skip(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflow
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
					return 0, ErrIntOverflow
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
					return 0, ErrIntOverflow
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
				return 0, ErrInvalidLength
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroup
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, errors.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLength
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}
