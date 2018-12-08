package orm

import (
	"database/sql"
	"github.com/envzo/zorm/db"
)

type PodUser struct {
	Id          int64
	Nickname    string
	Password    string
	MobilePhone string
	CreateDt    int64
	UpdateDt    int64
}

func NewPodUser() *PodUser {
	return &PodUser{}
}

type _PodUserMgr struct{}

var PodUserMgr = &_PodUserMgr{}

func (mgr *_PodUserMgr) IsNicknameExists(nickname string) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where nickname=?`,
		nickname)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) FindOneByNickname(nickname string) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, mobile_phone, create_dt, update_dt from pod.pod_user where nickname=?`,
		nickname)

	var id sql.NullInt64
	var nickname_1 sql.NullString
	var password sql.NullString
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var updateDt sql.NullInt64

	if err := row.Scan(&id, &nickname_1, &password, &mobilePhone, &createDt, &updateDt); err != nil {
		return nil, err
	}

	d := PodUser{
		Id:          id.Int64,
		Nickname:    nickname_1.String,
		Password:    password.String,
		MobilePhone: mobilePhone.String,
		CreateDt:    createDt.Int64,
		UpdateDt:    updateDt.Int64,
	}

	return &d, nil
}

func (mgr *_PodUserMgr) IsMobilePhoneExists(mobilePhone string) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where mobile_phone=?`,
		mobilePhone)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) FindOneByMobilePhone(mobilePhone string) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, mobile_phone, create_dt, update_dt from pod.pod_user where mobile_phone=?`,
		mobilePhone)

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var mobilePhone_1 sql.NullString
	var createDt sql.NullInt64
	var updateDt sql.NullInt64

	if err := row.Scan(&id, &nickname, &password, &mobilePhone_1, &createDt, &updateDt); err != nil {
		return nil, err
	}

	d := PodUser{
		Id:          id.Int64,
		Nickname:    nickname.String,
		Password:    password.String,
		MobilePhone: mobilePhone_1.String,
		CreateDt:    createDt.Int64,
		UpdateDt:    updateDt.Int64,
	}

	return &d, nil
}

func (mgr *_PodUserMgr) Create(d *PodUser) error {
	r, err := db.DB().Exec(`insert into pod.pod_user (nickname, password, mobile_phone, create_dt, update_dt) value (?,?,?,?,?)`, d.Nickname, d.Password, d.MobilePhone, d.CreateDt, d.UpdateDt)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = id
	return nil
}

func (mgr *_PodUserMgr) Update(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`update pod.pod_user set nickname=?, password=?, mobile_phone=?, create_dt=?, update_dt=? where id=?`, d.Id)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}
