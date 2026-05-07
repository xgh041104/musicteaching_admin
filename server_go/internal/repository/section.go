package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type SectionRepository interface {
	EditerUploadFile(files *multipart.Form) (string, error)
	AddSection(section model.SectionEdit, files *multipart.Form) error
	DelSection(section model.Section) error
	EditSection(section model.SectionEdit, files *multipart.Form) error
	QuerySectionBySectionId(SectionId int) (model.SectionView, error)
	QuerySectionByChapterId(chapterId int) ([]*model.SectionModel, error)
}
type sectionRepository struct {
	*BaseRepository
}

func NewSectionRepository(repository *BaseRepository) SectionRepository {
	return &sectionRepository{
		BaseRepository: repository,
	}
}

func (r *sectionRepository) AddSection(section model.SectionEdit, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {

		return err
	}

	tx.QueryRow("select  COALESCE( MAX(sectionOrder)+1,0) 'sectionOrder' from  section where chapterId=?", section.ChapterId).Scan(&section.SectionOrder)

	ret, err := tx.Exec("insert into section(chapterId,sectionTitle,sectionDesc,sectionType,sectionContent) values(?,?,?,?,?)", section.ChapterId, section.SectionTitle, section.SectionDesc, section.SectionType, section.SectionContent)
	if err != nil {
		tx.Rollback()
		return err
	}

	newsectionId, err := ret.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}
	if newsectionId == 0 {
		tx.Rollback()
		return errors.New("新建小节失败")
	}

	section.SectionId = int(newsectionId)
	//该切片需要在通道完成视频转换
	// var videofileinfo []*model.FileInfo
	// var officefileinfo []*model.OfficePdf

	videofileinfo, officefileinfo, err := common.AddFileAndAnnex(tx, section, files)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	fileinfochanmutex := common.LoadViderChan()
	for i := 0; i < len(videofileinfo); i++ {
		fileinfochanmutex.Fileinfochan <- videofileinfo[i]
	}
	go func() {
		if !fileinfochanmutex.Mutex.TryLock() {
			return
		}
		common.DealVideoChan(r.db)

		fileinfochanmutex.Mutex.Unlock()

	}()

	officechanmutex := common.LoadOfficeChan()
	for i := 0; i < len(officefileinfo); i++ {
		officechanmutex.OfficePdf <- officefileinfo[i]
	}
	go func() {
		if !officechanmutex.Mutex.TryLock() {
			return
		}
		common.OfficeToPdf(r.db)

		officechanmutex.Mutex.Unlock()

	}()

	return nil
}

// 富文本编辑器上传文件
func (r *sectionRepository) EditerUploadFile(files *multipart.Form) (string, error) {
	fileinfo := new(model.FileInfo)
	if len(files.File["files"]) > 0 {

		tx, err := r.db.Begin()

		if err != nil {
			return "", err
		}

		file := files.File["files"][0]

		fileinfo.FileType = path.Ext(file.Filename)

		fileinfo.FileName = uuid.GenUUID()
		isImg := common.In(fileinfo.FileType, common.ImgExtArr())
		if isImg {
			//是图片
			fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + fileinfo.FileType
			fileinfo.FileUseTo = "用于图文小节的富文本编辑器的图片"
		}
		fileId := common.AddFileInfo(tx, fileinfo, file)

		if fileId == 0 {
			tx.Rollback()
			return "", errors.New("文件上传错误")
		}

		err = tx.Commit()

		if err != nil {
			tx.Rollback()
			return "", errors.New("文件上传错误")
		}
	} else {
		return "", errors.New("没有上传文件")
	}
	return fileinfo.FilePath, nil

}

