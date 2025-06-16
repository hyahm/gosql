package gosql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 增加 default: 使用默认值
// 修改： force, counter  增量器
// 删除 修改 需要 primarykey

func fill(dest interface{}, rows *sql.Rows) error {
	value := reflect.ValueOf(dest)
	typ := reflect.TypeOf(dest)
	// cols := 0
	// // json.Unmarshal returns errors for these
	if typ.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	// stt 是数组基础数据结构

	typ = typ.Elem()
	// 判断是否是数组
	isArr := false
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
		isArr = true
	}
	// 标识最后的接受体是指针还是结构体
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	// ss 是切片
	ss := value.Elem()
	names := make(map[string]int)
	cls, _ := rows.Columns()
	for i, v := range cls {
		names[v] = i
	}

	vals := make([][]byte, len(cls))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cls))
	//这里scans引用vals，把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}
	for rows.Next() {
		// scan into the struct field pointers and append to our results
		err := rows.Scan(scans...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		new := reflect.New(typ)
		if !isPtr {
			new = new.Elem()
		}
		if new.Type().Kind() == reflect.Ptr {
			new = new.Elem()
		}
		for index := 0; index < typ.NumField(); index++ {
			dbname := typ.Field(index).Tag.Get("db")
			tags := strings.Split(dbname, ",")
			if len(tags) == 0 {
				continue
			}
			if tags[0] == "" {
				continue
			}

			if v, ok := names[tags[0]]; ok {
				if new.Field(index).CanSet() {
					// 判断这一列的值
					kind := new.Field(index).Kind()
					b := *(scans[v]).(*[]byte)
					switch kind {
					case reflect.String:
						new.Field(index).SetString(string(b))
					case reflect.Int64:
						i64, _ := strconv.ParseInt(string(b), 10, 64)
						new.Field(index).SetInt(i64)
					case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
						i, _ := strconv.Atoi(string(b))
						new.Field(index).Set(reflect.ValueOf(i))

					case reflect.Bool:
						t, _ := strconv.ParseBool(string(b))
						new.Field(index).SetBool(t)

					case reflect.Float32:
						f64, _ := strconv.ParseFloat(string(b), 32)
						new.Field(index).SetFloat(f64)

					case reflect.Float64:
						f64, _ := strconv.ParseFloat(string(b), 64)
						new.Field(index).SetFloat(f64)

					case reflect.Struct:
						if new.Field(index).Type().String() == "time.Time" {
							tv, err := time.Parse("2006-01-02 15:04:05", string(b))
							if err != nil {
								new.Field(index).Set(reflect.ValueOf(time.Time{}))
								continue
							}
							new.Field(index).Set(reflect.ValueOf(tv))
							continue
							// tv := value.Field(v).Interface().(time.Time).Format("2006-01-02 15:04:05")
							// keys = append(keys, signs[0])
							// values = append(values, tv)
							// placeholders = append(placeholders, "?")
							// continue
						}

						j := reflect.New(new.Field(index).Type())
						json.Unmarshal(b, j.Interface())
						new.Field(index).Set(j.Elem())

					case reflect.Slice, reflect.Interface:
						j := reflect.New(new.Field(index).Type())
						err = json.Unmarshal(b, j.Interface())
						if err != nil {
							new.Field(index).Set(reflect.MakeSlice(new.Field(index).Type(), 0, 0))
							continue
						}
						new.Field(index).Set(j.Elem())

					case reflect.Ptr:
						j := reflect.New(new.Field(index).Type())
						json.Unmarshal(b, j.Interface())
						new.Field(index).Set(j)
					default:
						log.Println("not support , you can make a issue to report in https://github.com/hyahm/gosql, kind: ", kind)
					}
				} else {
					fmt.Println("can not set: ", index)
				}
			}

		}
		if !isArr {
			if isPtr {
				value.Elem().Elem().Set(new)
			} else {
				value.Elem().Set(new)
			}

			return nil
		} else {
			if isPtr {
				ss = reflect.Append(ss, new.Addr())
			} else {
				ss = reflect.Append(ss, new)

			}
		}
	}
	value.Elem().Set(ss)
	return nil
}

