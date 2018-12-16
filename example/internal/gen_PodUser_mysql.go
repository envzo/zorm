// usage:
// FindByXXX will not return sql.ErrNoRows, so it's caller's ability to check error

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/envzo/zorm/db"
)

var _ = errors.New
var _ = fmt.Printf
var _ = sql.ErrNoRows

type PodUser struct {
	Id          int64
	Nickname    string
	Password    string
	Age         int32
	MobilePhone string
	CreateDt    int64
	IsBlocked   bool
	UpdateDt    int64

	baby bool
}

func NewPodUser() *PodUser {
	return &PodUser{baby: true}
}

type _PodUserMgr struct{}

var PodUserMgr = &_PodUserMgr{}

func (mgr *_PodUserMgr) IsNicknameExists(nickname string) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where nickname = ?`,
		nickname)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) UniFindByNickname(nickname string) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where nickname = ?`,
		nickname)

	var id sql.NullInt64
	var nickname_1 sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	if err := row.Scan(&id, &nickname_1, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt); err != nil {
		return nil, err
	}

	d := PodUser{
		Id:          id.Int64,
		Nickname:    nickname_1.String,
		Password:    password.String,
		Age:         int32(age.Int64),
		MobilePhone: mobilePhone.String,
		CreateDt:    createDt.Int64,
		IsBlocked:   isBlocked.Bool,
		UpdateDt:    updateDt.Int64,
	}

	return &d, nil
}

func (mgr *_PodUserMgr) IsMobilePhoneExists(mobilePhone string) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where mobile_phone = ?`,
		mobilePhone)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) UniFindByMobilePhone(mobilePhone string) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where mobile_phone = ?`,
		mobilePhone)

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone_1 sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	if err := row.Scan(&id, &nickname, &password, &age, &mobilePhone_1, &createDt, &isBlocked, &updateDt); err != nil {
		return nil, err
	}

	d := PodUser{
		Id:          id.Int64,
		Nickname:    nickname.String,
		Password:    password.String,
		Age:         int32(age.Int64),
		MobilePhone: mobilePhone_1.String,
		CreateDt:    createDt.Int64,
		IsBlocked:   isBlocked.Bool,
		UpdateDt:    updateDt.Int64,
	}

	return &d, nil
}

