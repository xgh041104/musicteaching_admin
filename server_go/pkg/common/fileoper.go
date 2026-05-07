package common

import (
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"server_go/internal/model"
	"server_go/pkg/helper/uuid"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//const PublicPath = "../../"

// TODO 发布需修改
const PublicPath = ""

func SaveFile(fileHeader *multipart.FileHeader, destination string) error {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建目标文件
	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 将文件内容拷贝到目标文件
	_, err = io.Copy(destinationFile, file)
	if err != nil {
		return err
	}

	return nil
}

func AddFileInfo(tx *sql.Tx, fileinfo *model.FileInfo, tempfile *multipart.FileHeader) int64 {

	ret, err := tx.Exec("insert into fileinfo(fileType,fileName,fileUseTo,filePath) values(?,?,?,?) ",
		fileinfo.FileType, fileinfo.FileName, fileinfo.FileUseTo, fileinfo.FilePath)

	if err != nil {
		return 0
	}
	fileId, err := ret.LastInsertId()

	if err != nil {
		return 0
	}

	if fileId < 0 {
		return 0
	}

	// common.PublicPath  解决 执行的文件不在根目录问题
	err = SaveFile(tempfile, PublicPath+fileinfo.FilePath)
	if err != nil {
		return 0
	}
	return fileId
}

func UpdateFileInfo(tx *sql.Tx, fileinfo *model.FileInfo, tempfile *multipart.FileHeader) error {

	oldfilepath := ""

	tx.QueryRow("select filePath from fileinfo where fileInfoId=? ", fileinfo.FileInfoId).Scan(&oldfilepath)

	if oldfilepath == "" {
		return errors.New("删除文件时文件路径错误")
	}
	ret, err := tx.Exec("update fileinfo set  filePath=?,fileType=?,fileName=? where fileInfoId=? ", fileinfo.FilePath, path.Ext(tempfile.Filename), fileinfo.FileName, fileinfo.FileInfoId)

	if err != nil {
		return err
	}
	fileId, err := ret.LastInsertId()

	if err != nil {
		return err
	}

	if fileId < 0 {
		return err
	}

	err = os.Remove(PublicPath + oldfilepath) //删除原图片
	if err != nil {
		return err
	}
	err = SaveFile(tempfile, PublicPath+fileinfo.FilePath) //保存新图片
	if err != nil {
		return err
	}

	return nil
}

func DelFileInfo(tx *sql.Tx, fileinfo *model.FileInfo) error {
	oldfilepath := ""
	tx.QueryRow("select filePath from fileinfo where fileInfoId=? ", fileinfo.FileInfoId).Scan(&oldfilepath)

	if oldfilepath == "" {
		return errors.New("删除文件时文件路径错误")
	}

	ret, err := tx.Exec("delete from   fileinfo  where  fileInfoId=?", fileinfo.FileInfoId)

	if err != nil {
		return err
	}
	fileId, err := ret.RowsAffected()

	if err != nil {
		return err
	}

	if fileId < 0 {
		return err
	}

	_ = os.Remove(PublicPath + oldfilepath) //删除原图片

	return nil
}

func EditCover(tx *sql.Tx, files *multipart.Form, fileUseTo string, infileid int, filedist string, IsDelImg int) (int64, error) {

	var outfileid int64
	fileinfo := new(model.FileInfo)
	if len(files.File["files"]) == 1 && infileid == 0 { //上传了图片 且原先没有图片
		tempfile := files.File["files"][0]
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = fileUseTo
		fileinfo.FilePath = filedist + fileinfo.FileName + path.Ext(tempfile.Filename)
		outfileid = AddFileInfo(tx, fileinfo, tempfile)

		if outfileid == 0 {
			return outfileid, errors.New("文件上传错误")
		}

	} else if len(files.File["files"]) == 1 && infileid > 0 { //上传了图片 且原先有图片 进行更改fileinfo

		outfileid = int64(infileid)

		tempfile := files.File["files"][0]
		fileinfo.FileInfoId = infileid
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FilePath = filedist + fileinfo.FileName + path.Ext(tempfile.Filename)
		err := UpdateFileInfo(tx, fileinfo, tempfile)

		if err != nil {
			return outfileid, errors.New("文件上传失败")
		}

	} else if len(files.File["files"]) == 0 && IsDelImg == 1 && infileid > 0 { //不上传了图片 且原先有图片  删除原先的图片]

		outfileid = 0
		fileinfo.FileInfoId = infileid
		err := DelFileInfo(tx, fileinfo)
		if err != nil {
			return outfileid, errors.New("文件删除失败")
		}

	} else if len(files.File["files"]) == 0 && IsDelImg == 0 && infileid > 0 { //不上传了图片 且原先有图片   不删图片

		return int64(infileid), nil

	} else if len(files.File["files"]) == 0 && infileid == 0 { // 原先没有封面  新传的也没有封面
		outfileid = 0
	}

	return outfileid, nil
}

// 删除小节的文件  附件和 源文件
func DelSectionFileInfo(tx *sql.Tx, fileinfo *model.FileInfo) error {

	if fileinfo.FilePath == "" {
		return errors.New("文件路径错误")
	}

	ret, err := tx.Exec("delete from   fileinfo  where  fileInfoId=?", fileinfo.FileInfoId)

	if err != nil {
		return err
	}
	fileId, err := ret.RowsAffected()

	if err != nil {
		return err
	}

	if fileId < 0 {

		return errors.New("文件表删除错误")
	}

	if In(fileinfo.FileType, OfficeExtArr()) { //判断是否是office文件类型

		_ = os.Remove(PublicPath + fileinfo.FilePath) //删除pdf文件

		//删除office文件
		{
			filenameall := path.Base(fileinfo.FilePath)
			filesuffix := path.Ext(filenameall)

			filename := filenameall[0 : len(filenameall)-len(filesuffix)]
			filedir := path.Dir(PublicPath + fileinfo.FilePath)
			_ = os.Remove(filedir + "/" + filename + fileinfo.FileType)
		}

	} else if In(fileinfo.FileType, VideoExtArr()) { //判断是否是视频文件类型
		filesuffix := path.Ext(fileinfo.FilePath)
		if filesuffix != ".m3u8" {
			return errors.New("文件还在审核，不能删除")
		}
		if !strings.Contains(path.Dir(fileinfo.FilePath), "Resources/Video") {
			return errors.New("文件路径不对，不能删除")
		}
		_ = os.RemoveAll(PublicPath + path.Dir(fileinfo.FilePath) + "/") //直接删除整个muu8的文件夹就行了

	} else if In(fileinfo.FileType, ImgExtArr()) || In(fileinfo.FileType, AudioExtArr()) { //如果是图片或者是音频文件  只需要删除源文件即可
		_ = os.Remove(PublicPath + fileinfo.FilePath) //删除原图片
	} else {
		_ = os.Remove(PublicPath + fileinfo.FilePath) //删除文件本身
	}

	return nil
}

// 删除富文本编辑器的文件
func DelHtmlResources(res string, tx *sql.Tx) bool {

	r := strings.NewReader(res)
	doc, err := goquery.NewDocumentFromReader(r)
	if err == nil {

		flag := true
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					_ = os.Remove(filepathtemp)

					ret, err := tx.Exec(" delete from fileinfo where FilePath=? ", filepathtemp)

					if err != nil {
						flag = false
					}
					n, err := ret.RowsAffected()
					if err != nil {
						flag = false
					}

					if n < 0 {
						flag = false
					}
				}
			}

		})

		if !flag {
			return flag
		}

		doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					_ = os.Remove(filepathtemp)
					ret, err := tx.Exec(" delete from fileinfo where FilePath=? ", filepathtemp)

					if err != nil {
						flag = false
					}
					n, err := ret.RowsAffected()
					if err != nil {
						flag = false
					}

					if n < 0 {
						flag = false
					}
				}
			}

		})
		if !flag {
			return flag
		}

		return true
	} else {
		return false
	}
}

