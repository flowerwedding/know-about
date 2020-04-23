package home

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"where-dream-continues/know-about/personal"
	"where-dream-continues/know-about/picture"
)

func Insert(c *gin.Context){
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	message := c.PostForm("message")
	label := c.PostForm("label")
	detail := c.PostForm("detail")
	browse := c.PostForm("browse")
	picturee:=picture.Upload(c,"picture")
	ni := c.PostForm("ni")
	niNew, _ :=strconv.Atoi(ni)
	var value string
	if browse == "-1"{value ="文章"}else{value ="问题"}
	id := SelectIdMax(2)+1
	if QuestionInsert(label,id,user_id,message,detail,browse,picturee)&&DynamicInsert(user_id,message,value){
		if niNew == -1 {
			c.JSON(200,gin.H{"status": http.StatusOK,"message":message,"user_id":"匿名"})
		} else{
			c.JSON(200,gin.H{"status": http.StatusOK,"message":message,"user_id": personal.User_idToName(user_id)})
		}
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}
}

func Delete(c *gin.Context){
	message:=c.PostForm("message")
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	wNew := MessageToId(message,0)
	if DeleteMesage(user_id,message)&&DynamicDelete(user_id,message)&&DeleteComment(wNew){
		c.JSON(200, gin.H{"Data": "删除成功"})
	}else {
		c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"数据库删除失败"})

	}
}

type Question struct {
	Message string
	Detail  string
	Follow  int
	Browse  int
	Answer  *personal.Answer
}

func JsonNest(questionSlice []Question,j int) []gin.H {
	var questionJsons []gin.H
	var questionJson gin.H
	if j == 0{
		for i:=0;i<3;i++ {
			questions := questionSlice[i]
			questionJson = gin.H{
				"message": questions.Message,
				"detail" : questions.Detail,
				"temperature": questions.Follow+questions.Browse,
			}
			questionJsons = append(questionJsons,questionJson)
		}
	}else{
		for _, questions := range questionSlice {
			questionJson = gin.H{
				"message": questions.Message,
				"answer" : questions.Answer,
			}
			questionJsons = append(questionJsons,questionJson)
		}
	}
	return questionJsons
}

func Follow(messageSlice []string) []gin.H {//messa是对方的id
	var messageJsons []gin.H
	var messageJson []gin.H
	for _, messa := range messageSlice {
		messageSlice,valueSlice:=personal.FinfDynamic(messa)
		for i, valueNew := range valueSlice {
			if valueNew == "问题"{
				newSlice :=[]string{messageSlice[i]}
				messageJson=personal.JsonNestll(newSlice,messa,1,valueNew)
			    if  messageJson != nil {messageJsons = append(messageJsons, messageJson[0])}
			}else if valueNew == "回答"{
				newSlice :=[]string{messageSlice[i]}
				messageJson=personal.JsonNestll(newSlice,messa,3,valueNew)
			    if  messageJson != nil {messageJsons = append(messageJsons, messageJson[0])}
			}else if valueNew == "文章"{
				name,headportrait,_,_:=personal.FindTest(messa)
				_,detail,follow,picturee:=personal.FindQuestionll("messa",messageSlice[i])
				messageJsons = append(messageJsons,gin.H{"name":name,"headportrait":headportrait,"message":messageSlice[i],"detail":detail,"follow":follow,"picture":picturee})
			}
		}
	}
	return messageJsons
}