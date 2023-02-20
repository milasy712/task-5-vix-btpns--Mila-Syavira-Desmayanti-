package formaterror

import (
	"errors"
	"strings"
)

func ErrorMessage (err string) error {
	
	if strings.Contains(err, "pkey") {
		return errors.New("user id already registered")
	}else if strings.Contains(err, "email_key") {
		return errors.New("email has been registered")
	}else if strings.Contains(err, "user not found") {
		return errors.New("email has not been registered")
	}else if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect password")
	}
	return errors.New(err)
}