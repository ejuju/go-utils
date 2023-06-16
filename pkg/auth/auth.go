package auth

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ejuju/go-utils/pkg/email"
	"github.com/ejuju/go-utils/pkg/uid"
	"github.com/ejuju/go-utils/pkg/validation"
)

type Sessions interface {
	Create(*Session) error
	Find(string) (*Session, error)
	Delete(string) error
}

type Session struct {
	ID        string
	CreatedAt time.Time
	UserID    string
}

type Users interface {
	FindByEmailAddress(string) (*User, error)
}

type User struct {
	ID           string
	EmailAddress string
}

type OTPs interface {
	Create(*OTP) error
	Find(string) (*OTP, error)
	Delete(string) error
}

type OTP struct {
	userID    string
	code      string
	createdAt time.Time
}

func NewOTP(userID string) *OTP {
	return &OTP{
		userID:    userID,
		code:      uid.MustNewID(12).Hex(),
		createdAt: time.Now(),
	}
}

var ErrNotFound = errors.New("not found")

func NewErrNotFound(id string) error { return fmt.Errorf("%q %w") }

type OTPAuthenticatorConfig struct {
	Host                string
	ConfirmLoginRoute   string
	SuccessfulLoginPath string
	CookieName          string
	Emailer             email.Emailer
	Users               Users
	OTPs                OTPs
	Sessions            Sessions
}

func (conf *OTPAuthenticatorConfig) validate() error {
	return validation.Validate(
		validation.CheckStringNotEmpty(conf.Host),
		validation.CheckStringNotEmpty(conf.ConfirmLoginRoute),
		validation.CheckStringNotEmpty(conf.SuccessfulLoginPath),
		validation.CheckStringNotEmpty(conf.CookieName),
		validation.CheckNotNil(conf.Emailer),
		validation.CheckNotNil(conf.Users),
		validation.CheckNotNil(conf.OTPs),
		validation.CheckNotNil(conf.Sessions),
	)
}

type OTPAuthenticator struct{ conf *OTPAuthenticatorConfig }

func NewOTPAuthenticator(config *OTPAuthenticatorConfig) *OTPAuthenticator {
	if err := config.validate(); err != nil {
		panic(err)
	}
	return &OTPAuthenticator{conf: config}
}

func (authr *OTPAuthenticator) SendLoginLinkByEmail(addr string) error {
	user, err := authr.conf.Users.FindByEmailAddress(addr)
	if err != nil {
		return err
	}
	otp := NewOTP(user.ID)
	err = authr.conf.OTPs.Create(otp)
	if err != nil {
		return err
	}
	link := fmt.Sprintf("%s%s?email-address=%s&code=%s",
		authr.conf.Host,
		authr.conf.ConfirmLoginRoute,
		url.QueryEscape(addr),
		url.QueryEscape(otp.code),
	)
	return authr.conf.Emailer(&email.Email{
		To:            []string{addr},
		Subject:       "New one-time password",
		PlainTextBody: link,
	})
}

func (authr *OTPAuthenticator) LoginWithLink(w http.ResponseWriter, r *http.Request, onInternalErr, onErr func(error)) {
	emailAddr := r.URL.Query().Get("email-address")
	otpCode := r.URL.Query().Get("code")

	// Check if email address exists
	user, err := authr.conf.Users.FindByEmailAddress(emailAddr)
	if err != nil {
		onInternalErr(err)
		return
	}

	// Get OTP for provided user
	otp, err := authr.conf.OTPs.Find(otpCode)
	if errors.Is(err, ErrNotFound) {
		onErr(fmt.Errorf("OTP not found for email address: %q", emailAddr))
		return
	} else if err != nil {
		onInternalErr(err)
		return
	}

	// Check OTP
	if subtle.ConstantTimeCompare([]byte(otp.code), []byte(otpCode)) == 0 {
		onErr(errors.New("invalid OTP"))
		return
	}
	// OK, valid credentials

	// Delete OTP
	err = authr.conf.OTPs.Delete(otpCode)
	if err != nil {
		onInternalErr(err)
		return
	}

	// Create session
	s := &Session{ID: uid.MustNewID(12).Hex(), UserID: user.ID, CreatedAt: time.Now()}
	err = authr.conf.Sessions.Create(s)
	if err != nil {
		onInternalErr(err)
		return
	}

	// Set auth cookie and redirect to page
	http.SetCookie(w, newAuthCookie(authr.conf.CookieName, s.ID, time.Hour))
	http.Redirect(w, r, authr.conf.SuccessfulLoginPath, http.StatusSeeOther)
}

func (authr *OTPAuthenticator) Authenticate(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// Get auth cookie
	cookie, err := r.Cookie(authr.conf.CookieName)
	if err != nil {
		return nil, nil
	}

	// Get session by ID
	session, err := authr.conf.Sessions.Find(cookie.Value)
	if err != nil {
		http.SetCookie(w, newAuthCookie(authr.conf.CookieName, "", 0))
		return nil, err
	}

	// Check if session is expired
	if time.Since(session.CreatedAt) > time.Hour {
		http.SetCookie(w, newAuthCookie(authr.conf.CookieName, "", 0))
		return nil, nil
	}

	return session, nil
}

func newAuthCookie(name, value string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		HttpOnly: true,
		Value:    value,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
		Secure:   true,
	}
}

type MockUsers map[string]*User

func NewMockUsers(emailAddrs ...string) MockUsers {
	out := MockUsers{}
	for _, addr := range emailAddrs {
		id := uid.MustNewID(12).Hex()
		out[id] = &User{ID: id, EmailAddress: addr}
	}
	return out
}

func (users MockUsers) FindByEmailAddress(addr string) (*User, error) {
	for _, u := range users {
		if u.EmailAddress == addr {
			return u, nil
		}
	}
	return nil, NewErrNotFound(addr)
}

type MockSessions map[string]*Session

func (sessions MockSessions) Create(s *Session) error { sessions[s.ID] = s; return nil }
func (sessions MockSessions) Find(id string) (*Session, error) {
	out, ok := sessions[id]
	if !ok {
		return nil, NewErrNotFound(id)
	}
	return out, nil
}
func (sessions MockSessions) Delete(id string) error { delete(sessions, id); return nil }

type MockOTPs map[string]*OTP

func (otps MockOTPs) Create(otp *OTP) error { otps[otp.code] = otp; return nil }
func (otps MockOTPs) Find(id string) (*OTP, error) {
	out, ok := otps[id]
	if !ok {
		return nil, NewErrNotFound(id)
	}
	return out, nil
}
func (otps MockOTPs) Delete(id string) error { delete(otps, id); return nil }
