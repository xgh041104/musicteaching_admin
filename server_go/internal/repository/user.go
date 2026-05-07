package repository

import (
	"errors"
	"server_go/internal/model"
	"server_go/pkg/helper/md5"
	tostring "server_go/pkg/helper/toString"
	"server_go/pkg/verify"
)

type UserRepository interface {
	LoginUser(userparam model.LoginReq) (*model.LoginUser, error)
	EditUserPwd(userparam model.LoginUser) error
	GetUserBySchoolId(schoolId int64) ([]*model.CommonUser, error)
	AddCommonUser(commonUser model.CommonUser) error
	UpdateCommonUser(commonUser model.CommonUser) error
	DelCommonUser(commonUserId int64) error
}
type userRepository struct {
	*BaseRepository
}

func NewUserRepository(repository *BaseRepository) UserRepository {
	return &userRepository{
		BaseRepository: repository,
	}
}

func (r *userRepository) LoginUser(userparam model.LoginReq) (*model.LoginUser, error) {

	newpwd := md5.Md5(userparam.UserPwd)

	loginUser := new(model.LoginUser)
	var err error
	sqlloggin := ""
	if userparam.UserType == 2 && userparam.AccessToken == "" { // 普通用户
		sqlloggin = "select commonUserId,commonUserTrueName,commonUserAccount,commonUserPwd,schoolId,commonUserType from  commonuser   where commonUserAccount=? and commonUserPwd=? and commonUserType=?  "
	} else if userparam.UserType == 0 || userparam.UserType == 1 { // 管理员和超级管理员
		sqlloggin = "select adminId,adminTrueName,adminAccount,adminPwd,schoolId,adminType from  admin  where adminAccount=? and adminPwd=?  and adminType=? "
	} else { // 学生
		sqlloggin = "select s.studentId, studentName, studentAccount, schoolId, className, userType " +
			" from studentuser s left join student_teacher st on s.studentId = st.studentId " +
			" where st.teacherId=? and studentAccount=? and studentPwd=?"
	}

	if userparam.AccessToken == "" {
		err = r.db.QueryRow(sqlloggin, userparam.UserAccount, newpwd, userparam.UserType).
			Scan(&loginUser.UserId, &loginUser.UserTrueName, &loginUser.UserAccount, &loginUser.UserPwd, &loginUser.SchoolId, &loginUser.UserType)
	} else {
		err = r.db.QueryRow(sqlloggin, userparam.TeacherId, userparam.UserAccount, newpwd).
			Scan(&loginUser.UserId, &loginUser.UserTrueName, &loginUser.UserAccount, &loginUser.SchoolId, &loginUser.ClassName, &loginUser.UserType)
	}

	if loginUser.UserId == 0 {
		return nil, nil // 表示没有查询到数据 登陆密码错误
	}
	if err != nil {
		return nil, err
	}

	if userparam.UserType == 2 && userparam.AccessToken == "" { //验证普通用户
		dataMap := make(map[string]string)
		dataMap["onlyMark"] = verify.SystemName + loginUser.UserAccount
		dataMap["serialNum"] = userparam.SerialNum
		dataMap["hostMAC"] = userparam.HostMAC

		loginUser.AccessToken, err = verify.SecretVerifyUser(dataMap)
		if err != nil {
			return nil, err
		}
	} else if userparam.AccessToken != "" { //验证二级用户
		var count int
		err = r.db.QueryRow("SELECT COUNT(*) FROM student_teacher WHERE teacherId = ? AND studentId = ?", loginUser.UserId, userparam.TeacherId).Scan(&count)
		if err != nil {
			return nil, errors.New("该学生不属于该老师，无法登录")
		}

		var account string
		err = r.db.QueryRow("SELECT commonUserAccount FROM commonuser WHERE commonUserId = ?", userparam.TeacherId).Scan(&account)
		if err != nil {
			return nil, errors.New("系统错误")
		}

		err = verify.SecretVerifyStudent(verify.SystemName+account, userparam.AccessToken)
		if err != nil {
			return nil, err
		}
	}

	return loginUser, nil
}

