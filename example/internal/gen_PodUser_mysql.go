// usage:
// FindByXXX will not return sql.ErrNoRows, so it's caller's ability to check error

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/envzo/zorm/db"
	"github.com/envzo/zorm/util"
)

var _ = errors.New
var _ = fmt.Printf
var _ = strings.Trim
var _ = sql.ErrNoRows
var _ = util.I64
var _ = time.Nanosecond

type PodUser struct {
	//
	Id int64
	// 昵称，非姓名
	Nickname string
	//
	Password string
	//
	Age int32
	// 手机号
	MobilePhone string
	//
	CreateDt int64
	//
	IsBlocked bool
	// 更新时间
	UpdateDt int64
	//
	StatsDt *time.Time
	//
	Dt *time.Time

	// 调用Upsert方法时，baby为true则insert，反之update
	baby bool
}

func NewPodUser() *PodUser {
	return &PodUser{baby: true}
}

type _PodUserMgr struct{}

type IxEntityPodUserCreateDtAge struct {
	CreateDt int64
	Age      int32
}

type IxEntityPodUserUpdateDt struct {
	UpdateDt int64
}

var PodUserMgr = &_PodUserMgr{}

func (mgr *_PodUserMgr) IsNicknameMobilePhoneExists(nickname string, mobilePhone string) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where nickname = ? and mobile_phone = ?`,
		nickname, mobilePhone)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) UniFindByNicknameMobilePhone(nickname string, mobilePhone string) (*PodUser, error) {
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where nickname = ? and mobile_phone = ?`,
		nickname, mobilePhone)

	var id sql.NullInt64
	var nickname_1 sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone_1 sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id, &nickname_1, &password, &age, &mobilePhone_1, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id.Int64
	d.Nickname = nickname_1.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone_1.String
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	return &d, nil
}

