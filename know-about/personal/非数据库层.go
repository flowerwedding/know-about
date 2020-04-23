package personal

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"where-dream-continues/know-about/picture"
)

func Update(c *gin.Context){
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	t:=c.PostForm("t")
	if t != "1"{
		name:=c.PostForm("name")
		gender:=c.PostForm("gender")
		background:=picture.Upload(c,"background")
		headportrait:=picture.Upload(c,"headportrait")
		introduce:=c.PostForm("introduce")
		address:=c.PostForm("address")
		industry:=c.PostForm("industry")
		occupation:=c.PostForm("occupation")
		education:=c.PostForm("education")
		synopsis:=c.PostForm("synopsis")
		if UpdateMesage(user_id,name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis){
			c.JSON(200, gin.H{"Data": "信息更新成功"})
		}else{
			c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"数据库更新失败"})
		}
	}else{
		name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis:=See(user_id)
		c.JSON(200,gin.H{"user_id":user_id,"name":name,"gender":gender,"background":background,"headportrait":headportrait,"introduce":introduce,"address":address,"industry":industry,"occupation":occupation,"education":education,"synopsis":synopsis})
	}
}

func Follow(messageSlice []string) []gin.H {//message是对方的id
	var messageJsons []gin.H
	var messageJson gin.H
	for _, message := range messageSlice {
		name, headportrait, _, _ := FindTest(message)
		count_guanzhuwoderen, _ := FinfCount(message)
		count_huida, _, count_wenzhang, _, _, _, _, _, _, _, _ := FindCount(message)
		messageJson = gin.H{"name": name, "headportrait": headportrait, "count_guanzhuwoderen": count_guanzhuwoderen, "count_huida": count_huida, "count_wenzhang": count_wenzhang,}
	messageJsons = append(messageJsons, messageJson)
	}
	return messageJsons
}

type Answer struct {
	Praise int
	UserId string
	Message string
	Comment int
}

func JsonNestll(messageSlice []string,user_id string,i int,value string) []gin.H {
	var messageJsons []gin.H
	var messageJson gin.H
	if i == 1 {
		for _, messages := range messageSlice {
			id, _, _, _ := FindQuestionll(user_id, messages)
			messageJson = gin.H{
				"value" : value,
				"question": messages,
				"answer":   FindMaxAnswerll(id),
			}
			messageJsons = append(messageJsons, messageJson)
		}
	}else if i == 3{
		for _, messages := range messageSlice {
			answer,now:=FindAnswerll(messages)
			_,headportrait,_,_:=FindTest(answer.UserId)
			messageJson = gin.H{
				"name"    : answer.UserId,
				"headportrait":headportrait,
				"question":FindWentill(now),
				"message":messages,
				"praise" :answer.Praise,
				"comment":answer.Comment,
				"value" : value,
			}
			messageJsons = append(messageJsons, messageJson)
		}
	}
	return messageJsons
}