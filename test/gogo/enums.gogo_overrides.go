package test

import (
	"fmt"
	"strconv"
)

// MarshalJSON implements json.Marshaler interface.
func (v CustomEnum) MarshalJSON() ([]byte, error) {
	switch v {
	case CustomEnum_CUSTOM_UNKNOWN:
		return []byte(strconv.Quote("CUSTOM_UNKNOWN")), nil
	case CustomEnum_CUSTOM_V1_0:
		return []byte(strconv.Quote("1.0")), nil
	case CustomEnum_CUSTOM_V1_0_1:
		return []byte(strconv.Quote("1.0.1")), nil
	}
	return []byte(strconv.Itoa(int(v))), nil
}

func init() {
	for k, v := range CustomEnum_customvalue {
		CustomEnum_value[k] = v
	}
}

// UnmarshalJSON implements json.Marshaler interface.
func (v *CustomEnum) UnmarshalJSON(b []byte) error {
	if unquoted, err := strconv.Unquote(string(b)); err == nil {
		switch unquoted {
		case "0", "CUSTOM_UNKNOWN", "UNKNOWN":
			*v = CustomEnum_CUSTOM_UNKNOWN
		case "1.0", "1.0.0", "CUSTOM_V1_0", "V1_0":
			*v = CustomEnum_CUSTOM_V1_0
		case "1.0.1", "CUSTOM_V1_0_1", "V1_0_1":
			*v = CustomEnum_CUSTOM_V1_0
		default:
			return fmt.Errorf("invalid value: %q", unquoted)
		}
		return nil
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*v = CustomEnum(n)
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (v CustomEnumValue) MarshalJSON() ([]byte, error) {
	return v.Value.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (v *CustomEnumValue) UnmarshalJSON(b []byte) error {
	var vv CustomEnum
	if err := vv.UnmarshalJSON(b); err != nil {
		return err
	}
	*v = CustomEnumValue{
		Value: vv,
	}
	return nil
}
