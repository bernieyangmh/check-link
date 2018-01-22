package check_link

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

var (
	MongoSession, err = mgo.Dial("127.0.0.1")
	DB                = "worktest"
	CheckWebsite      = "check_website"
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
	err := MongoSession.DB(DB).C(CheckWebsite).EnsureIndex(crawlurlIndex)
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
	updateAt    time.Time     `json:"-" bson:"created_at"`
	QueryError  string        `json:"QueryError" bson:"query_error"`
}

func (cu *CUrl) Insert() error {

	session := Session()
	defer session.Close()
	c := session.DB(DB).C(CheckWebsite)
	cu.Id = bson.NewObjectId()
	cu.updateAt = time.Now()

	return c.Insert(cu)

}

func (cu *CUrl) Update() error {
	log.Print("Update mongo		" + cu.CrawlUrl)

	session := Session()
	defer session.Close()
	c := session.DB(DB).C(CheckWebsite)

	selector := bson.M{"crawl_url": cu.CrawlUrl}
	data := bson.M{"crawl_url": "testupdate1111111111"}
	return c.Update(selector, data)

}
