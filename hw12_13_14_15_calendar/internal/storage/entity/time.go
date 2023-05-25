package entity

import (
	"fmt"
	"time"
)

const (
	timeLayout = time.RFC3339
)

type MyTime time.Time

func NewMyTime(dateTime time.Time) MyTime {
	return MyTime(dateTime)
}

func (m *MyTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}

	dateTime, err := time.Parse(`"`+timeLayout+`"`, str)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*m = MyTime(dateTime)

	return nil
}

func (m *MyTime) Time() time.Time {
	return time.Time(*m)
}

func (m *MyTime) IsEmpty() bool {
	return m.Time().IsZero()
}

func (m *MyTime) MarshalJSON() ([]byte, error) {
	if m.IsEmpty() {
		return []byte(`""`), nil
	}

	return []byte(`"` + m.Time().Format(timeLayout) + `"`), nil
}
