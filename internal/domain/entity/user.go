package entity

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

//easyjson:json
type Users []User

type CreateUser struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type UpdateUser struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Error struct {
	Message string `json:"message"`
}
