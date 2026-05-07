package service

import (
	"errors"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type CommentService interface {
	AddComment(comment model.Comment) error
	QueryComment() ([]*model.CommenntView, error)
	DelComment(commentId int64) error
	QueryCommentByCommentId(commentId int64) (*model.CommenntView, error)
}

type commentService struct {
	*Service
	commentRepository repository.CommentRepository
}

func NewCommentService(service *Service, commentRepository repository.CommentRepository) CommentService {
	return &commentService{
		Service:           service,
		commentRepository: commentRepository,
	}
}

func (s *commentService) AddComment(comment model.Comment) error {
	return s.commentRepository.AddComment(comment)
}

func (s *commentService) QueryComment() ([]*model.CommenntView, error) {
	return s.commentRepository.QueryComment()
}

func (s *commentService) DelComment(commentId int64) error {
	if commentId == 0 {
		return errors.New("参数错误")
	}

	return s.commentRepository.DelComment(commentId)
}

func (s *commentService) QueryCommentByCommentId(commentId int64) (*model.CommenntView, error) {
	return s.commentRepository.QueryCommentByCommentId(commentId)
}