func (r *sectionRepository) DelSection(section model.Section) error {

	tx, err := r.db.Begin()
	if err != nil {

		return err
	}
	rows, err := tx.Query("select b.fileInfoId,b.fileType,b.fileName,b.fileUseTo,b.filePath from sectionrelationfile a "+
		" left join fileinfo b on a.fileInfoId=b.fileInfoId "+
		" where  sectionId=? and isResourceTable=0", section.SectionId)
	if err != nil {
		tx.Rollback()
		return err
	}
	//  删除office时 需要先删除原本的文件  再删除.pdf文件
	//删除视频时  自己删除视频的文件夹就行
	//删除图片也只需要 删除源文件就行
	//删除图文课 需要把图文课里面的内容全部删掉
	fileinfoViewarr := make([]*model.FileInfo, 0)
	for rows.Next() {
		tempmodel := new(model.FileInfo)
		rows.Scan(&tempmodel.FileInfoId, &tempmodel.FileType, &tempmodel.FileName, &tempmodel.FileUseTo, &tempmodel.FilePath)
		if tempmodel.FileInfoId == 0 {
			continue
		}
		fileinfoViewarr = append(fileinfoViewarr, tempmodel)
	}

	for i := 0; i < len(fileinfoViewarr); i++ {
		err = common.DelSectionFileInfo(tx, fileinfoViewarr[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if section.SectionType == 0 { //删除图文课
		var htmtcontent string
		tx.QueryRow("select sectionContent from  section   where sectionId=?", section.SectionId).Scan(&htmtcontent)

		flag := common.DelHtmlResources(htmtcontent, tx)
		if !flag {
			tx.Rollback()
			return errors.New("删除图文课失败")
		}
	}

	//删除小节数据
	ret, err := tx.Exec("delete from section where sectionId=?", section.SectionId)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if n == 0 {
		tx.Rollback()
		return errors.New("删除小节数据错误")
	}
	//删除小节文件关联表数据
	ret, err = tx.Exec("delete from sectionrelationfile where sectionId=?", section.SectionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	//删除文件表数据
	return nil
}

func (r *sectionRepository) EditSection(section model.SectionEdit, files *multipart.Form) error {

	tx, err := r.db.Begin()
	if err != nil {

		return err
	}

	if section.SectionType == 0 { //图文课特殊处理
		var htmtcontent string
		tx.QueryRow("select sectionContent from  section   where sectionId=?", section.SectionId).Scan(&htmtcontent)

		flag := common.UpdateHtmlResources(htmtcontent, section.SectionContent, tx)
		if !flag {
			tx.Rollback()
			return errors.New("删除图文课失败")
		}
	}

	ret, err := tx.Exec("update section set sectionTitle=?,sectionDesc=?,sectionContent=? where sectionId=?", section.SectionTitle, section.SectionDesc, section.SectionContent, section.SectionId)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err := ret.RowsAffected()

	if err != nil {
		tx.Rollback()

		return err
	}
	if n < 0 {
		tx.Rollback()
		return err
	}

	//该切片需要在通道完成视频转换
	// var videofileinfo []*model.FileInfo
	// var officefileinfo []*model.OfficePdf

	videofileinfo, officefileinfo, err := common.AddFileAndAnnex(tx, section, files)

	if err != nil {
		tx.Rollback()
		return err
	}

	var IsResourceTable int

	if len(section.RemoveFile) > 0 { //修改小节时 删除的文件目录
		for i := 0; i < len(section.RemoveFile); i++ {
			if section.RemoveFile[i] == 0 {
				continue
			}

			tx.QueryRow("select IsResourceTable from sectionrelationfile  where fileInfoId=? and sectionId=?", section.RemoveFile[i], section.SectionId).
				Scan(&IsResourceTable)

			ret, err := tx.Exec("delete from sectionrelationfile where fileInfoId=? and sectionId=?", section.RemoveFile[i], section.SectionId)
			if err != nil {
				tx.Rollback()

				return err
			}
			n, err := ret.RowsAffected()

			if err != nil {
				tx.Rollback()

				return err
			}
			if n == 0 {
				tx.Rollback()

				return err
			}

			if IsResourceTable == 0 {
				filemodel := new(model.FileInfo)
				tx.QueryRow("select fileInfoId,fileType,fileName,fileUseTo,filePath from fileinfo  "+
					" where fileInfoId=?", section.RemoveFile[i]).
					Scan(&filemodel.FileInfoId, &filemodel.FileType, &filemodel.FileName, &filemodel.FileUseTo, &filemodel.FilePath)

				if filemodel.FileInfoId < 0 {

					tx.Rollback()
					return errors.New("删除的文件找不到")
				}

				err = common.DelSectionFileInfo(tx, filemodel)
				if err != nil {

					tx.Rollback()
					return err
				}
			}

		}

	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	fileinfochanmutex := common.LoadViderChan()
	for i := 0; i < len(videofileinfo); i++ {
		fileinfochanmutex.Fileinfochan <- videofileinfo[i]
	}
	go func() {
		if !fileinfochanmutex.Mutex.TryLock() {
			return
		}
		common.DealVideoChan(r.db)

		fileinfochanmutex.Mutex.Unlock()

	}()

	officechanmutex := common.LoadOfficeChan()
	for i := 0; i < len(officefileinfo); i++ {
		officechanmutex.OfficePdf <- officefileinfo[i]
	}
	go func() {
		if !officechanmutex.Mutex.TryLock() {
			return
		}
		common.OfficeToPdf(r.db)

		officechanmutex.Mutex.Unlock()

	}()
	return nil
}

func (r *sectionRepository) QuerySectionBySectionId(SectionId int) (model.SectionView, error) {

	querysql := "select sectionId,a.chapterId,sectionTitle,sectionDesc,sectionType,sectionContent,sectionOrder,b.chapterTitle  from section a  " +
		" left join chapter  b on a.chapterId=b.chapterId 	where sectionId=?"

	var returnmodel model.SectionView
	r.db.QueryRow(querysql, SectionId).Scan(&returnmodel.SectionId, &returnmodel.ChapterId, &returnmodel.SectionTitle, &returnmodel.SectionDesc, &returnmodel.SectionType, &returnmodel.SectionContent, &returnmodel.SectionOrder, &returnmodel.ChapterTitle)

	if returnmodel.SectionId < 0 {
		return returnmodel, errors.New("找不到小节数据")
	}

	filerows, err := r.db.Query("select b.fileInfoId,b.fileType,b.fileName,b.fileUseTo,b.filePath,a.sectionrelationId,a.sectionFileInfoType,a.sectionFileOrder from sectionrelationfile a "+
		" left join fileinfo b on a.fileInfoId=b.fileInfoId "+
		" where sectionId=?  order by a.sectionFileInfoType,a.sectionFileOrder ", SectionId)
	if err != nil {
		return returnmodel, err
	}

	defer filerows.Close()
	for filerows.Next() {
		tempmodel := new(model.SectionrelationFileView)
		filerows.Scan(&tempmodel.FileInfoId, &tempmodel.FileType, &tempmodel.FileName, &tempmodel.FileUseTo, &tempmodel.FilePath, &tempmodel.SectionrelationId, &tempmodel.SectionFileInfoType, &tempmodel.SectionFileOrder)
		if tempmodel.FileInfoId == 0 {
			continue
		}

		if tempmodel.SectionFileInfoType == 0 {
			returnmodel.FileAnnex = append(returnmodel.FileAnnex, tempmodel)
		} else {
			returnmodel.FileContent = append(returnmodel.FileContent, tempmodel)
		}
	}

	return returnmodel, nil
}

func (r *sectionRepository) QuerySectionByChapterId(chapterId int) ([]*model.SectionModel, error) {

	var returnmarr []*model.SectionModel
	querysectionsql := "select sectionId,a.chapterId,sectionTitle,sectionDesc,sectionType,sectionContent,sectionOrder,b.chapterTitle  from section a  " +
		" left join chapter  b on a.chapterId=b.chapterId 	where a.chapterId=? order by sectionOrder"

	rows, err := r.db.Query(querysectionsql, chapterId)
	if err != nil {
		return returnmarr, err
	}
	defer rows.Close()
	for rows.Next() {
		tempmodel := new(model.SectionModel)
		rows.Scan(&tempmodel.SectionId, &tempmodel.ChapterId,
			&tempmodel.SectionTitle, &tempmodel.SectionDesc,
			&tempmodel.SectionType, &tempmodel.SectionContent,
			&tempmodel.SectionOrder, &tempmodel.ChapterTitle)
		returnmarr = append(returnmarr, tempmodel)
	}
	return returnmarr, nil
}

//todo  单独写修改附件顺序 和修改文件内容顺序
