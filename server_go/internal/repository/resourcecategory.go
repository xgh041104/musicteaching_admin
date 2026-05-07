package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type ResourceCategoryRepository interface {
	AddResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error
	EditResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error
	DelResourceCategory(resourceCategory model.ResourcesCategory) error
	QueryResourceCategoryTree(schoolId int) ([]*model.ResourcesCategoryTreeList, error)
	QueryResourceCategoryParentNodeByParentId(resourceCategoryParentId int) (*model.ResourcesCategoryView, error)
	QueryResourceCategoryChildNodesById(resourceCategoryId int, schoolId int) ([]*model.ResourcesCategoryView, error)
}
type resourceCategoryRepository struct {
	*BaseRepository
}

func NewResourceCategoryRepository(repository *BaseRepository) ResourceCategoryRepository {
	return &resourceCategoryRepository{
		BaseRepository: repository,
	}
}

func (r *resourceCategoryRepository) AddResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var fileId int64                   //用户接受新增文件表之后fileid的存储
	if len(files.File["files"]) == 1 { // 这里只获取有一个文件的时候  如果files传了多张图片   这里则用for循环了  限定只能传img格式
		tempfile := files.File["files"][0]

		fileinfo := new(model.FileInfo)
		fileinfo.FileType = path.Ext(tempfile.Filename)
		fileinfo.FileName = uuid.GenUUID()
		fileinfo.FileUseTo = "用于资源类目封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {

			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec(" insert into resourcecategory(resourceCategoryName,resourceCategoryDesc,resourceCategoryParentId,resourceCategoryImgFileId,schoolId) values (?,?,?,?,?)", resourceCategory.ResourceCategoryName, resourceCategory.ResourceCategoryDesc, resourceCategory.ResourceCategoryParentId, fileId, resourceCategory.SchoolId)
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
		return errors.New("操作错误,资源类目插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

func (r *resourceCategoryRepository) EditResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempResourcesCategory := new(model.ResourcesCategory)
	err = tx.QueryRow("select resourceCategoryImgFileId from resourcecategory where resourceCategoryId=? ", resourceCategory.ResourceCategoryId).Scan(&tempResourcesCategory.ResourceCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	//用于接受新增文件表之后fileid的存储

	fileId, err := common.EditCover(tx, files, "用于资源类目封面", tempResourcesCategory.ResourceCategoryImgFileId, "Resources/Img/", resourceCategory.IsDelImg)

	if err != nil {

		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("update resourcecategory set   resourceCategoryName=?,resourceCategoryDesc=?,schoolId=?,resourceCategoryImgFileId=?    where resourceCategoryId=?",
		resourceCategory.ResourceCategoryName, resourceCategory.ResourceCategoryDesc, resourceCategory.SchoolId, fileId, resourceCategory.ResourceCategoryId)
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
		return errors.New("操作错误,资源类目插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}
func (r *resourceCategoryRepository) DelResourceCategory(resourceCategory model.ResourcesCategory) error {

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	resourceCategorynum := 0

	tx.QueryRow("  select count(1) from  resourcecategory   where resourceCategoryParentId=?", resourceCategory.ResourceCategoryId).Scan(&resourceCategorynum)
	if resourceCategorynum > 0 {
		tx.Rollback()
		return errors.New("该类目不能删除，请先删除该类目的子目录")
	}
	resourcenum := 0
	tx.QueryRow("  select count(1) from  resource   where resourceCategoryId=?", resourceCategory.ResourceCategoryId).Scan(&resourcenum)
	if resourcenum > 0 {
		tx.Rollback()
		return errors.New("该类目不能删除，请先删除该类目的的资源")
	}

	tempResourcesCategory := new(model.ResourcesCategory)
	err = tx.QueryRow("select resourceCategoryImgFileId from resourcecategory where resourceCategoryId=? ", resourceCategory.ResourceCategoryId).Scan(&tempResourcesCategory.ResourceCategoryImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempResourcesCategory.ResourceCategoryImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempResourcesCategory.ResourceCategoryImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  resourcecategory    where resourceCategoryId=?", resourceCategory.ResourceCategoryId)
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
		return errors.New("操作错误,资源类目删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

func (r *resourceCategoryRepository) QueryResourceCategoryTree(schoolId int) ([]*model.ResourcesCategoryTreeList, error) {

	querytree := "	WITH RECURSIVE cte AS (   " +
		"	SELECT resourceCategoryId,resourceCategoryName,resourceCategoryDesc,resourceCategoryParentId,schoolId,resourceCategoryImgFileId,COALESCE(b.filePath,'') 'filePath', 0 AS level    " +
		"	FROM resourcecategory  a    " +
		"	 LEFT JOIN fileinfo b on a.resourceCategoryImgFileId=b.fileInfoId   " +
		"	WHERE resourceCategoryParentId=0     " +

		"	UNION ALL     " +
		"	SELECT c.resourceCategoryId,c.resourceCategoryName,c.resourceCategoryDesc,c.resourceCategoryParentId,c.schoolId,c.resourceCategoryImgFileId,COALESCE(d.filePath,'') 'filePath', cte.level+1    " +
		"	FROM resourcecategory c      " +
		"	JOIN cte ON c.resourceCategoryParentId = cte.resourceCategoryId     " +
		"				 LEFT JOIN fileinfo d on c.resourceCategoryImgFileId=d.fileInfoId  " +
		"	)  SELECT * FROM cte where (schoolId=?  or schoolId=0) "

	rows, err := r.db.Query(querytree, schoolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*model.ResourcesCategoryTreeList
	for rows.Next() {

		tempmodel := new(model.ResourcesCategoryTreeList)
		err = rows.Scan(&tempmodel.ResourceCategoryId, &tempmodel.ResourceCategoryName, &tempmodel.ResourceCategoryDesc, &tempmodel.ResourceCategoryParentId, &tempmodel.SchoolId, &tempmodel.ResourceCategoryImgFileId, &tempmodel.FilePath, &tempmodel.Level)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, tempmodel)
	}

	treeLists := ResourceCategoryTree(nodes, 0)

	return treeLists, nil

}

// 资源类目树形递归
func ResourceCategoryTree(node []*model.ResourcesCategoryTreeList, resourcesCategoryParentId int) []*model.ResourcesCategoryTreeList {
	res := make([]*model.ResourcesCategoryTreeList, 0)
	for _, v := range node {
		if v.ResourceCategoryParentId == resourcesCategoryParentId {
			v.Children = ResourceCategoryTree(node, v.ResourceCategoryId)
			res = append(res, v)
		}
	}
	return res
}

//根据当前目录id查询父节点数据

func (r *resourceCategoryRepository) QueryResourceCategoryParentNodeByParentId(resourceCategoryParentId int) (*model.ResourcesCategoryView, error) {

	if resourceCategoryParentId == 0 {
		return nil, errors.New("该节点已为根节点")
	}

	tempmodel := new(model.ResourcesCategoryView)

	r.db.QueryRow("		SELECT resourceCategoryId,resourceCategoryName,resourceCategoryDesc,resourceCategoryParentId,schoolId, "+
		" resourceCategoryImgFileId,COALESCE(b.filePath,'') 'filePath' "+
		"   FROM resourcecategory  a   		 LEFT JOIN fileinfo b on a.resourceCategoryImgFileId=b.fileInfoId     "+
		"	 WHERE resourceCategoryId=?  ", resourceCategoryParentId).
		Scan(&tempmodel.ResourceCategoryId, &tempmodel.ResourceCategoryName, &tempmodel.ResourceCategoryDesc, &tempmodel.ResourceCategoryParentId, &tempmodel.SchoolId,
			&tempmodel.ResourceCategoryImgFileId, &tempmodel.FilePath)

	if tempmodel.ResourceCategoryId == 0 {
		return nil, errors.New("未查询到数据")
	}
	return tempmodel, nil
}

// 根据当前id查询子节点数据
func (r *resourceCategoryRepository) QueryResourceCategoryChildNodesById(resourceCategoryId int, schoolId int) ([]*model.ResourcesCategoryView, error) {

	// if resourceCategoryId == 0 {
	// 	return nil, errors.New("该节点已为根节点")
	// }

	temparr := make([]*model.ResourcesCategoryView, 0)

	rows, err := r.db.Query("		SELECT resourceCategoryId,resourceCategoryName,resourceCategoryDesc,resourceCategoryParentId,schoolId, "+
		" resourceCategoryImgFileId,COALESCE(b.filePath,'') 'filePath' "+
		"   FROM resourcecategory  a   		 LEFT JOIN fileinfo b on a.resourceCategoryImgFileId=b.fileInfoId     "+
		"	 WHERE resourceCategoryParentId=?  and ( schoolId=? or  schoolId=0) ", resourceCategoryId, schoolId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tempmodel := new(model.ResourcesCategoryView)

		err = rows.Scan(&tempmodel.ResourceCategoryId, &tempmodel.ResourceCategoryName, &tempmodel.ResourceCategoryDesc, &tempmodel.ResourceCategoryParentId, &tempmodel.SchoolId,
			&tempmodel.ResourceCategoryImgFileId, &tempmodel.FilePath)

		if err != nil {
			return nil, err
		}
		temparr = append(temparr, tempmodel)
	}

	return temparr, nil
}
