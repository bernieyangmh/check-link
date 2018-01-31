package check_link

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var (
	MongoSession, err = mgo.Dial("127.0.0.1")
	DB                = "worktest"
	CheckUrl          = "check_url"
)

func init() {
	//group coll
	crawlurlIndex := mgo.Index{
		Key:        []string{"crawl_url"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     false,
	}
	err := MongoSession.DB(DB).C(CheckUrl).EnsureIndex(crawlurlIndex)
	if err != nil {
		panic(err)
	}
}

func Session() *mgo.Session {
	return MongoSession.Copy()
}

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

func (cu *CUrl) Insert() error {

	session := Session()
	defer session.Close()
	c := session.DB(DB).C(CheckUrl)
	cu.Id = bson.NewObjectId()
	cu.updateAt = time.Now()

	return c.Insert(cu)

}

func (cu *CUrl) Update() error {
	log.Print("Update mongo		" + cu.CrawlUrl)

	session := Session()
	defer session.Close()
	c := session.DB(DB).C(CheckUrl)

	selector := bson.M{"crawl_url": cu.CrawlUrl}
	data := bson.M{
		"crawl_url":    cu.CrawlUrl,
		"status_code":  cu.StatusCode,
		"origin":       cu.Origin,
		"domain":       cu.Domain,
		"ref_url":      cu.RefUrl,
		"content_type": cu.ContentType,
		"update_at":    time.Now(),
		"query_error":  cu.QueryError,
		"context":      cu.Context,
	}
	return c.Update(selector, data)

}

//todo 不直接返回，抽象出来
func GetIterUrl() *mgo.Iter {
	session := Session()
	c := session.DB(DB).C(CheckUrl)
	find := c.Find(bson.M{}).Select(bson.M{
		"crawl_url":    1,
		"ref_url":      1,
		"status_code":  1,
		"context":      1,
		"query_error":  1,
		"update_at":    1,
		"content_type": 1,
	})
	items := find.Iter()
	return items

}
