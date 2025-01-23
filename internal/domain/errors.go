package domain

// Define errors here for using in the project

type ErrInvalidCredentials struct{}

func (e *ErrInvalidCredentials) Error() string {
	return "invalid credentials"
}
