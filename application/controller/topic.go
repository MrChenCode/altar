package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Topic struct{}

func (_ *Topic) GetTopicInfo(c *cctx.ControllerContext) {
	tid := "5a0bdc104cecc186538b4567"

	objID, _ := primitive.ObjectIDFromHex(tid)
	where := bson.M{
		"_id": objID,
		"$or": []bson.M{
			{
				"is_limit": 0,
			},
			{
				"is_limit":   1,
				"start_time": bson.M{"$lte": time.Now().Format("2006-01-02 15:04:05")},
				"end_time":   bson.M{"$gte": time.Now().Format("2006-01-02 15:04:05")},
			},
		},
	}
	fields := bson.M{
		"_id": 0, "title": 1, "desc": 1, "image": 1,
		"con": 1, "type": 1, "similarity": 1, "platform_type": 1,
		"add_time": 1, "wxqq_title": 1, "wxqq_desc": 1, "weibo_desc": 1,
	}
	opt := options.FindOne()
	opt.SetProjection(fields)
	res := c.Mongo.FindOne("topic_subject", where, opt)
	temp := make(map[string]interface{})
	if err := res.Decode(&temp); err != nil {
		c.JSON(200, gin.H{
			"code":   10000,
			"msg":    err.Error(),
			"result": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":   0,
		"msg":    "",
		"result": temp,
	})
}
