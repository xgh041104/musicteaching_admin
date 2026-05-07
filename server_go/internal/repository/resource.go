package repository

import (
	"database/sql"
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type ResourceRepository interface {
	AddResource(resource model.Resources, files *multipart.Form) error
	DelResource(resource model.Resources) error
	QueryResourceByResourceCategoryId(schoolId int, lecturerCommonUserId int, ResourceCategoryId int) ([]*model.ResourcesView, error)
	QueryResourceById(resourceId int) (*model.ResourcesView, error)
}
type resourceRepository struct {
	*BaseRepository
}

func NewResourceRepository(repository *BaseRepository) ResourceRepository {
	return &resourceRepository{
		BaseRepository: repository,
	}
}

func (r *resourceRepository) AddResource(resource model.Resources, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var videofileinfo []*model.FileInfo

	var officefileinfo []*model.OfficePdf

	if len(files.File["files"]) > 0 { //资源文件
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = tempfile.Filename

		filename := uuid.GenUUID()
		resource.ResourceType = fileinfo.FileType
		if common.In(fileinfo.FileType, common.PPTExtArr()) {
			fileinfo.FilePath = "Resources/Ppt/" + filename + fileinfo.FileType
			fileinfo.FileUseTo = "资源文件Office"
		} else if common.In(fileinfo.FileType, common.OfficeNotFoundPPTExtArr()) { // office 文件不包括ppt
			fileinfo.FilePath = "Resources/Annex/" + filename + fileinfo.FileType
			fileinfo.FileUseTo = "资源文件Office"
		} else if common.In(fileinfo.FileType, common.AudioExtArr()) { //音频 文件不包括ppt
			fileinfo.FilePath = "Resources/Audio/" + filename + fileinfo.FileType
			fileinfo.FileUseTo = "资源文件 音频文件"
		} else if common.In(fileinfo.FileType, common.VideoExtArr()) { //视频文件
			fileinfo.FileUseTo = "资源文件 视频文件"
			fileinfo.FilePath = "Resources/Video/" + filename + fileinfo.FileType
		} else if common.In(fileinfo.FileType, common.ImgExtArr()) { //图片文件
			fileinfo.FileUseTo = "资源文件 图片"
			fileinfo.FilePath = "Resources/Img/" + filename + fileinfo.FileType
		} else { //其他文件
			fileinfo.FileUseTo = "资源文件  其他文件"
			fileinfo.FilePath = "Resources/Annex/" + filename + fileinfo.FileType
		}

		resource.ResourceFileId = int(common.AddFileInfo(tx, fileinfo, tempfile))

		if resource.ResourceFileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

		fileinfo.FileInfoId = resource.ResourceFileId           //新文件id赋值到文件对象
		if common.In(fileinfo.FileType, common.VideoExtArr()) { //视频文件
			videofileinfo = append(videofileinfo, fileinfo)
		} else if common.In(fileinfo.FileType, common.OfficeExtArr()) {
			officeinfo := new(model.OfficePdf)
			officeinfo.InputFile = fileinfo.FilePath
			officeinfo.OutputFile = "Resources/Annex"

			officeinfo.FileName = fileinfo.FileName
			officeinfo.FileInfoId = fileinfo.FileInfoId
			officefileinfo = append(officefileinfo, officeinfo)
		}

	}

	if len(files.File["cover"]) > 0 { //资源资源封面
		tempfile := files.File["cover"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于资源封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		resource.ResourceImgFileId = int(common.AddFileInfo(tx, fileinfo, tempfile))

		if resource.ResourceImgFileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec("insert into resource(resourceCategoryId,resourceName,resourceType,resourceFileId,schoolId,resourceDesc,resourceImgFileId,lecturerCommonUserId) values(?,?,?,?,?,?,?,?)",
		resource.ResourceCategoryId, resource.ResourceName, resource.ResourceType, resource.ResourceFileId, resource.SchoolId, resource.ResourceDesc, resource.ResourceImgFileId, resource.LecturerCommonUserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	num, err := ret.RowsAffected()

	if err != nil {
		tx.Rollback()
		return err
	}
	if num <= 0 {
		tx.Rollback()
		return errors.New("新建资源失败")
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
func (r *resourceRepository) DelResource(resource model.Resources) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// err = tx.QueryRow("select count(1) from sectionrelationresource where resourceId=? ", resource.ResourceId).Scan(&isdel)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// if isdel > 0 {
	// 	tx.Rollback()
	// 	return errors.New("当前资源不能删除  ，资源已被使用")
	// }

	resourcefileid := 0
	tx.QueryRow("select  resourceFileId from resource where resourceId=? ", resource.ResourceId).Scan(&resourcefileid)

	isdel := 0
	err = tx.QueryRow("select count(1) from sectionrelationfile where isResourceTable=1  and fileInfoId=?", resourcefileid).Scan(&isdel)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isdel > 0 {
		tx.Rollback()
		return errors.New("该资源已被使用，不能删除")
	}

	err = tx.QueryRow("select resourceFileId,resourceImgFileId from resource where resourceId=?    ", resource.ResourceId).Scan(&resource.ResourceFileId, &resource.ResourceImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if resource.ResourceFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = resource.ResourceFileId
		err = tx.QueryRow("select fileInfoId,fileType,fileName,fileUseTo,filePath from fileinfo where  fileInfoId=? ", resource.ResourceFileId).
			Scan(&fileinfo.FileInfoId, &fileinfo.FileType, &fileinfo.FileName, &fileinfo.FileUseTo, &fileinfo.FilePath)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}

		if fileinfo.FileInfoId == 0 {
			tx.Rollback()
			return errors.New("文件未找到")
		}
		err = common.DelSectionFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	if resource.ResourceImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = resource.ResourceImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  resource    where resourceId=?", resource.ResourceId)
	if err != nil {
		tx.Rollback()
		return err
	}
	num, err := ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if num < 0 {
		tx.Rollback()
		return errors.New("操作错误,学校删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}
func (r *resourceRepository) QueryResourceByResourceCategoryId(schoolId int, lecturerCommonUserId int, ResourceCategoryId int) ([]*model.ResourcesView, error) {
	resourceviewarr := make([]*model.ResourcesView, 0)
	queryCoyrseSql := " SELECT 		a.resourceId,a.resourceCategoryId,a.resourceName,a.resourceType,a.resourceFileId, " +
		" a.schoolId,a.resourceDesc,a.resourceImgFileId,a.lecturerCommonUserId, " +
		"	COALESCE(resourcefile.filePath,'') 'resourcefile',			COALESCE(resourcefile.fileName,'') 'resourceFileName'," +
		" COALESCE(imgfile.filePath,'') 'imgfilePath', 		  " +
		" COALESCE(c.resourceCategoryName,'') 'resourceCategoryName', " +
		" COALESCE(d.schoolName,'') 'schoolName',COALESCE(e.commonUserTrueName,'') 'commonUserTrueName' " +
		" FROM resource a " +
		" left join fileinfo imgfile on a.resourceImgFileId=imgfile.fileInfoId " +
		" left join fileinfo resourcefile on a.resourceFileId=resourcefile.fileInfoId " +
		" left join resourcecategory c on a.resourceCategoryId=c.resourceCategoryId " +
		" LEFT JOIN school d on a.schoolId=d.schoolId " +
		" LEFT JOIN commonuser e on a.lecturerCommonUserId=e.commonUserId WHERE a.resourceCategoryId=? "

	var rows *sql.Rows
	var err error
	if lecturerCommonUserId == 0 {
		queryCoyrseSql += "  and  (a.schoolId=?  or a.schoolId=0 ) and  a.lecturerCommonUserId=0 "
		rows, err = r.db.Query(queryCoyrseSql, ResourceCategoryId, schoolId)

		if err != nil {
			return resourceviewarr, err
		}
	} else {
		queryCoyrseSql += "  and   (a.lecturerCommonUserId=? or (a.schoolId=0 or a.schoolId=?))  "
		rows, err = r.db.Query(queryCoyrseSql, ResourceCategoryId, lecturerCommonUserId, schoolId)
		if err != nil {
			return resourceviewarr, err
		}
	}

	defer rows.Close()
	for rows.Next() {

		tempmodel := new(model.ResourcesView)

		rows.Scan(&tempmodel.ResourceId, &tempmodel.ResourceCategoryId, &tempmodel.ResourceName, &tempmodel.ResourceType, &tempmodel.ResourceFileId, &tempmodel.SchoolId,
			&tempmodel.ResourceDesc, &tempmodel.ResourceImgFileId, &tempmodel.LecturerCommonUserId, &tempmodel.ResourceFilePath, &tempmodel.ResourceFileName,
			&tempmodel.ImgfilePath, &tempmodel.ResourceCategoryName, &tempmodel.SchoolName, &tempmodel.CommonUserTrueName)
		resourceviewarr = append(resourceviewarr, tempmodel)
	}
	return resourceviewarr, nil
}
func (r *resourceRepository) QueryResourceById(resourceId int) (*model.ResourcesView, error) {
	tempmodel := new(model.ResourcesView)
	queryCoyrseSql := " SELECT 		a.resourceId,a.resourceCategoryId,a.resourceName,a.resourceType,a.resourceFileId, " +
		" a.schoolId,a.resourceDesc,a.resourceImgFileId,a.lecturerCommonUserId, " +
		"	COALESCE(resourcefile.filePath,'') 'resourcefile',			COALESCE(resourcefile.fileName,'') 'resourceFileName'," +
		" COALESCE(imgfile.filePath,'') 'imgfilePath', 		  " +
		" COALESCE(c.resourceCategoryName,'') 'resourceCategoryName', " +
		" COALESCE(d.schoolName,'') 'schoolName',COALESCE(e.commonUserTrueName,'') 'commonUserTrueName' " +
		" FROM resource a " +
		" left join fileinfo imgfile on a.resourceImgFileId=imgfile.fileInfoId " +
		" left join fileinfo resourcefile on a.resourceFileId=resourcefile.fileInfoId " +
		" left join resourcecategory c on a.resourceCategoryId=c.resourceCategoryId " +
		" LEFT JOIN school d on a.schoolId=d.schoolId " +
		" LEFT JOIN commonuser e on a.lecturerCommonUserId=e.commonUserId WHERE a.resourceId=? "

	err := r.db.QueryRow(queryCoyrseSql, resourceId).Scan(&tempmodel.ResourceId, &tempmodel.ResourceCategoryId, &tempmodel.ResourceName, &tempmodel.ResourceType, &tempmodel.ResourceFileId, &tempmodel.SchoolId,
		&tempmodel.ResourceDesc, &tempmodel.ResourceImgFileId, &tempmodel.LecturerCommonUserId, &tempmodel.ResourceFilePath, &tempmodel.ResourceFileName,
		&tempmodel.ImgfilePath, &tempmodel.ResourceCategoryName, &tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

	if err != nil {
		return nil, err
	}

	if tempmodel.ResourceId == 0 {
		return nil, err
	}
	return tempmodel, nil
}
