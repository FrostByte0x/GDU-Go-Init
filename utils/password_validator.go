package utils

import (
	"errors"
	"regexp"
)

func Validatepassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be 8 caracters minimum length")
	}
	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	num := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#%\$\^\&\*\.]`)
	if !upper.MatchString(password) {
		return errors.New("Password must contain at least one upper case caracter")
	}
	if !lower.MatchString(password) {
		return errors.New("Password must contain at least one lower case caracter")
	}
	if !num.MatchString(password) {
		return errors.New("Password must contain at least one numerical caracter")
	}
	if !special.MatchString(password) {
		return errors.New("Password must contain at least one special caracter: ! @ # %")
	}
	return nil
}
