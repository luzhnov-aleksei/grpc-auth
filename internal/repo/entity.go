package repo

type UserRegister struct {
	Email        string
	Username     string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	Role         string
}
