package common

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"server_go/internal/model"
	tostring "server_go/pkg/helper/toString"
	"sync"
	"time"
)

func In(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func VideoExtArr() []string {

	return []string{".avi", ".mp4", ".mov", ".wmv", ".flv", ".mkv", ".mpg", ".rmvb"}

}
func ImgExtArr() []string {

	return []string{".jpeg", ".jpg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".heif"}

}

func PPTExtArr() []string {

	return []string{".ppt", ".pptx"}

}
func PdfExtArr() []string {

	return []string{".pdf"}

}
func OfficeExtArr() []string {

	return []string{".xlsx", ".xls", ".docx", ".doc", ".pptx", ".ppt"}

}

func OfficeNotFoundPPTExtArr() []string {

	return []string{".xlsx", ".xls", ".docx", ".doc"}

}

func AudioExtArr() []string {
	return []string{".mp3", ".wav", ".ogg"}

}

func DealVideoChan(db *sql.DB) {

	fileinfochanmutex := LoadViderChan()

	for {

		select {
		case fileinfo := <-fileinfochanmutex.Fileinfochan:
			fmt.Println(fileinfo)
			filepath, err := VideoToM3u8(fileinfo)
			if err != nil {
				fmt.Println(err)
			}
			ret, err := db.Exec(" update fileinfo set FilePath=? where fileInfoId=?", filepath, fileinfo.FileInfoId)
			if err != nil {
				fmt.Println(err)
			}

			n, err := ret.RowsAffected()
			if err != nil {
				fmt.Println(err)
			}
			if n > 0 {
				fmt.Println("操作成功")
			} else {
				fmt.Println("操作失败")
			}
		case <-time.After(time.Second * 2):
			return
		}
	}
}

func VideoToM3u8(fileinfo *model.FileInfo) (string, error) {

	logicalCPU := int(runtime.NumCPU() / 2)

	filenameall := path.Base(fileinfo.FilePath)
	filesuffix := path.Ext(filenameall)

	filename := filenameall[0 : len(filenameall)-len(filesuffix)]

	filedir := path.Dir(PublicPath + fileinfo.FilePath)
	err := os.MkdirAll(filedir+"/"+filename, 0777)
	if err != nil {
		return "", err

	} else {
		fmt.Println("Successfully created directories")
	}

	cmd := exec.Command("/usr/bin/ffmpeg/ffmpeg",
		"-i", PublicPath+fileinfo.FilePath,
		"-threads", tostring.Strval(logicalCPU),
		"-c:v", "copy",
		"-c:a", "copy",
		"-y",
		filedir+"/"+filename+"/"+filename+".mp4")

	err = cmd.Run()
	if err != nil {
		return "", err
	}
	err = os.Remove(PublicPath + fileinfo.FilePath)
	if err != nil {
		return "", err
	}
	cmd = exec.Command("/usr/bin/ffmpeg/ffmpeg",
		"-i", filedir+"/"+filename+"/"+filename+".mp4",
		//"-profile:v", "baseline",
		"-acodec", "copy",
		"-level", "3.0",
		"-start_number", "0",
		"-hls_time", "20",
		"-hls_list_size", "0",
		"-threads", tostring.Strval(logicalCPU),
		"-f", "hls", filedir+"/"+filename+"/"+filename+".m3u8")
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	err = os.Remove(filedir + "/" + filename + "/" + filename + ".mp4")
	if err != nil {
		return "", err
	}
	return path.Dir(fileinfo.FilePath) + "/" + filename + "/" + filename + ".m3u8", err
}

var onceChan sync.Once
var fileinfochanmutex *MutexFileChan

type MutexFileChan struct {
	Mutex        sync.Mutex
	Fileinfochan chan *model.FileInfo
}

func LoadViderChan() *MutexFileChan {
	onceChan.Do(func() {
		fileinfochanmutex = new(MutexFileChan)
		fileinfochanmutex.Fileinfochan = make(chan *model.FileInfo, 500)

	})
	return fileinfochanmutex
}

type MutexOfficeChan struct {
	Mutex     sync.Mutex
	OfficePdf chan *model.OfficePdf
}

var onceofficeChan sync.Once
var officechanmutex *MutexOfficeChan

func LoadOfficeChan() *MutexOfficeChan {
	onceofficeChan.Do(func() {
		officechanmutex = new(MutexOfficeChan)
		officechanmutex.OfficePdf = make(chan *model.OfficePdf, 500)

	})
	return officechanmutex

}

func OfficeToPdf(db *sql.DB) {

	officechanmutex := LoadOfficeChan()
	for {

		select {
		case officePdf := <-officechanmutex.OfficePdf:
			//dir, _ := os.Getwd()

			filenameall := path.Base(officePdf.InputFile)
			filesuffix := path.Ext(filenameall)

			filename := filenameall[0 : len(filenameall)-len(filesuffix)]

			cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf:writer_pdf_Export", PublicPath+officePdf.InputFile, "--outdir", PublicPath+officePdf.OutputFile)
			err := cmd.Run()
			if err != nil {
				fmt.Println("转换失败:", err)
				return
			}

			fmt.Println("转换成功！")

			ret, err := db.Exec(" update fileinfo set FilePath=? where fileInfoId=?", officePdf.OutputFile+"/"+filename+".pdf", officePdf.FileInfoId)
			if err != nil {
				fmt.Println(err)
			}

			n, err := ret.RowsAffected()
			if err != nil {
				fmt.Println(err)
			}
			if n > 0 {
				fmt.Println("操作成功")
			} else {
				fmt.Println("操作失败")
			}

		case <-time.After(time.Second * 2):
			return
		}
	}
}
