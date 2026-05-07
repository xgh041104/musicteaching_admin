package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"

	"github.com/xuri/excelize/v2"
)

type GiveScoreRepository interface {
	QueryGiveScore(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveTeacherScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveOneselfScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveMutualScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	AddGiveScore(giveScore model.GiveScore, files *multipart.Form) error
	EditGiveScore(giveScore model.GiveScore) error
	DelGiveSCore(giveScoreId int64) error
}

type giveScoreRepository struct {
	*BaseRepository
}

func NewGiveScoreRepository(repository *BaseRepository) GiveScoreRepository {
	return &giveScoreRepository{
		BaseRepository: repository,
	}
}

func (r *giveScoreRepository) QueryGiveScore(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	giveScores := make([]*model.GiveScore, 0)

	rows, err := r.db.Query("SELECT giveScoreId, giveScoreThemeId, studentName, studentGrade, studentSubject, teacherScore, oneselfScore, mutualScore FROM givescore WHERE giveScoreThemeId = ?", giveScoreThemeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScore := new(model.GiveScore)

		err := rows.Scan(&giveScore.GeScoreId, &giveScore.GeScoreThemeId, &giveScore.StudentName, &giveScore.StudentGrade, &giveScore.StudentSubject, &giveScore.TeacherScore, &giveScore.OneselfScore, &giveScore.MutualScore)
		if err != nil {
			return nil, err
		}

		giveScores = append(giveScores, giveScore)
	}

	return giveScores, nil
}

func (r *giveScoreRepository) QueryGiveTeacherScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	giveScores := make([]*model.GiveScore, 0)

	rows, err := r.db.Query("SELECT giveScoreId, giveScoreThemeId, studentName, studentGrade, studentSubject, teacherScore FROM givescore WHERE giveScoreThemeId = ?", giveScoreThemeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScore := new(model.GiveScore)

		err := rows.Scan(&giveScore.GeScoreId, &giveScore.GeScoreThemeId, &giveScore.StudentName, &giveScore.StudentGrade, &giveScore.StudentSubject, &giveScore.TeacherScore)
		if err != nil {
			return nil, err
		}

		giveScores = append(giveScores, giveScore)
	}

	return giveScores, nil
}

func (r *giveScoreRepository) QueryGiveOneselfScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	giveScores := make([]*model.GiveScore, 0)

	rows, err := r.db.Query("SELECT giveScoreId, giveScoreThemeId, studentName, studentGrade, studentSubject, oneselfScore FROM givescore WHERE giveScoreThemeId = ?", giveScoreThemeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScore := new(model.GiveScore)

		err := rows.Scan(&giveScore.GeScoreId, &giveScore.GeScoreThemeId, &giveScore.StudentName, &giveScore.StudentGrade, &giveScore.StudentSubject, &giveScore.OneselfScore)
		if err != nil {
			return nil, err
		}

		giveScores = append(giveScores, giveScore)
	}

	return giveScores, nil
}

func (r *giveScoreRepository) QueryGiveMutualScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	giveScores := make([]*model.GiveScore, 0)

	rows, err := r.db.Query("SELECT giveScoreId, giveScoreThemeId, studentName, studentGrade, studentSubject, mutualScore FROM givescore WHERE giveScoreThemeId = ?", giveScoreThemeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScore := new(model.GiveScore)

		err := rows.Scan(&giveScore.GeScoreId, &giveScore.GeScoreThemeId, &giveScore.StudentName, &giveScore.StudentGrade, &giveScore.StudentSubject, &giveScore.MutualScore)
		if err != nil {
			return nil, err
		}

		giveScores = append(giveScores, giveScore)
	}

	return giveScores, nil
}

func (r *giveScoreRepository) AddGiveScore(giveScore model.GiveScore, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	// 从multipart.Form中获取Excel文件
	if len(files.File["files"]) == 1 {
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)

		if !common.In(fileinfo.FileType, []string{".xlsx", ".xls"}) {
			tx.Rollback()
			return errors.New("不是Excel文件")
		}

		filePart, err := tempfile.Open()
		if err != nil {
			tx.Rollback()
			return err
		}

		defer filePart.Close()

		// 使用excelize解析Excel文件 解析名称，年纪，科目，教室评分，自评分，互评分
		f, err := excelize.OpenReader(filePart)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer f.Close()

		// 获取工作表的所有行
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			tx.Rollback()
			return err
		}

		// 遍历每一行，从第二行开始解析数据
		for i, row := range rows {
			// 跳过第一行（标题行）
			if i == 0 {
				continue
			}

			if len(row) < 6 {
				continue
			}

			// 获取前六列的数据
			studentName := row[0]
			studentGrade := row[1]
			studentSubject := row[2]
			teacherScore := row[3]
			oneselfScore := row[4]
			mutualScore := row[5]

			ret, err := tx.Exec("INSERT INTO givescore(giveScoreThemeId, studentName, studentGrade, studentSubject, teacherScore, oneselfScore, mutualScore) VALUES(?, ?, ?, ?, ?, ?, ?)", giveScore.GeScoreThemeId, studentName, studentGrade, studentSubject, teacherScore, oneselfScore, mutualScore)
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
				return errors.New("操作错误,评分表插入错误")
			}

		}

	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}

	return nil
}

func (r *giveScoreRepository) EditGiveScore(giveScore model.GiveScore) error {
	ret, err := r.db.Exec(" UPDATE giveScore SET teacherScore = ?, oneselfScore = ?, mutualScore = ? WHERE giveScoreId = ?",
		giveScore.TeacherScore, giveScore.OneselfScore, giveScore.MutualScore, giveScore.GeScoreId)

	if err != nil {
		return err
	}

	updatenum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if updatenum < 0 {
		return errors.New("操作失败, 评分插入失败")
	}

	return nil
}

func (r *giveScoreRepository) DelGiveSCore(giveScoreId int64) error {
	ret, err := r.db.Exec("DELETE FROM giveScore WHERE giveScoreId = ?", giveScoreId)
	if err != nil {
		return err
	}

	delnum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if delnum < 0 {
		return errors.New("操作失败，评分删除失败")
	}

	return nil
}