func insertInterfaceSql(dest interface{}, cmd string, args ...interface{}) (string, []interface{}, error) {
	// 插入到args之前  dest 是struct或切片的指针
	if !strings.Contains(cmd, "$key") {
		return "", nil, errors.New("not found placeholders $key")
	}

	if !strings.Contains(cmd, "$value") {
		return "", nil, errors.New("not found placeholders $value")
	}
	values := make([]interface{}, 0)
	keys := make([]string, 0)
	// ？号
	placeholders := make([]string, 0)
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)

	if typ.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Struct {
		// 如果是struct， 执行插入
		for i := 0; i < value.NumField(); i++ {
			tag := typ.Field(i).Tag.Get("db")
			if tag == "" {
				continue
			}
			signs := strings.Split(tag, ",")
			kind := value.Field(i).Kind()
			switch kind {
			case reflect.String:
				if value.Field(i) == reflect.ValueOf("") && strings.Contains(tag, "default") {
					continue
				}
				keys = append(keys, signs[0])
				// placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())

			case reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
				if value.Field(i).Int() == 0 && !strings.Contains(tag, "force") {
					continue
				}
				keys = append(keys, signs[0])
				// placeholders = append(placeholders, "?")
				values = append(values, value.Field(i).Interface())
			case reflect.Float32, reflect.Float64:
				if value.Field(i).Float() == 0 && !strings.Contains(tag, "force") {
					continue
				}
				keys = append(keys, signs[0])
				values = append(values, value.Field(i).Interface())
			case reflect.Uint64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32:
				if value.Field(i).Uint() == 0 && !strings.Contains(tag, "force") {
					continue
				}
				keys = append(keys, signs[0])
				values = append(values, value.Field(i).Interface())
			case reflect.Bool:
				keys = append(keys, signs[0])
				values = append(values, value.Field(i).Interface())
			case reflect.Slice:
				if value.Field(i).IsNil() {
					keys = append(keys, signs[0])
					// placeholders = append(placeholders, "?")
					values = append(values, "[]")
				} else {
					if value.Field(i).Len() == 0 && !strings.Contains(tag, "default") {
						continue
					}

					keys = append(keys, signs[0])
					// placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "[]")
						placeholders = append(placeholders, "?")
						continue
					}
					values = append(values, send)
				}
			case reflect.Ptr:
				if value.Field(i).IsNil() {
					if !strings.Contains(tag, "default") {
						continue
					}
					keys = append(keys, signs[0])
					// placeholders = append(placeholders, "?")
					values = append(values, "{}")
				} else {
					keys = append(keys, signs[0])
					// placeholders = append(placeholders, "?")
					send, err := json.Marshal(value.Field(i).Interface())
					if err != nil {
						values = append(values, "{}")
						placeholders = append(placeholders, "?")
						continue
					}
					values = append(values, send)
				}
			case reflect.Struct, reflect.Interface:
				if typ.Field(i).Type.String() == "time.Time" {
					if !value.Field(i).Interface().(time.Time).IsZero() {
						tv := value.Field(i).Interface().(time.Time).Format("2006-01-02 15:04:05")
						keys = append(keys, signs[0])
						values = append(values, tv)
						placeholders = append(placeholders, "?")
						continue
					} else {
						if strings.Contains(tag, "default") {
							keys = append(keys, signs[0])
							values = append(values, time.Time{}.Format("2006-01-02 15:04:05"))
						}
						placeholders = append(placeholders, "?")
						continue
					}

				}
				keys = append(keys, signs[0])
				// placeholders = append(placeholders, "?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "{}")
					placeholders = append(placeholders, "?")
					continue
				}
				values = append(values, send)
			default:
				return "", nil, errors.New("not support , you can add issue: " + kind.String())
			}
			placeholders = append(placeholders, "?")
		}
	}

	cmd = strings.Replace(cmd, "$key", strings.Join(keys, ","), 1)
	cmd = strings.Replace(cmd, "$value", strings.Join(placeholders, ","), 1)
	newargs := append(values, args...)
	return cmd, newargs, nil
}

