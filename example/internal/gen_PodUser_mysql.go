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
	Id          int64
	Nickname    string
	Password    string
	Age         int32
	MobilePhone string
	Sequence    int64
	CreateDt    int64
	IsBlocked   bool
	UpdateDt    int64
	StatsDt     *time.Time
	Dt          *time.Time

	baby bool
}

func NewPodUser() *PodUser {
	return &PodUser{baby: true}
}

type _PodUserMgr struct{}

var PodUserMgr = &_PodUserMgr{}

func (mgr *_PodUserMgr) IsNicknameMobilePhoneExists(nickname string, mobilePhone string) (bool, error) {
	util.Log(`pod.pod_user`, `IsNicknameMobilePhoneExists`)
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where nickname = ? and mobile_phone = ?`,
		nickname, mobilePhone)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	util.Log(`pod.pod_user`, `IsNicknameMobilePhoneExists ... done`)
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) UniFindByNicknameMobilePhone(nickname string, mobilePhone string) (*PodUser, error) {
	util.Log(`pod.pod_user`, `UniFindByNicknameMobilePhone`)
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where nickname = ? and mobile_phone = ?`,
		nickname, mobilePhone)

	var id sql.NullInt64
	var nickname_1 sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone_1 sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id, &nickname_1, &password, &age, &mobilePhone_1, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id.Int64
	d.Nickname = nickname_1.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone_1.String
	d.Sequence = sequence.Int64
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	util.Log(`pod.pod_user`, `UniFindByNicknameMobilePhone ... done`)
	return &d, nil
}

func (mgr *_PodUserMgr) UpdateByNicknameMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UpdateByNicknameMobilePhone`)
	r, err := db.DB().Exec(`update pod.pod_user set password = ?, age = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where nickname = ? and mobile_phone = ?`, d.Password, d.Age, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UpdateByNicknameMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) TxUpdateByNicknameMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UpdateByNicknameMobilePhone`)
	r, err := Zotx.Exec(`update pod.pod_user set password = ?, age = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where nickname = ? and mobile_phone = ?`, d.Password, d.Age, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UpdateByNicknameMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) UniRmByNicknameMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UniRmByNicknameMobilePhone`)
	r, err := db.DB().Exec(`delete from pod.pod_user where nickname = ? and mobile_phone = ?`, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UniRmByNicknameMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) TxUniRmByNicknameMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UniRmByNicknameMobilePhone`)
	r, err := Zotx.Exec(`delete from pod.pod_user where nickname = ? and mobile_phone = ?`, d.Nickname, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UniRmByNicknameMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) IsMobilePhoneExists(mobilePhone string) (bool, error) {
	util.Log(`pod.pod_user`, `IsMobilePhoneExists`)
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where mobile_phone = ?`,
		mobilePhone)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	util.Log(`pod.pod_user`, `IsMobilePhoneExists ... done`)
	return c.Int64 > 0, nil
}

func (mgr *_PodUserMgr) UniFindByMobilePhone(mobilePhone string) (*PodUser, error) {
	util.Log(`pod.pod_user`, `UniFindByMobilePhone`)
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where mobile_phone = ?`,
		mobilePhone)

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone_1 sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id, &nickname, &password, &age, &mobilePhone_1, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id.Int64
	d.Nickname = nickname.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone_1.String
	d.Sequence = sequence.Int64
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	util.Log(`pod.pod_user`, `UniFindByMobilePhone ... done`)
	return &d, nil
}

