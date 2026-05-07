package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type CourseService interface {
	AddCourse(course model.Course, files *multipart.Form) error
	EditCourse(course model.Course, files *multipart.Form) error
	DelCourse(course model.Course) error
	QueryCourse(schoolId int, lecturerCommonUserId int) ([]*model.CourseView, error)
	QueryCourseByCategoryId(schoolId int, lecturerCommonUserId int, page int, pageSize int, categoryId int, titleWords string, trueSchoolId int, userId int) (*model.CourseResult, error)
	QueryCourseDirectory(CourseId int, commonUserId int, schoolId int) (model.CourseDirectory, error)
	CollectCourse(courseId int, userId int, isCollect bool) error
	QueryCourseIsCollected(courseId int, userId int) (bool, error)
	QueryPaidCourse() ([]*model.Course, error)
	QueryCollectCourse(userId int, schoolId int, page int, pageSize int, titleWords string) (*model.CourseResult, error)
	QueryPublicCourse(userId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
	QuerySchoolCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
	QueryMyCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error)
}

type courseService struct {
	*Service
	courseRepository repository.CourseRepository
}

func NewCourseService(service *Service, courseRepository repository.CourseRepository) CourseService {
	return &courseService{
		Service:          service,
		courseRepository: courseRepository,
	}
}

func (s *courseService) AddCourse(course model.Course, files *multipart.Form) error {

	return s.courseRepository.AddCourse(course, files)
}

func (s *courseService) EditCourse(course model.Course, files *multipart.Form) error {

	if course.CourseId == 0 {
		return errors.New("参数错误")
	}
	return s.courseRepository.EditCourse(course, files)
}

func (s *courseService) DelCourse(course model.Course) error {
	if course.CourseId == 0 {
		return errors.New("参数错误")
	}
	return s.courseRepository.DelCourse(course)
}

func (s *courseService) QueryCourse(schoolId int, lecturerCommonUserId int) ([]*model.CourseView, error) {

	return s.courseRepository.QueryCourse(schoolId, lecturerCommonUserId)
}
func (s *courseService) QueryCourseDirectory(CourseId int, commonUserId int, schoolId int) (model.CourseDirectory, error) {

	return s.courseRepository.QueryCourseDirectory(CourseId, commonUserId, schoolId)
}

func (s *courseService) QueryCourseByCategoryId(schoolId int, lecturerCommonUserId int, page int, pageSize int, categoryId int, titleWords string, trueSchool int, userId int) (*model.CourseResult, error) {
	return s.courseRepository.QueryCourseByCategoryId(schoolId, lecturerCommonUserId, page, pageSize, categoryId, titleWords, trueSchool, userId)
}
func (s *courseService) CollectCourse(courseId int, userId int, isCollect bool) error {
	return s.courseRepository.CollectCourse(courseId, userId, isCollect)
}

func (s *courseService) QueryCourseIsCollected(courseId int, userId int) (bool, error) {
	return s.courseRepository.QueryCourseIsCollected(courseId, userId)
}

func (s *courseService) QueryPaidCourse() ([]*model.Course, error) {
	return s.courseRepository.QueryPaidCourse()
}

func (s *courseService) QueryCollectCourse(userId int, schoolId int, page int, pageSize int, titleWords string) (*model.CourseResult, error) {
	return s.courseRepository.QueryCollectCourse(userId, schoolId, page, pageSize, titleWords)
}

func (s *courseService) QueryPublicCourse(userId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	return s.courseRepository.QueryPublicCourse(userId, page, pageSize, titleWords, categoryId)
}
func (s *courseService) QuerySchoolCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	return s.courseRepository.QuerySchoolCourse(userId, schoolId, page, pageSize, titleWords, categoryId)
}
func (s *courseService) QueryMyCourse(userId int, schoolId int, page int, pageSize int, titleWords string, categoryId int) (*model.CourseResult, error) {
	return s.courseRepository.QueryMyCourse(userId, schoolId, page, pageSize, titleWords, categoryId)
}
