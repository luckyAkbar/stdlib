package mail

import "errors"

var (
	// ErrMailgunNotActivated is returned when mailgun is not activated by configuration
	ErrMailgunNotActivated = errors.New("mailgun is not activated by configuration")

	// ErrSendInBlueNotActivated is returned when sendinblue is not activated by configuration
	ErrSendInBlueNotActivated = errors.New("sendinblue is not activated by configuration")
)
