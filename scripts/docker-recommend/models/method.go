package models

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id          int64  `orm:"auto"`
	ObjectId    string `orm:"size(64)"`  //ha_id
	Signature   string `orm:"size(64)"`  //签名
	Lights      string `orm:"size(128)"` //轻应用
	Phonenum    string `orm:"size(32)"`  //手机号
	Addressbook string `orm:"size(32)"`  //手机号
	CreatedAt   int64  //创建时间
	UpdatedAt   int64  //更新时间
}

func init() {
	orm.RegisterModel(new(User))
}

//遍历轻应用user,获取手机号和通讯录

func GetUserPhonenum() []orm.Params {

	o := orm.NewOrm()

	var users []orm.Params
	_, err := o.QueryTable("user").Values(&users, "object_id", "addressbook")
	if err != nil {
		beego.Error(err)
	}
	return users
}

//phones := models.GetUserPhonenum()
//		for _, m := range phones {
//			fmt.Println(m["ObjectId"], m["Phonenum"])
//			// map 中的数据都是展开的，没有复杂的嵌套
//		}

//验证user是否存在
func ExistUser(phonenum string) bool {
	o := orm.NewOrm()

	qs := o.QueryTable("user")
	exist := qs.Filter("phonenum", phonenum).Exist()

	if exist {
		return true
	} else {
		return false
	}
}

//获取收藏的轻应用
func GetUserLights(phonenum string) (string, error) {
	o := orm.NewOrm()
	user := new(User)

	qs := o.QueryTable("user")
	err := qs.Filter("phonenum", phonenum).One(user)

	if err == nil {
		return user.Lights, nil
	} else {
		return "", err
	}
}

//冒泡排序，对收藏对轻应用长度进行排序
func BubbleSort(vector []string) []string {
	for i := 0; i < len(vector); i++ {
		for j := 0; j < len(vector)-i-1; j++ {
			if len(vector[j]) > len(vector[j+1]) {
				vector[j], vector[j+1] = vector[j+1], vector[j]
			}
		}
	}
	return vector
}

//对收藏对轻应用进行数量统计
func CountLights(vector []string) (map[int]string, map[int]int) {
	counts := make(map[int]int)    //保存对应轻应用的数量
	lights := make(map[int]string) //保存统计对轻应用id

	//初始化
	lights_init := strings.Split(vector[0], ",")
	for i := 0; i < len(lights_init)-1; i++ {
		lights[i] = lights_init[i]
	}
	for j := 0; j < len(lights_init)-1; j++ {
		counts[j] = 1
	}

	//统计
	for i := 1; i < len(vector); i++ {
		lights_all := strings.Split(vector[i], ",")

		for j := 0; j < len(lights_all)-1; j++ {

			for x := 0; x < len(lights); x++ {
				if strings.EqualFold(lights[x], lights_all[j]) {

					counts[x] = counts[x] + 1

					break
				} else {

					if x == len(lights)-1 {
						lights[x+1] = lights_all[j]
						counts[x+1] = 1

						break
					}

				}
			}
		}
	}

	return lights, counts //返回数据
}

//获取好友收藏数量较多的轻应用
func GetRecommend_3(counts map[int]int, lights map[int]string) ([3]string, [3]int) {

	for i := 0; i < len(counts); i++ {
		for j := 0; j < len(counts)-i-1; j++ {
			if counts[j] > counts[j+1] {
				counts[j], counts[j+1] = counts[j+1], counts[j]
				lights[j], lights[j+1] = lights[j+1], lights[j]
			}
		}
	}

	rec_counts := [3]int{counts[len(counts)-1], counts[len(counts)-2], counts[len(counts)-3]}
	rec_lights := [3]string{lights[len(lights)-1], lights[len(lights)-2], lights[len(lights)-3]}

	return rec_lights, rec_counts

}

func GetRecommend_5(counts map[int]int, lights map[int]string) ([5]string, [5]int) {

	for i := 0; i < len(counts); i++ {
		for j := 0; j < len(counts)-i-1; j++ {
			if counts[j] > counts[j+1] {
				counts[j], counts[j+1] = counts[j+1], counts[j]
				lights[j], lights[j+1] = lights[j+1], lights[j]
			}
		}
	}

	rec_counts := [5]int{counts[len(counts)-1], counts[len(counts)-2], counts[len(counts)-3], counts[len(counts)-4], counts[len(counts)-5]}
	rec_lights := [5]string{lights[len(lights)-1], lights[len(lights)-2], lights[len(lights)-3], lights[len(lights)-4], lights[len(lights)-5]}

	return rec_lights, rec_counts
}
