package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type GiveScoreCategoryRepository interface {
	QueryGiveScoreCategory(schoolId int64) ([]*model.GiveScoreCategoryView, error)
	QueryGiveScoreCategoryTree(schoolId int64) ([]*model.GiveScoreCategoryTreeList, error)
	QueryGiveScoreCategoryParentNodeByParentId(giveScoreCategoryParentId int64) (*model.GiveScoreCategoryView, error)
	QueryGiveScoreCategoryChildNodesById(giveScoreCategoryId int64, schoolId int64) ([]*model.GiveScoreCategoryView, error)
	AddGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error
	EditGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error
	DelGiveScoreCategory(giveScoreCategoryId int64) error
}

type giveScoreCategoryRepository struct {
	*BaseRepository
}

func NewGiveScoreCategoryRepository(repository *BaseRepository) GiveScoreCategoryRepository {
	return &giveScoreCategoryRepository{
		BaseRepository: repository,
	}
}

func (r *giveScoreCategoryRepository) AddGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var fileId int64
	if len(files.File["files"]) == 1 {
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于评分类目封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec("INSERT INTO givescorecategory(giveScoreCategoryName, giveScoreCategoryDesc, giveScoreCategoryParentId, schoolId, giveScoreCategoryImgFileId) VALUES(?,?,?,?,?)",
		giveScoreCategory.GiveScoreCategoryName, giveScoreCategory.GiveScoreCategoryDesc, giveScoreCategory.GiveScoreCategoryParentId,
		giveScoreCategory.SchoolId, fileId)
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
		return errors.New("操作错误,评分类目插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}

	return nil
}

func (r *giveScoreCategoryRepository) DelGiveScoreCategory(giveScoreCategoryId int64) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempgiveScoreCategory := new(model.GiveScoreCategory)
	err = tx.QueryRow("SELECT giveScoreCategoryImgFileId FROM giveScoreCategory WHERE giveScoreCategoryId = ?", giveScoreCategoryId).Scan(&tempgiveScoreCategory.GiveScoreCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempgiveScoreCategory.GiveScoreCategoryImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempgiveScoreCategory.GiveScoreCategoryImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("DELETE FROM giveScoreCategory WHERE giveScoreCategoryId = ?", giveScoreCategoryId)
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
		return errors.New("操作错误,评分类目删除失败")
	}

	ret, err = tx.Exec("DELETE FROM giveScoreCategory WHERE giveScoreCategoryParentId = ?", giveScoreCategoryId)
	if err != nil {
		tx.Rollback()
		return err
	}

	num, err = ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if num < 0 {
		tx.Rollback()
		return errors.New("操作错误,评分类目删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}

	return nil
}

func (r *giveScoreCategoryRepository) EditGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempgiveScoreCategory := new(model.GiveScoreCategory)
	err = tx.QueryRow("SELECT giveScoreCategoryImgFileId FROM giveScoreCategory WHERE giveScoreCategoryId = ?", giveScoreCategory.GiveScoreCategoryId).Scan(&tempgiveScoreCategory.GiveScoreCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	fileId, err := common.EditCover(tx, files, "用于评分类目封面", tempgiveScoreCategory.GiveScoreCategoryImgFileId, "Resources/Img/", tempgiveScoreCategory.IsDelImg)
	if err != nil {
		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("UPDATE giveScoreCategory SET giveScoreCategoryName = ?, giveScoreCategoryDesc = ?, giveScoreCategoryParentId = ?, schoolId = ?, giveScoreCategoryImgFileId = ? WHERE giveScoreCategoryId = ?", giveScoreCategory.GiveScoreCategoryName, giveScoreCategory.GiveScoreCategoryDesc, giveScoreCategory.GiveScoreCategoryParentId, giveScoreCategory.SchoolId, fileId, giveScoreCategory.GiveScoreCategoryId)
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
		return errors.New("操作错误， 评分类目插入错误")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}

	return nil
}

// TODO未完成
func (r *giveScoreCategoryRepository) QueryGiveScoreCategory(schoolId int64) ([]*model.GiveScoreCategoryView, error) {

	giveScoreCategoryViews := make([]*model.GiveScoreCategoryView, 0)

	query := "	SELECT gs.giveScoreCategoryId, gs.giveScoreCategoryName, gs.giveScoreCategoryDesc, gs.giveScoreCategoryParentId, gs.schoolId, gs.giveScoreCategoryImgFileId, COALESCE(f.FilePath, '') 'giveScoreCategoryImgFilePath'	" +
		"	FROM givescorecategory gs LEFT JOIN fileinfo f ON gs.giveScoreCategoryImgFileId = f.fileInfoId WHERE gs.giveScoreCategoryParentId = 0 AND schoolId = ?"

	rows, err := r.db.Query(query, schoolId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		giveScoreCategoryView := new(model.GiveScoreCategoryView)

		err = rows.Scan(&giveScoreCategoryView.GiveScoreCategoryId, &giveScoreCategoryView.GiveScoreCategoryName, &giveScoreCategoryView.GiveScoreCategoryDesc,
			&giveScoreCategoryView.GiveScoreCategoryParentId, &giveScoreCategoryView.SchoolId, &giveScoreCategoryView.GiveScoreCategoryImgFileId,
			&giveScoreCategoryView.FilePath)

		if err != nil {
			return nil, err
		}

		giveScoreCategoryViews = append(giveScoreCategoryViews, giveScoreCategoryView)
	}

	return giveScoreCategoryViews, nil
}

func (r *giveScoreCategoryRepository) QueryGiveScoreCategoryTree(schoolId int64) ([]*model.GiveScoreCategoryTreeList, error) {
	querytree := "	WITH RECURSIVE cte AS (	" +
		"	SELECT giveScoreCategoryId, giveScoreCategoryName, giveScoreCategoryDesc, giveScoreCategoryParentId, schoolId, giveScoreCategoryImgFileId, COALESCE(f1.filePath,'') 'giveScoreCategoryImgFilePath', 0 AS level	" +
		"	FROM givescorecategory gs1 LEFT JOIN fileinfo f1 on gs1.giveScoreCategoryImgFileId = f1.fileInfoId	" +
		"	WHERE gs1.giveScoreCategoryParentId = 0	" +

		"	UNION ALL	" +
		"	SELECT gs2.giveScoreCategoryId, gs2.giveScoreCategoryName, gs2.giveScoreCategoryDesc, gs2.giveScoreCategoryParentId, gs2.schoolId, gs2.giveScoreCategoryImgFileId, COALESCE(f2.filePath,'') 'giveScoreCategoryImgFilePath', cte.level + 1	" +
		"	FROM givescorecategory gs2	" +
		"	JOIN cte ON gs2.giveScoreCategoryParentId = cte.giveScoreCategoryId	" +
		"	LEFT JOIN fileinfo f2  ON gs2.giveScoreCategoryImgFileId = f2.fileInfoId )	" +
		"	SELECT * FROM cte where schoolId = ?"

	rows, err := r.db.Query(querytree, schoolId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nodes []*model.GiveScoreCategoryTreeList

	for rows.Next() {
		tempmodel := new(model.GiveScoreCategoryTreeList)

		err := rows.Scan(&tempmodel.GiveScoreCategoryId, &tempmodel.GiveScoreCategoryName, &tempmodel.GiveScoreCategoryDesc, &tempmodel.GiveScoreCategoryParentId, &tempmodel.SchoolId, &tempmodel.GiveScoreCategoryImgFileId, &tempmodel.FilePath, &tempmodel.Level)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, tempmodel)
	}

	treeLists := GiveScoreCategoryTree(nodes, 0)

	return treeLists, nil
}

func GiveScoreCategoryTree(node []*model.GiveScoreCategoryTreeList, giveScoreCategoryParentId int) []*model.GiveScoreCategoryTreeList {
	res := make([]*model.GiveScoreCategoryTreeList, 0)

	for _, v := range node {
		if v.GiveScoreCategoryParentId == giveScoreCategoryParentId {
			v.Children = GiveScoreCategoryTree(node, v.GiveScoreCategoryId)
			res = append(res, v)
		}
	}

	return res
}

func (r *giveScoreCategoryRepository) QueryGiveScoreCategoryParentNodeByParentId(giveScoreCategoryParentId int64) (*model.GiveScoreCategoryView, error) {
	if giveScoreCategoryParentId == 0 {
		return nil, errors.New("该节点已为根节点")
	}

	tempmodel := new(model.GiveScoreCategoryView)
	r.db.QueryRow("	SELECT gs.giveScoreCategoryId, gs.giveScoreCategoryName, gs.giveScoreCategoryDesc, gs.giveScoreCategoryParentId, gs.schoolId, "+
		"	gs.giveScoreCategoryImgFileId, COALESCE(f.filePath,'') 'giveScoreCategoryImgFilePath' "+
		"	FROM givescorecategory gs     LEFT JOIN fileinfo f on gs.giveScoreCategoryImgFileId = f.fileInfoId     "+
		"	WHERE gs.giveScoreCategoryId = ?  ", giveScoreCategoryParentId).
		Scan(&tempmodel.GiveScoreCategoryId, &tempmodel.GiveScoreCategoryName, &tempmodel.GiveScoreCategoryDesc,
			&tempmodel.GiveScoreCategoryParentId, &tempmodel.SchoolId, &tempmodel.GiveScoreCategoryImgFileId, &tempmodel.FilePath)

	if tempmodel.GiveScoreCategoryId == 0 {
		return nil, errors.New("未查询到数据")
	}

	return tempmodel, nil
}

func (r *giveScoreCategoryRepository) QueryGiveScoreCategoryChildNodesById(giveScoreCategoryId int64, schoolId int64) ([]*model.GiveScoreCategoryView, error) {

	temparr := make([]*model.GiveScoreCategoryView, 0)

	rows, err := r.db.Query("	SELECT gs.giveScoreCategoryId, gs.giveScoreCategoryName, gs.giveScoreCategoryDesc, gs.giveScoreCategoryParentId, gs.schoolId, "+
		"	gs.giveScoreCategoryImgFileId, COALESCE(f.filePath,'') 'giveScoreCategoryImgFilePath' "+
		"	FROM givescorecategory gs    LEFT JOIN fileinfo f ON gs.giveScoreCategoryImgFileId = f.fileInfoId     "+
		"	WHERE gs.giveScoreCategoryParentId = ? AND ( schoolId = ? or schoolId = 0) ", giveScoreCategoryId, schoolId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		tempmodel := new(model.GiveScoreCategoryView)

		err = rows.Scan(&tempmodel.GiveScoreCategoryId, &tempmodel.GiveScoreCategoryName, &tempmodel.GiveScoreCategoryDesc,
			&tempmodel.GiveScoreCategoryParentId, &tempmodel.SchoolId, &tempmodel.GiveScoreCategoryImgFileId, &tempmodel.FilePath)

		if err != nil {
			return nil, err
		}

		temparr = append(temparr, tempmodel)
	}

	return temparr, nil
}
