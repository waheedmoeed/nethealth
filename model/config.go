package model

type Config struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	NeedToPouplateUsers bool   `json:"needToPouplateUsers"`
}
