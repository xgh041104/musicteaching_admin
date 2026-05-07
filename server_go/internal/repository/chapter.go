package repository

import (
	"errors"
	"server_go/internal/model"
)

type ChapterRepository interface {
	AddChapter(chapter model.Chapter) error
	EditChapter(chapter model.Chapter) error
	DelChapter(chapter model.Chapter) error
	QueryChapterByCourseId(CourseId int) ([]*model.ChapterView, error)
}

type chapterRepository struct {
	*BaseRepository
}

func NewChapterRepository(repository *BaseRepository) ChapterRepository {
	return &chapterRepository{
		BaseRepository: repository,
	}
}

func (r *chapterRepository) AddChapter(chapter model.Chapter) error {

	err := r.db.QueryRow("select  COALESCE( MAX(ChapterOrder)+1,0) 'ChapterOrder' from  chapter where CourseId=?", chapter.CourseId).Scan(&chapter.ChapterOrder)

	if err != nil {
		return err
	}
	ret, err := r.db.Exec("insert into chapter(courseId,chapterTitle,chapterOrder) values(?,?,?) ", chapter.CourseId, chapter.ChapterTitle, chapter.ChapterOrder)
	if err != nil {
		return err
	}

	num, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if num <= 0 {
		return errors.New("章节新增失败")
	}
	return nil
}

func (r *chapterRepository) EditChapter(chapter model.Chapter) error {

	ret, err := r.db.Exec("update chapter set chapterTitle=? where chapterId=? ", chapter.ChapterTitle, chapter.ChapterId)
	if err != nil {
		return err
	}

	num, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if num <= 0 {
		return errors.New("章节修改失败")
	}
	return nil
}

func (r *chapterRepository) DelChapter(chapter model.Chapter) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	sectionnum := 0
	tx.QueryRow("select count(1) from section where chapterId=?", chapter.ChapterId).Scan(&sectionnum)

	if sectionnum > 0 {
		tx.Rollback()
		return errors.New("该章节不能删除，请先删除章节中的小节")
	}

	ret, err := tx.Exec("delete from  chapter   where chapterId=? ", chapter.ChapterId)
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
		return errors.New("章节删除失败")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *chapterRepository) QueryChapterByCourseId(CourseId int) ([]*model.ChapterView, error) {

	chapterarr := make([]*model.ChapterView, 0)
	querychaptersql := "select a.chapterId,a.chapterTitle,a.courseId,a.chapterOrder,COALESCE(b.courseTitle,'') 'courseTitle' from chapter a " +
		"left join course b on a.courseId=b.courseId  " +
		"where a.courseId=? order by a.chapterOrder  "
	rows, err := r.db.Query(querychaptersql, CourseId)

	if err != nil {
		return chapterarr, err
	}

	defer rows.Close()

	for rows.Next() {
		tempmodel := new(model.ChapterView)

		err = rows.Scan(&tempmodel.ChapterId, &tempmodel.ChapterTitle, &tempmodel.CourseId, &tempmodel.ChapterOrder, &tempmodel.CourseTitle)
		if err != nil {
			return chapterarr, err
		}
		chapterarr = append(chapterarr, tempmodel)
	}

	return chapterarr, nil
}
