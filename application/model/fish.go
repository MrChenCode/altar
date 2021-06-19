package model

import (
	"errors"
)

type FishModel struct {
	*BasicModel
}

type FishList struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	Img     string `json:"img"`
	Weight  string `json:"weight"`
	Lenght  string `json:"lenght"`
	Id      string `json:"id"`
}

func (f *FishModel) GetFishList(pageId, pageSize, userId int) (map[string][]FishList, int, error) {
	sql := "select id, title, address, user_id, img, weight,lenght from fish where user_id = ? order by id desc limit ?, ? "
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
			Lenght:  v["lenght"],
			Id:      v["id"],
		})
	}
	return m, 1006, nil

}
