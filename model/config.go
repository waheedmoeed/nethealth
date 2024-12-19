package model

type Config struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	ManualEmail         string `json:"manualEmail"`
	ManualPassword      string `json:"manualPassword"`
	NeedToPouplateUsers bool   `json:"needToPouplateUsers"`
	DownloaderUser      string `json:"downloaderUser"`
	DownloaderPassword  string `json:"downloaderPassword"`
	Entity              string `json:"entity"`
	Debug               bool   `json:"debug"`
	Headless            bool   `json:"headless"`
}

//37 "Ageility at 75 State Street"