// 修改富文本编辑器里面的文件内容
func UpdateHtmlResources(res string, newres string, tx *sql.Tx) bool {

	var resarr []string    //原先的图片路径数组
	var newresarr []string //修改之后的图片路径数组

	//var removearr []string
	r := strings.NewReader(res)
	doc, err := goquery.NewDocumentFromReader(r)
	if err == nil {
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					resarr = append(resarr, filepathtemp)
				}
			}

		})
		doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					resarr = append(resarr, filepathtemp)
				}
			}

		})
	}

	newr := strings.NewReader(newres)
	newdoc, err := goquery.NewDocumentFromReader(newr)
	if err == nil {
		newdoc.Find("img").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					newresarr = append(newresarr, filepathtemp)
				}
			}

		})
		newdoc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			//解析<div>标签
			//h,err := s.Html()
			v, t := s.Attr("src")
			v = strings.Replace(v, "\n", "", -1)
			if t {
				myUrl, _ := url.Parse(v)
				params, _ := url.ParseQuery(myUrl.RawQuery)
				filepathtemp := params.Get("filename")
				if filepathtemp != "" {
					newresarr = append(newresarr, filepathtemp)
				}
			}

		})
	}

	for i := 0; i < len(resarr); i++ {
		flag := true
		for j := 0; j < len(newresarr); j++ {
			if resarr[i] == newresarr[j] {
				continue
			}
			flag = false

		}
		if !flag {
			//removearr = append(removearr, resarr[i])

			ret, err := tx.Exec(" delete from fileinfo where FilePath=? ", resarr[i])

			if err != nil {
				return false
			}
			n, err := ret.RowsAffected()
			if err != nil {
				return false
			}

			if n < 0 {
				return false
			}
			_ = os.Remove(resarr[i])

		}

	}
	return true
}

