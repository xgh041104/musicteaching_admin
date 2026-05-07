package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/md5"
	"strings"

	"github.com/xuri/excelize/v2"
)

type StudentUserRepository interface {
	QueryStudentByTeacherId(teacherId, schoolId int64) ([]*model.Student, error)
	QueryStudentList(schoolId int64) ([]*model.StudentUser, error)
	AddStudent(student *model.StudentUser, teacherId int64) error
	AddStudentBatch(fileHeader *multipart.FileHeader, teacherId int64, schoolId int64) error
	EditStudent(student *model.Student) error
	DelStudent(studentId, schoolId int64) error
}
type studentUserRepository struct {
	*BaseRepository
}

func NewStudentUserRepository(repository *BaseRepository) StudentUserRepository {
	return &studentUserRepository{
		BaseRepository: repository,
	}
}

func (r *studentUserRepository) QueryStudentByTeacherId(teacherId, schoolId int64) ([]*model.Student, error) {
	var students []*model.Student

	rows, err := r.db.Query("SELECT s.studentId, studentName, studentAccount, s.schoolId, s.className, s.userType "+
		" FROM  studentuser s "+
		" LEFT JOIN student_teacher st ON s.studentId = st.studentId "+
		" WHERE st.teacherId = ? AND s.schoolId = ?", teacherId, schoolId)
	if err != nil {
		return nil, errors.New("系统错误")
	}

	defer rows.Close()

	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.StudentId, &student.StudentName, &student.StudentAccount, &student.SchoolId, &student.ClassName, &student.UserType)
		if err != nil {
			return nil, errors.New("系统错误")
		}
		students = append(students, &student)
	}

	for _, student := range students {
		rows, err := r.db.Query("SELECT teacherId FROM student_teacher WHERE studentId = ?", student.StudentId)
		if err != nil {
			return nil, errors.New("系统错误")
		}

		defer rows.Close()

		for rows.Next() {
			var teacherId int
			err := rows.Scan(&teacherId)
			if err != nil {
				return nil, errors.New("系统错误")
			}
			student.TeacherIdList = append(student.TeacherIdList, teacherId)
		}
	}

	return students, nil
}

func (r *studentUserRepository) QueryStudentList(schoolId int64) ([]*model.StudentUser, error) {
	var students []*model.StudentUser

	rows, err := r.db.Query("SELECT studentId, studentName, studentAccount, schoolId, className, userType "+
		" FROM  studentuser s WHERE schoolId = ?", schoolId)
	if err != nil {
		return nil, errors.New("系统错误")
	}

	defer rows.Close()

	for rows.Next() {
		var student model.StudentUser
		err := rows.Scan(&student.StudentId, &student.StudentName, &student.StudentAccount, &student.SchoolId, &student.ClassName, &student.UserType)
		if err != nil {
			return nil, errors.New("系统错误")
		}
		students = append(students, &student)
	}
	return students, nil
}

func (r *studentUserRepository) AddStudent(student *model.StudentUser, teacherId int64) error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM studentuser WHERE studentAccount = ?", student.StudentAccount).Scan(&count)
	if err != nil {
		return errors.New("系统错误")
	}

	if count > 0 {
		return errors.New("用户名重复")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("系统错误")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newpwd := md5.Md5(student.StudentPwd)
	student.UserType = 3

	ret, err := tx.Exec("INSERT INTO studentuser(studentName, studentAccount, studentPwd, schoolId, className, userType) VALUES(?,?,?,?,?,?) ",
		student.StudentName, student.StudentAccount, newpwd, student.SchoolId, student.ClassName, student.UserType)

	if err != nil {
		return errors.New("系统错误")
	}

	studentId, err := ret.LastInsertId()
	if err != nil {
		return errors.New("系统错误")
	}

	_, err = tx.Exec("INSERT INTO student_teacher(teacherId, studentId) VALUES (?, ?)", teacherId, studentId)
	if err != nil {
		return errors.New("系统错误")
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("系统错误")
	}

	return nil
}

