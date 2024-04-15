package vtcompare

// EqualVT is a message with a EqualVT function (VTProtobuf).
type EqualVT[T comparable] interface {
	comparable
	// EqualVT compares against the other message for equality.
	EqualVT(other T) bool
}

// CompareEqualVT returns a compare function to compare two VTProtobuf messages.
func CompareEqualVT[T EqualVT[T]]() func(t1, t2 T) bool {
	return func(t1, t2 T) bool {
		return IsEqualVT[T](t1, t2)
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
