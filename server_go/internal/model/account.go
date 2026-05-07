package model

type Account struct {
	AccountId        int    `json:"accountId"`
	AccountTitle     string `json:"accountTitle"`
	AccountDesc      string `json:"accountDesc"`
	AccountUrl       string `json:"accountUrl"`
	AccountImgFileId int    `json:"accountImgFileId"`
	IsDelImg         int    `json:"isDelImg" `
}

type AccountView struct {
	AccountId          int    `json:"accountId"`
	AccountTitle       string `json:"accountTitle"`
	AccountDesc        string `json:"accountDesc"`
	AccountUrl         string `json:"accountUrl"`
	AccountImgFileId   int    `json:"accountImgFileId"`
	AccountImgFileName string `json:"accountImgFileName"`
	AccountImgPath     string `json:"accountImgPath"`
}
