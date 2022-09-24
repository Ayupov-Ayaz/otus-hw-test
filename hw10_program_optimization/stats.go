package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
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
	scan := bufio.NewScanner(r)
	scan.Split(bufio.ScanLines)

	var (
		user User
		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	i := 0
	for scan.Scan() {
		if err = json.Unmarshal(scan.Bytes(), &user); err != nil {
			return
		}

		result[i] = user
		i++
		user.Reset()
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for i := 0; i < len(u); i++ {
		user := u[i]
		if user.Email == "" {
			break
		}

		if strings.HasSuffix(user.Email, domain) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}
	return result, nil
}
