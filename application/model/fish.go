package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
	"time"
)

var UploadDir string = "/upload/img/"

type FishModel struct {
	*BasicModel
}

type FishList struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	Img     string `json:"img"`
	Weight  string `json:"weight"`
	Lenght  string `json:"length"`
	Id      string `json:"id"`
}

func (f *FishModel) GetFishList(pageId, pageSize, userId int) (map[string][]FishList, int, error) {
	sql := "select id, title, address, user_id, img, weight,length from fish where user_id = ? order by id desc limit ?, ? "
	res, err := f.ctx.Mysql.Query(sql, userId, f.library.Func.GetPage(pageId, pageSize), pageSize)
	if err != nil {
		return nil, 1000, errors.New("服务错误")
	}
	//fmt.Println(res)
	if len(res) == 0 {
		return nil, 1005, errors.New("暂无数据")
	}
	m := make(map[string][]FishList)
	m["fish_list"] = make([]FishList, 0, len(res))
	for _, v := range res {
		m["fish_list"] = append(m["fish_list"], FishList{
			Title:   v["title"],
			Address: v["address"],
			Img:     v["img"],
			Weight:  v["weight"],
			Lenght:  v["length"],
			Id:      v["id"],
		})
	}
	return m, 1006, nil
}

func (f *FishModel) UploadImgAndFishInfo(ctx *gin.Context, title, weight, length, address, userId string) error {
	Img, err := ctx.FormFile("imgfile")
	if err != nil {
		return errors.New("上传文件有问题")
	}
	fileExt := strings.ToLower(path.Ext(Img.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		return errors.New("上传失败!只允许png,jpg,gif,jpeg文件")
	}
	dirName, _ := os.Getwd()
	fileName := f.library.Func.Md5(fmt.Sprintf("%s%s", Img.Filename, time.Now().String()))
	fildDir := fmt.Sprintf("%s%s%d%s/", dirName, UploadDir, time.Now().Year(), time.Now().Month().String())
	isExist := f.BasicModel.library.Func.ExistsDir(fildDir)
	if !isExist {
		err = os.MkdirAll(fildDir, os.ModePerm)
		if err != nil {
			return errors.New("创建文件失败")
		}
	}
	filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, ".png")
	err = ctx.SaveUploadedFile(Img, filepath)
	if err != nil {
		return err
	}
	DBpath := fmt.Sprintf("%s%d%s/%s%s", UploadDir, time.Now().Year(), time.Now().Month().String(), fileName, ".png")
	inserSql := "insert into fish (title, address,user_id,img,weight,length) values (?,?,?,?,?,?)"
	_, err = f.ctx.Mysql.Exec(inserSql, title, address, userId, DBpath, weight, length)
	if err != nil {
		f.ctx.Log.Error("error", err)
		return err
	}

	return nil
}
