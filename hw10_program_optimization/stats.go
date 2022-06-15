package hw10programoptimization

import (
	"bufio"
	"errors"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	lines, readErrors := readLines(r)
	emails := extractEmail(lines)
	domains := extractEmailDomain(emails, domain)

	ds := make(DomainStat)
	for {
		select {
		case err := <-readErrors:
			if err != nil {
				return nil, err
			}
		case d, ok := <-domains:
			if !ok {
				return ds, nil
			}
			ds[d]++
		}
	}
}

func readLines(r io.Reader) (<-chan []byte, <-chan error) {
	in := make(chan []byte)
	errCh := make(chan error)
	rd := bufio.NewReader(r)
	go func() {
		defer func() {
			close(in)
			close(errCh)
		}()
		for {
			line, _, err := rd.ReadLine()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				errCh <- err
				return
			}
			b := make([]byte, len(line))
			copy(b, line)
			in <- b
		}
	}()
	return in, errCh
}

func extractEmail(lines <-chan []byte) <-chan string {
	in := make(chan string)
	go func() {
		defer close(in)
		for line := range lines {
			email := jsoniter.Get(line, "Email").ToString()
			if email != "" {
				in <- email
			}
		}
	}()
	return in
}

func extractEmailDomain(emails <-chan string, domain string) <-chan string {
	in := make(chan string)
	go func() {
		defer close(in)
		for email := range emails {
			if strings.HasSuffix(email, domain) {
				parts := strings.SplitAfterN(email, "@", 2)
				if len(parts) > 1 {
					in <- strings.ToLower(parts[1])
				}
			}
		}
	}()
	return in
}
