package entity

import (
	"fmt"
	"time"
)

const timeLayout = time.RFC3339

type MyTime time.Time

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
