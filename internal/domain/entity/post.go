package entity

type Post struct {
	Id       int    `json:"id"`
	Parent   int    `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}

type CreatePost struct {
	Parent  int    `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type PostDetails struct {
	DPost   *Post   `json:"post"`
	DAuthor *User   `json:"author"`
	DThread *Thread `json:"thread"`
	DForum  *Forum  `json:"forum"`
}

type UpdatePost struct {
	Message string `json:"message"`
}

type PostWithoutEdited struct {
	Id      int    `json:"id"`
	Parent  int    `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Forum   string `json:"forum"`
	Thread  int    `json:"thread"`
	Created string `json:"created"`
}