func (mgr *_PodUserMgr) UpdateByMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UpdateByMobilePhone`)
	r, err := db.DB().Exec(`update pod.pod_user set password = ?, age = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where mobile_phone = ?`, d.Password, d.Age, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UpdateByMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) TxUpdateByMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UpdateByMobilePhone`)
	r, err := Zotx.Exec(`update pod.pod_user set password = ?, age = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where mobile_phone = ?`, d.Password, d.Age, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UpdateByMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) UniRmByMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UniRmByMobilePhone`)
	r, err := db.DB().Exec(`delete from pod.pod_user where mobile_phone = ?`, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UniRmByMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) TxUniRmByMobilePhone(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `UniRmByMobilePhone`)
	r, err := Zotx.Exec(`delete from pod.pod_user where mobile_phone = ?`, d.MobilePhone)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `UniRmByMobilePhone ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) FindByCreateDt(createDt int64, order []string, offset, limit int64) ([]*PodUser, error) {
	util.Log(`pod.pod_user`, `FindByCreateDt`)
	query := `select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where create_dt = ?`
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

	rows, err := db.DB().Query(query, createDt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt_1 sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &sequence, &createDt_1, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.Sequence = sequence.Int64
		d.CreateDt = createDt_1.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	util.Log(`pod.pod_user`, `FindByCreateDt ... done`)
	return ret, nil
}

func (mgr *_PodUserMgr) CountByCreateDt(createDt int64) (int64, error) {
	util.Log(`pod.pod_user`, `CountByCreateDt`)
	query := `select count(1) from pod.pod_user where create_dt = ?`
	row := db.DB().QueryRow(query, createDt)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	util.Log(`pod.pod_user`, `CountByCreateDt ... done`)
	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByUpdateDt(updateDt int64, order []string, offset, limit int64) ([]*PodUser, error) {
	util.Log(`pod.pod_user`, `FindByUpdateDt`)
	query := `select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where update_dt = ?`
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
	defer rows.Close()

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt_1 sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt_1, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.Sequence = sequence.Int64
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt_1.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	util.Log(`pod.pod_user`, `FindByUpdateDt ... done`)
	return ret, nil
}

func (mgr *_PodUserMgr) CountByUpdateDt(updateDt int64) (int64, error) {
	util.Log(`pod.pod_user`, `CountByUpdateDt`)
	query := `select count(1) from pod.pod_user where update_dt = ?`
	row := db.DB().QueryRow(query, updateDt)

	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return 0, err
	}

	util.Log(`pod.pod_user`, `CountByUpdateDt ... done`)
	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByMultiJoin(joins []db.Join, where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	util.Log(`pod.pod_user`, `FindByMultiJoin`)
	var params []interface{}

	query := `select pod_user.id, pod_user.nickname, pod_user.password, pod_user.age, pod_user.mobile_phone, pod_user.sequence, pod_user.create_dt, pod_user.is_blocked, pod_user.update_dt, pod_user.stats_dt, pod_user.dt from pod.pod_user`
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
	defer rows.Close()

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.Sequence = sequence.Int64
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	util.Log(`pod.pod_user`, `FindByMultiJoin ... done`)
	return ret, nil
}

func (mgr *_PodUserMgr) CountByMultiJoin(joins []db.Join, where []db.Rule) (int64, error) {
	util.Log(`pod.pod_user`, `CountByMultiJoin`)
	var params []interface{}

	query := `select count(1) from (select pod_user.id, pod_user.nickname, pod_user.password, pod_user.age, pod_user.mobile_phone, pod_user.sequence, pod_user.create_dt, pod_user.is_blocked, pod_user.update_dt, pod_user.stats_dt, pod_user.dt from pod.pod_user`
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
	util.Log(`pod.pod_user`, `CountByMultiJoin ... done`)
	return c.Int64, nil
}

func (mgr *_PodUserMgr) FindByJoin(t string, on, where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	return mgr.FindByMultiJoin([]db.Join{
		{T: t, Rule: on},
	}, where, order, offset, limit)
}
func (mgr *_PodUserMgr) FindByCond(where []db.Rule, order []string, offset, limit int64) ([]*PodUser, error) {
	util.Log(`pod.pod_user`, `FindByCond`)
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where `
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
	defer rows.Close()

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.Sequence = sequence.Int64
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
	util.Log(`pod.pod_user`, `FindAllByCond`)
	var params []interface{}

	query := `select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user `
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
	defer rows.Close()

	var id sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	var ret []*PodUser

	for rows.Next() {
		if err = rows.Scan(&id, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
			return nil, err
		}

		d := PodUser{}
		d.Id = id.Int64
		d.Nickname = nickname.String
		d.Password = password.String
		d.Age = int32(age.Int64)
		d.MobilePhone = mobilePhone.String
		d.Sequence = sequence.Int64
		d.CreateDt = createDt.Int64
		d.IsBlocked = isBlocked.Bool
		d.UpdateDt = updateDt.Int64
		d.StatsDt = util.SafeParseDateStr(statsDt.String)
		d.Dt = util.SafeParseDateTimeStr(dt.String)
		ret = append(ret, &d)
	}
	util.Log(`pod.pod_user`, `FindAllByCond ... done`)
	return ret, nil
}

func (mgr *_PodUserMgr) Create(d *PodUser) error {
	util.Log(`pod.pod_user`, `Create`)
	r, err := db.DB().Exec(`insert into pod.pod_user (nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt) value (?,?,?,?,?,?,?,?,?,?)`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = id
	util.Log(`pod.pod_user`, `Create ... done`)
	return nil
}

func (mgr *_PodUserMgr) TxCreate(d *PodUser) error {
	util.Log(`pod.pod_user`, `TxCreate`)
	r, err := Zotx.Exec(`insert into pod.pod_user (nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt) value (?,?,?,?,?,?,?,?,?,?)`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = id
	util.Log(`pod.pod_user`, `TxCreate ... done`)
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
	util.Log(`pod.pod_user`, `CountByRule`)
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

	util.Log(`pod.pod_user`, `CountByRule ... done`)
	return c.Int64, nil
}

func (mgr *_PodUserMgr) RmByRule(rules ...db.Rule) (int64, error) {
	util.Log(`pod.pod_user`, `RmByRule`)
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
	util.Log(`pod.pod_user`, `RmByRule ... done`)
	return n, nil
}
func (mgr *_PodUserMgr) TxRmByRule(rules ...db.Rule) (int64, error) {
	util.Log(`pod.pod_user`, `RmByRule`)
	query := "delete from pod.pod_user where "
	var p []interface{}
	for i, r := range rules {
		if i > 0 {
			query += " and "
		}
		query += r.S
		p = append(p, r.P)
	}
	r, err := Zotx.Exec(query, p...)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `RmByRule ... done`)
	return n, nil
}
func (mgr *_PodUserMgr) UniFindByPK(id int64) (*PodUser, error) {
	util.Log(`pod.pod_user`, `UniFindByPK`)
	row := db.DB().QueryRow(`select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where id = ?`, id)

	var id_1 sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id_1, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id_1.Int64
	d.Nickname = nickname.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone.String
	d.Sequence = sequence.Int64
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	util.Log(`pod.pod_user`, `UniFindByPK ... done`)
	return &d, nil
}

func (mgr *_PodUserMgr) TxUniFindByPK(id int64) (*PodUser, error) {
	util.Log(`pod.pod_user`, `UniFindByPK`)
	row := Zotx.QueryRow(`select id, nickname, password, age, mobile_phone, sequence, create_dt, is_blocked, update_dt, stats_dt, dt from pod.pod_user where id = ?`, id)

	var id_1 sql.NullInt64
	var nickname sql.NullString
	var password sql.NullString
	var age sql.NullInt64
	var mobilePhone sql.NullString
	var sequence sql.NullInt64
	var createDt sql.NullInt64
	var isBlocked sql.NullBool
	var updateDt sql.NullInt64
	var statsDt sql.NullString
	var dt sql.NullString

	if err := row.Scan(&id_1, &nickname, &password, &age, &mobilePhone, &sequence, &createDt, &isBlocked, &updateDt, &statsDt, &dt); err != nil {
		return nil, err
	}

	d := PodUser{}
	d.Id = id_1.Int64
	d.Nickname = nickname.String
	d.Password = password.String
	d.Age = int32(age.Int64)
	d.MobilePhone = mobilePhone.String
	d.Sequence = sequence.Int64
	d.CreateDt = createDt.Int64
	d.IsBlocked = isBlocked.Bool
	d.UpdateDt = updateDt.Int64
	d.StatsDt = util.SafeParseDateStr(statsDt.String)
	d.Dt = util.SafeParseDateTimeStr(dt.String)
	util.Log(`pod.pod_user`, `UniFindByPK ... done`)
	return &d, nil
}

func (mgr *_PodUserMgr) Update(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `Update`)
	r, err := db.DB().Exec(`update pod.pod_user set nickname = ?, password = ?, age = ?, mobile_phone = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where id = ?`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Id)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `Update ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) TxUpdate(d *PodUser) (int64, error) {
	util.Log(`pod.pod_user`, `Update`)
	r, err := Zotx.Exec(`update pod.pod_user set nickname = ?, password = ?, age = ?, mobile_phone = ?, sequence = ?, create_dt = ?, is_blocked = ?, update_dt = ?, stats_dt = ?, dt = ? where id = ?`, d.Nickname, d.Password, d.Age, d.MobilePhone, d.Sequence, d.CreateDt, d.IsBlocked, d.UpdateDt, d.StatsDt, d.Dt, d.Id)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `Update ... done`)
	return n, nil
}

func (mgr *_PodUserMgr) RmByPK(pk int64) (int64, error) {
	util.Log(`pod.pod_user`, `RmByPK`)
	query := "delete from pod.pod_user where id = ?"
	r, err := db.DB().Exec(query, pk)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `RmByPK ... done`)
	return n, nil
}
func (mgr *_PodUserMgr) TxRmByPK(pk int64) (int64, error) {
	util.Log(`pod.pod_user`, `RmByPK`)
	query := "delete from pod.pod_user where id = ?"
	r, err := Zotx.Exec(query, pk)
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	util.Log(`pod.pod_user`, `RmByPK ... done`)
	return n, nil
}
func (mgr *_PodUserMgr) IsExistsByPK(pk int64) (bool, error) {
	util.Log(`pod.pod_user`, `IsExistsByPK`)
	row := db.DB().QueryRow(`select count(1) from pod.pod_user where id = ?`, pk)
	var c sql.NullInt64

	if err := row.Scan(&c); err != nil {
		return false, err
	}
	util.Log(`pod.pod_user`, `IsExistsByPK ... done`)
	return c.Int64 > 0, nil
}
