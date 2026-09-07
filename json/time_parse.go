package json

import "time"

func parseJSONTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05.999999999Z",
		time.RFC3339Nano,
		time.RFC3339,
	}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}
