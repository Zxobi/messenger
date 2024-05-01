package model

type User struct {
	Id    []byte `bson:"_id"`
	Email string `bson:"email"`
	Bio   string `bson:"bio"`
}

type UserCredentials struct {
	Id       []byte `bson:"_id"`
	Email    string `bson:"email"`
	PassHash []byte `bson:"pass_hash"`
}
