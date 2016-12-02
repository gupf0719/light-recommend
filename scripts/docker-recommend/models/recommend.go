package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Recommend struct {
	Id        int64  `orm:"auto"`
	ObjectId  string `orm:"size(32)"` //object_id
	Lights    string `orm:"size(64)"` //根据好友推荐的轻应用
	Counts    string `orm:"size(64)"` //推荐的轻应用收藏的好友数
	CreatedAt int64  //创建时间
	UpdatedAt int64  //更新时间
}

func init() {
	orm.RegisterModel(new(Recommend))
}

func AddRec(m *Recommend) error {
	o := orm.NewOrm()
	recommend := new(Recommend)

	qs := o.QueryTable("recommend")
	err := qs.Filter("object_id", m.ObjectId).One(recommend)

	if err == nil {
		recommend.ObjectId = m.ObjectId
		recommend.Lights = m.Lights
		recommend.Counts = m.Counts
		recommend.UpdatedAt = time.Now().Unix()

		o.Update(recommend)
		return err
	} else {
		_, err := o.Insert(m)
		return err
	}
}
