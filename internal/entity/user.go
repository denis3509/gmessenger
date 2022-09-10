package entity

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	HashedPassword string `json:"hashedPassword"`
	Email string `json:"email"`
}
