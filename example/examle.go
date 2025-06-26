package main

import (
	"log"
	"slices"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hyahm/gosql"
)

type MpMessageList struct {
	ID              int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id" `
	TypeID          int64     `gorm:"notNull;column:type_id" json:"type_id" `
	Title           string    `gorm:"type:varchar(255);notNull;default:'';column:title" json:"title"`
	Summary         string    `gorm:"type:varchar(255);notNull;default:'';column:summary" json:"summary" ` // 摘要
	UpdateAt        time.Time `gorm:"type:datetime;notNull;column:update_at" json:"update_at" `
	CreateAt        time.Time `gorm:"type:datetime;notNull;column:create_at" json:"create_at" db:"created"`
	ContentType     int       `gorm:"type:int;notNull;default:0;column:content_type" json:"content_type" `
	Platform        int       `gorm:"type:int;notNull;default:0;column:platform" json:"platform" ` // 0: 空  1: all  2: android  3: ios
	Scope           int       `gorm:"type:int;notNull;default:0;column:scope" json:"scope" `       // 推广范围
	Content         string    `gorm:"type:longtext;notNull;default:'';column:content" json:"content" `
	AuditStatus     int       `gorm:"type:int;notNull;default:0;column:audit_status" json:"audit_status" db:"audit_status,force"` // 0: 无 1: 待审核 2: 审核通过 3: 审核不通过
	AuditTime       time.Time `gorm:"type:datetime;notNull;column:audit_time" json:"audit_time" `
	SendTime        time.Time `gorm:"type:datetime;notNull;column:send_time" json:"send_time"`
	SendType        int       `gorm:"type:int;notNull;default:0;column:send_type" json:"send_type" `     //   0： 空 1:立即发送   2:定时发送
	SendStatus      int       `gorm:"type:int;notNull;default:0;column:send_status" json:"send_status" ` // 0: 空 1: 已发送 2: 发送失败
	SendCount       int       `gorm:"type:int;notNull;default:0;column:send_count" json:"send_count" `
	ReceptionVolume int       `gorm:"type:int;notNull;default:0;column:reception_volume" json:"reception_volume" `
	ClickVolume     int       `gorm:"type:int;notNull;default:0;column:click_volume" json:"click_volume" db:"click_volume,default"`
	UID             int64     `gorm:"type:bigint;notNull;default:0;column:uid" json:"uid" `
	MsgId           string    `gorm:"type:nvarchar(20);notNull;default:0;column:msg_id" json:"msg_id" `
	Deleted         bool      `gorm:"column:deleted" json:"deleted" `
}

var sql = `
CREATE TABLE mp_message_list (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    type_id bigint NOT NULL,
    title VARCHAR(255) not null default '',
    summary nvarchar(255) not null default '',
    update_at datetime  default null,
    platform int not null default 0,
	scope int not null default 0,
    create_at datetime not null,
    content_type int not null default 0,
    content LONGTEXT not null ,
    audit_time datetime default null,
    audit_status int not null default 0,    --  0: 无   1： 待审核 2： 审核不通过  3： 审核通过
    push_time datetime default null,
	click_volume int not null default 0,
    send_time datetime default null ,
    send_status int not null default 0,   --   0:  空   1： 已发送  
    send_count int not null default 0,
    reception_volume int not null default 0,
    send_type INT NOT NULL DEFAULT 0,    --   0:  空   1： 立即发送  2： 邮件  3： 微信  4： 短信
	uid bigint not null default 0,
    deleted boolean not null default false,
    msg_id nvarchar(20) not null default ''
) comment "消息推送列表";

`

// type MpTypeTable struct {
// 	ID       int64     `gorm:"column:id;primaryKey" json:"id" db:"id"`
// 	UID      int64     `gorm:"column:uid;primaryKey" json:"uid" db:"uid"`
// 	TypeName string    `gorm:"column:type_name;not null" json:"type_name" db:"type_name"`
// 	Icon     string    `gorm:"column:icon" json:"icon" db:"icon"`
// 	UpdateAt time.Time `gorm:"column:update_at" json:"update_at" db:"updated"`
// 	ParentID int64     `gorm:"column:parent_id" json:"parent_id,omitempty" db:"parent_id"` // 使用指针类型以支持 NULL 值
// 	Deleted  bool      `gorm:"column:deleted" json:"-" db:"deleted"`
// }

func sort() {
	a := []int{1, 2, 3, 4, 5}
	for i := 0; i < len(a); i++ {
		....
	}
	a = nil
}

func main() {

	var conf = &gosql.Sqlconfig{
		Host:     "192.168.3.110",
		Port:     3306,
		UserName: "cander",
		Password: "123456",
		Debug:    true,
	}
	db, err := conf.CreateDB("dxzg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(sql)
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok {
			if v.Number != 1050 {
				log.Fatal(err)
			}
		}
	}

	ps := &MpMessageList{
		ID:          38,
		TypeID:      1,
		ContentType: 1,
		Content:     "图片地址",
		Title:       "测试",
		SendTime:    time.Now(),
	}
	// err = db.InsertInterfaceWithID(ps, "insert into mp_message_list($key) values($value)").Err

	// if err != nil {

	// 	log.Fatal(err)
	// }
	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	// res := db.UpdateInterface(ps, "update mp_message_list set $set where id=?", ps.ID)
	// if res.Err != nil {
	// 	log.Fatal(res.Err)
	// }
	err = db.Select(ps, "select * from mp_message_list where id=?", 38).Err
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ps.SendTime)
	// fmt.Printf("%v\n", ps.Icon.Format("2006-01-02 15:04:05.000000000"))
	// // cate := &User{}
	// res = db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// // _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	// if res.Err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res.LastInsertId)
}
