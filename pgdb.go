package gosql

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGConn struct {
	*pgxpool.Pool
	conf  string
	Ctx   context.Context
	sc    *Sqlconfig
	mu    *sync.RWMutex
	debug bool
}

func (d *PGConn) InsertInterfaceWithID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// db.InsertInterfaceWithID(&value, "insert into test($key)  values($value)")
	res := Result{
		LastInsertIds: make([]int64, 0),
	}

	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	if typ.Kind() == reflect.Struct {
		return d.insertInterface(dest, cmd, args...)
	}
	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return d.insertInterface(dest, cmd, args...)
		}
		for i := 0; i < length; i++ {
			result := d.insertInterface(value.Index(i).Interface(), cmd, args...)
			res.Sql += ";" + result.Sql
			if result.Err != nil {
				return result
			}
			res.LastInsertIds = append(res.LastInsertIds, result.LastInsertId)
		}
	} else {
		res.Err = ErrNotSupport
	}
	return res
}

func (d *PGConn) insertInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// 插入到args之前  dest 是struct或切片的指针
	newcmd, newargs, err := insertInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Insert(newcmd, newargs...)
}

func (d *PGConn) insertWithoutInterface(dest interface{}, cmd string, args ...interface{}) Result {
	// 插入到args之前  dest 是struct或切片的指针
	newcmd, newargs, err := insertInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.Insert(newcmd, newargs...)
}

func (d *PGConn) Insert(cmd string, args ...interface{}) Result {
	res := Result{
		Sql: ToSql(cmd, args...),
	}
	for i := 1; i <= len(args); i++ {
		cmd = strings.Replace(cmd, "?", fmt.Sprintf("$%d", i), 1)
	}

	result := d.QueryRow(d.Ctx, cmd, args...)
	err := result.Scan(&res.LastInsertId)
	if err != nil {
		res.Err = err
		return res
	}
	return res
}

func (d *PGConn) InsertWithoutId(cmd string, args ...interface{}) Result {
	res := Result{
		Sql: ToSql(cmd, args...),
	}
	for i := 1; i <= len(args); i++ {
		cmd = strings.Replace(cmd, "?", fmt.Sprintf("$%d", i), 1)
	}

	result, err := d.Exec(d.Ctx, cmd, args...)
	res.RowsAffected = result.RowsAffected()
	res.Err = err
	return res
}

func (d *PGConn) UpdateInterface(dest interface{}, cmd string, args ...interface{}) Result {
	newcmd, newargs, err := updateInterfaceSql(dest, cmd, args...)
	if err != nil {
		return Result{Err: err}
	}

	return d.Update(newcmd, newargs...)
}

func (d *PGConn) Update(cmd string, args ...interface{}) Result {
	res := Result{}
	if d.debug {
		res.Sql = ToPGSql(cmd, args...)
	}

	for i := 1; i <= len(args); i++ {
		cmd = strings.Replace(cmd, "?", fmt.Sprintf("$%d", i), 1)
	}

	tags, err := d.Exec(d.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	res.RowsAffected, res.Err = tags.RowsAffected(), err
	return res
}

func (d *PGConn) InsertInterfaceWithoutID(dest interface{}, cmd string, args ...interface{}) Result {
	// $key 和 $value 固定位置固定值
	// ID 自增的必须设置 default
	// db.InsertInterfaceWithoutID(&value, "insert into test($key)  values($value)")
	res := Result{}
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}
	if typ.Kind() == reflect.Struct {
		return d.insertWithoutInterface(dest, cmd, args...)
	}

	if typ.Kind() == reflect.Slice {
		// 如果是切片， 那么每个值都做一次处理
		length := value.Len()
		if length == 1 {
			return d.insertWithoutInterface(dest, cmd, args...)
		}

		arguments := make([]interface{}, 0)
		for i := 0; i < length; i++ {
			newcmd, newargs, err := insertInterfaceSql(value.Index(i).Interface(), cmd, args...)
			if err != nil {
				res.Err = err
				return res
			}
			cmd = newcmd
			arguments = append(arguments, newargs...)
		}
		return d.InsertMany(cmd, arguments...)
	} else {
		res.Err = ErrNotSupport
	}
	return res
}

func (d *PGConn) InsertMany(cmd string, args ...interface{}) Result {
	// sql: insert into test(id, name) values(?,?)  args: interface{}...  1,'t1', 2, 't2', 3, 't3'
	// 每次返回的是第一次插入的id
	if args == nil {
		return d.InsertWithoutId(cmd)
	}
	newcmd, err := formatSql(cmd, args...)
	if err != nil {
		return Result{Err: err}
	}
	return d.InsertWithoutId(newcmd, args...)
}

func (d *PGConn) Select(dest interface{}, cmd string, args ...interface{}) Result {
	// db.Select(&value, "select * from test")
	// 传入切片的地址， 根据tag 的 db 自动补充，
	// 最求性能建议还是使用 GetRows or GetOne
	for i := 1; i <= len(args); i++ {
		cmd = strings.Replace(cmd, "?", fmt.Sprintf("$%d", i), 1)
	}
	res := Result{}
	if d.debug {
		res.Sql = ToPGSql(cmd, args...)
	}
	rows, err := d.Query(d.Ctx, cmd, args...)
	if err != nil {
		res.Err = err
		return res
	}
	defer rows.Close()
	// 需要设置的值
	res.Err = pgfill(dest, rows)
	return res
}

func (d *PGConn) Delete(cmd string, args ...interface{}) Result {
	return d.Update(cmd, args...)
}
