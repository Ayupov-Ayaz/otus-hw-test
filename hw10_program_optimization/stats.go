package hw10programoptimization

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func (u *User) Reset() {
	*u = User{}
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r) // дорого
	if err != nil {
		return
	}

	var (
		user User
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	lines := strings.Split(string(content), "\n") // сильное выделение памяти
	for i, line := range lines {
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
		user.Reset()
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	rg, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		if rg.Match([]byte(user.Email)) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
