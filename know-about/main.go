package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"where-dream-continues/know-about/home"
	"where-dream-continues/know-about/lorey"
	"where-dream-continues/know-about/personal"
	"where-dream-continues/know-about/picture"
	"where-dream-continues/know-about/problem"
)

func main(){
	router:=gin.Default()
	router.Use(cors.Default())

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("know-about-session",store))

	router.POST("/hello/first",lorey.Before)
	router.POST("/hello/next", lorey.Registe)
	router.POST("/hello/last", lorey.Login)

	router.POST("/cure/recommend",func(c *gin.Context){
		res:= home.FindQuestion(1)
		c.JSON(200,home.JsonNest(res,1))
	})
	router.POST("/cure/follow",func(c *gin.Context){
		session :=sessions.Default(c)
		user_id := session.Get("username").(string)
		messageSlice:=personal.FindDynamic(user_id,"我关注的人")
		if messageSlice != nil {c.JSON(200,home.Follow(messageSlice))
		}else{
			var mezaJsons []gin.H
			nilSlice:= home.Guanzhunil(user_id)
			for _, meza := range nilSlice {
				name, headportrait, _, _ := personal.FindTest(meza)
				count_guanzhuwoderen, _ := personal.FinfCount(meza)
				count_huida, _, _, _, _, _, _, _, _, _, _ := personal.FindCount(meza)
				if count_huida != 0{
					zhaohuida:=personal.FindDynamic(meza,"回答")
					mezaJson := gin.H{"name": name, "headportrait": headportrait, "count_guanzhuwoderen": count_guanzhuwoderen, "count_huida": count_huida,"question":personal.FindWentill(home.MessageToId(zhaohuida[0],3)),"answer":zhaohuida[0]}
					mezaJsons = append(mezaJsons, mezaJson)
				}
			}
			c.JSON(200,mezaJsons)
		}
	})
	router.POST("/cure/hostlist",func(c *gin.Context){
		res:= home.FindQuestion(0)
		c.JSON(200,home.JsonNest(res,0))
	})

	router.POST("/cure/insert",home.Insert)
	router.DELETE("/cure/delete",home.Delete)

	router.POST("/cure/fuzzysearch",func(c *gin.Context){
		session :=sessions.Default(c)
		user_id := session.Get("username").(string)
		message:=c.PostForm("message")
		t:=c.PostForm("t")
		if t == "1" {
			var questionSlice []home.Question
			searchSlice:=home.Fuzzysearch(message)
			for _, search := range searchSlice {
				var questions home.Question
				id:=home.MessageToId(search,0)
				child:= home.FindAnswer(id)
				questions.Message = search
				questions.Answer = &child
				questionSlice = append(questionSlice,questions)
			}
			c.JSON(200,home.JsonNest(questionSlice,1))
			if home.DynamicInsert(user_id,message,"历史记录"){
				c.JSON(200, gin.H{"status": http.StatusOK, "message": "搜索成功"})
			}else{
				c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
			}
		}else{
			history:=personal.FindDynamic(user_id,"历史记录")
			c.JSON(200,gin.H{"history":history})
			search:=home.Fuzzysearch(message)
			c.JSON(200,gin.H{"search":search})
		}
	})

	router.POST("/sun/findcomment",func(c *gin.Context){//查看第几个问题的第几个回答的评论
		pid:=c.PostForm("pid")
		w:=c.PostForm("w")
		f:=c.PostForm("f")
		wNew := home.MessageToId(w,0)
		var pidNew int
		if pid == w { pidNew = 0}else{pidNew = home.MessageToId(pid,1)}
		fNew, _ :=strconv.Atoi(f)
		res:= problem.FindMessageByPid(pidNew,wNew,fNew)
		var ress []problem.Answer
		ress= problem.MToA(res,ress)
		if fNew == 0 && cap(ress)>3 {
			problem.Qsort(ress,0,cap(res)-1)
			c.JSON(200,problem.JsonNeste( problem.MaxThree(ress)))
		}
		c.JSON(200,problem.JsonNested(res))
	})

	router.POST("/sun/findanswer",func(c *gin.Context){//查看第几个问题的全部回答
		w:=c.PostForm("w")
		f:=c.PostForm("f")
		wNew := home.MessageToId(w,0)
		fNew, _ :=strconv.Atoi(f)
		res,_:= problem.FindAnswer(wNew,fNew)
		c.JSON(200,problem.JsonNeste(res))
	})

	router.POST("/sun/findall",func(c *gin.Context){//查看所有
		w:=c.PostForm("w")
		wNew := home.MessageToId(w,0)
		ress:= problem.FindQuestion(wNew)
		c.JSON(200,gin.H{
			"message":ress.Message,
			"detail":ress.Detail,
			"label":ress.Label,
			"follow":ress.Follow,
			"browse":ress.Browse,
			"comment": ress.Comment,
		})
		res,_:= problem.FindAnswer(wNew,2)
		_,count:= problem.FindAnswer(wNew,1)
		//要第一个回答的关于作者：头像、昵称、回答数、文章数、关注者数，第一个回答为res[0]
		name, headportrait, _, _ := personal.FindTest(strconv.Itoa(res[0].UserId))
		count_guanzhuwoderen, _ := personal.FinfCount(strconv.Itoa(res[0].UserId))
		count_huida, _, count_wenzhang, _, _, _, _, _, _, _, _ := personal.FindCount(strconv.Itoa(res[0].UserId))
		c.JSON(200,gin.H{"name": name, "headportrait": headportrait, "count_guanzhuwoderen": count_guanzhuwoderen, "count_huida": count_huida, "count_wenzhang": count_wenzhang,
		})
		c.JSON(200,gin.H{"count":count})
		c.JSON(200,problem.JsonNeste(res))
	})

	router.POST("/sun/insert",problem.Insert)
	router.DELETE("/sun/delete",problem.Delete)
	router.POST("/sun/update",problem.Update)

	router.POST("/smile/dynamic",func(c *gin.Context){
		session :=sessions.Default(c)
		user_id := session.Get("username").(string)
		value:=c.PostForm("value")
		if value == "文章"{
			messageSlice:=personal.FindDynamic(user_id,value)
			for _, message := range messageSlice {
				name,headportrait,_,_:=personal.FindTest(user_id)
				_,detail,follow,pictures:=personal.FindQuestionll(user_id,message)
				c.JSON(200,gin.H{"name":name,"headportrait":headportrait,"message":message,"detail":detail,"follow":follow,"picture":pictures})
			}
		}else if value =="问题"{
			messageSlice:=personal.FindDynamic(user_id,value)
			c.JSON(200,personal.JsonNestll(messageSlice,user_id,1,value))
		}else if value =="回答"{
			messageSlice:=personal.FindDynamic(user_id,value)
			c.JSON(200,personal.JsonNestll(messageSlice,user_id,3,value))
		}else if value =="动态"{
			messageSlice,valueSlice:=personal.FinfDynamic(user_id)
			for i, valueNew := range valueSlice {
				if valueNew == "问题"{newSlice :=[]string{messageSlice[i]};c.JSON(200, personal.JsonNestll(newSlice,user_id,1,valueNew))
				}else if valueNew == "赞"|| valueNew== "踩"||valueNew == "回答"{newSlice :=[]string{messageSlice[i]};c.JSON(200,personal.JsonNestll(newSlice,user_id,3,valueNew))
				}else if valueNew == "我关注的问题"{c.JSON(200, gin.H{"valueNew": valueNew, "message": messageSlice[i]})}
			}
		}else if value =="收藏"{
		}else if value =="我关注的人"{
			messageSlice:=personal.FindDynamic(user_id,value)
			c.JSON(200,personal.Follow(messageSlice))
		}else if value =="关注我的人"{
			_,guanzhuwoderenSlice:=personal.FinfCount(user_id)
			c.JSON(200,personal.Follow(guanzhuwoderenSlice))
		}else {
			messageSlice:=personal.FindDynamic(user_id,value)
			for _, message := range messageSlice {
				c.JSON(200, gin.H{"value": value, "message": message})
			}
		}
	})
	router.POST("/smile/others",func(c *gin.Context){
		session :=sessions.Default(c)
		user_id := session.Get("username").(string)
		name,headportrait,background,browse:=personal.FindTest(user_id)
		count_huida,count_wenti,count_wenzhang,count_zhuanlan,count_xiangfa,count_woguanzhuderen,count_woguanzhudehuati,count_woguanzhudezhuanlan,count_woguanzhudewenti,count_woguanzhudeshoucangjia,count_shoucang:=personal.FindCount(user_id)
		count_guanzhuwoderen,_:=personal.FinfCount(user_id)
		c.JSON(200,gin.H{"name":name,"headportrait":headportrait,"background":background,"browse":browse})
		c.JSON(200,gin.H{"count_huida":count_huida,"count_wenti":count_wenti,"count_wenzhang":count_wenzhang,"count_zhuanlan":count_zhuanlan,"count_xiangfa":count_xiangfa, "count_woguanzhuderen":count_woguanzhuderen,"count_shoucang":count_shoucang,
			"count_woguanzhudehuati":count_woguanzhudehuati,"count_woguanzhudezhuanlan":count_woguanzhudezhuanlan,"count_woguanzhudewenti":count_woguanzhudewenti,"count_woguanzhudeshoucangjia":count_woguanzhudeshoucangjia,"count_guanzhuwoderen":count_guanzhuwoderen})
	})

	router.POST("/smile/update",personal.Update)
	router.POST("/smile/logout",picture.Logout)

	_ = router.Run(":8080")
}