func (mgr *_PodUserMgr) UpdateByNicknameMobilePhone(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`update pod.pod_user set password = ?, age = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where nickname = ? and mobile_phone = ?`, d.Password, d.Age, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (mgr *_PodUserMgr) UniRmByNicknameMobilePhone(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`delete from pod.pod_user where nickname = ? and mobile_phone = ?`, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
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
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where mobile_phone = ?`,
		mobilePhone)

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone_1 sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id, &nickname, &password, &age, &mobilePhone_1, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id.Int64
	d.Nickname = nickname.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone_1.String
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	return &d, nil
}

func (mgr *_PodUserMgr) UpdateByMobilePhone(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`update pod.pod_user set password = ?, age = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where mobile_phone = ?`, d.Password, d.Age, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (mgr *_PodUserMgr) UniRmByMobilePhone(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`delete from pod.pod_user where mobile_phone = ?`, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (mgr *_PodUserMgr) FindByCreateDtAge(createDt int64, age int32, order []string, offset, limit int64) ([]*PodUser, error) {
	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where create_dt = ? and age = ?`
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	rows, err := db.DB().Query(query, createDt, age)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age_1 sql.NullInt64
	var mobilePhone sql.NullString
	var createDt_1 sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age_1, &mobilePhone, &createDt_1, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age_1.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt_1.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, nil
}

// 通过索引数组查询
func (mgr *_PodUserMgr) FindByCreateDtAgeArray(entities []*IxEntityPodUserCreateDtAge, order []string, offset, limit int64) ([]*PodUser, string, error) {
	if len(entities) == 0 {
		return nil, "", errors.New("input entities empty. ")
	}

	str := "(?,?)" + strings.Repeat(",(?,?)", len(entities)-1)
	query := fmt.Sprintf("select `id`, `nickname`, `password`, `age`, `mobile_phone`, `create_dt`, `is_blocked`, `update_dt`, `stats_dt`, `dt` from pod.pod_user where (`create_dt`, `age`) in (%s)", str)
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	params := make([]interface{}, 0, 2*len(entities))
	for _, entity := range entities {
		params = append(params, entity.CreateDt)
		params = append(params, entity.Age)
	}
	rows, err := db.DB().Query(query, params...)
	if err != nil {
		return nil, query, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, query, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, query, nil
}

func (mgr *_PodUserMgr) CountByCreateDtAge(createDt int64, age int32) (int64, error) {
	query := `select count(1) from pod.pod_user where create_dt = ? and age = ?`
	row := db.DB().QueryRow(query, createDt, age)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByUpdateDt(updateDt int64, order []string, offset, limit int64) ([]*PodUser, error) {
	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where update_dt = ?`
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
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
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt_1, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt_1.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, nil
}

// 通过索引数组查询
func (mgr *_PodUserMgr) FindByUpdateDtArray(entities []*IxEntityPodUserUpdateDt, order []string, offset, limit int64) ([]*PodUser, string, error) {
	if len(entities) == 0 {
		return nil, "", errors.New("input entities empty. ")
	}

	str := "?" + strings.Repeat(",?", len(entities)-1)
	query := fmt.Sprintf("select `id`, `nickname`, `password`, `age`, `mobile_phone`, `create_dt`, `is_blocked`, `update_dt`, `stats_dt`, `dt` from pod.pod_user where (`update_dt`) in (%s)", str)
	for i, o := range order {
		if i == 0 {
			query += " order by "
		} else {
			query += ", "
		}
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
		if o[0] == '-' {
			query += " desc"
		}
	}
	if offset != -1 && limit != -1 {
		query += fmt.Sprintf(" limit %d, %d", offset, limit)
	}

	params := make([]interface{}, 0, 1*len(entities))
	for _, entity := range entities {
		params = append(params, entity.UpdateDt)
	}
	rows, err := db.DB().Query(query, params...)
	if err != nil {
		return nil, query, err
	}

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, query, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, query, nil
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

func (mgr *_PodUserMgr) FindByMultiJoin(joins []db.Join, where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	var params []interface{}

	query := `select pod_user.id, pod_user.nickname, pod_user.password, pod_user.age, pod_user.mobile_phone, pod_user.create_dt, pod_user.is_blocked, pod_user.update_dt, pod_user.stats_dt, pod_user.dt from pod.pod_user`
	for _, join := range joins {
		query += ` join pod.` + join.T + ` on `
		for i, v := range join.Rule {
			if i > 0 {
				query += " and "
			}
			query += v.S
			if v.P != nil {
				params = append(params, v.P)
			}
		}
	}
	for i, v := range where {
		if i == 0 {
			query += " where "
		} else {
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
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
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
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) CountByMultiJoin(joins []db.Join, where []db.Rule) (int64, error) {
	var params []interface{}

	query := `select count(1) from (select pod_user.id, pod_user.nickname, pod_user.password, pod_user.age, pod_user.mobile_phone, pod_user.create_dt, pod_user.is_blocked, pod_user.update_dt, pod_user.stats_dt, pod_user.dt from pod.pod_user`
	for _, join := range joins {
		query += ` join pod.` + join.T + ` on `
		for i, v := range join.Rule {
			if i > 0 {
				query += " and "
			}
			query += v.S
			if v.P != nil {
				params = append(params, v.P)
			}
		}
	}
	for i, v := range where {
		if i == 0 {
			query += " where "
		} else {
			query += " and "
		}
		query += v.S
		if v.P != nil {
			params = append(params, v.P)
		}
	}
	query += ") t"

	row := db.DB().QueryRow(query, params...)
	var c sql.NullInt64
	if err := row.Scan(&c); err != nil {
		return 0, err
	}
	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByJoin(t string, on, where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	return mgr.FindByMultiJoin([]db.Join{
		{T: t, Rule: on},
	}, where, order, offset, limit)
}
func (mgr *_PodUserMgr) FindByCond(where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where `
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
		} else {
			query += ", "
		}
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
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
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) FindAllByCond(where []db.Rule, order []string) ([]*PodUser, error) {
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user `
	for i, v := range where {
		if i == 0 {
			query += " where "
		} else if i > 0 {
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
		if strings.HasPrefix(o, "-") {
			query += o[1:]
		} else {
			query += o
		}
		if o[0] == '-' {
			query += " desc"
		}
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
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	return ret, nil
}

func (mgr *_PodUserMgr) Create(d *PodUser) error {
	r, err := db.DB().Exec(`insert into pod.pod_user (nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt) value (?,?,?,?,?,?,?,?,?)`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt)
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
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where id = ?`, id)

	var id_1 sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id_1, &nickname, &password, &age, &mobilePhone, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id_1.Int64
	d.Nickname = nickname.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone.String
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	return &d, nil
}

func (mgr *_PodUserMgr) Update(d *PodUser) (int64, error) {
	r, err := db.DB().Exec(`update pod.pod_user set nickname = ?, password = ?, age = ?, mobile_phone = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where id = ?`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Id)
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
func (mgr *_PodUserMgr) IsExistsByPK(pk int64) (bool, error) {
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where id = ?`, pk)
	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c.Int64 > 0, nil
}
