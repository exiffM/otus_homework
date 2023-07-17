package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

// easyjson:json
type User struct {
	ID       int    `json:"id"`
	Name     string `jsosn:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
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
		err = easyjson.Unmarshal(scaner.Bytes(), &result[i])
		if err != nil {
			return
		}
		i++
	}
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
		ok := strings.Contains(user.Email, sb.String()) && strings.Contains(user.Email, "@")

		if ok {
			num := result[strings.ToLower(strings.Split(user.Email, "@")[1])]
			num++
			result[strings.ToLower(strings.Split(user.Email, "@")[1])] = num
		}
	}
	return result, nil
}
