package repository

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
	"strconv"
	"strings"
)

type SchoolRepository interface {
	AddSchool(school model.School, files *multipart.Form) error
	EditSchool(school model.School, files *multipart.Form) error
	DelSchool(school model.School) error
	QuerySchoolAll() ([]*model.SchoolView, error)
}

type schoolRepository struct {
	*BaseRepository
}

func NewSchoolRepository(repository *BaseRepository) SchoolRepository {
	return &schoolRepository{
		BaseRepository: repository,
	}
}

func (r *schoolRepository) AddSchool(school model.School, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var fileId int64                   //用户接受新增文件表之后fileid的存储
	if len(files.File["files"]) == 1 { // 这里只获取有一个文件的时候  如果files传了多张图片   这里则用for循环了  限定只能传img格式
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于学校封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec(" insert into school( schoolName,address,schoolImgFileId) values (?,?,?)", school.SchoolName, school.Address, fileId)
	if err != nil {
		tx.Rollback()
		return err
	}
	schoolId, err := ret.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	if schoolId < 0 {
		tx.Rollback()
		return errors.New("操作错误,学校插入错误")
	}

	if len(school.PaidCourseIds) != 0 {
		var placeholders []string
		var values []interface{}

		for _, courseId := range school.PaidCourseIds {
			placeholders = append(placeholders, "(?, ?)")
			values = append(values, courseId, schoolId)
		}

		query := fmt.Sprintf("INSERT INTO coursepaid (courseId, schoolId) VALUES %s", strings.Join(placeholders, ", "))

		// 执行批量插入
		_, err = tx.Exec(query, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil

}

func (r *schoolRepository) EditSchool(school model.School, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempschool := new(model.School)
	err = tx.QueryRow("select schoolImgFileId from school where schoolId=? ", school.SchoolId).Scan(&tempschool.SchoolImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	fileId, err := common.EditCover(tx, files, "用于学校封面", tempschool.SchoolImgFileId, "Resources/Img/", school.IsDelImg)

	if err != nil {

		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("  update school set  schoolName=?,address=?,schoolImgFileId=?  where schoolId=?", school.SchoolName, school.Address, fileId, school.SchoolId)
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
		return errors.New("操作错误,学校插入错误")
	}

	_, err = tx.Exec("DELETE FROM coursepaid WHERE schoolId = ?", school.SchoolId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(school.PaidCourseIds) != 0 {
		var placeholders []string
		var values []interface{}

		for _, courseId := range school.PaidCourseIds {
			placeholders = append(placeholders, "(?, ?)")
			values = append(values, courseId, school.SchoolId)
		}

		query := fmt.Sprintf("INSERT INTO coursepaid (courseId, schoolId) VALUES %s", strings.Join(placeholders, ", "))

		// 执行批量插入
		_, err = tx.Exec(query, values...)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

func (r *schoolRepository) DelSchool(school model.School) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempschool := new(model.School)
	err = tx.QueryRow("select schoolImgFileId from school where schoolId=? ", school.SchoolId).Scan(&tempschool.SchoolImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempschool.SchoolImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempschool.SchoolImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  school    where schoolId=?", school.SchoolId)
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

func (r *schoolRepository) QuerySchoolAll() ([]*model.SchoolView, error) {

	schoolviewarr := make([]*model.SchoolView, 0)

	query := "SELECT COALESCE(GROUP_CONCAT(c.courseId), 0) 'courseIds', COALESCE(GROUP_CONCAT(c.courseTitle), '') 'courseTitles', s.schoolId, s.schoolName, s.address, s.schoolImgFileId, COALESCE(f.fileName , '')'schoolImgFileName', COALESCE(f.filePath , '') 'schoolImgPath' " +
		" FROM school s LEFT JOIN fileinfo f ON s.schoolImgFileId = f.fileInfoId " +
		" LEFT JOIN coursepaid cp ON cp.schoolId = s.schoolId " +
		" LEFT JOIN course c ON c.courseId = cp.courseId " +
		" WHERE c.status is null or (c.status = 0 ) GROUP BY s.schoolId"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		tempschoolview := new(model.SchoolView)
		var courseIds, courseTitles string

		err = rows.Scan(&courseIds, &courseTitles, &tempschoolview.SchoolId, &tempschoolview.SchoolName, &tempschoolview.Address, &tempschoolview.SchoolImgFileId, &tempschoolview.SchoolImgFileName, &tempschoolview.SchoolImgPath)
		if err != nil {
			return nil, err
		}

		courseIdList := strings.Split(courseIds, ",")
		tempschoolview.PaidCourseTitles = strings.Split(courseTitles, ",")

		for _, idString := range courseIdList {
			id, err := strconv.Atoi(idString)
			if err != nil {
				return nil, err
			}
			tempschoolview.PaidCourseIds = append(tempschoolview.PaidCourseIds, id)
		}

		schoolviewarr = append(schoolviewarr, tempschoolview)
	}

	for i := 0; i < len(schoolviewarr); i++ {
		r.db.QueryRow("SELECT COUNT(*) FROM commonuser WHERE schoolId = ?", schoolviewarr[i].SchoolId).Scan(&schoolviewarr[i].UserNum)
	}

	return schoolviewarr, nil
}
