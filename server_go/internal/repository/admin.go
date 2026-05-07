package repository

import (
	"errors"
	"server_go/internal/model"
	"server_go/pkg/helper/md5"
)

type AdminRepository interface {
	GetSchoolAdmin() ([]*model.Admin, error)
	AddSchoolAdmin(admin model.Admin) error
	DelSchoolAdmin(adminId int64) error
	UpdateSchoolAdmin(admin model.Admin) error
}

type adminRepository struct {
	*BaseRepository
}

func NewAdminRepository(repository *BaseRepository) AdminRepository {
	return &adminRepository{
		BaseRepository: repository,
	}
}

func (r *adminRepository) AddSchoolAdmin(admin model.Admin) error {

	var count int
	err := r.db.QueryRow("select count(*) from admin where adminAccount = ?", admin.AdminAccount).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 如果已经存在相同的 userAccount 记录，则返回错误信息
		return errors.New("用户名重复")
	}

	//进行md5加密
	newpwd := md5.Md5(admin.AdminPwd)

	//执行插入
	_, err = r.db.Exec("insert into admin(adminTrueName, adminAccount,adminPwd, schoolId, adminType) values(?,?,?,?,?)", admin.AdminTrueName, admin.AdminAccount, newpwd, admin.SchoolId, admin.AdminType)
	if err != nil {
		return err
	}
	return nil
	// r.db.
	// err := r.db.Debug().Select("Username", "Pwd").Create(&admin).Error
	// if err != nil {
	// 	return err
	// }
}

func (r *adminRepository) GetSchoolAdmin() ([]*model.Admin, error) {
	// 存储多个 Admin 结构体指针
	var admins []*model.Admin

	rows, err := r.db.Query("select a.adminId, a.adminTrueName, a.adminAccount, a.adminPwd, a.schoolId, a.adminType, COALESCE(s.schoolName,'') 'schoolName' from admin a " +
		"left join school s on a.schoolId = s.schoolId WHERE   a.adminType = 1;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var admin model.Admin
		err := rows.Scan(&admin.AdminId, &admin.AdminTrueName, &admin.AdminAccount, &admin.AdminPwd, &admin.SchoolId, &admin.AdminType, &admin.SchoolName)
		if err != nil {
			return nil, err
		}
		admins = append(admins, &admin)
	}

	if err := rows.Err(); err != nil {
		// 处理迭代器错误
		return nil, err
	}

	// 返回存储了从数据库中检索到的多个 Admin 结构体的指针的切片
	return admins, nil
}

func (r *adminRepository) DelSchoolAdmin(adminId int64) error {

	schoolnum := 0
	r.db.QueryRow("select count(1) from  admin where schoolId in  (select  schoolId from  admin where adminId = ?)").
		Scan(&schoolnum)

	if schoolnum <= 1 {
		return errors.New("这个管理员是该学校最后一个管理员,不能删除")
	}
	ret, err := r.db.Exec("delete from admin where adminId = ? and adminType = 1", adminId)
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

func (r *adminRepository) UpdateSchoolAdmin(admin model.Admin) error {

	var count int
	err := r.db.QueryRow("select count(*) from admin where adminAccount = ? and adminId != ?", admin.AdminAccount, admin.AdminId).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 如果已经存在相同的 userAccount 记录，则返回错误信息
		return errors.New("用户名重复")
	}

	ret, err := r.db.Exec(" update admin set adminTrueName = ?, adminAccount = ?, schoolId = ?, adminType = ? where adminId = ?",
		admin.AdminTrueName, admin.AdminAccount, admin.SchoolId, admin.AdminType, admin.AdminId)
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
