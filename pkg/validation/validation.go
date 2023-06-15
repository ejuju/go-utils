package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"unicode/utf8"
)

type Check func() error

func Validate(checks ...Check) error { return CheckMultiple(checks...)() }

// Wrapper to check errors.
func CheckErr(err error) Check { return func() error { return err } }

func CheckMultiple(checks ...Check) Check {
	return func() error {
		for i, v := range checks {
			err := v()
			if err != nil {
				return fmt.Errorf("check failed [%d/%d]: %w", i, len(checks), err)
			}
		}
		return nil
	}
}

func CheckWhen(ok bool, v Check) Check {
	return func() error {
		if !ok {
			return nil
		}
		return v()
	}
}

func CheckEmailAddress(addr string) Check {
	return func() error {
		_, err := mail.ParseAddress(addr)
		return err
	}
}

// Inclusive
func CheckUTF8StringMinLength(in string, min int) Check {
	return func() error {
		length := utf8.RuneCountInString(in)
		if length < min {
			return fmt.Errorf("got string length %d but want minimum %d", length, min)
		}
		return nil
	}
}

// Exclusive
func CheckUTF8StringMaxLength(in string, max int) Check {
	return func() error {
		length := utf8.RuneCountInString(in)
		if length >= max {
			return fmt.Errorf("got string length %d but want maximum %d", length, max)
		}
		return nil
	}
}

func CheckStringIs(in string, match string) Check {
	return func() error {
		if in != match {
			return fmt.Errorf("string %q doesn't match %q", in, match)
		}
		return nil
	}
}

func CheckStringIsEither(in string, options ...string) Check {
	return func() error {
		for _, opt := range options {
			if in == opt {
				return nil
			}
		}
		return fmt.Errorf("string %q doesn't match any of options: %q", in, options)
	}
}

func CheckNotNil(in any) Check {
	return func() error {
		if in == nil {
			return errors.New("should not be nil")
		}
		return nil
	}
}

func CheckNetworkPort(in int) Check {
	return func() error {
		if in < 0 || in > 65535 {
			return fmt.Errorf("invalid port number: %d", in)
		}
		return nil
	}
}
