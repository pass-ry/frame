package email

import (
	stdMail "net/mail"

	"github.com/pkg/errors"
	otherSMTP "github.com/scorredoira/email"
)

func SendWithAttachment(filter Filter,
	title, subject, body string,
	to []string, cc []string,
	attach map[string][]byte) error {
	if filter == nil {
		return errors.Errorf("NIL Filter")
	}
	if filter(title, subject, body) {
		return nil
	}
	if len(to) == 0 {
		return errors.Errorf("nil Tos")
	}

	mail := otherSMTP.NewHTMLMessage(subject, body)

	mail.From = stdMail.Address{
		Name:    title,
		Address: user,
	}

	for _, t := range to {
		mail.AddTo(stdMail.Address{Address: t})
	}

	for _, c := range cc {
		mail.AddCc(stdMail.Address{Address: c})
	}

	for aName, a := range attach {
		err := mail.AttachBuffer(aName, a, false)
		if err != nil {
			return errors.Wrapf(err, "AttachBuffer %s", aName)
		}
	}

	mail.AddHeader("X-CUSTOMER-id", "xxxxx")

	return otherSMTP.Send("service.cheng95.com:25",
		&smtpAuth{
			username: user,
			password: password,
		}, mail)
}
