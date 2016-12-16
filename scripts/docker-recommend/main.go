package main

import (
	"fmt"
	_ "light-recommend/docs"
	"light-recommend/models"
	_ "light-recommend/routers"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	dbuser := beego.AppConfig.String("dbuser")
	dbpwd := beego.AppConfig.String("dbpwd")
	dbname := beego.AppConfig.String("dbname")
	dbhost := beego.AppConfig.String("dbhost")
	maxIdle, _ := beego.AppConfig.Int("maxidle")
	maxConn, _ := beego.AppConfig.Int("maxconn")

	orm.RegisterDriver("mysql", orm.DRMySQL)

	conn := dbuser + ":" + dbpwd + "@tcp(" + dbhost + ")/" + dbname + "?charset=utf8"
	err := orm.RegisterDataBase("default", "mysql", conn, maxIdle, maxConn)
	if err != nil {
		beego.Error(err.Error)
	}
	orm.Debug = true

	//自动建表
	name := "default"                         //数据库别名
	force := false                            //不强制建数据库
	verbose := true                           //打印建表过程
	err = orm.RunSyncdb(name, force, verbose) //建表

	if err != nil {
		beego.Error(err)
	}
	beego.Debug("初始化")

}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	//好友轻应用推荐
	counts_rec := 1
	for counts_rec > 0 {
		//获取全部user手机号和通讯录
		users := models.GetUserPhonenum()
		for _, u := range users {

			object_id := u["ObjectId"].(string)
			addressbook := u["Addressbook"].(string)

			//通讯录不为空
			if !strings.EqualFold(addressbook, "") {

				//通讯录string转array
				addressbook_array := strings.Split(addressbook, ",")

				//通讯录中号码数量大于4
				addressbook_array_length := len(addressbook_array)

				if addressbook_array_length > 4 {
					lights_all := make([]string, 0, addressbook_array_length)

					//验证通讯录是否注册轻应用
					for _, v := range addressbook_array {
						//获取该手机号的轻应用
						exist := models.ExistUser(v) //用户是否存在
						if exist {
							//获取该手机号收藏的轻应用
							lights_user, err := models.GetUserLights(v)
							if err != nil {
								beego.Error(err)
							}
							//不为空
							if !strings.EqualFold(lights_user, "") {
								fmt.Println(lights_user)
								lights_all = append(lights_all, lights_user)
							}

						}
					}

					lights_all = models.BubbleSort(lights_all)
					fmt.Println("------------lighs_all------------------------------")
					fmt.Println(lights_all)

					lights, counts := models.CountLights(lights_all)

					//收藏的轻应用数量大于2或大于4
					if len(counts) > 2 && len(counts) < 5 {

						rec_lights, rec_counts := models.GetRecommend_3(counts, lights)

						rec_lights_s := rec_lights[0] + "," + rec_lights[1] + "," + rec_lights[2]
						rec_counts_s := strconv.Itoa(rec_counts[0]) + "," + strconv.Itoa(rec_counts[1]) + "," + strconv.Itoa(rec_counts[2])

						recommend := models.Recommend{ObjectId: object_id, Lights: rec_lights_s, Counts: rec_counts_s, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()}
						merr := models.AddRec(&recommend) //插入

						if merr != nil {
							beego.Error(merr)
						}
					} else if len(counts) > 4 {
						rec_lights, rec_counts := models.GetRecommend_5(counts, lights)

						rec_lights_s := rec_lights[0] + "," + rec_lights[1] + "," + rec_lights[2] + "," + rec_lights[3] + "," + rec_lights[4]
						rec_counts_s := strconv.Itoa(rec_counts[0]) + "," + strconv.Itoa(rec_counts[1]) + "," + strconv.Itoa(rec_counts[2]) + "," + strconv.Itoa(rec_counts[3]) + "," + strconv.Itoa(rec_counts[4])

						recommend := models.Recommend{ObjectId: object_id, Lights: rec_lights_s, Counts: rec_counts_s, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()}
						merr := models.AddRec(&recommend) //插入

						if merr != nil {
							beego.Error(merr)
						}
					}
				}
			}
		}

		fmt.Println("")
		fmt.Println(counts_rec)
		fmt.Println("------------------------------------------------------------------------")
		counts_rec++
		time.Sleep(time.Second * 60) //暂停60秒
	}
	beego.Run()
}
