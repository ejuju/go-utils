package auth

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ejuju/go-utils/pkg/email"
	"github.com/ejuju/go-utils/pkg/uid"
)

type Session struct {
	ID        string
	CreatedAt time.Time
	UserID    string
}

type User struct {
	ID           string
	EmailAddress string
}

type OTP struct {
	userID    string
	code      string
	createdAt time.Time
}

func NewAuthCookie(name, value string, ttl time.Duration) *http.Cookie {
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

type StaticUsers map[string]*User

func NewStaticUsers(emailAddrs ...string) StaticUsers {
	out := StaticUsers{}
	for _, addr := range emailAddrs {
		id := uid.MustNewID(12).Hex()
		out[id] = &User{ID: id, EmailAddress: addr}
	}
	return out
}

func (users StaticUsers) FindByEmailAddress(addr string) (*User, error) {
	for _, u := range users {
		if u.EmailAddress == addr {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with email address %q not found", addr)
}

func (users StaticUsers) FindByID(id string) *User { return users[id] }

func NewOTP(userID string) *OTP {
	return &OTP{
		userID:    userID,
		code:      uid.MustNewID(12).Hex(),
		createdAt: time.Now(),
	}
}

type OTPAuthenticatorConfig struct {
	Host                string
	ConfirmLoginRoute   string
	SuccessfulLoginPath string
	Users               StaticUsers
	CookieName          string
	Emailer             email.Emailer
}

type OTPAuthenticator struct {
	conf     *OTPAuthenticatorConfig
	otps     map[string]*OTP
	sessions map[string]*Session
}

func NewOTPAuthenticator(config *OTPAuthenticatorConfig) *OTPAuthenticator {
	return &OTPAuthenticator{
		conf:     config,
		sessions: map[string]*Session{},
		otps:     map[string]*OTP{},
	}
}

func (authr *OTPAuthenticator) SendLoginLinkByEmail(addr string) error {
	user, err := authr.conf.Users.FindByEmailAddress(addr)
	if err != nil {
		return err
	}
	otp := NewOTP(user.ID)
	authr.otps[user.ID] = otp
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
	otp, ok := authr.otps[user.ID]
	if !ok {
		onErr(fmt.Errorf("OTP not found for email address: %q", emailAddr))
		return
	}

	// Check OTP
	if subtle.ConstantTimeCompare([]byte(otp.code), []byte(otpCode)) == 0 {
		onErr(errors.New("invalid OTP"))
		return
	}

	// OK
	authr.sessions[user.ID] = &Session{ID: uid.MustNewID(12).Hex(), UserID: user.ID, CreatedAt: time.Now()}
	http.SetCookie(w, NewAuthCookie(authr.conf.CookieName, authr.sessions[user.ID].ID, time.Hour))
	http.Redirect(w, r, authr.conf.SuccessfulLoginPath, http.StatusSeeOther)
}

func (authr *OTPAuthenticator) Authenticate(w http.ResponseWriter, r *http.Request) *Session {
	// Get auth cookie
	cookie, err := r.Cookie(authr.conf.CookieName)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(cookie.Value)
	// Get session by ID
	session, ok := authr.sessions[cookie.Value]
	if !ok {
		http.SetCookie(w, NewAuthCookie(authr.conf.CookieName, "", 0))
		return nil
	}
	log.Println(session)
	// Check if session is expired
	if time.Since(session.CreatedAt) > time.Hour {
		http.SetCookie(w, NewAuthCookie(authr.conf.CookieName, "", 0))
		return nil
	}
	log.Println("non expired")
	return session
}
