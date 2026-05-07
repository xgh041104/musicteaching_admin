package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type CourseRepository interface {
	AddCourse(course model.Course, files *multipart.Form) error
	EditCourse(course model.Course, files *multipart.Form) error
	DelCourse(course model.Course) error
	QueryCourse(shcoolId int, lecturerCommonUserId int) ([]*model.CourseView, error)
	QueryCourseByCategoryId(shcoolId int, lecturerCommonUserId int, page int, pageSize int, categoryId int, titleWords string, trueSchoolId int, userId int) (*model.CourseResult, error)
	QueryCourseDirectory(courseId int, commonUserId int, schoolId int) (model.CourseDirectory, error)
	CollectCourse(courseId int, userId int, isCollect bool) error
	QueryCourseIsCollected(courseId int, userId int) (bool, error)
	QueryPaidCourse() ([]*model.Course, error)
	QueryCollectCourse(userId int, schoolId int, page int, pageSize int, titleWords string) (*model.CourseResult, error)
	QueryPublicCourse(userId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
	QuerySchoolCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
	QueryMyCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
}

type courseRepository struct {
	*BaseRepository
}

func NewCourseRepository(repository *BaseRepository) CourseRepository {
	return &courseRepository{
		BaseRepository: repository,
	}
}

func (r *courseRepository) AddCourse(course model.Course, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var count int
	err = r.db.QueryRow("select count(*) from course where courseTitle = ? and schoolId = ?", course.CourseTitle, course.SchoolId).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("课程标题重复")
	}

	var fileId int64                   //用户接受新增文件表之后fileid的存储
	if len(files.File["files"]) == 1 { // 这里只获取有一个文件的时候  如果files传了多张图片   这里则用for循环了
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于课程封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec("insert into course(courseCategoryId,courseTitle,courseDesc,schoolId,lecturerCommonUserId,courseUpdateTime,courseImgFileId, courseType, status) values(?,?,?,?,?,now(),?,?,?)", course.CourseCategoryId, course.CourseTitle, course.CourseDesc, course.SchoolId, course.LecturerCommonUserId, fileId, course.CourseType, course.Status)
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

func (r *courseRepository) EditCourse(course model.Course, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempCourse := new(model.Course)
	err = tx.QueryRow("select courseImgFileId from course where courseId=? ", course.CourseId).Scan(&tempCourse.CourseImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	//用于接受新增文件表之后fileid的存储

	fileId, err := common.EditCover(tx, files, "用于课程封面", tempCourse.CourseImgFileId, "Resources/Img/", course.IsDelImg)

	if err != nil {

		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("update course set   courseCategoryId=?,courseTitle=?,courseDesc=?,schoolId=?,lecturerCommonUserId=?,courseUpdateTime=now(),courseImgFileId=?, courseType = ?, status = ?   where courseId=?", course.CourseCategoryId, course.CourseTitle, course.CourseDesc, course.SchoolId, course.LecturerCommonUserId, fileId, course.CourseType, course.Status, course.CourseId)
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

func (r *courseRepository) DelCourse(course model.Course) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempcourse := new(model.Course)
	err = tx.QueryRow("select courseImgFileId from course where courseId=? ", course.CourseId).Scan(&tempcourse.CourseImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempcourse.CourseImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempcourse.CourseImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  course    where courseId=?", course.CourseId)
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
		return errors.New("操作错误,课程删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

// 管理中心查询课程
func (r *courseRepository) QueryCourse(schoolId int, lecturerCommonUserId int) ([]*model.CourseView, error) {

	coursevierarr := make([]*model.CourseView, 0)
	queryCourseSql := "SELECT a.courseId,a.courseCategoryId,a.courseType,a.status,a.courseTitle,a.courseDesc,a.schoolId,a.lecturerCommonUserId, " +
		" a.courseUpdateTime,a.courseImgFileId,a.viewNum,a.likeNum,  " +
		" COALESCE(b.filePath,'') 'filePath',COALESCE(c.courseCategoryName,'') 'courseCategoryName', " +
		" COALESCE(d.schoolName,'') 'schoolName',COALESCE(e.commonUserTrueName,'') 'commonUserTrueName' " +
		" FROM course a " +
		" left join fileinfo b on a.courseImgFileId=b.fileInfoId " +
		" left join coursecategory c on a.courseCategoryId=c.courseCategoryId " +
		" LEFT JOIN school d on a.schoolId=d.schoolId " +
		" LEFT JOIN commonuser e on a.lecturerCommonUserId=e.commonUserId "
	var rows *sql.Rows
	var err error

	if lecturerCommonUserId == 0 {
		if schoolId != 0 { //管理员查询学校和公共课程
			queryCourseSql += "WHERE a.schoolId = ? AND a.lecturerCommonUserId = 0 or (a.status = 0 AND (a.courseType = 0 AND a.schoolId = 0 OR (a.courseType = 1 AND a.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))))"
			rows, err = r.db.Query(queryCourseSql, schoolId, schoolId)

			if err != nil {
				return coursevierarr, err
			}
		} else { //超管查询公共课程
			queryCourseSql += " where (a.schoolId=?  or a.schoolId=0 ) and    a.lecturerCommonUserId=0 "
			rows, err = r.db.Query(queryCourseSql, schoolId)

			if err != nil {
				return coursevierarr, err
			}
		}

	} else { //用户查询个人课程，学校，公共课程
		queryCourseSql += "WHERE a.lecturerCommonUserId=? or a.schoolId = ? or (a.status = 0 AND (a.courseType = 0 AND a.schoolId = 0 OR (a.courseType = 1 AND a.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))))"

		rows, err = r.db.Query(queryCourseSql, lecturerCommonUserId, schoolId, schoolId)
		if err != nil {
			return coursevierarr, err
		}
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.CourseType, &tempmodel.Status, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return coursevierarr, err
		}

		coursevierarr = append(coursevierarr, tempmodel)

	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId in (select  courseCategoryParentId from  coursecategory where courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}
	return coursevierarr, nil
}

func (r *courseRepository) QueryCourseDirectory(courseId int, commonUserId int, schoolId int) (model.CourseDirectory, error) {
	var returnModel model.CourseDirectory

	var courseType, status int
	err := r.db.QueryRow("SELECT courseType, status FROM course WHERE courseId = ?", courseId).Scan(&courseType, &status)
	if err != nil {
		return returnModel, err
	}

	if schoolId != 0 {
		if status == 1 {
			return returnModel, errors.New("课程已下架")
		} else if courseType == 1 {
			var count int
			err = r.db.QueryRow("select count(*) from coursepaid where courseId = ? AND schoolId = ?", courseId, schoolId).Scan(&count)

			if err != nil {
				return returnModel, err
			}

			if count == 0 {
				return returnModel, errors.New("课程未购买")
			}
		}
	}

	queryCourseSql := "SELECT a.courseId,a.courseCategoryId,a.courseTitle,CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = a.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected,a.courseDesc,a.schoolId,a.lecturerCommonUserId, " +
		" a.courseUpdateTime,a.courseImgFileId,a.viewNum,a.likeNum,  " +
		" COALESCE(b.filePath,'') 'filePath',COALESCE(c.courseCategoryName,'') 'courseCategoryName', " +
		" COALESCE(d.schoolName,'') 'schoolName',COALESCE(e.commonUserTrueName,'') 'commonUserTrueName' " +
		" FROM course a " +
		" left join fileinfo b on a.courseImgFileId=b.fileInfoId " +
		" left join coursecategory c on a.courseCategoryId=c.courseCategoryId " +
		" LEFT JOIN school d on a.schoolId=d.schoolId " +
		" LEFT JOIN commonuser e on a.lecturerCommonUserId=e.commonUserId " +
		"  where a.courseId=? "

	r.db.QueryRow(queryCourseSql, commonUserId, courseId).Scan(&returnModel.CourseId, &returnModel.CourseCategoryId, &returnModel.CourseTitle, &returnModel.IsCollected,
		&returnModel.CourseDesc, &returnModel.SchoolId, &returnModel.LecturerCommonUserId, &returnModel.CourseUpdateTime, &returnModel.CourseImgFileId,
		&returnModel.ViewNum, &returnModel.LikeNum, &returnModel.FilePath, &returnModel.CourseCategoryName, &returnModel.SchoolName, &returnModel.CommonUserTrueName)

	if returnModel.CourseId == 0 {
		return returnModel, errors.New("查询不到数据")
	}

	r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId=?", returnModel.CourseCategoryId).
		Scan(&returnModel.CourseCategoryParentName, &returnModel.CourseCategoryParentId)

	querychaptersql := "select a.chapterId,a.chapterTitle,a.courseId,a.chapterOrder  from chapter a " +
		"where a.courseId=? order by a.chapterOrder  "
	rows, err := r.db.Query(querychaptersql, courseId)

	if err != nil {
		return returnModel, err
	}

	for rows.Next() {
		tempmodel := new(model.ChapterDirectory)

		err = rows.Scan(&tempmodel.ChapterId, &tempmodel.ChapterTitle, &tempmodel.CourseId, &tempmodel.ChapterOrder)
		if err != nil {
			return returnModel, err
		}
		returnModel.ChapterArr = append(returnModel.ChapterArr, tempmodel)
	}
	rows.Close()

	querysectionsql := "select sectionId,chapterId,sectionTitle,sectionDesc,sectionType,sectionContent,sectionOrder  " +
		"	from section  where chapterId=?  order by sectionOrder "
	for i := 0; i < len(returnModel.ChapterArr); i++ {

		rows, err = r.db.Query(querysectionsql, returnModel.ChapterArr[i].ChapterId)
		if err != nil {
			return returnModel, err
		}
		for rows.Next() {
			tempmodel := new(model.Section)
			rows.Scan(&tempmodel.SectionId, &tempmodel.ChapterId, &tempmodel.SectionTitle, &tempmodel.SectionDesc, &tempmodel.SectionType, &tempmodel.SectionContent, &tempmodel.SectionOrder)
			returnModel.ChapterArr[i].SectionArr = append(returnModel.ChapterArr[i].SectionArr, tempmodel)
		}
		rows.Close()

	}

	return returnModel, nil
}

// 模糊查询
func (r *courseRepository) QueryCourseByCategoryId(schoolId int, lecturerCommonUserId int, page int, pageSize int, categoryId int, titleWords string, trueSchoolId int, userId int) (*model.CourseResult, error) {
	courseResult := new(model.CourseResult)
	var total int

	coursevierarr := make([]*model.CourseView, 0)

	queryTotal := "SELECT COUNT(*) FROM course a "
	queryCourseSql := ""

	var queryParams []interface{}
	var queryParamsTotal []interface{}
	offset := (page - 1) * pageSize

	baseQuery := "SELECT a.courseId,a.courseCategoryId, CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = a.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected, a.courseTitle, a.courseDesc, a.schoolId, a.lecturerCommonUserId, a.courseUpdateTime, a.courseImgFileId, a.viewNum, a.likeNum, " +
		" COALESCE(b.filePath, '') 'filePath', COALESCE(c.courseCategoryName, '') 'courseCategoryName' ,COALESCE(d.schoolName, '') 'schoolName',COALESCE(e.commonUserTrueName, '') 'commonUserTrueName' " +
		" FROM course a LEFT JOIN fileinfo b ON a.courseImgFileId = b.fileInfoId " +
		" LEFT JOIN coursecategory c ON a.courseCategoryId = c.courseCategoryId " +
		" LEFT JOIN school d ON a.schoolId = d.schoolId " +
		" LEFT JOIN commonuser e ON a.lecturerCommonUserId = e.commonUserId "

	if lecturerCommonUserId == 0 && schoolId == 0 { //查询公共课程
		baseQuery += " WHERE a.schoolId = 0 AND a.lecturerCommonUserId = 0 "
		queryTotal += " WHERE a.schoolId = 0 AND a.lecturerCommonUserId = 0 "

		baseQuery2 := " SELECT a.courseId,a.courseCategoryId, CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = a.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected, a.courseTitle,a.courseDesc,a.schoolId,a.lecturerCommonUserId,a.courseUpdateTime,a.courseImgFileId,a.viewNum,a.likeNum, " +
			" COALESCE(b.filePath, '') AS filePath, COALESCE(c.courseCategoryName, '') AS courseCategoryName, COALESCE(d.schoolName, '') AS schoolName, COALESCE(e.commonUserTrueName, '') AS commonUserTrueName " +
			" FROM course a LEFT JOIN fileinfo b ON a.courseImgFileId = b.fileInfoId " +
			" LEFT JOIN coursecategory c ON a.courseCategoryId = c.courseCategoryId " +
			" LEFT JOIN school d ON a.schoolId = d.schoolId " +
			" LEFT JOIN commonuser e ON a.lecturerCommonUserId = e.commonUserId " +
			" JOIN coursepaid f ON a.courseId = f.courseId " +
			" WHERE f.schoolId = ?"

		if trueSchoolId != 0 {
			baseQuery += " AND a.courseType = 0 AND a.status = 0"
			baseQuery2 += " AND a.status = 0"
			queryTotal += " AND a.courseType = 0 AND a.status = 0"
		}

		if categoryId != 0 {
			baseQuery += " AND a.courseCategoryId = ? "
			baseQuery2 += " AND a.courseCategoryId = ? "
			queryTotal += " AND a.courseCategoryId = ? "
			queryParams = append(queryParams, categoryId)
			queryParamsTotal = append(queryParamsTotal, categoryId)
		}

		if titleWords != "" {
			baseQuery += " AND a.courseTitle LIKE ?"
			baseQuery2 += " AND a.courseTitle LIKE ?"
			queryTotal += " AND a.courseTitle LIKE ?"
			queryParams = append(queryParams, "%"+titleWords+"%")
			queryParamsTotal = append(queryParamsTotal, "%"+titleWords+"%")
		}

		tempParams := queryParams
		queryParams = append([]interface{}{userId}, queryParams...)

		queryParams = append(queryParams, userId, trueSchoolId)

		queryParams = append(queryParams, tempParams...)

		queryCourseSql = fmt.Sprintf("SELECT * FROM ( %s UNION %s ) AS combined_results ", baseQuery, baseQuery2)
	} else { //个人，学校课程
		queryParams = append(queryParams, userId)
		if lecturerCommonUserId == 0 {
			baseQuery += " WHERE a.schoolId=? AND a.lecturerCommonUserId=0"
			queryTotal += " WHERE a.schoolId=? AND a.lecturerCommonUserId=0"
			queryParams = append(queryParams, schoolId)
			queryParamsTotal = append(queryParamsTotal, schoolId)
		} else { //查询个人课程
			baseQuery += " WHERE a.lecturerCommonUserId=? "
			queryTotal += " WHERE a.lecturerCommonUserId=? "
			queryParams = append(queryParams, lecturerCommonUserId)
			queryParamsTotal = append(queryParamsTotal, lecturerCommonUserId)
		}

		if categoryId != 0 {
			baseQuery += " AND a.courseCategoryId=?"
			queryTotal += " AND a.courseCategoryId=?"
			queryParams = append(queryParams, categoryId)
			queryParamsTotal = append(queryParamsTotal, categoryId)
		}

		if titleWords != "" {
			baseQuery += " AND a.courseTitle LIKE ?"
			queryParams = append(queryParams, "%"+titleWords+"%")
			queryParamsTotal = append(queryParamsTotal, "%"+titleWords+"%")
		}

		queryCourseSql = baseQuery
	}

	err := r.db.QueryRow(queryTotal, queryParamsTotal...).Scan(&total)
	if err != nil {
		return courseResult, err
	}

	queryCourseSql += " LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize, offset)

	// 执行查询
	rows, err := r.db.Query(queryCourseSql, queryParams...)
	if err != nil {
		return courseResult, err
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.IsCollected, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return courseResult, err
		}

		coursevierarr = append(coursevierarr, tempmodel)
	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId in (select  courseCategoryParentId from  coursecategory where courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}

	courseResult.Records = coursevierarr
	courseResult.Total = total

	return courseResult, nil
}

// 收藏课程
func (r *courseRepository) CollectCourse(courseId int, userId int, isCollect bool) error {
	var err error

	if isCollect {
		var count int
		err = r.db.QueryRow("SELECT COUNT(*) FROM collect WHERE courseId = ? AND commonUserId = ?", courseId, userId).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("该课程已收藏")
		}

		_, err = r.db.Exec("INSERT INTO collect (courseId, commonUserId, collectTime) VALUES (?, ?, now())", courseId, userId)
		if err != nil {
			return err
		}

	} else {
		_, err = r.db.Exec("DELETE FROM collect WHERE courseId = ? AND commonUserId = ?", courseId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

// 查询课程是否被收藏
func (r *courseRepository) QueryCourseIsCollected(courseId int, userId int) (bool, error) {
	var count int
	var isCollected bool

	err := r.db.QueryRow(" SELECT 1 FROM collect WHERE courseId = ? AND commonUserId = ?", courseId, userId).Scan(&count)
	if err != nil {
		return false, nil
	}

	if count > 0 {
		isCollected = true
	} else {
		isCollected = false
	}
	return isCollected, nil
}

// 查询全部付费课程
func (r *courseRepository) QueryPaidCourse() ([]*model.Course, error) {
	courses := make([]*model.Course, 0)

	rows, err := r.db.Query("SELECT course.courseId, course.courseTitle FROM course WHERE courseType = 1 AND status = 0")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var course model.Course
		err = rows.Scan(&course.CourseId, &course.CourseTitle)
		if err != nil {
			return nil, err
		}

		courses = append(courses, &course)
	}

	return courses, nil
}

// 查询收藏课程
func (r *courseRepository) QueryCollectCourse(userId int, schoolId int, page int, pageSize int, titleWords string) (*model.CourseResult, error) {
	courseResult := new(model.CourseResult)
	var total int

	coursevierarr := make([]*model.CourseView, 0)

	queryTotal := "SELECT COUNT(*) FROM course a LEFT JOIN collect ct ON ct.courseId = a.courseId " +
		" WHERE a.status = 0 AND ct.commonUserId = ? AND (a.courseType = 0 OR (a.courseType = 1 AND a.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?)))"

	var queryParams []interface{}
	offset := (page - 1) * pageSize

	queryCourseSql := "SELECT a.courseId,a.courseCategoryId, 'true', a.courseTitle, a.courseDesc, a.schoolId, a.lecturerCommonUserId, a.courseUpdateTime, a.courseImgFileId, a.viewNum, a.likeNum, " +
		" COALESCE(b.filePath, '') 'filePath', COALESCE(c.courseCategoryName, '') 'courseCategoryName' ,COALESCE(d.schoolName, '') 'schoolName',COALESCE(e.commonUserTrueName, '') 'commonUserTrueName' " +
		" FROM course a LEFT JOIN fileinfo b ON a.courseImgFileId = b.fileInfoId " +
		" LEFT JOIN coursecategory c ON a.courseCategoryId = c.courseCategoryId " +
		" LEFT JOIN school d ON a.schoolId = d.schoolId " +
		" LEFT JOIN commonuser e ON a.lecturerCommonUserId = e.commonUserId " +
		" LEFT JOIN collect ct ON ct.courseId = a.courseId " +
		" WHERE  ct.commonUserId = ? AND (a.status = 0 AND (a.courseType = 0 OR (a.courseType = 1  AND a.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))))"

	queryParams = append(queryParams, userId, schoolId)

	if titleWords != "" {
		queryCourseSql += " AND a.courseTitle LIKE ? "
		queryTotal += " AND a.courseTitle LIKE ? "
		queryParams = append(queryParams, "%"+titleWords+"%")
	}

	err := r.db.QueryRow(queryTotal, queryParams...).Scan(&total)
	if err != nil {
		return courseResult, err
	}

	queryCourseSql += " LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize, offset)

	// 执行查询
	rows, err := r.db.Query(queryCourseSql, queryParams...)
	if err != nil {
		return courseResult, err
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.IsCollected, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return courseResult, err
		}

		coursevierarr = append(coursevierarr, tempmodel)
	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId in (select  courseCategoryParentId from  coursecategory where courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}

	courseResult.Records = coursevierarr
	courseResult.Total = total

	return courseResult, nil
}

// 查询公共课程
func (r *courseRepository) QueryPublicCourse(userId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	courseResult := new(model.CourseResult)
	var total int

	coursevierarr := make([]*model.CourseView, 0)

	queryTotal := "SELECT COUNT(*) FROM course c WHERE c.status = 0 AND c.schoolId = 0 AND c.courseType = 0"

	var queryParams []interface{}
	var queryParamsTotal []interface{}
	offset := (page - 1) * pageSize

	queryCourseSql := "SELECT c.courseId, cg.courseCategoryId, CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = c.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected, c.courseTitle, c.courseDesc, c.schoolId, c.lecturerCommonUserId, c.courseUpdateTime, c.courseImgFileId, c.viewNum, c.likeNum, " +
		" COALESCE(f.filePath, '') 'filePath', COALESCE(cg.courseCategoryName, '') 'courseCategoryName' , COALESCE(s.schoolName, '') 'schoolName', COALESCE(u.commonUserTrueName, '') 'commonUserTrueName' " +
		" FROM course c LEFT JOIN fileinfo f ON c.courseImgFileId = f.fileInfoId " +
		" LEFT JOIN coursecategory cg ON c.courseCategoryId = cg.courseCategoryId " +
		" LEFT JOIN school s ON c.schoolId = s.schoolId " +
		" LEFT JOIN commonuser u ON c.lecturerCommonUserId = u.commonUserId " +
		" WHERE c.status = 0 AND c.schoolId = 0 AND c.courseType = 0"

	queryParams = append(queryParams, userId)
	if titleWords != "" {
		queryCourseSql += " AND c.courseTitle LIKE ? "
		queryTotal += " AND c.courseTitle LIKE ? "
		queryParams = append(queryParams, "%"+titleWords+"%")
	}

	if categoryId != 0 {
		queryCourseSql += " AND c.courseCategoryId = ? "
		queryTotal += " AND c.courseCategoryId = ? "
		queryParams = append(queryParams, categoryId)
	}

	queryParamsTotal = queryParams[1:]

	err := r.db.QueryRow(queryTotal, queryParamsTotal...).Scan(&total)
	if err != nil {
		return courseResult, err
	}

	queryCourseSql += " LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize, offset)

	// 执行查询
	rows, err := r.db.Query(queryCourseSql, queryParams...)
	if err != nil {
		return courseResult, err
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.IsCollected, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return courseResult, err
		}

		coursevierarr = append(coursevierarr, tempmodel)
	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("SELECT courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' FROM coursecategory WHERE courseCategoryId IN (SELECT  courseCategoryParentId FROM coursecategory WHERE courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}

	courseResult.Records = coursevierarr
	courseResult.Total = total

	return courseResult, nil
}

// 查询学校课程(含付费)
func (r *courseRepository) QuerySchoolCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	courseResult := new(model.CourseResult)
	var total int

	coursevierarr := make([]*model.CourseView, 0)

	queryTotal := "SELECT COUNT(*) FROM course c WHERE c.schoolId = ? AND c.lecturerCommonUserId = 0 OR (c.status = 0 AND c.courseType = 1  AND c.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))"
	queryCourseSql := ""

	var queryParams []interface{}
	var queryParamsTotal []interface{}
	offset := (page - 1) * pageSize

	//查询学校课程和已经购买的付费课程
	queryCourseSql = "SELECT c.courseId, c.courseCategoryId, CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = c.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected, c.courseTitle, c.courseDesc, c.schoolId, c.lecturerCommonUserId, c.courseUpdateTime, c.courseImgFileId, c.viewNum, c.likeNum, " +
		" COALESCE(f.filePath, '') 'filePath', COALESCE(cg.courseCategoryName, '') 'courseCategoryName' , COALESCE(s.schoolName, '') 'schoolName',COALESCE(u.commonUserTrueName, '') 'commonUserTrueName' " +
		" FROM course c LEFT JOIN fileinfo f ON c.courseImgFileId = f.fileInfoId " +
		" LEFT JOIN coursecategory cg ON c.courseCategoryId = cg.courseCategoryId " +
		" LEFT JOIN school s ON c.schoolId = s.schoolId " +
		" LEFT JOIN commonuser u ON c.lecturerCommonUserId = u.commonUserId " +
		" WHERE c.schoolId = ? " +
		" AND c.lecturerCommonUserId = 0 " +
		" OR ( " +
		" 	c.status = 0 " +
		"	AND c.courseType = 1 " +
		"	AND c.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))"

	queryParams = append(queryParams, userId, schoolId, schoolId)

	if categoryId != 0 {
		queryCourseSql += " AND c.courseCategoryId = ? "
		queryTotal += " AND c.courseCategoryId = ? "
		queryParams = append(queryParams, categoryId)
	}

	if titleWords != "" {
		queryCourseSql += " AND c.courseTitle LIKE ? "
		queryTotal += " AND c.courseTitle LIKE ?"
		queryParams = append(queryParams, "%"+titleWords+"%")
	}

	queryParamsTotal = queryParams[1:]

	err := r.db.QueryRow(queryTotal, queryParamsTotal...).Scan(&total)
	if err != nil {
		return courseResult, err
	}

	queryCourseSql += " LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize, offset)

	// 执行查询
	rows, err := r.db.Query(queryCourseSql, queryParams...)
	if err != nil {
		return courseResult, err
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.IsCollected, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return courseResult, err
		}

		coursevierarr = append(coursevierarr, tempmodel)
	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId in (select  courseCategoryParentId from  coursecategory where courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}

	courseResult.Records = coursevierarr
	courseResult.Total = total

	return courseResult, nil
}

// 查询个人课程(含收藏课程)
func (r *courseRepository) QueryMyCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	courseResult := new(model.CourseResult)
	var total int

	coursevierarr := make([]*model.CourseView, 0)

	queryTotal := "SELECT COUNT(*) FROM course c LEFT JOIN collect ct ON ct.courseId = c.courseId " +
		" WHERE c.lecturerCommonUserId = ? " +
		" OR ( " +
		" 	ct.commonuserId = ? " +
		"	AND c.status = 0 " +
		"	AND ( " +
		"		c.courseType = 0 " +
		"		OR ( " +
		" 			c.courseType = 1  " +
		"			AND c.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))))"

	var queryParams []interface{}
	var queryParamsTotal []interface{}
	offset := (page - 1) * pageSize

	queryCourseSql := "SELECT c.courseId, c.courseCategoryId, CASE WHEN EXISTS (SELECT 1 FROM collect WHERE courseId = c.courseId AND commonUserId = ?) THEN 'true' ELSE 'false' END as is_collected, c.courseTitle, c.courseDesc, c.schoolId, c.lecturerCommonUserId, c.courseUpdateTime, c.courseImgFileId, c.viewNum, c.likeNum, " +
		" COALESCE(f.filePath, '') 'filePath', COALESCE(cg.courseCategoryName, '') 'courseCategoryName' , COALESCE(s.schoolName, '') 'schoolName',COALESCE(u.commonUserTrueName, '') 'commonUserTrueName' " +
		" FROM course c LEFT JOIN fileinfo f ON c.courseImgFileId = f.fileInfoId " +
		" LEFT JOIN coursecategory cg ON c.courseCategoryId = cg.courseCategoryId " +
		" LEFT JOIN school s ON c.schoolId = s.schoolId " +
		" LEFT JOIN commonuser u ON c.lecturerCommonUserId = u.commonUserId " +
		" LEFT JOIN collect ct ON c.courseId = ct.courseId " +
		" WHERE c.lecturerCommonUserId = ? " +
		" OR ( " +
		" 	ct.commonuserId = ? " +
		"	AND c.status = 0 " +
		"	AND ( " +
		"		c.courseType = 0 " +
		"		OR ( " +
		" 			c.courseType = 1  " +
		"			AND c.courseId IN (SELECT courseId FROM coursepaid WHERE schoolId = ?))))"

	queryParams = append(queryParams, userId, userId, userId, schoolId)

	if titleWords != "" {
		queryCourseSql += " AND c.courseTitle LIKE ? "
		queryTotal += " AND c.courseTitle LIKE ? "
		queryParams = append(queryParams, "%"+titleWords+"%")
	}

	if categoryId != 0 {
		queryCourseSql += " AND c.courseTitle LIKE ? "
		queryTotal += " AND c.courseCategoryId=?"
		queryParams = append(queryParams, categoryId)
	}

	queryParamsTotal = queryParams[1:]

	err := r.db.QueryRow(queryTotal, queryParamsTotal...).Scan(&total)
	if err != nil {
		return courseResult, err
	}

	queryCourseSql += " LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize, offset)

	// 执行查询
	rows, err := r.db.Query(queryCourseSql, queryParams...)
	if err != nil {
		return courseResult, err
	}

	for rows.Next() {
		tempmodel := new(model.CourseView)
		err = rows.Scan(&tempmodel.CourseId, &tempmodel.CourseCategoryId, &tempmodel.IsCollected, &tempmodel.CourseTitle, &tempmodel.CourseDesc, &tempmodel.SchoolId, &tempmodel.LecturerCommonUserId,
			&tempmodel.CourseUpdateTime, &tempmodel.CourseImgFileId, &tempmodel.ViewNum, &tempmodel.LikeNum, &tempmodel.FilePath, &tempmodel.CourseCategoryName,
			&tempmodel.SchoolName, &tempmodel.CommonUserTrueName)

		if err != nil {
			return courseResult, err
		}

		coursevierarr = append(coursevierarr, tempmodel)
	}

	for i := 0; i < len(coursevierarr); i++ {
		r.db.QueryRow("select  courseCategoryName 'CourseCategoryParentName',courseCategoryId 'CourseCategoryParentId' from  coursecategory where courseCategoryId in (select  courseCategoryParentId from  coursecategory where courseCategoryId=?) ", coursevierarr[i].CourseCategoryId).
			Scan(&coursevierarr[i].CourseCategoryParentName, &coursevierarr[i].CourseCategoryParentId)
	}

	courseResult.Records = coursevierarr
	courseResult.Total = total

	return courseResult, nil
}
