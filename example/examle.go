package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hyahm/gosql"
)

type MpMessageList struct {
	ID              int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id" db:"id,omitempty"`
	TypeID          int64     `gorm:"notNull;column:type_id" json:"type_id" db:"type_id,default"`
	Title           string    `gorm:"type:varchar(255);notNull;default:'';column:title" json:"title" db:"title,default"`
	Summary         string    `gorm:"type:varchar(255);notNull;default:'';column:summary" json:"summary" db:"summary,default"` // 摘要
	UpdateAt        time.Time `gorm:"type:datetime;notNull;column:update_at" json:"update_at" db:"update_at,default"`
	CreateAt        time.Time `gorm:"type:datetime;notNull;column:create_at" json:"create_at" db:"create_at,default"`
	ContentType     int       `gorm:"type:int;notNull;default:0;column:content_type" json:"content_type" db:"content_type,default"`
	Platform        int       `gorm:"type:int;notNull;default:0;column:platform" json:"platform" db:"platform,default"` // 0: 空  1: all  2: android  3: ios
	Scope           int       `gorm:"type:int;notNull;default:0;column:scope" json:"scope" db:"scope,default"`          // 推广范围
	Content         string    `gorm:"type:longtext;notNull;default:'';column:content" json:"content" db:"content,default"`
	AuditStatus     int       `gorm:"type:int;notNull;default:0;column:audit_status" json:"audit_status" db:"audit_status,force"` // 0: 无 1: 待审核 2: 审核通过 3: 审核不通过
	AuditTime       time.Time `gorm:"type:datetime;notNull;column:audit_time" json:"audit_time" db:"audit_time,default"`
	SendTime        time.Time `gorm:"type:datetime;notNull;default:CURRENT_TIMESTAMP;column:send_time" json:"send_time" db:"send_time,default"`
	SendType        int       `gorm:"type:int;notNull;default:0;column:send_type" json:"send_type" db:"send_type,default"`       //   0： 空 1:立即发送   2:定时发送
	SendStatus      int       `gorm:"type:int;notNull;default:0;column:send_status" json:"send_status" db:"send_status,default"` // 0: 空 1: 已发送 2: 发送失败
	SendCount       int       `gorm:"type:int;notNull;default:0;column:send_count" json:"send_count" db:"send_count,default"`
	ReceptionVolume int       `gorm:"type:int;notNull;default:0;column:reception_volume" json:"reception_volume" db:"reception_volume,default"`
	ClickVolume     int       `gorm:"type:int;notNull;default:0;column:click_volume" json:"click_volume" db:"click_volume,default"`
	UID             int64     `gorm:"type:bigint;notNull;default:0;column:uid" json:"uid" db:"uid,default"`
	MsgId           string    `gorm:"type:nvarchar(20);notNull;default:0;column:msg_id" json:"msg_id" db:"msg_id,default"`
	Deleted         bool      `gorm:"column:deleted" json:"deleted" db:"deleted,default"`
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
    content LONGTEXT not null default '',
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

type HengshengZsFundflow struct {
	Code                  string    `gorm:"primaryKey;column:code" json:"code"`                                                // 指数代码
	SecuMarket            int       `gorm:"notNull;default:0;column:secu_market" json:"secu_market"`                           // 指数代码市场
	Date                  int       `gorm:"notNull;column:date" json:"date"`                                                   // 日期
	MainNetTurnover       int64     `gorm:"notNull;default:0;column:main_net_turnover" json:"main_net_turnover"`               // 主力资金净额
	NetTurnover           int64     `gorm:"notNull;default:0;column:net_turnover" json:"net_turnover"`                         // 资金净额
	SuperAmountIn         int64     `gorm:"notNull;default:0;column:super_amount_in" json:"super_amount_in"`                   // 超大单流入成交量
	SuperAmountOut        int64     `gorm:"notNull;default:0;column:super_amount_out" json:"super_amount_out"`                 // 超大单流出成交量
	SuperTurnoverIn       int64     `gorm:"notNull;default:0;column:super_turnover_in" json:"super_turnover_in"`               // 超大单流入成交额
	SuperTurnoverOut      int64     `gorm:"notNull;default:0;column:super_turnover_out" json:"super_turnover_out"`             // 超大单流出成交额
	SuperCountIn          int64     `gorm:"notNull;default:0;column:super_count_in" json:"super_count_in"`                     // 超大单流入成交笔数
	SuperCountOut         int64     `gorm:"notNull;default:0;column:super_count_out" json:"super_count_out"`                   // 超大单流出成交笔数
	SuperEntrustCountIn   int64     `gorm:"notNull;default:0;column:super_entrust_count_in" json:"super_entrust_count_in"`     // 超大单流入委托单数
	SuperEntrustCountOut  int64     `gorm:"notNull;default:0;column:super_entrust_count_out" json:"super_entrust_count_out"`   // 超大单流出委托单数
	LargeAmountIn         int64     `gorm:"notNull;default:0;column:large_amount_in" json:"large_amount_in"`                   // 大单流入成交量
	LargeAmountOut        int64     `gorm:"notNull;default:0;column:large_amount_out" json:"large_amount_out"`                 // 大单流出成交量
	LargeTurnoverIn       int64     `gorm:"notNull;default:0;column:large_turnover_in" json:"large_turnover_in"`               // 大单流入成交额
	LargeTurnoverOut      int64     `gorm:"notNull;default:0;column:large_turnover_out" json:"large_turnover_out"`             // 大单流出成交额
	LargeCountIn          int64     `gorm:"notNull;default:0;column:large_count_in" json:"large_count_in"`                     // 大单流入成交笔数
	LargeCountOut         int64     `gorm:"notNull;default:0;column:large_count_out" json:"large_count_out"`                   // 大单流出成交笔数
	LargeEntrustCountIn   int64     `gorm:"notNull;default:0;column:large_entrust_count_in" json:"large_entrust_count_in"`     // 大单流入委托单数
	LargeEntrustCountOut  int64     `gorm:"notNull;default:0;column:large_entrust_count_out" json:"large_entrust_count_out"`   // 大单流出委托单数
	MediumAmountIn        int64     `gorm:"notNull;default:0;column:medium_amount_in" json:"medium_amount_in"`                 // 中单流入成交量
	MediumAmountOut       int64     `gorm:"notNull;default:0;column:medium_amount_out" json:"medium_amount_out"`               // 中单流出成交量
	MediumTurnoverIn      int64     `gorm:"notNull;default:0;column:medium_turnover_in" json:"medium_turnover_in"`             // 中单流入成交额
	MediumTurnoverOut     int64     `gorm:"notNull;default:0;column:medium_turnover_out" json:"medium_turnover_out"`           // 中单流出成交额
	MediumCountIn         int64     `gorm:"notNull;default:0;column:medium_count_in" json:"medium_count_in"`                   // 中单流入成交笔数
	MediumCountOut        int64     `gorm:"notNull;default:0;column:medium_count_out" json:"medium_count_out"`                 // 中单流出成交笔数
	MediumEntrustCountIn  int64     `gorm:"notNull;default:0;column:medium_entrust_count_in" json:"medium_entrust_count_in"`   // 中单流入委托单数
	MediumEntrustCountOut int64     `gorm:"notNull;default:0;column:medium_entrust_count_out" json:"medium_entrust_count_out"` // 中单流出委托单数
	LittleAmountIn        int64     `gorm:"notNull;default:0;column:little_amount_in" json:"little_amount_in"`                 // 小单流入成交量
	LittleAmountOut       int64     `gorm:"notNull;default:0;column:little_amount_out" json:"little_amount_out"`               // 小单流出成交量
	LittleTurnoverIn      int64     `gorm:"notNull;default:0;column:little_turnover_in" json:"little_turnover_in"`             // 小单流入成交额
	LittleTurnoverOut     int64     `gorm:"notNull;default:0;column:little_turnover_out" json:"little_turnover_out"`           // 小单流出成交额
	LittleCountIn         int       `gorm:"notNull;default:0;column:little_count_in" json:"little_count_in"`                   // 小单流入成交笔数
	LittleCountOut        int       `gorm:"notNull;default:0;column:little_count_out" json:"little_count_out"`                 // 小单流出成交笔数
	LittleEntrustCountIn  int       `gorm:"notNull;default:0;column:little_entrust_count_in" json:"little_entrust_count_in"`   // 小单流入委托单数
	LittleEntrustCountOut int       `gorm:"notNull;default:0;column:little_entrust_count_out" json:"little_entrust_count_out"` // 小单流出委托单数
	Created               time.Time `gorm:"notNull;default:current_timestamp();column:created" json:"created"`                 // 创建时间
}

func main() {
	var conf = &gosql.Sqlconfig{
		Host:     "192.168.56.128",
		Port:     3306,
		UserName: "cander",
		Password: "123456",
		Debug:    true,
	}
	db, err := conf.CreateDB("test")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(sql)
	ps := &MpMessageList{
		TypeID:      1,
		UID:         1,
		Title:       "test",
		Summary:     "test",
		Content:     "test",
		SendTime:    time.Now(),
		AuditTime:   time.Now(),
		ClickVolume: 1,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		Deleted:     false,
		MsgId:       "test",
		Platform:    1,
	}
	// err = db.InsertInterfaceWithID(ps, "insert into mp_message_list($key) values($value)").Err

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// $key  $value 是固定占位符
	// omitempty: 如果为空， 那么为数据库的默认值
	// struct, 指针， 切片 默认值为 ""
	// $set
	// res := db.UpdateInterface(ps, "update mp_message_list set $set where id=?", 1)
	// if res.Err != nil {
	// 	log.Fatal(res.Err)
	// }
	err = db.Select(ps, "select * from mp_message_list where id=?", 1).Err
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", ps.AuditTime.Format("2006-01-02 15:04:05.000000000"))
	// cate := &User{}
	// res := db.Insert("INSERT INTO user (username, password) VALUES ('77tom', '123') ON DUPLICATE KEY UPDATE username='tom', password='123';")
	// // _, err = db.ReplaceInterface(&cate, "INSERT INTO user ($key) VALUES ($value) ON DUPLICATE KEY UPDATE $set")
	// if res.Err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(res.LastInsertId)
}
