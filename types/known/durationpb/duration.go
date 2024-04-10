package durationpb

import (
	"strconv"
	"strings"
)

// String formats the duration to a string.
func (d *Duration) String() string {
	var out strings.Builder
	secs, nanos := d.GetSeconds(), d.GetNanos()
	if secs != 0 {
		_, _ = out.WriteString("seconds:")
		_, _ = out.WriteString(strconv.FormatInt(secs, 10))
	}
	if nanos != 0 {
		if out.Len() != 0 {
			_, _ = out.WriteString(" ")
		}
		_, _ = out.WriteString("nanos:")
		_, _ = out.WriteString(strconv.Itoa(int(nanos)))
	}
	return out.String()
}
