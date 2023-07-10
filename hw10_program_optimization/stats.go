package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

// type users []User

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
