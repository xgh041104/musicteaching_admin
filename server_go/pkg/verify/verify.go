package verify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"server_go/pkg/helper/aes"
	tostring "server_go/pkg/helper/toString"
	"server_go/pkg/helper/uuid"
	"time"
)

// 账号验证
const (
	baseURL = "http://127.0.0.1:9878"
	// baseURL                = "http://47.116.207.219/testverify"  //服务器地址
	secretAddURL           = baseURL + "/verify/user/add"        //添加用户
	secretVerifyURL        = baseURL + "/verify/user/verifyuser" //验证用户
	secretVerifyStudentURL = baseURL + "/verify/children/verify" //验证二级用户
	secretDelURL           = baseURL + "/verify/user/del"        //删除用户
	aesKey                 = "1951EC8DA5"                        //密钥字符串
	SystemName             = "subjectcourse:"                    //系统名
)

func generateEncryptedRequest(dataMap map[string]string, urlStr string) ([]byte, string, error) {
	if len(dataMap) == 0 || urlStr == "" {
		return nil, "", errors.New("参数缺失")
	}

	//生成密钥
	key := aes.GenerateKey(aesKey)

	//加密
	userInfo := make(map[string]string)

	for k, v := range dataMap {
		data, err := aes.Encrypt(key, []byte(v))
		if err != nil {
			return nil, "", errors.New("系统错误，请稍等重试")
		}

		userInfo[k] = data
	}

	req, err := json.Marshal(userInfo)
	if err != nil {
		return nil, "", errors.New("系统错误，请稍等重试")
	}

	//生成timestamp、nonce、sign
	timestamp := tostring.Strval(time.Now().Unix())

	nonce := uuid.GenUUID()

	sign, err := aes.Encrypt(key, []byte(timestamp+nonce))
	if err != nil {
		return nil, "", errors.New("系统错误，请稍等重试")
	}

	//生成url，形式：xxx?timestamp=xx
	u, _ := url.ParseRequestURI(urlStr)

	data := url.Values{}
	data.Set("nonce", nonce)
	data.Set("timestamp", timestamp)
	data.Set("sign", sign)

	u.RawQuery = data.Encode()

	return req, u.String(), nil
}

func SecretAddUser(dataMap map[string]string) error {
	req, url, err := generateEncryptedRequest(dataMap, secretAddURL)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	if resp.Code != 0 {
		return errors.New(resp.Message)
	}

	return nil
}

func SecretVerifyUser(dataMap map[string]string) (string, error) {
	req, url, err := generateEncryptedRequest(dataMap, secretVerifyURL)
	if err != nil {
		return "", err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return "", errors.New("系统错误，请稍等重试")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("系统错误，请稍等重试")
	}

	var resp struct {
		Code        int    `json:"code"`
		Message     string `json:"message"`
		AccessToken string `json:"accessToken"`
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return "", errors.New("系统错误，请稍等重试")
	}

	if resp.Code != 0 {
		return "", errors.New(resp.Message)
	}

	return resp.AccessToken, nil
}

func SecretVerifyStudent(parentMark, accessToken string) error {
	key := aes.GenerateKey(aesKey)

	mark, err := aes.Encrypt(key, []byte(parentMark))
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	dataMap := make(map[string]string)
	dataMap["parentMark"] = mark

	reqData, err := json.Marshal(dataMap)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	timestamp := tostring.Strval(time.Now().Unix())

	nonce := uuid.GenUUID()

	sign, err := aes.Encrypt(key, []byte(timestamp+nonce))
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	u, _ := url.ParseRequestURI(secretVerifyStudentURL)

	data := url.Values{}
	data.Set("nonce", nonce)
	data.Set("timestamp", timestamp)
	data.Set("sign", sign)

	u.RawQuery = data.Encode()

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqData))
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	if resp.Code != 0 {
		return errors.New(resp.Message)
	}

	return nil
}

func SecretDelUser(dataMap map[string]string) error {
	req, url, err := generateEncryptedRequest(dataMap, secretDelURL)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return errors.New("系统错误，请稍等重试")
	}

	if resp.Code != 0 {
		return errors.New(resp.Message)
	}

	return nil
}
