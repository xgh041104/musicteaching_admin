package repository

import (
	"errors"
	"server_go/internal/model"
	"time"
)

type GiveScoreThemeRepository interface {
	QueryGiveScoreThemeByGiveScoreCategoryId(giveScoreCatgeoryId int64) ([]*model.GiveScoreTheme, error)
	QueryGiveScoreTheme(schoolId int64) ([]*model.GiveScoreThemeView, error)
	AddGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error
	EditGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error
	DelGiveScoreTheme(giveScoreThemeId int64) error
}

type giveScoreThemeRepository struct {
	*BaseRepository
}

func NewGiveScoreThemeRepository(repository *BaseRepository) GiveScoreThemeRepository {
	return &giveScoreThemeRepository{
		BaseRepository: repository,
	}
}

func (r *giveScoreThemeRepository) QueryGiveScoreThemeByGiveScoreCategoryId(giveScoreCatgeoryId int64) ([]*model.GiveScoreTheme, error) {
	giveScoreThemes := make([]*model.GiveScoreTheme, 0)

	query := "SELECT giveScoreThemeId, giveScoreCategoryId, giveScoreThemeTitle, createTime, schoolId FROM giveScoreTheme WHERE giveScoreCategoryId = ?"

	rows, err := r.db.Query(query, giveScoreCatgeoryId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScoreTheme := new(model.GiveScoreTheme)

		err := rows.Scan(&giveScoreTheme.GiveScoreThemeId, &giveScoreTheme.GiveScoreCategoryId, &giveScoreTheme.GiveScoreThemeTitle, &giveScoreTheme.CreateTime, &giveScoreTheme.SchoolId)
		if err != nil {
			return nil, err
		}

		giveScoreThemes = append(giveScoreThemes, giveScoreTheme)
	}

	return giveScoreThemes, nil
}

func (r *giveScoreThemeRepository) QueryGiveScoreTheme(schoolId int64) ([]*model.GiveScoreThemeView, error) {
	giveScoreThemes := make([]*model.GiveScoreThemeView, 0)

	query := " SELECT gt.giveScoreThemeId, gt.giveScoreCategoryId, gt.giveScoreThemeTitle, gt.createTime, " +
		" gt.schoolId,COALESCE(gc.giveScoreCategoryName,'') 'giveScoreCategoryName'," +
		" gc.giveScoreCategoryParentId,COALESCE(gcp.giveScoreCategoryName,'') 'giveScoreCategorParentyName' " +
		" FROM givescoretheme gt " +
		" LEFT JOIN givescorecategory gc ON gt.giveScoreCategoryId = gc.giveScoreCategoryId " +
		"  LEFT JOIN  givescorecategory gcp on gcp.giveScoreCategoryId=gc.giveScoreCategoryParentId " +
		" WHERE (gt.schoolId = ? or gt.schoolId=0) "

	rows, err := r.db.Query(query, schoolId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		giveScoreTheme := new(model.GiveScoreThemeView)

		err := rows.Scan(&giveScoreTheme.GiveScoreThemeId, &giveScoreTheme.GiveScoreCategoryId, &giveScoreTheme.GiveScoreThemeTitle,
			&giveScoreTheme.CreateTime, &giveScoreTheme.SchoolId, &giveScoreTheme.GiveScoreCategoryName, &giveScoreTheme.GiveScoreCategoryParentId,
			&giveScoreTheme.GiveScoreCategoryParentName)
		if err != nil {
			return nil, err
		}

		giveScoreThemes = append(giveScoreThemes, giveScoreTheme)
	}

	return giveScoreThemes, nil
}

func (r *giveScoreThemeRepository) AddGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error {
	giveScoreTheme.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	_, err := r.db.Exec("INSERT INTO givescoretheme(giveScoreCategoryId, giveScoreThemeTitle, createTime, schoolId) VALUES(?,?,?,?)", giveScoreTheme.GiveScoreCategoryId, giveScoreTheme.GiveScoreThemeTitle, giveScoreTheme.CreateTime, giveScoreTheme.SchoolId)
	if err != nil {
		return err
	}

	return nil
}

func (r *giveScoreThemeRepository) EditGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error {

	ret, err := r.db.Exec("UPDATE givescoretheme set giveScoreCategoryId = ?, giveScoreThemeTitle = ?,  schoolId = ? WHERE giveScoreThemeId = ?",
		giveScoreTheme.GiveScoreCategoryId, giveScoreTheme.GiveScoreThemeTitle, giveScoreTheme.SchoolId, giveScoreTheme.GiveScoreThemeId)

	if err != nil {
		return err
	}

	updateNum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if updateNum == 0 {
		return errors.New("操作无效，数据无变化")
	}

	return nil
}

func (r *giveScoreThemeRepository) DelGiveScoreTheme(giveScoreThemeId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	ret, err := tx.Exec("DELETE FROM givescoretheme WHERE giveScoreThemeId = ?", giveScoreThemeId)
	if err != nil {
		tx.Rollback()
		return err
	}

	delnum, err := ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if delnum < 0 {
		tx.Rollback()
		return errors.New("操作错误,评分主题删除失败")
	}

	ret, err = tx.Exec("DELETE FROM givescore WHERE giveScoreThemeId = ?", giveScoreThemeId)
	if err != nil {
		tx.Rollback()
		return err
	}

	delnum, err = ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if delnum < 0 {
		tx.Rollback()
		return errors.New("操作错误,评分主题删除失败")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}

	return nil
}
