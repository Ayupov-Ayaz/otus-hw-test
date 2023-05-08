package entity

import (
	"fmt"
	"time"
)

const (
	timeLayout  = time.RFC3339
	mysqlLayout = "2006-01-02 15:04:05"
)

type MyTime time.Time

func NewMyTime(dateTime time.Time) MyTime {
	return MyTime(dateTime)
}

func ParseTime(dateTime string) (MyTime, error) {
	t, err := time.Parse(timeLayout, dateTime)
	if err != nil {
		return MyTime{}, err
	}

	return NewMyTime(t), nil
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

func (m *MyTime) MySQLFormat() string {
	return m.Time().Format(mysqlLayout)
}

func (m *MyTime) Time() time.Time {
	return time.Time(*m)
}

func (m *MyTime) IsEmpty() bool {
	return m.Time().IsZero()
}
