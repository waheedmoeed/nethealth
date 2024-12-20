package model

type Config struct {
	Accounts            []Cred `json:"accounts"`
	NeedToPouplateUsers bool   `json:"needToPouplateUsers"`
	DownloaderUser      string `json:"downloaderUser"`
	DownloaderPassword  string `json:"downloaderPassword"`
	Entity              string `json:"entity"`
	Debug               bool   `json:"debug"`
	Headless            bool   `json:"headless"`
}

type Cred struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//37 "Ageility at 75 State Street"
