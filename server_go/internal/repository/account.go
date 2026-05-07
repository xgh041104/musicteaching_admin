package repository

import (
	"errors"
	"mime/multipart"
	"path"
	"server_go/internal/model"
	"server_go/pkg/common"
	"server_go/pkg/helper/uuid"
)

type AccountRepository interface {
	AddAccount(account model.Account, files *multipart.Form) error
	EditAccount(account model.Account, files *multipart.Form) error
	DelAccount(account model.Account) error
	QueryAllAccount() ([]*model.AccountView, error)
}

type accountRepository struct {
	*BaseRepository
}

func NewAccountRepository(repository *BaseRepository) AccountRepository {
	return &accountRepository{
		BaseRepository: repository,
	}
}

func (r *accountRepository) AddAccount(account model.Account, files *multipart.Form) error {

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
		fileinfo.FileUseTo = "用于公众号封面"
		fileinfo.FilePath = "Resources/Img/" + fileinfo.FileName + path.Ext(tempfile.Filename)
		fileId = common.AddFileInfo(tx, fileinfo, tempfile)

		if fileId == 0 {
			tx.Rollback()
			return errors.New("文件上传错误")
		}

	}

	ret, err := tx.Exec(" insert into account(accountTitle, accountDesc, accountUrl, accountImgFileId) values (?,?,?,?)", account.AccountTitle, account.AccountDesc, account.AccountUrl, fileId)
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
		return errors.New("操作错误,公众号插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil

}

func (r *accountRepository) EditAccount(account model.Account, files *multipart.Form) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempaccount := new(model.Account)
	err = tx.QueryRow("select accountImgFileId from account where accountId=? ", account.AccountId).Scan(&tempaccount.AccountImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	fileId, err := common.EditCover(tx, files, "用于公众号封面", tempaccount.AccountImgFileId, "Resources/Img/", account.IsDelImg)

	if err != nil {

		tx.Rollback()
		return err
	}

	ret, err := tx.Exec("  update account set accountTitle=?, accountDesc=?, accountUrl=?,accountImgFileId=?  where accountId=?", account.AccountTitle, account.AccountDesc, account.AccountUrl, fileId, account.AccountId)
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
		return errors.New("操作错误,公众号插入错误")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

func (r *accountRepository) DelAccount(account model.Account) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	tempaccount := new(model.Account)
	err = tx.QueryRow("select accountImgFileId from account where accountId=? ", account.AccountId).Scan(&tempaccount.AccountImgFileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if tempaccount.AccountImgFileId > 0 {
		fileinfo := new(model.FileInfo)
		fileinfo.FileInfoId = tempaccount.AccountImgFileId
		err = common.DelFileInfo(tx, fileinfo)
		if err != nil {
			tx.Rollback()
			return errors.New("文件删除失败")
		}
	}

	ret, err := tx.Exec("  delete from  account  where accountId=?", account.AccountId)
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
		return errors.New("操作错误,公众号删除失败")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errors.New("事务提交失败")
	}
	return nil
}

func (r *accountRepository) QueryAllAccount() ([]*model.AccountView, error) {

	accountviewarr := make([]*model.AccountView, 0)

	rows, err := r.db.Query("select a.accountId,a.accountTitle,COALESCE(a.accountDesc, '') 'accountDesc', a.accountUrl, a.accountImgFileId,COALESCE(b.fileName , '')'accountImgFileName', COALESCE(b.filePath , '') 'accountImgPath' from account a   left join fileinfo b on a.accountImgFileId=b.fileInfoId ")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		tempaccountview := new(model.AccountView)
		err = rows.Scan(&tempaccountview.AccountId, &tempaccountview.AccountTitle, &tempaccountview.AccountDesc, &tempaccountview.AccountUrl, &tempaccountview.AccountImgFileId, &tempaccountview.AccountImgFileName, &tempaccountview.AccountImgPath)
		if err != nil {
			return nil, err
		}
		accountviewarr = append(accountviewarr, tempaccountview)
	}

	return accountviewarr, nil
}
