package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type Filter = func(title, subject, body string) (skip bool)

var (
	DefaultFilter Filter = func(title, subject, body string) (skip bool) {
		if len(body) == 0 {
			// skip nil body
			return true
		}
		now := time.Now()
		if 0 <= now.Hour() && now.Hour() <= 7 {
			// skip 0-7
			return true
		}
		return false
	}

	user     = "notice@cheng95.com"
	password = "jackking"
)

func Send(filter Filter,
	title, subject, body string, to ...string) error {
	if filter == nil {
		return fmt.Errorf("NIL Filter")
	}
	if filter(title, subject, body) {
		return nil
	}
	if len(to) == 0 {
		return fmt.Errorf("nil To")
	}

	header := map[string]string{
		"To": func() string {
			toLocal := make([]string, len(to))
			for i, t := range to {
				toLocal[i] = fmt.Sprintf("<%s>", t)
			}
			return strings.Join(toLocal, ",")
		}(),
		"From": fmt.Sprintf("%s<%s>", title, user),
		"Subject": fmt.Sprintf("=?UTF-8?B?%s?=",
			base64.StdEncoding.EncodeToString([]byte(subject))),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/html;chartset=UTF-8",
		"Content-Transfer-Encoding": "base64",
	}

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return smtp.SendMail("service.cheng95.com:25",
		&smtpAuth{
			username: user,
			password: password,
		},
		user, to, []byte(msg))
}

var _ smtp.Auth = (*smtpAuth)(nil)

type smtpAuth struct {
	username string
	password string
}

func (a *smtpAuth) Start(server *smtp.ServerInfo) (string, []byte, error) { return "LOGIN", nil, nil }

func (a *smtpAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)
	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}
