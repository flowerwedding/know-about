package problem

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"where-dream-continues/know-about/home"
	"where-dream-continues/know-about/personal"
)

func Insert(c *gin.Context){
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	message := c.PostForm("message")
	pid := c.PostForm("pid")
	w:= c.PostForm("w")
	wNew := home.MessageToId(w,0)
	user_idNew, _ :=strconv.Atoi(user_id)
	var pidNew int
	var value string
	if pid == w { pidNew = -1;value = "回答"}else{pidNew = home.MessageToId(pid,1);value = "评论"}
	id := home.SelectIdMax(1)+1
	if MessageInsert(pidNew,id,user_idNew,wNew,message)&&home.DynamicInsert(user_id,message,value){
		if pidNew == -1 {
			c.JSON(200, gin.H{"status": http.StatusOK, "message": "写回答成功"})
		} else{
			c.JSON(200, gin.H{"status": http.StatusOK, "message": "写评论成功"})
		}
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}
}

func Delete(c *gin.Context){
	message:=c.PostForm("message")
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	if DeleteMesage(user_id,message)&&home.DynamicDelete(user_id,message){
		c.JSON(200, gin.H{"Data": "删除成功"})
	}else {
		c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"数据库删除失败"})
	}
}

func Update(c *gin.Context){
	message:=c.PostForm("message")
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	id:=c.PostForm("id")
	f:=c.PostForm("f")
	fNew, _ := strconv.Atoi(f)
	var value string
	if fNew == 0 {value="我关注的问题"}else if fNew == 1 {value= "赞"}else if fNew == 2 {value="踩"}else if fNew == -1 {value ="我关注的人"}else {value="浏览"}
	if id == user_id {
		c.JSON(500, gin.H{"status": http.StatusInternalServerError, "message": "不能赞、踩、收藏自己"})
	}else if UpdateMesage(message, id, fNew) && value == "浏览" {
		c.JSON(200, gin.H{"Data": "浏览成功"})
	}else if home.DynamicInsert(user_id,id,value) && value == "我关注的人" {
		c.JSON(200, gin.H{"Data": "关注人成功"})
	}else if UpdateMesage(message,id,fNew)&&home.DynamicInsert(user_id,message,value){
		if fNew == 0{
			c.JSON(200, gin.H{"Data": "关注问题成功"})
		}else if fNew == 1 {
			c.JSON(200, gin.H{"Data": "赞成功"})
		}else if fNew == 2 {
			c.JSON(200, gin.H{"Data": "踩成功"})
		}
	}else {
		c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"数据库删除失败"})

	}
}

type Message struct {
	Praise int
	Stamp int
	UserId int
	Message string
	time string
	ChildMessage *[]Message
}

type Answer struct {
	Praise int
	Stamp int
	UserId int
	time string
	Message string
	Comment int
}

type Question struct{
	Message string
	Detail string
	Label string
	Follow int
	Browse int
	Comment int
}

func JsonNested(messageSlice []Message) []gin.H {
	var messageJsons []gin.H
	var messageJson gin.H
	for _, messages := range messageSlice {
		message := *messages.ChildMessage
		if messages.ChildMessage != nil {
			messageJson = gin.H{
				"praise" :         messages.Praise,
				"stamp"  :         messages.Stamp,
				"user_id":         personal.User_idToName(strconv.Itoa(messages.UserId)),
				"message":         messages.Message,
				"time"   :         messages.time,
				"ChildrenMessage": JsonNested(message),
			}
		} else {
			messageJson = gin.H{
				"praise" :         messages.Praise,
				"stamp"  :         messages.Stamp,
				"user_id":         personal.User_idToName(strconv.Itoa(messages.UserId)),
				"message":         messages.Message,
				"time"   :         messages.time,
				"ChildrenMessage": "null",
			}
		}
		messageJsons = append(messageJsons, messageJson)
	}
	return messageJsons
}

func JsonNeste(answerSlice []Answer) []gin.H {
	var answerJsons []gin.H
	var answerJson gin.H
	for _, answers := range answerSlice {
		answerJson = gin.H{
			    "praise" : answers.Praise,
			    "stamp"  : answers.Stamp,
			    "user_id": personal.User_idToName(strconv.Itoa(answers.UserId)),
				"message": answers.Message,
			    "time"   : answers.time,
				"comment": answers.Comment,
		}
		answerJsons = append(answerJsons,answerJson)
	}
	return answerJsons
}