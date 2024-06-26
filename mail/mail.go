package mail

import (
	"github.com/sendinblue/APIv3-go-library/lib"
	"gopkg.in/guregu/null.v4"
)

// GenericReceipient is a generic receipient format to be used in mail
type GenericReceipient struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

// Mail is datatype for mail
type Mail struct {
	ID          string              `json:"id"`
	To          []GenericReceipient `json:"to"`
	Cc          []GenericReceipient `json:"cc,omitempty"`
	Bcc         []GenericReceipient `json:"bcc,omitempty"`
	HTMLContent string              `json:"html_content"`
	Subject     string              `json:"subject"`
	Metadata    null.String         `json:"metadata,omitempty"`
}

// SendInBlueTo get send in blue SendSmtpEmailTo
func (m *Mail) SendInBlueTo() []lib.SendSmtpEmailTo {
	var to []lib.SendSmtpEmailTo

	for _, t := range m.To {
		to = append(to, lib.SendSmtpEmailTo{
			Email: t.Email,
			Name:  t.Name,
		})
	}

	return to
}

// SendInBlueCc get send in blue SendSmtpEmailCc
func (m *Mail) SendInBlueCc() []lib.SendSmtpEmailCc {
	var cc []lib.SendSmtpEmailCc

	for _, c := range m.Cc {
		cc = append(cc, lib.SendSmtpEmailCc{
			Email: c.Email,
			Name:  c.Name,
		})
	}

	return cc
}

// SendInBlueBcc get send in blue SendSmtpEmailBcc
func (m *Mail) SendInBlueBcc() []lib.SendSmtpEmailBcc {
	var bcc []lib.SendSmtpEmailBcc

	for _, b := range m.Bcc {
		bcc = append(bcc, lib.SendSmtpEmailBcc{
			Email: b.Email,
			Name:  b.Name,
		})
	}

	return bcc
}

// MailgunTo convert to to mailgun compatible to
func (m *Mail) MailgunTo() []string {
	var res []string
	for _, t := range m.To {
		res = append(res, t.Email)
	}

	return res
}

// MailgunCC convert cc to mailgun compatible cc
func (m *Mail) MailgunCC() []string {
	var cc []string
	for _, c := range m.Cc {
		cc = append(cc, c.Email)
	}

	return cc
}

// MailgunBCC convert bcc to mailgun compatible bcc
func (m *Mail) MailgunBCC() []string {
	var bcc []string
	for _, c := range m.Bcc {
		bcc = append(bcc, c.Email)
	}

	return bcc
}

func (m *Mail) ResendTo() []string {
	var tos []string
	for _, to := range m.To {
		tos = append(tos, to.Email)
	}

	return tos
}

func (m *Mail) ResendCc() []string {
	var ccs []string
	for _, cc := range m.Cc {
		ccs = append(ccs, cc.Email)
	}

	return ccs
}

func (m *Mail) ResendBcc() []string {
	var bccs []string
	for _, bcc := range m.Bcc {
		bccs = append(bccs, bcc.Email)
	}

	return bccs
}
