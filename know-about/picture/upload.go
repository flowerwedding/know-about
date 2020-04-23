package picture

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Upload(c *gin.Context,t string)string {
	file, _:= c.FormFile(t)
	log.Println(file.Filename)
	filename := file.Filename
	if filename[len(filename)-4:]!=".jpg"{
		fmt.Println(filename[len(filename)-4:])
		c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"格式错误"})
	}else {
		timee := time.Now().Unix()
		filename = "./imagine/" + strconv.FormatInt(timee, 10) + ".jpg"
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			fmt.Println(err.Error())
		}
		return filename
	}
	return "null"
}

func Logout(c *gin.Context){
	session:=sessions.Default(c)
	session.Set("user_id","~")
	_ = session.Save()
	c.String(401,"byebye")
}