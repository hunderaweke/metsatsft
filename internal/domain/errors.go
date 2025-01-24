package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type ErrInvalidCredentials struct{}

func (e *ErrInvalidCredentials) Error() string {
	return "invalid credentials"
}

type ErrDuplicateUser struct {
	Message string
}

func (e *ErrDuplicateUser) Error() string {
	fields := parseDuplicateKeyError(e.Message)
	if fields != nil {
		return fmt.Sprintf("duplicate user with %v", fields)
	}
	return "duplicate user"
}

func parseDuplicateKeyError(message string) map[string]string {
	re := regexp.MustCompile(`dup key: {([^}]+)}`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		dupFields := strings.Split(matches[1], ", ")
		fields := make(map[string]string)
		for _, field := range dupFields {
			parts := strings.Split(field, ": ")
			if len(parts) == 2 {
				value := strings.Trim(parts[1], "\"")
				value = strings.TrimSpace(value)
				if len(value) > 0 {
					fields[parts[0]] = value
				}
			}
		}
		return fields
	}
	return nil
}

type ErrUserNotFound struct{}

func (e *ErrUserNotFound) Error() string {
	return "user not found"
}
