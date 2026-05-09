package service

import (
	"ai_summary_project/internal/model"
	"ai_summary_project/internal/repository"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ----------------- 旧AI总结接口鉴权参数 -------------------
const (
	oldAppID     = "d0fdf554"
	oldSecretKey = "ca448a46b7a702184318066cd101e0b6"
	xModelURL    = "https://spark-api-open.xf-yun.com/v2/chat/completions"
	apiPassword  = "CRVqMJlUrJlfZNCoLsOd:XbJoLmPesQfxtnlDUGoe"
)

// ----------------- 新语音转写接口鉴权参数 -------------------
const (
	newAppID      = "8c782774"
	newApiKey     = "ed6f9d8c44dc6331f7f8e871906227c7"
	newApiSecret  = "MzBlMWEyYjY5YTc1ZjlmODQ2N2JlOTA4"
	uploadBaseURL = "https://upload-ost-api.xfyun.cn/file"
	openBaseURL   = "https://ost-api.xfyun.cn/v2"
	sliceSize     = 15 * 1024 * 1024 // 15MB
)

// 业务层接口定义
type CourseService interface {
	CreateCourseWithFiles(ctx *gin.Context, bookID uint, title, summary string, mp4Header, mp3Header *multipart.FileHeader) error
	GenerateSummaryFromAudio(ctx *gin.Context, audioFile *multipart.FileHeader, prompt string) (string, error)
	GetCoursesByBookID(ctx *gin.Context, bookID uint) ([]*model.Course, error)
	DeleteCourse(ctx *gin.Context, courseID uint) error
	GetAllBooks(ctx *gin.Context) ([]*model.Book, error)
}

// 业务层实现结构体（假设model和repository包存在）
type courseService struct {
	*Service
	courseRepository repository.CourseRepository
	bookRepository   repository.BookRepository
}

func NewCourseService(svc *Service, courseRepo repository.CourseRepository, bookRepo repository.BookRepository) CourseService {
	return &courseService{
		Service:          svc,
		courseRepository: courseRepo,
		bookRepository:   bookRepo,
	}
}

// ------------------- 旧课程相关接口代码 -------------------

func (s *courseService) CreateCourseWithFiles(
	ctx *gin.Context,
	bookID uint,
	title, summary string,
	mp4Header, mp3Header *multipart.FileHeader,
) error {
	now := time.Now()
	dir := filepath.Join("static", "media", now.Format("2006"), now.Format("01"), now.Format("02"))

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	uuid := generateUUID()
	videoPath := filepath.Join(dir, uuid+".mp4")
	audioPath := filepath.Join(dir, uuid+".mp3")

	if err := ctx.SaveUploadedFile(mp4Header, videoPath); err != nil {
		return fmt.Errorf("保存视频失败: %w", err)
	}

	if err := ctx.SaveUploadedFile(mp3Header, audioPath); err != nil {
		_ = os.Remove(videoPath)
		return fmt.Errorf("保存音频失败: %w", err)
	}

	course := &model.Course{
		Title:      title,
		Summary:    summary,
		BookID:     bookID,
		VideoPath:  videoPath,
		RecordPath: audioPath,
	}

	if err := s.courseRepository.Create(ctx.Request.Context(), course); err != nil {
		_ = os.Remove(videoPath)
		_ = os.Remove(audioPath)
		return fmt.Errorf("课程写入数据库失败: %w", err)
	}

	return s.bookRepository.AddCourseCount(ctx.Request.Context(), bookID, +1)
}

func (s *courseService) GetCoursesByBookID(ctx *gin.Context, bookID uint) ([]*model.Course, error) {
	return s.courseRepository.ListByBookID(ctx.Request.Context(), bookID)
}

func (s *courseService) DeleteCourse(ctx *gin.Context, courseID uint) error {
	course, err := s.courseRepository.GetByID(ctx.Request.Context(), courseID)
	if err != nil {
		return err
	}

	if course.VideoPath != "" {
		_ = os.Remove(course.VideoPath)
	}
	if course.RecordPath != "" {
		_ = os.Remove(course.RecordPath)
	}

	if err := s.courseRepository.Delete(ctx.Request.Context(), courseID); err != nil {
		return err
	}

	return s.bookRepository.AddCourseCount(ctx.Request.Context(), course.BookID, -1)
}

func (s *courseService) GetAllBooks(ctx *gin.Context) ([]*model.Book, error) {
	return s.bookRepository.List(ctx.Request.Context())
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// 旧签名生成（保留，未变）
func generateSigna(ts string) string {
	base := oldAppID + ts
	md5Sum := md5.Sum([]byte(base))
	h := hmac.New(sha1.New, []byte(oldSecretKey))
	h.Write([]byte(hex.EncodeToString(md5Sum[:])))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ------------------ 新的语音转写 + AI总结实现 -----------------

// 生成临时文件名
func tmpFileName(prefix, original string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s", prefix, original))
}

// 入口函数，替代旧的GenerateSummaryFromAudio
func (s *courseService) GenerateSummaryFromAudio(ctx *gin.Context, audioFile *multipart.FileHeader, prompt string) (string, error) {
	// 保存上传文件
	tmpPath := tmpFileName(generateUUID(), audioFile.Filename)
	if err := ctx.SaveUploadedFile(audioFile, tmpPath); err != nil {
		return "", fmt.Errorf("保存上传文件失败: %w", err)
	}
	defer os.Remove(tmpPath)

	// mp3转pcm16k
	pcmFile, err := convertMp3ToPcm16k(tmpPath)
	if err != nil {
		return "", err
	}
	if pcmFile != tmpPath {
		defer os.Remove(pcmFile)
	}

	// 初始化上传
	fi, err := os.Stat(pcmFile)
	if err != nil {
		return "", err
	}
	uploadID, err := initMultipartUpload(filepath.Base(pcmFile), fi.Size())
	if err != nil {
		return "", err
	}

	// 分片上传
	if err = uploadFileInSlices(pcmFile, uploadID); err != nil {
		return "", err
	}

	// 完成上传
	fileURL, err := completeUpload(uploadID)
	if err != nil {
		return "", err
	}

	// 创建转写任务
	taskID, err := createTranscribeTask(fileURL)
	if err != nil {
		return "", err
	}

	// 查询转写结果
	text, err := queryTranscribeResult(taskID)
	if err != nil {
		return "", err
	}

	// 调用旧AI总结接口
	return callXModel(prompt, text)
}

// --------- mp3转pcm，分片上传相关函数 -----------

func convertMp3ToPcm16k(inputFile string) (string, error) {
	ext := strings.ToLower(filepath.Ext(inputFile))
	if ext != ".mp3" {
		return inputFile, nil
	}

	outputFile := filepath.Join(os.TempDir(), generateUUID()+".pcm")
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, "-f", "s16le", "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", outputFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg 转换失败: %s 详细: %s", err, string(out))
	}
	return outputFile, nil
}

func initMultipartUpload(fileName string, fileSize int64) (string, error) {
	url := uploadBaseURL + "/mpupload/init"
	reqBody := map[string]interface{}{
		"app_id":     newAppID,
		"file_name":  fileName,
		"file_size":  fileSize,
		"request_id": strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	respBytes, err := postJSONWithAuth(url, reqBody)
	if err != nil {
		return "", err
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UploadID string `json:"upload_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("初始化上传失败 code=%d msg=%s", resp.Code, resp.Message)
	}
	return resp.Data.UploadID, nil
}

func uploadSlice(uploadID string, sliceID int, data []byte) error {
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	_ = writer.WriteField("app_id", newAppID)
	_ = writer.WriteField("upload_id", uploadID)
	_ = writer.WriteField("slice_id", strconv.Itoa(sliceID))
	_ = writer.WriteField("request_id", strconv.FormatInt(time.Now().UnixNano(), 10))

	part, err := writer.CreateFormFile("data", fmt.Sprintf("slice_%d", sliceID))
	if err != nil {
		return err
	}
	if _, err = part.Write(data); err != nil {
		return err
	}

	_ = writer.Close()

	req, err := http.NewRequest("POST", uploadBaseURL+"/mpupload/upload", bodyBuf)
	if err != nil {
		return err
	}

	addAuthHeaders(req, bodyBuf.Bytes(), writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, _ := io.ReadAll(resp.Body)
	fmt.Printf("分片%d上传返回：%s\n", sliceID, string(bs))
	return nil
}

func uploadFileInSlices(filePath, uploadID string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, sliceSize)
	sliceID := 1

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if err := uploadSlice(uploadID, sliceID, buf[:n]); err != nil {
			return err
		}
		sliceID++
		if err == io.EOF {
			break
		}
	}
	return nil
}

func completeUpload(uploadID string) (string, error) {
	url := uploadBaseURL + "/mpupload/complete"
	reqBody := map[string]interface{}{
		"app_id":     newAppID,
		"upload_id":  uploadID,
		"request_id": strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	respBytes, err := postJSONWithAuth(url, reqBody)
	if err != nil {
		return "", err
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Url string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("完成上传失败 code=%d msg=%s", resp.Code, resp.Message)
	}
	return resp.Data.Url, nil
}

func createTranscribeTask(audioURL string) (string, error) {
	url := openBaseURL + "/ost/pro_create"
	reqBody := map[string]interface{}{
		"common": map[string]string{"app_id": newAppID},
		"business": map[string]string{
			"request_id": strconv.FormatInt(time.Now().UnixNano(), 10),
			"language":   "zh_cn",
			"accent":     "mandarin",
			"domain":     "pro_ost_ed",
		},
		"data": map[string]string{
			"audio_url": audioURL,
			"encoding":  "raw",
			"format":    "audio/L16;rate=16000",
			"audio_src": "http",
		},
	}
	respBytes, err := postJSONWithAuth(url, reqBody)
	if err != nil {
		return "", err
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TaskID string `json:"task_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("创建转写任务失败 code=%d msg=%s", resp.Code, resp.Message)
	}
	return resp.Data.TaskID, nil
}

func queryTranscribeResult(taskID string) (string, error) {
	for i := 0; i < 20; i++ {
		time.Sleep(2 * time.Second)
		url := openBaseURL + "/ost/query"
		reqBody := map[string]interface{}{
			"common": map[string]string{"app_id": newAppID},
			"business": map[string]string{
				"task_id": taskID,
			},
		}
		respBytes, err := postJSONWithAuth(url, reqBody)
		if err != nil {
			fmt.Printf("查询失败，第%d次，错误：%v\n", i+1, err)
			continue
		}
		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				TaskStatus string `json:"task_status"`
				Result     struct {
					Lattice []struct {
						JSON1Best struct {
							St struct {
								Rt []struct {
									Ws []struct {
										Cw []struct {
											W string `json:"w"`
										} `json:"cw"`
									} `json:"ws"`
								} `json:"rt"`
							} `json:"st"`
						} `json:"json_1best"`
					} `json:"lattice"`
				} `json:"result"`
			} `json:"data"`
		}
		if err := json.Unmarshal(respBytes, &resp); err != nil {
			fmt.Printf("解析失败，第%d次，错误：%v\n", i+1, err)
			continue
		}
		fmt.Printf("查询第%d次，状态: %s\n", i+1, resp.Data.TaskStatus)
		if resp.Data.TaskStatus == "4" || resp.Data.TaskStatus == "3" {
			var sb strings.Builder
			for _, lat := range resp.Data.Result.Lattice {
				for _, rt := range lat.JSON1Best.St.Rt {
					for _, ws := range rt.Ws {
						for _, cw := range ws.Cw {
							sb.WriteString(cw.W)
						}
					}
				}
			}
			return sb.String(), nil
		}
	}
	return "", fmt.Errorf("转写任务超时未完成")
}

// 公共请求，带新转写鉴权头
func postJSONWithAuth(url string, body interface{}) ([]byte, error) {
	jsonBytes, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	addAuthHeaders(req, jsonBytes, "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// 新转写鉴权头生成
func addAuthHeaders(req *http.Request, body []byte, contentType string) {
	u := req.URL
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	hash := sha256.Sum256(body)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])

	signStr := fmt.Sprintf("host: %s\ndate: %s\nPOST %s HTTP/1.1\ndigest: %s", u.Host, date, u.Path, digest)

	mac := hmac.New(sha256.New, []byte(newApiSecret))
	mac.Write([]byte(signStr))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authHeader := fmt.Sprintf(`hmac username="%s", algorithm="hmac-sha256", headers="host date request-line digest", signature="%s"`, newApiKey, signature)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Date", date)
	req.Header.Set("Digest", digest)
	req.Header.Set("Authorization", authHeader)
}

// 调用旧AI总结接口
func callXModel(theme string, text string) (string, error) {
	body := map[string]interface{}{
		"user":  "课程助手",
		"model": "x1",
		"messages": []map[string]string{{
			"role": "system",
			"content": `你是一名教学总结助手，请对以下上课文本内容进行总结。总结时必须遵守以下规则：
严格根据原始文本内容进行总结，不进行扩展、联想或补充；
省略思考过程，直接返回结果
若文本内容非常少或无实质信息，请只返回：“上课内容” 或 “无”；
若文本内容充实，请以书面总结格式输出，内容要求简洁，结构如下：
课程要点：
用简洁语句概括课堂讲授的核心内容。
不得对原文进行扩展以及无效添加
主要知识点：
分条列出本节课所涵盖的重要知识点；
每条仅陈述知识本身，不加解释、比喻或建议。
难点提醒：
不加拓展，仅列出原文中实际提及或强调的难点。
分条列出学生在本节课中容易混淆或出错的地方，需针对课堂实际情况；
仅列出原文中明确提及或重点强调的难点，不得进行拓展。
列出课堂上老师强调的需要学生注意的事项
注意事项：（若无内容请整段省略）

师生互动：（若无内容请整段省略）

作业布置：（若无内容请整段省略）

其他备注：（若无内容请整段省略）

请总结以下文本内容：主题：` + theme + "内容为：" + text,
		}},
		"stream": false,
	}

	jsonBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", xModelURL, bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respParsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	raw, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(raw, &respParsed); err != nil {
		return "", fmt.Errorf("解析大模型响应失败: %v", err)
	}
	if len(respParsed.Choices) == 0 {
		return "", fmt.Errorf("大模型未返回内容")
	}
	return respParsed.Choices[0].Message.Content, nil
}