// 添加小节的文件  附件和 源文件
func AddFileAndAnnex(tx *sql.Tx, section model.SectionEdit, files *multipart.Form) ([]*model.FileInfo, []*model.OfficePdf, error) {
	var videofileinfo []*model.FileInfo

	var officefileinfo []*model.OfficePdf

	annernum := 0
	filecontentnum := 0
	tx.QueryRow("select  COALESCE( MAX(sectionFileOrder)+1,0) 'sectionFileOrder' from  sectionrelationfile where sectionId=? and sectionFileInfoType=0", section.SectionId).
		Scan(&annernum)

	tx.QueryRow("select  COALESCE( MAX(sectionFileOrder)+1,0) 'sectionFileOrder' from  sectionrelationfile where sectionId=?  and sectionFileInfoType=1", section.SectionId).
		Scan(&filecontentnum)

	if len(files.File["files"]) > 0 {
		for i := 0; i < len(files.File["files"]); i++ {
			file := files.File["files"][i]
			contentfileinfo := new(model.FileInfo)
			contentfileinfo.FileName = file.Filename
			contentfileinfo.FileType = path.Ext(file.Filename)
			pathfilename := uuid.GenUUID()

			if In(contentfileinfo.FileType, PPTExtArr()) {
				contentfileinfo.FilePath = "Resources/Ppt/" + pathfilename + contentfileinfo.FileType
				contentfileinfo.FileUseTo = "小节文件 ppt课小节文件"
			} else if In(contentfileinfo.FileType, ImgExtArr()) {
				contentfileinfo.FilePath = "Resources/Img/" + pathfilename + contentfileinfo.FileType
				contentfileinfo.FileUseTo = "小节文件 图片课小节文件"
			} else if In(contentfileinfo.FileType, VideoExtArr()) {
				contentfileinfo.FileUseTo = "小节文件 视频课小节文件"
				contentfileinfo.FilePath = "Resources/Video/" + pathfilename + contentfileinfo.FileType
			} else if In(contentfileinfo.FileType, PdfExtArr()) {
				contentfileinfo.FilePath = "Resources/Pdf/" + pathfilename + contentfileinfo.FileType
				contentfileinfo.FileUseTo = "小节文件 ppt课小节文件 pdf"
			}
			newfileid := AddFileInfo(tx, contentfileinfo, file)
			if newfileid == 0 {
				return videofileinfo, officefileinfo, errors.New("新建小节失败")
			}
			contentfileinfo.FileInfoId = int(newfileid)

			if section.SectionType == 1 && contentfileinfo.FileInfoId > 0 {
				videofileinfo = append(videofileinfo, contentfileinfo)
			}

			ret, err := tx.Exec("insert into sectionrelationfile(sectionId,fileInfoId,sectionFileOrder,sectionFileInfoType)  values (?,?,?,?)", section.SectionId, newfileid, filecontentnum, 1)
			filecontentnum++
			if err != nil {
				return videofileinfo, officefileinfo, err
			}
			num, err := ret.RowsAffected()
			if err != nil {
				return videofileinfo, officefileinfo, err
			}
			if num == 0 {
				return videofileinfo, officefileinfo, errors.New("新建小节失败")
			}
		}
	}

	if len(files.File["annex"]) > 0 { //附件save
		for i := 0; i < len(files.File["annex"]); i++ {
			file := files.File["annex"][i]
			contentfileinfo := new(model.FileInfo)
			contentfileinfo.FileName = file.Filename
			contentfileinfo.FileType = path.Ext(file.Filename)
			pathfilename := uuid.GenUUID()

			if In(contentfileinfo.FileType, OfficeExtArr()) { // office 附件
				contentfileinfo.FilePath = "Resources/Annex/" + pathfilename + contentfileinfo.FileType
				contentfileinfo.FileUseTo = "小节文件附件Office"
			} else if In(contentfileinfo.FileType, AudioExtArr()) { //音频附件
				contentfileinfo.FilePath = "Resources/Audio/" + pathfilename + contentfileinfo.FileType
				contentfileinfo.FileUseTo = "小节附件 音频附件文件"
			} else if In(contentfileinfo.FileType, VideoExtArr()) { //视频附件
				contentfileinfo.FileUseTo = "小节附件 视频课小节文件"
				contentfileinfo.FilePath = "Resources/Video/" + pathfilename + contentfileinfo.FileType
			} else {
				contentfileinfo.FileUseTo = "小节文件 其他文件"
				contentfileinfo.FilePath = "Resources/Annex/" + pathfilename + contentfileinfo.FileType

			}
			newfileid := AddFileInfo(tx, contentfileinfo, file)
			if newfileid == 0 {
				return videofileinfo, officefileinfo, errors.New("新建小节附件失败")
			}
			contentfileinfo.FileInfoId = int(newfileid)

			ret, err := tx.Exec("insert into sectionrelationfile(sectionId,fileInfoId,sectionFileOrder,sectionFileInfoType)  values (?,?,?,?)", section.SectionId, newfileid, annernum, 0)

			annernum++
			if err != nil {
				return videofileinfo, officefileinfo, err
			}
			num, err := ret.RowsAffected()
			if err != nil {
				return videofileinfo, officefileinfo, err
			}
			if num == 0 {
				return videofileinfo, officefileinfo, errors.New("新建小节附件失败")
			}

			if In(contentfileinfo.FileType, VideoExtArr()) {
				videofileinfo = append(videofileinfo, contentfileinfo)
			}

			if In(contentfileinfo.FileType, OfficeExtArr()) {
				officeinfo := new(model.OfficePdf)
				officeinfo.InputFile = contentfileinfo.FilePath
				officeinfo.OutputFile = "Resources/Annex"

				//TODO 这里可能需要把后缀名去掉
				officeinfo.FileName = contentfileinfo.FileName
				officeinfo.FileInfoId = contentfileinfo.FileInfoId
				officefileinfo = append(officefileinfo, officeinfo)
			}

		}

	}

	//如果是资源列表添加则按资源列表的方式添加
	for i := 0; i < len(section.NewResourceFiles); i++ {

		tempmodel := section.NewResourceFiles[i]
		ret, err := tx.Exec("insert into sectionrelationfile(sectionId,fileInfoId,sectionFileOrder,sectionFileInfoType,isResourceTable)  values (?,?,?,?,?)", section.SectionId, tempmodel.ResourceId, tempmodel.Position, tempmodel.Resourcetype, 1)
		if err != nil {
			return videofileinfo, officefileinfo, err
		}
		num, err := ret.RowsAffected()
		if err != nil {
			return videofileinfo, officefileinfo, err
		}
		if num == 0 {
			return videofileinfo, officefileinfo, errors.New("新建小节 资源列表插入失败")
		}

	}

	return videofileinfo, officefileinfo, nil
}
