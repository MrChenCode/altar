package model

import (
	"altar/application/config"
	"altar/application/library/util"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TopicModel struct {
	*BasicModel
}

type TopicDetail struct {
	Tid          primitive.ObjectID `bson:"_id" json:"tid"`
	Title        string             `bson:"title" json:"title"`
	Desc         string             `bson:"desc" json:"desc"`
	Type         int                `bson:"type" json:"type"`
	Image        string             `bson:"image" json:"image"`
	Icon         string             `bson:"icon" json:"icon"`
	PlatformType int                `bson:"platform_type" json:"platform_type"`
	WxqqTitle    string             `bson:"wxqq_title" json:"wxqq_title"`
	WxqqDesc     string             `bson:"wxqq_desc" json:"wxqq_desc"`
	WeiboDesc    string             `bson:"weibo_desc" json:"weibo_desc"`
	AddTime      string             `bson:"add_time" json:"add_time"`

	List interface{} `json:"list"`
}

//获取自定义专题详情
//tid 专题id，必须是有效的mongo的id
//platformType 平台类型，1-app 2-H5 会返回本平台相似的专题
func (t *TopicModel) GetTopicDetail(tid string, platformType int) (*TopicDetail, error) {
	mtid, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		return nil, err
	}
	nowDate := time.Now().Format("2006-01-02 15:04:05")
	where := bson.M{
		"_id":    mtid,
		"status": 1,
		"$or": bson.A{
			bson.D{{"is_limit", 0}},
			bson.D{
				{"is_limit", 1},
				{"start_time", bson.M{"$lte": nowDate}},
				{"end_time", bson.M{"$gte": nowDate}},
			},
		},
	}
	fields := bson.M{
		"_id": 1, "title": 1, "desc": 1, "image": 1,
		"icon": 1, "type": 1, "similarity": 1, "platform_type": 1,
		"add_time": 1, "wxqq_title": 1, "wxqq_desc": 1, "weibo_desc": 1,
	}
	opt := options.FindOne().SetProjection(fields)
	detail := &TopicDetail{}
	if err := t.ctx.Mongo.FindOne("topic_subject", where, opt).Decode(detail); err != nil {
		//没有查询到具体的mongo数据, 如果想要返回其他的错误，可以在这里指定
		if err == mongo.ErrNoDocuments {
			//return detail, nil
		}
		return nil, err
	}

	//获取专题详情
	where = bson.M{"tid": tid, "status": 1}
	detailOpt := options.Find().SetSort(bson.M{"sort_num": -1})
	cursor, err := t.ctx.Mongo.Find("topic_fragment", where, detailOpt)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.Background())
	}(cursor)

	var list []map[string]interface{}
	if err := cursor.All(context.Background(), &list); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return detail, nil
	}
	detailList := make([]map[string]interface{}, 0, len(list))
	for _, dv := range list {
		ret := make(map[string]interface{})
		ret["title"] = dv["title"]
		ret["desc"] = dv["desc"]
		ret["type"] = dv["type"]
		ret["sort_num"] = util.Int(dv["sort_num"])
		ret["add_time"] = dv["add_time"]
		if _, ok := dv["content"]; !ok {
			continue
		}
		contentData, ok := dv["content"].(primitive.A)
		if !ok {
			continue
		}
		content := make([]map[string]interface{}, 0, len(contentData))
		for _, cvt := range contentData {
			cv, ok := cvt.(map[string]interface{})
			if !ok {
				continue
			}
			switch ret["type"] {
			case "carda", "cardb", "lista", "listb", "bookwall":
				ret["show_rec"] = 1
				ret["buy"] = 1
				ret["isfree"] = 0
				if rec, ok := dv["show_rec"]; ok {
					ret["show_rec"] = util.Int(rec)
				}
				if buy, ok := dv["buy"]; ok {
					ret["buy"] = util.Int(buy)
				}
				if isfree, ok := dv["is_free"]; ok {
					ret["isfree"] = util.Int(isfree)
				}
				cont := map[string]interface{}{
					"bookid":       util.Int(cv["bookid"]),
					"bookname":     util.TrimString(cv["bookname"]),
					"booktype":     util.GetBookType(cv["bookid"]),
					"bookdesc":     util.TrimString(cv["bookdesc"]),
					"authorname":   util.TrimString(cv["authorname"]),
					"frontcover":   util.FbkImg(cv["frontcover_6"]),
					"booksize":     util.Float64(cv["booksize"]),
					"booktypename": cv["booktypename"],
					"book_score":   (float64(int(util.Float64(cv["book_score"]) * 2 * 10))) / 10,
					"recom":        cv["recommend"],
					"sort_num":     util.Int(cv["sort_num"]),
				}
				if ret["type"] == "carda" || ret["type"] == "cardb" {
					cont["chapterid"] = util.Int(cv["chapter_id"])
				}
				if bookstatus, ok := cv["bookstatus"]; ok {
					cont["bookstatus"] = bookstatus
				}
				//如果是音频书籍，则把图片替换为原图（如果原图存在的话）
				if cont["booktype"] == config.BookTypeAudio {
					if f1000, ok := cv["frontcover_1000"]; ok {
						cont["frontcover"] = util.FbkImg(f1000)
					}
				}
				content = append(content, cont)
			case "video":
				content = append(content, map[string]interface{}{
					"video_type": util.TrimString(cv["video_type"]),
					"video_url":  util.TrimString(cv["video_url"]),
					"video_pic":  util.TrimString(cv["front_cover"]),
					"recom":      util.TrimString(cv["video_recom"]),
					"sort_num":   util.Int(cv["sort_num"]),
				})
			case "news":
				content = append(content, map[string]interface{}{
					"title":    util.TrimString(cv["title"]),
					"content":  util.TrimString(cv["content"]),
					"recom":    util.TrimString(cv["info_recom"]),
					"sort_num": util.Int(cv["sort_num"]),
				})
			case "coupon":
				objID := ""
				if id, ok := dv["_id"].(primitive.ObjectID); ok {
					objID = id.Hex()
				}
				content = append(content, map[string]interface{}{
					"coupon_id": fmt.Sprintf(
						"topic_coupon_%v_%v_%v_%v",
						objID, cv["coupon_num"], cv["coupon_day"], cv["index"]),
					"coupon_name": util.TrimString(cv["coupon_name"]),
					"coupon_num":  util.Int(cv["coupon_num"]),
					"coupon_day":  util.Int(cv["coupon_day"]),
					"start_time":  util.TrimString(cv["coupon_start_time"]),
					"end_time":    util.TrimString(cv["coupon_end_time"]),
					"sort_num":    util.Int(cv["sort_num"]),
				})
			case "welfare_book":
				objID := ""
				if id, ok := dv["_id"].(primitive.ObjectID); ok {
					objID = id.Hex()
				}
				ret["day"] = util.Int(dv["day"])
				ret["start_time"] = util.TrimString(dv["start_time"])
				ret["end_time"] = util.TrimString(dv["end_time"])
				ret["welfare_book_id"] = fmt.Sprintf("topic_welfare_book_%v_%v", objID, dv["day"])
				ret["free_type"] = 1
				if freetype, ok := dv["free_type"]; ok {
					ret["free_type"] = util.Int(freetype)
				}
				frontcover := cv["frontcover_6"]
				booktype := util.GetBookType(cv["bookid"])
				if booktype == config.BookTypeAudio {
					if f1000, ok := cv["frontcover_1000"]; ok {
						frontcover = f1000
					}
				}
				frontcover = util.FbkImg(frontcover)
				content = append(content, map[string]interface{}{
					"bookid":       util.Int(cv["bookid"]),
					"bookname":     util.TrimString(cv["bookname"]),
					"booktype":     booktype,
					"bookdesc":     util.TrimString(cv["bookdesc"]),
					"authorname":   util.TrimString(cv["authorname"]),
					"frontcover":   frontcover,
					"booksize":     util.Float64(cv["booksize"]),
					"booktypename": cv["booktypename"],
					"book_score":   (float64(int(util.Float64(cv["book_score"]) * 2 * 10))) / 10,
					"recom":        util.TrimString(cv["recommend"]),
					"sort_num":     util.Int(cv["sort_num"]),
				})
			case "author":
				content = append(content, map[string]interface{}{
					"author_id":   util.Int(cv["author_id"]),
					"author_name": util.TrimString(cv["author_name"]),
					"author_desc": util.TrimString(cv["author_desc"]),
					"author_pic":  util.TrimString(cv["author_pic"]),
					"sort_num":    util.Int(cv["sort_num"]),
				})
			case "actor":
				content = append(content, map[string]interface{}{
					"actor_title":    util.TrimString(cv["actor_title"]),
					"actor_subtitle": util.TrimString(cv["actor_subtitle"]),
					"actor_desc":     util.TrimString(cv["actor_desc"]),
					"actor_pic":      util.TrimString(cv["actor_pic"]),
					"sort_num":       util.Int(cv["sort_num"]),
				})
			case "text_href":
				if cv["text"] != nil && cv["url"] != nil {
					content = append(content, map[string]interface{}{
						"text":     util.TrimString(cv["text"]),
						"url":      util.TrimString(cv["url"]),
						"sort_num": util.Int(cv["sort_num"]),
					})
				}
			case "img_href":
				if cv["img_url"] != nil && cv["url"] != nil {
					content = append(content, map[string]interface{}{
						"img":      util.TrimString(cv["img_url"]),
						"url":      util.TrimString(cv["url"]),
						"sort_num": util.Int(cv["sort_num"]),
					})
				}
			}
		}
		util.SortSliceMapStringInterface(&content, "sort_num", util.OrderDesc)
		ret["content"] = content
		detailList = append(detailList, ret)
	}
	detail.List = detailList
	return detail, nil
}
