package hw10programoptimization

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

const testDomain = "gmail.com"

func makeUser(i int) User {
	str := strconv.Itoa(i)

	return User{
		ID:       i,
		Name:     str,
		Username: str,
		Email:    str + "@" + testDomain,
		Phone:    str,
		Password: str,
		Address:  str,
	}
}

func BenchmarkCountDomains(b *testing.B) {
	const (
		count = 100
	)

	var testUsers users
	for i := 0; i < count; i++ {
		testUsers[i] = makeUser(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stat, err := countDomains(testUsers, testDomain)
		if err != nil {
			b.Fatal(err)
		}

		if stat[testDomain] != count {
			b.Fatal("invalid")
		}
	}
}

func BenchmarkGetUsers(b *testing.B) {
	const count = 100

	var builder strings.Builder
	for i := 0; i < count; i++ {
		data, err := json.Marshal(makeUser(i))
		if err != nil {
			b.Fatal(err)
		}

		user := string(data)

		builder.WriteString(user)
		if i != count-1 {
			builder.WriteString("\n")
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	var buff *bytes.Buffer

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buff = bytes.NewBufferString(builder.String())
		b.StartTimer()

		_, err := getUsers(buff)
		if err != nil {
			b.Fatal(err)
		}
	}
}
