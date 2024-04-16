package model

type User struct {
	Id    []byte
	Email string
	Bio   string
}

type UserCredentials struct {
	Id       []byte
	Email    string
	PassHash []byte
}
