package entity

type CreateThread struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

type Thread struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

type ThreadResponse struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes"`
	Created string `json:"created"`
}
