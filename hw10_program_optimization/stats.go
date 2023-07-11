package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int    `json:"Id"`
	Name     string `jsosn:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
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
	scaner := bufio.NewScanner(r)
	scaner.Split(bufio.ScanLines)

	i := 0
	for scaner.Scan() {
		if err = easyjson.Unmarshal(scaner.Bytes(), &result[i]); err != nil {
			return
		}
		i++
	}

	// easyjson.UnmarshalFromReader(r, &result)

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	sb := strings.Builder{}
	if _, err := sb.WriteString("."); err != nil {
		return nil, err
	}

	if _, err := sb.WriteString(domain); err != nil {
		return nil, err
	}

	for _, user := range u {
		ok := strings.Contains(user.Email, sb.String())

		if ok {
			num := result[strings.ToLower(strings.Split(user.Email, "@")[1])]
			num++
			result[strings.ToLower(strings.Split(user.Email, "@")[1])] = num
		}
	}
	return result, nil
}