func (r *userRepository) EditUserPwd(userparam model.LoginUser) error {

	newpwd := md5.Md5(userparam.UserPwd)

	sql := ""
	if userparam.UserType == 2 { // 普通用户
		sql = "UPDATE commonuser set commonUserPwd = ?  where commonUserId = ? and commonUserType = ?  "
	} else if userparam.UserType == 1 || userparam.UserType == 0 { // 管理员和超级管理员
		sql = "UPDATE admin set adminPwd = ?  where adminId = ? and adminType = ?  "
	}

	ret, err := r.db.Exec(sql, newpwd, userparam.UserId, userparam.UserType)
	if err != nil {
		return err
	}

	updatenum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if updatenum == 0 {
		return errors.New("操作无效，数据无变化")
	}

	return nil
}

func (r *userRepository) GetUserBySchoolId(schoolId int64) ([]*model.CommonUser, error) {
	var commonUsers []*model.CommonUser

	rows, err := r.db.Query("select commonUserId,commonUserTrueName,commonUserAccount,commonUserPwd,schoolId,commonUserType from  commonuser  where schoolId = ?", schoolId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var commonUser model.CommonUser
		err := rows.Scan(&commonUser.CommonUserId, &commonUser.CommonUserTrueName, &commonUser.CommonUserAccount, &commonUser.CommonUserPwd, &commonUser.SchoolId, &commonUser.CommonUserType)
		if err != nil {
			return nil, err
		}
		commonUsers = append(commonUsers, &commonUser)
	}
	return commonUsers, nil
}

func (r *userRepository) AddCommonUser(commonUser model.CommonUser) error {
	var count int
	err := r.db.QueryRow("select count(*) from commonuser where commonUserAccount = ?", commonUser.CommonUserAccount).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 如果已经存在相同的 userAccount 记录，则返回错误信息
		return errors.New("用户名重复")
	}

	dataMap := make(map[string]string)
	dataMap["onlyMark"] = verify.SystemName + commonUser.CommonUserAccount
	dataMap["userName"] = commonUser.CommonUserTrueName
	dataMap["userAccount"] = commonUser.CommonUserAccount
	dataMap["userTypeId"] = tostring.Strval(1)

	err = verify.SecretAddUser(dataMap)
	if err != nil {
		return err
	}

	//进行md5加密
	newpwd := md5.Md5(commonUser.CommonUserPwd)

	//执行插入
	_, err = r.db.Exec("insert into commonuser(commonUserTrueName, commonUserAccount, commonUserPwd, schoolId, commonUserType) values(?,?,?,?,?) ", commonUser.CommonUserTrueName, commonUser.CommonUserAccount, newpwd, commonUser.SchoolId, commonUser.CommonUserType)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateCommonUser(commonUser model.CommonUser) error {
	// var count int
	// err := r.db.QueryRow("select count(*) from commonuser where commonUserAccount = ?", commonUser.CommonUserAccount).Scan(&count)
	// if err != nil {
	// 	return err
	// }

	// , commonUserAccount = ?, schoolId = ?, commonUserType = ?

	// , commonUser.CommonUserAccount, commonUser.SchoolId, commonUser.CommonUserType
	ret, err := r.db.Exec(" update commonuser set commonUserTrueName = ? where commonUserId = ?",
		commonUser.CommonUserTrueName, commonUser.CommonUserId)
	if err != nil {
		return err
	}

	updatenum, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if updatenum == 0 {
		return errors.New("操作无效，数据无变化")
	}
	return nil
}

func (r *userRepository) DelCommonUser(commonUserId int64) error {
	var account string
	err := r.db.QueryRow("select commonUserAccount from commonuser where commonUserId = ? ", commonUserId).Scan(&account)
	if err != nil {
		return errors.New("系统错误")
	}

	dataMap := make(map[string]string)
	dataMap["onlyMark"] = verify.SystemName + account

	err = verify.SecretDelUser(dataMap)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("系统错误")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ret, err := tx.Exec("delete from commonuser where commonUserId = ? ", commonUserId)
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

	_, err = tx.Exec("delete from student_teacher where teacherId = ? ", commonUserId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("系统错误")
	}

	return nil
}
