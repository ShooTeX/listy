package auth

import (
	"encoding/json"
	"time"
)

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*t = UnixTime(time.Time{})
		return nil
	}
	var ts int64
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}
	*t = UnixTime(time.Unix(ts, 0))
	return nil
}

func (t UnixTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Unix()
	return json.Marshal(ts)
}
