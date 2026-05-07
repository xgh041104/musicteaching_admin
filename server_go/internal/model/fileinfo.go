package model

type FileInfo struct {
	FileInfoId int    `json:"fileInfoId"`
	FileType   string `json:"fileType"`
	FileName   string `json:"fileName"`
	FileUseTo  string `json:"fileUseTo"`
	FilePath   string `json:"filePath"`
}

type OfficePdf struct {
	InputFile  string
	OutputFile string
	FileInfoId int
	FileName   string
}
