package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)

	var user User
	var err error
	var emailDomain string

	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, domain) {
			emailDomain = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[emailDomain]++
		}
	}

	return result, nil
}
