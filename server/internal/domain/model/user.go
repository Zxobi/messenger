package model

type User struct {
	Id    []byte `bson:"_id"`
	Email string `bson:"email"`
	Bio   string `bson:"bio"`
}

type UserCredentials struct {
	Id       []byte
	Email    string
	PassHash []byte
}