func (mgr *_PodUserMgr) FindByCreateDt(createDt int64, order []string, offset, limit int64) ([]*PodUser, error) {
	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where create_dt = ?`
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		query += o[1:]
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	rows, err := db.DB().Query(query, createDt)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt_1 sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt_1, &isBlocked, &updateDt); err != nil {
			return nil, err
		}

		d := PodUser{
			Id:          id.Int64,
			Nickname:    nickname.String,
			Password:    password.String,
			Age:         int32(age.Int64),
			MobilePhone: mobilePhone.String,
			CreateDt:    createDt_1.Int64,
			IsBlocked:   isBlocked.Bool,
			UpdateDt:    updateDt.Int64,
		}
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) CountByCreateDt(createDt int64) (int64, error) {
	query := `select count(1) from pod.pod_user where create_dt = ?`
	row := db.DB().QueryRow(query, createDt)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByUpdateDt(updateDt int64, order []string, offset, limit int64) ([]*PodUser, error) {
	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where update_dt = ?`
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		query += o[1:]
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	rows, err := db.DB().Query(query, updateDt)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt_1 sql.NullInt64

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt_1); err != nil {
			return nil, err
		}

		d := PodUser{
			Id:          id.Int64,
			Nickname:    nickname.String,
			Password:    password.String,
			Age:         int32(age.Int64),
			MobilePhone: mobilePhone.String,
			CreateDt:    createDt.Int64,
			IsBlocked:   isBlocked.Bool,
			UpdateDt:    updateDt_1.Int64,
		}
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) CountByUpdateDt(updateDt int64) (int64, error) {
	query := `select count(1) from pod.pod_user where update_dt = ?`
	row := db.DB().QueryRow(query, updateDt)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByJoin(t string, on, where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user join t on `
	for i, v := range on {
		if i > 0 {
			query += " and "
		}
		query += v.S
		if v.P != nil {
			params = append(params, v.P)
		}
	}
	for i, v := range where {
		if i == 0 {
			query += " where "
		} else if i != len(where)-1 {
			query += " and "
		}
		query += v.S
		if v.P != nil {
			params = append(params, v.P)
		}
	}
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		query += o[1:]
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	rows, err := db.DB().Query(query, params...)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt); err != nil {
			return nil, err
		}

		d := PodUser{
			Id:          id.Int64,
			Nickname:    nickname.String,
			Password:    password.String,
			Age:         int32(age.Int64),
			MobilePhone: mobilePhone.String,
			CreateDt:    createDt.Int64,
			IsBlocked:   isBlocked.Bool,
			UpdateDt:    updateDt.Int64,
		}
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) FindByCond(where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where `
	for i, v := range where {
		if i > 0 {
			query += " and "
		}
		query += v.S
		if v.P != nil {
			params = append(params, v.P)
		}
	}
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else if i != len(order)-1 {
			query += ", "
		}
		query += o
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	rows, err := db.DB().Query(query, params...)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt); err != nil {
			return nil, err
		}

		d := PodUser{
			Id:          id.Int64,
			Nickname:    nickname.String,
			Password:    password.String,
			Age:         int32(age.Int64),
			MobilePhone: mobilePhone.String,
			CreateDt:    createDt.Int64,
			IsBlocked:   isBlocked.Bool,
			UpdateDt:    updateDt.Int64,
		}
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) Create(d *PodUser) error {
	r, err := db.DB().Exec(`insert into pod.pod_user (nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt) value (?,?,?,?,?,?,?)`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.CreateDt, d.IsBlocked, d.UpdateDt)
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

func (mgr *_PodUserMgr) Upsert(d *PodUser) error {
	if d.baby {
		return mgr.Create(d)
	}
	_, err := mgr.Update(d)
	return err
}

func (mgr *_PodUserMgr) CountByRule(rules ...db.Rule) (int64, error) {
	var p []interface{}
	query := `select count(1) from pod.pod_user where `
	for i, rule := range rules {
		if i > 0 {
			query += " and "
		}
		query += rule.S
		if rule.P != nil {
			p = append(p, rule.P)
		}

	}

	row := db.DB().QueryRow(query, p...)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	return c.Int64, nil
}

func (mgr *_PodUserMgr) RmByRule(rules ...db.Rule) (int64, error) {
	query := "delete from pod.pod_user where "
	var p []interface{}
	for i, r := range rules {
		if i > 0 {
			query += " and "
		}
		query += r.S
		p = append(p, r.P)
	}
	r, err := db.DB().Exec(query, p...)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (mgr *_PodUserMgr) UniFindByPK(id int64) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt from pod.pod_user where id = ?`, id)

	var id_1 sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64

	if err := row.Scan(&id_1, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt); err != nil {
		return nil, err
	}

	d := PodUser{
		Id:          id_1.Int64,
		Nickname:    nickname.String,
		Password:    password.String,
		Age:         int32(age.Int64),
		MobilePhone: mobilePhone.String,
		CreateDt:    createDt.Int64,
		IsBlocked:   isBlocked.Bool,
		UpdateDt:    updateDt.Int64,
	}

	return &d, nil
}

func (mgr *_PodUserMgr) Update(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`update pod.pod_user set nickname = ?, password = ?, age = ?, mobile_phone = ?, create_dt = ?, is_blocked = ?, update_dt = ? where id = ?`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.CreateDt, d.IsBlocked, d.UpdateDt, d.Id)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (mgr *_PodUserMgr) RmByPK(pk int64) (int64, error) {
	query := "delete from pod.pod_user where id = ?"
	r, err := db.DB().Exec(query, pk)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}
