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
	SaveNew(*FormData) error
}

type FormData struct {
	ID           string
	At           time.Time
	EmailAddress string
	Message      string
}

func ParseAndValidateForm(r *http.Request, emailFieldName, messageFieldName string) (*FormData, error) {
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
	return &FormData{
		ID:           uid.MustNewID(12).Hex(),
		At:           time.Now(),
		EmailAddress: emailAddr.Address,
		Message:      message,
	}, nil
}

type MockDB map[string]*FormData

func (db MockDB) SaveNew(f *FormData) error { db[f.ID] = f; return nil }
