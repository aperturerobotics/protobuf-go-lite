package json

import "time"

func parseJSONTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
}
