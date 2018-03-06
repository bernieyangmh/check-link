package check_link

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//-1		链接放入管道未爬取
//-2		http请求报错
//-3		读取管道超时，一般为没有新链接放入管道，自动结束
type CUrl struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	CrawlUrl    string        `json:"CrawlUrl" bson:"crawl_url"`
	StatusCode  int           `json:"StatusCode" bson:"status_code"`
	Origin      string        `json:"Origin" bson:"origin"`
	Domain      string        `json:"Domain" bson:"domain"`
	RefUrl      string        `json:"RefUrl" bson:"ref_url"`
	ContentType string        `json:"ContentType" bson:"content_type"`
	updateAt    time.Time     `json:"-" bson:"update_at"`
	QueryError  string        `json:"QueryError" bson:"query_error"`
	Context     string        `json:"Context" bson:"context"`
}
