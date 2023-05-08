package entity

import (
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Duration time.Duration

func NewSecondsDuration(seconds int) Duration {
	return Duration(time.Duration(seconds) * time.Second)
}

func NewDuration(duration time.Duration) Duration {
	return Duration(duration)
}

func (d Duration) DurationInSec() int {
	return int(time.Duration(d).Seconds())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := jsoniter.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(value)
		return nil
	case string:
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(duration)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
