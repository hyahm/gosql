# gosql

已经不再维护， 建议使用xorm
mysql 只是简单封装
 - [x] 支持高并发
 - [x] 支持更新和删除失败的日志记录
 - [x] 支持驱动自带的连接池，
 - [x] 避免连接过多导致的失败
 - [x] 支持in的操作

# mysql Curd
```go

// 高级方法
// 高级方法只有第一级的db有效，后面的都无视
type MeStruct struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// 使用高级方法的第一个是对应数据库的字段
// default: 插入时候，如果没有传入值将使用数据库default的值， 如果没写就是默认值
// force： 修改的时候， 如果设置了force， 那么强制修改字段的值， 如果没写， 零值的时候不会修改值
// 主键必须设置 omitempty 并且不能有force
type Person struct {
	ID        int64    `db:"id,default"`
	FirstName string   `db:"first_name,force"`
	LastName  string   `db:"last_name"`
	Email     string   `db:"email,default,force"`
	Me        MeStruct `db:"me"`
	
}


func main() {
	conf := gosql.Sqlconfig{
		UserName:        "test",
		Password:        "123456",
		Port:            3306,
		DbName:          "test",
		Host:            "192.168.101.4",
		MultiStatements: true,
	}
	// 如果要使用 postgresql   pg, err := conf.NewPGPool()
	db, err := conf.NewMysqlDb()
	if err != nil {
		t.Fatal(err)
	}
	
	// 插入
	ps := &Person{
		FirstName: "cander",
		LastName:  "biao",
		Email:     "aaaaa@eaml.com",
		Me: MeStruct{
			X: 10,
			Y: 20,
			Z: 30,
		},
	}
	pss := make([]*Person, 0)
	for i := 0; i < 20; i++ {
		pss = append(pss, ps)
	}
	// $key  $value 是插入事的固定占位符， 在这个占位符之前不能有参数占位符？，如果有的话请使用 Insert处理
	// default: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// 传入的 dest 值 可以是指针，可以是数据，可以是结构体
	result = db.InsertInterfaceWithoutID(pss, "insert into person($key) values($value)")
	if result.Err != nil {
		t.Fatal(err)
	}
	// 将会生成20条数据

	// 修改
	updateps := &Person{
		FirstName: "what is it",
		LastName:  "hyahm.com",
		Email:     "aaaaa@eaml.com",
		Me: MeStruct{
			X: 10,
			Y: 20,
			Z: 30,
		},
	}

	// $set 是固定占位符, 前面也必须没有参数占位符 ?
	// omitempty: 如果为空， 那么为数据库的默认值
	// 传入的值必须是指针或结构体
	result = db.UpdateInterface(updateps, "update person set $set where id=?", 1)
	if result.Err != nil {
		t.Fatal(err)
	}
	// 执行后会修改id为1的行
	persons := make([]*Category, 0)

	result = db.Select(&persons, "select * from Person")
	if result.Err != nil {
		fmt.Println(err)
	}

	for _, v := range persons {
		fmt.Printf("%#v", *v)
	}
}
```




