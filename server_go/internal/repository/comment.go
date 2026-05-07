package repository

import (
	"errors"
	"server_go/internal/model"
	"time"
)

type CommentRepository interface {
	AddComment(comment model.Comment) error
	QueryComment() ([]*model.CommenntView, error)
	DelComment(commentId int64) error
	QueryCommentByCommentId(commentId int64) (*model.CommenntView, error)
}
type commentRepository struct {
	*BaseRepository
}

func NewCommentRepository(repository *BaseRepository) CommentRepository {
	return &commentRepository{
		BaseRepository: repository,
	}
}

func (r *commentRepository) AddComment(comment model.Comment) error {

	comment.CommentTime = time.Now().Format("2006-01-02 15:04:05")

	_, err := r.db.Exec("insert into comment(courseId, chapterId, sectionId, commentContent, commentCommonUserId, commentTime) values(?,?,?,?,?,?)", comment.CourseId, comment.ChapterId, comment.SectionId, comment.CommentContent, comment.CommonUserId, comment.CommentTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *commentRepository) QueryComment() ([]*model.CommenntView, error) {
	// 存储多个 Commennt 结构体指针
	commentViews := make([]*model.CommenntView, 0)

	query := "select c.commentId, c.courseId, c.chapterId, c.sectionId, c.commentContent, c.commentCommonUserId, c.commentTime, COALESCE(cu.commonUserTrueName, ''), COALESCE(cr.courseTitle, ''), COALESCE(ch.chapterTitle, ''), COALESCE(s.sectionTitle, '') " +
		" from comment c " +
		" left join commonuser cu on c.commentCommonUserId = cu.commonUserId " +
		" left join course cr on c.courseId = cr.courseId " +
		" left join chapter ch on c.chapterId = ch.chapterId " +
		" left join section s on c.sectionId = s.sectionId"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		commentView := new(model.CommenntView)
		err := rows.Scan(&commentView.CommentId, &commentView.CourseId, &commentView.ChapterId, &commentView.SectionId, &commentView.CommentContent, &commentView.CommonUserId, &commentView.CommentTime, &commentView.CommonUserTrueName, &commentView.CourseTitle, &commentView.ChapterTitle, &commentView.SectionTitle)
		if err != nil {
			return nil, err
		}
		commentViews = append(commentViews, commentView)
	}

	if err := rows.Err(); err != nil {
		// 处理迭代器错误
		return nil, err
	}

	// 返回存储了从数据库中检索到的多个 comment 结构体的指针的切片
	return commentViews, nil
}

func (r *commentRepository) DelComment(commentId int64) error {
	ret, err := r.db.Exec("delete from comment where commentId = ?", commentId)
	if err != nil {
		return err
	}

	delnum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if delnum == 0 {
		return errors.New("操作无效，数据已删除")
	}

	return nil
}

func (r *commentRepository) QueryCommentByCommentId(commentId int64) (*model.CommenntView, error) {
	// 存储多个 Commennt 结构体指针
	commentView := new(model.CommenntView)

	query := "select c.commentId, c.courseId, c.chapterId, c.sectionId, c.commentContent, c.commentCommonUserId, c.commentTime, COALESCE(cu.commonUserTrueName, ''), COALESCE(cr.courseTitle, ''), COALESCE(ch.chapterTitle, ''), COALESCE(s.sectionTitle, '') " +
		" from comment c " +
		" left join commonuser cu on c.commentCommonUserId = cu.commonUserId " +
		" left join course cr on c.courseId = cr.courseId " +
		" left join chapter ch on c.chapterId = ch.chapterId " +
		" left join section s on c.sectionId = s.sectionId" +
		" where c.commentId = ?"

	err := r.db.QueryRow(query, commentId).Scan(&commentView.CommentId, &commentView.CourseId, &commentView.ChapterId, &commentView.SectionId, &commentView.CommentContent, &commentView.CommonUserId, &commentView.CommentTime, &commentView.CommonUserTrueName, &commentView.CourseTitle, &commentView.ChapterTitle, &commentView.SectionTitle)
	if err != nil {
		return nil, err
	}

	if commentView.CommentId <= 0 {
		return nil, errors.New("未查询到数据")
	}
	// 返回存储了从数据库中检索到的多个 comment 结构体的指针的切片
	return commentView, nil
}