func (r *studentUserRepository) AddStudentBatch(fileHeader *multipart.FileHeader, teacherId int64, schoolId int64) error {
	if !common.In(path.Ext(fileHeader.Filename), []string{".xlsx", ".xls"}) {
		return errors.New("文件格式错误, 只接受excel文件")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("文件打开失败")
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return errors.New("文件读取失败")
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return errors.New("文件读取失败")
	}

	students := make([]*model.StudentUser, 0)

	for _, row := range rows[1:] {
		if len(row) != 4 {
			return errors.New("文件内容缺失")
		}

		student := &model.StudentUser{
			StudentAccount: row[0],
			StudentPwd:     row[1],
			StudentName:    row[2],
			ClassName:      row[3],
		}
		students = append(students, student)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("系统错误")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, student := range students {
		var studentId int64

		err = tx.QueryRow("SELECT studentId FROM studentuser WHERE studentAccount = ?", student.StudentAccount).Scan(&studentId)
		if err != nil && err != sql.ErrNoRows {
			return errors.New("系统错误")
		}

		if studentId == 0 {
			nwpd := md5.Md5(student.StudentPwd)
			ret, err := tx.Exec("INSERT INTO studentuser (studentAccount, studentPwd, studentName, classname, schoolId, userType) "+
				" VALUES (?, ?, ?, ?, ?, ?) ", student.StudentAccount, nwpd, student.StudentName, student.ClassName, schoolId, 3)

			if err != nil {
				return errors.New("系统错误")
			}

			studentId, err = ret.LastInsertId()
			if err != nil {
				return errors.New("系统错误")
			}
		}

		var count int
		err = r.db.QueryRow("SELECT COUNT(*) FROM student_teacher WHERE teacherId = ? AND studentId = ?", teacherId, studentId).Scan(&count)
		if err != nil {
			return errors.New("系统错误")
		}

		if count != 0 {
			continue
		}

		_, err = tx.Exec("INSERT INTO student_teacher (teacherId, studentId) VALUES(?, ?)", teacherId, studentId)
		if err != nil {
			return errors.New("系统错误")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("事务提交失败")
	}

	return nil
}

func (r *studentUserRepository) EditStudent(student *model.Student) error {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM studentuser WHERE studentId = ?", student.StudentId).Scan(&count)
	if err != nil {
		return errors.New("系统错误")
	}

	if count == 0 {
		return errors.New("学生不存在")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("系统错误")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE studentuser SET studentName = ?, className = ? WHERE studentId = ?", student.StudentName, student.ClassName, student.StudentId)
	if err != nil {
		return errors.New("系统错误")
	}

	_, err = tx.Exec("DELETE FROM student_teacher WHERE studentId = ?", student.StudentId)
	if err != nil {
		return errors.New("系统错误")
	}

	var values []interface{}
	var valueStrs []string

	for _, teacherId := range student.TeacherIdList {
		values = append(values, teacherId, student.StudentId)
		valueStrs = append(valueStrs, "(?, ?)")
	}

	query := fmt.Sprintf("INSERT INTO student_teacher (teacherId, studentId) VALUES %s", strings.Join(valueStrs, ", "))

	_, err = tx.Exec(query, values...)
	if err != nil {
		return errors.New("系统错误")
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("系统错误")
	}

	return nil
}

func (r *studentUserRepository) DelStudent(studentId, schoolId int64) error {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM studentuser WHERE studentId = ? AND schoolId = ?", studentId, schoolId).Scan(&count)
	if err != nil {
		return errors.New("系统错误")
	}

	if count == 0 {
		return errors.New("学生不存在")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("系统错误")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM studentuser WHERE studentId = ?", studentId)
	if err != nil {
		return errors.New("系统错误")
	}

	_, err = tx.Exec("DELETE FROM student_teacher WHERE studentId = ?", studentId)
	if err != nil {
		return errors.New("系统错误")

	}

	err = tx.Commit()
	if err != nil {
		return errors.New("系统错误")
	}

	return nil
}
