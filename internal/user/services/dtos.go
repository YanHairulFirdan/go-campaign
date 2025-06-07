package services

type CreateUserDTO struct {
	Name     string
	Email    string
	Password string
}

type UserDTO struct {
	ID       int64
	Name     string
	Email    string
	Password string
}
