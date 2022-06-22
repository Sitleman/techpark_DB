package entity

type Vote struct {
	IdThread   int    `json:"idThread"`
	SlugThread string `json:"slugThread"`
	Nickname   string `json:"nickname"`
	Voice      int    `json:"voice"`
}