func updateInterfaceSql(dest interface{}, cmd string, args ...interface{}) (string, []interface{}, error) {
	// $set 固定位置固定值
	// db.UpdateInterface(&value, "update test set $set where id=1")
	// 插入到args之前  dest 是struct或切片的指针
	if !strings.Contains(cmd, "$set") {
		return "", nil, errors.New("not found placeholders $set")
	}

	// ？号
	typ := reflect.TypeOf(dest)
	value := reflect.ValueOf(dest)

	if typ.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return "", nil, errors.New("dest must ptr or struct")
	}
	values := make([]interface{}, 0)
	keys := make([]string, 0)
	// 如果是struct， 执行插入
	for i := 0; i < value.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("db")
		if tag == "" {
			continue
		}
		signs := strings.Split(tag, ",")
		// 获取字段名
		kind := value.Field(i).Kind()
		switch kind {

		case reflect.String:
			if value.Field(i).Interface().(string) == "" && !strings.Contains(tag, "force") {
				continue
			}

			values = append(values, value.Field(i).Interface())
		case reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32:
			if value.Field(i).Int() == 0 && !strings.Contains(tag, "force") {
				continue
			}
			values = append(values, value.Field(i).Interface())

		case reflect.Float32, reflect.Float64:
			if value.Field(i).Float() == 0 && !strings.Contains(tag, "force") {
				continue
			}

			values = append(values, value.Field(i).Interface())
		case reflect.Uint64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32:
			if value.Field(i).Uint() == 0 && !strings.Contains(tag, "force") {
				continue
			}

			values = append(values, value.Field(i).Interface())
		case reflect.Bool:
			if !value.Field(i).Bool() && !strings.Contains(tag, "force") {
				continue
			}
			// keys = append(keys, signs[0]+"=?")
			values = append(values, value.Field(i).Interface())
		case reflect.Slice:
			if value.Field(i).IsNil() {
				if !strings.Contains(tag, "force") {
					continue
				}
				// keys = append(keys, signs[0]+"=?")
				values = append(values, "")
			} else {
				if value.Field(i).Len() == 0 && !strings.Contains(tag, "force") {
					continue
				}
				// keys = append(keys, signs[0]+"=?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "")
					goto end
				}
				values = append(values, string(send))
			}
		case reflect.Ptr:
			if value.Field(i).IsNil() {
				if !strings.Contains(tag, "force") {
					continue
				}
				// keys = append(keys, signs[0]+"=?")
				values = append(values, "")
			} else {
				// keys = append(keys, signs[0]+"=?")
				send, err := json.Marshal(value.Field(i).Interface())
				if err != nil {
					values = append(values, "")
					goto end
				}
				values = append(values, string(send))
			}
		case reflect.Struct, reflect.Interface:
			empty := reflect.New(reflect.TypeOf(value.Field(i).Interface())).Elem().Interface()
			if typ.Field(i).Type.String() == "time.Time" {
				if value.Field(i).Interface().(time.Time).IsZero() && !strings.Contains(tag, "force") {
					continue
				}
				if value.Field(i).Interface().(time.Time).IsZero() && strings.Contains(tag, "force") {
					values = append(values, value.Field(i).Interface().(time.Time).Format("2006-01-02 15:04:05"))
					goto end
				}
				tv := value.Field(i).Interface().(time.Time).Format("2006-01-02 15:04:05")
				values = append(values, tv)
				goto end
			}
			if reflect.DeepEqual(value.Field(i).Interface(), empty) {
				continue
			}
			// keys = append(keys, signs[0]+"=?")
			send, err := json.Marshal(value.Field(i).Interface())
			if err != nil {
				values = append(values, "")
				goto end
			}
			values = append(values, string(send))
		default:
			return "", nil, errors.New("not support , you can add issue: " + kind.String())
		}
	end:
		if strings.Contains(tag, "counter") {
			keys = append(keys, fmt.Sprintf("%s=%s+?", signs[0], signs[0]))
		} else {
			keys = append(keys, signs[0]+"=?")
		}
	}

	cmd = strings.Replace(cmd, "$set", strings.Join(keys, ","), 1)

	newargs := append(values, args...)
	return cmd, newargs, nil
}
