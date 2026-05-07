package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type CourseCategoryRepository interface {
	AddCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error
	EditCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error
	DelCourseCategory(courseCategory model.CourseCategory) error
	QueryCourseCategoryTree(schoolId int) ([]*model.CourseCategoryTreeList, error)
}

type courseCategoryRepository struct {
	*BaseRepository
}

func NewCourseCategoryRepository(repository *BaseRepository) CourseCategoryRepository {
	return &courseCategoryRepository{
		BaseRepository: repository,
	}
}
func (r *courseCategoryRepository) AddCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var fileId int64                   //用户接受新增文件表之后fileid的存储
	if len(files.File["files"]) == 1 { // 这里只获取有一个文件的时候  如果files传了多张图片   这里则用for循环了
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于课程类目封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec("insert into coursecategory(courseCategoryName,courseCategoryParentId,schoolId,courseCategoryImgFileId) values(?,?,?,?)", courseCategory.CourseCategoryName, courseCategory.CourseCategoryParentId, courseCategory.SchoolId, fileId)
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
		return errors.New("操作错误,课程类目插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}
func (r *courseCategoryRepository) EditCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempCourseCategory := new(model.CourseCategory)
	err = tx.QueryRow("select courseCategoryImgFileId from coursecategory where courseCategoryId=? ", courseCategory.CourseCategoryId).Scan(&tempCourseCategory.CourseCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	//用于接受新增文件表之后fileid的存储

	fileId, err := common.EditCover(tx, files, "用于课程类目封面", tempCourseCategory.CourseCategoryImgFileId, "Resources/Img/", courseCategory.IsDelImg)

	if err != nil {

		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("update coursecategory set   courseCategoryName=?,courseCategoryParentId=?,schoolId=?,courseCategoryImgFileId=?    where courseCategoryId=?",
		courseCategory.CourseCategoryName, courseCategory.CourseCategoryParentId, courseCategory.SchoolId, fileId, courseCategory.CourseCategoryId)
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
		return errors.New("操作错误,课程插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}
func (r *courseCategoryRepository) DelCourseCategory(courseCategory model.CourseCategory) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}
	//   如果有课程则需要先删除挂载的课程或者子目录挂载的课程  删除子级目录

	courseCategorynum := 0

	tx.QueryRow("  select count(1) from  coursecategory   where courseCategoryParentId=?", courseCategory.CourseCategoryId).Scan(&courseCategorynum)
	if courseCategorynum > 0 {
		tx.Rollback()
		return errors.New("该类目不能删除，请先删除该类目的子类目")
	}
	coursenum := 0
	tx.QueryRow("  select count(1) from course where courseCategoryId=? ", courseCategory.CourseCategoryId).Scan(&coursenum)
	if coursenum > 0 {
		tx.Rollback()
		return errors.New("该类目不能删除，请先删除该类目的的课程")
	}

	tempCourseCategory := new(model.CourseCategory)
	err = tx.QueryRow("select courseCategoryImgFileId from coursecategory where courseCategoryId=? ", courseCategory.CourseCategoryId).Scan(&tempCourseCategory.CourseCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempCourseCategory.CourseCategoryImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempCourseCategory.CourseCategoryImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  coursecategory    where courseCategoryId=? or courseCategoryParentId=? ", courseCategory.CourseCategoryId, courseCategory.CourseCategoryId)
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
		return errors.New("操作错误,课程类目删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}
func (r *courseCategoryRepository) QueryCourseCategoryTree(schoolId int) ([]*model.CourseCategoryTreeList, error) {

	querytree := "	WITH RECURSIVE cte AS (   " +
		"		SELECT courseCategoryId, courseCategoryName, courseCategoryParentId,schoolId,courseCategoryImgFileId,COALESCE(b.filePath,'') 'filePath', 0 AS level  " +
		"		FROM coursecategory  a  " +
		"		 LEFT JOIN fileinfo b on a.courseCategoryImgFileId=b.fileInfoId " +
		"		WHERE courseCategoryParentId=0   " +
		"		UNION ALL   " +
		"		SELECT c.courseCategoryId, c.courseCategoryName, c.courseCategoryParentId,c.schoolId,c.courseCategoryImgFileId,COALESCE(d.filePath,'') 'filePath',cte.level+1  " +
		"		FROM coursecategory c    " +
		"		JOIN cte ON c.courseCategoryParentId = cte.courseCategoryId   " +
		"					 LEFT JOIN fileinfo d on c.courseCategoryImgFileId=d.fileInfoId" +
		"		)  SELECT * FROM cte where (schoolId=?  or schoolId=0)"

	rows, err := r.db.Query(querytree, schoolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*model.CourseCategoryTreeList
	for rows.Next() {

		tempmodel := new(model.CourseCategoryTreeList)
		err = rows.Scan(&tempmodel.CourseCategoryId, &tempmodel.CourseCategoryName, &tempmodel.CourseCategoryParentId, &tempmodel.SchoolId, &tempmodel.CourseCategoryImgFileId, &tempmodel.FilePath, &tempmodel.Level)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, tempmodel)
	}

	treeLists := Tree(nodes, 0)

	return treeLists, nil

}

// 课程类目树形递归
func Tree(node []*model.CourseCategoryTreeList, courseCategoryParentId int) []*model.CourseCategoryTreeList {
	res := make([]*model.CourseCategoryTreeList, 0)
	for _, v := range node {
		if v.CourseCategoryParentId == courseCategoryParentId {
			v.Children = Tree(node, v.CourseCategoryId)
			res = append(res, v)
		}
	}
	return res
}
