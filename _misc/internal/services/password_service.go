package services

type PasswordService interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
	Validate(password string) bool
}
