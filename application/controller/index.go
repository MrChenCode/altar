package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

type Index struct{}

func (_ *Index) Index(ctx *cctx.ControllerContext) {
	ctx.JSON(200, gin.H{"msg": "welcome to altar"})
}
func (_ *Index) Draw(ctx *cctx.ControllerContext) {
	class, _ := strconv.Atoi(ctx.Query("class"))
	switch class {
	case 1:
		class = 0
	case 2:
		class = 1
	case 5:
		class = 2
	default:
		class = 0
	}
	name := getStudentName(class)
	subject := getSubjectNum()
	ctx.Log.Info("中奖同学", name, "中奖题号", subject)

	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name": name,
		"subject": subject,
	})

	//ctx.JSON(200, gin.H{"同学名字": name, "题目序号": subject})
}

func getSubjectNum () int {
	num := rand.Intn( 9)
	return num + 1
}


func getStudentName(i int) string {
	names := [][]string{
		{
			"林欣雨",
			"宋恩龙",
			"李文博",
			"刘厶睿",
			"李雨彤",
			"杨悦",
			"姜博轩",
			"周志文",
			"赵鑫宝",
			"秦浩翔",
			"张宇坤",
			"单一晨",
			"杨鹏程",
			"张钟升",
			"吕佳莹",
			"臧家鹏",
			"卢研松",
			"李方桐",
			"于欣桐",
			"董欣蕊",
			"尤嘉浩",
			"杨明宇",
			"任海鸥",
			"刘歆蕊",
			"张烁炎",
			"张佳硕",
			"海兴锐",
			"丁胜楠",
			"李卓航",
			"邵思莹",
			"王建辉",
			"杨思影",
			"徐智恩",
			"辛瑀",
			"常佳乐",
			"辛依然",
			"吴瑞希",
			"李佳昊",
			"田家旺",
			"纪雨柔",
			"任兴龙",
			"吴承章",
			"李士济",
			"孙琸然",
			"周允涵",
			"袁新博",
			"袁野",
			"夏中博",
			"吴铭玉",
			"闫思琪",
			"张宇轩",
			"吴欣然",
			"冯椿媛",
			"谷欣怡",
			"何晶晶"},
		{
			"韩文康",
			"于然",
			"王宇哲",
			"宋依辉",
			"王猛",
			"王衍斌",
			"王子轩",
			"徐景奇",
			"孙禹尚昊",
			"奚楚贺",
			"韩河金子",
			"孙靖婷",
			"房梓睿",
			"邹畅",
			"尹元泽",
			"杨子睿",
			"时雨涵",
			"吕佳畅",
			"王海龙",
			"李月航",
			"李嘉鑫",
			"张鑫悦",
			"张莹月",
			"张琪",
			"冯任举",
			"梁馨予",
			"李志娜",
			"郑涵予",
			"顾桐阁",
			"许圣涵",
			"王麒智",
			"尹佳馨",
			"孟令达",
			"张广慧",
			"陈薪同",
			"程程",
			"孟美琪",
			"金禹含",
			"苏浩宁",
			"徐贺",
			"刘伟鑫",
			"黄境",
			"杨奕博",
			"迟浩森",
			"候金鹏",
			"蔺雷",
			"牛鑫月",
			"张新悦",
			"雍思睿",
			"彭天瑞",
			"王紫涵",
			"户泽伟",
			"丛山",
			"宿婉茹",
		},
		{
			"李超",
			"张东旭",
			"赵雨可轩",
			"宋慧子",
			"郑玉桐",
			"谷雨舒",
			"刘伟池",
			"李珍妮",
			"周金宇",
			"关佳欣",
			"左婉婷",
			"姜悦",
			"韩天泽",
			"魏秋蕊",
			"郑佳奇",
			"林昊",
			"李东昊",
			"焦婉婷",
			"孙浩淳",
			"乔妍",
			"邹佳露",
			"刘鑫宇",
			"姜智勇",
			"王野",
			"宋思琪",
			"韩英棋",
			"蔺欣然",
			"张坤",
			"李浩然",
			"丛耀辉",
			"苑雅馨",
			"杨文文",
			"朱子涵",
			"朱禹",
			"薛佳欣",
			"林天琦",
			"梁鑫宇",
			"连永森",
			"李金虹",
			"崔文吉",
			"李名可旭",
			"刘梓涵",
			"王忠娜",
			"许雯萱",
			"王瑞",
			"张紫钊",
			"张思彤",
			"周梦涵",
			"李润泽",
			"成蕾",
			"李博",
			"田忠凱",
			"王欣然",
			"张安楠",
			"王艺梦",
			"卢宇彤",
			"王俊杰",
			"李文博",
			"于可欣",
			"龙立辉",
		},


	}
	var randNum int
	if i == 0 {
		randNum = rand.Intn(54)
	} else if i == 1 {
		randNum = rand.Intn(53)
	} else {
		randNum = rand.Intn(59)
	}
	return names[i][randNum]
}
