package contact

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"
	"unicode/utf8"

	"github.com/ejuju/go-utils/pkg/uid"
)

var MaxMessageLength = 800

type Forms interface {
	SaveNew(*Form) error
}

type Form struct {
	ID           string
	At           time.Time
	VisitorHash  string
	RequestID    string
	EmailAddress string
	Message      string
}

func ParseAndValidateForm(r *http.Request, emailFieldName, messageFieldName string) (*Form, error) {
	// Validate email address and message length
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	emailAddr, err := mail.ParseAddress(r.FormValue(emailFieldName))
	if err != nil {
		return nil, fmt.Errorf("%s (%w)", "Invalid email address", err)
	}
	message := r.FormValue(messageFieldName)
	messageLength := utf8.RuneCountInString(message)
	if messageLength > MaxMessageLength {
		return nil, fmt.Errorf("message (%d chars) exceeds max length %d", len(message), MaxMessageLength)
	}

	// Generate unique ID
	return &Form{
		ID:           uid.MustNewID(12).Hex(),
		At:           time.Now(),
		VisitorHash:  r.Header.Get("Visitor-Hash"),
		RequestID:    r.Header.Get("Request-ID"),
		EmailAddress: emailAddr.Address,
		Message:      message,
	}, nil
}

type MockDB map[string]*Form

func (db MockDB) SaveNew(f *Form) error { db[f.ID] = f; return nil }
