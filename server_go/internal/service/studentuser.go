package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type StudentUserService interface {
	QueryStudentListByTeacherId(teacherId, schoolId int64) ([]*model.Student, error)
	QueryStudentList(schoolId int64) ([]*model.StudentUser, error)
	AddStudent(student *model.StudentUser, teacherId int64) error
	AddStudentBatch(fileHeader *multipart.FileHeader, teacherId int64, schoolId int64) error
	EditStudent(student *model.Student) error
	DelStudent(studentId, schoolId int64) error
}

type studentUserService struct {
	*Service
	studentUserRepository repository.StudentUserRepository
}

func NewStudentUserService(service *Service, studentUserRepository repository.StudentUserRepository) StudentUserService {
	return &studentUserService{
		Service:               service,
		studentUserRepository: studentUserRepository,
	}
}

func (s *studentUserService) QueryStudentList(schoolId int64) ([]*model.StudentUser, error) {
	return s.studentUserRepository.QueryStudentList(schoolId)
}

func (s *studentUserService) QueryStudentListByTeacherId(teacherId, schoolId int64) ([]*model.Student, error) {
	if teacherId == 0 {
		return nil, errors.New("参数错误")
	}
	return s.studentUserRepository.QueryStudentByTeacherId(teacherId, schoolId)
}

func (s *studentUserService) AddStudent(student *model.StudentUser, teacherId int64) error {
	return s.studentUserRepository.AddStudent(student, teacherId)
}

func (s *studentUserService) AddStudentBatch(fileHeader *multipart.FileHeader, teacherId int64, schoolId int64) error {
	return s.studentUserRepository.AddStudentBatch(fileHeader, teacherId, schoolId)
}

func (s *studentUserService) EditStudent(student *model.Student) error {
	return s.studentUserRepository.EditStudent(student)
}

func (s *studentUserService) DelStudent(studentId, schoolId int64) error {
	return s.studentUserRepository.DelStudent(studentId, schoolId)
}
