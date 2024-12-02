package model

type Job struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	Download bool   `json:"download"`
	PDFLink  string `json:"pdfLink"`
}
