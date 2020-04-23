package lorey

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Before(c *gin.Context){
	user_id :=c.PostForm("user_id")
	password :=c.PostForm("password")
	verification := c.PostForm("verification")
	if UserEmpty(user_id){ c.JSON(200, gin.H{"status": http.StatusOK, "message": "请输入手机号"})
	}else if UserEmpty(password){c.JSON(200, gin.H{"status": http.StatusOK, "message": "请输入验证码"})
	}else if password != verification {c.JSON(200, gin.H{"status": http.StatusOK, "message": "验证码错误"})
	}else if UserExist(user_id){
		session:=sessions.Default(c)
		session.Set("user_id",user_id)
		_ = session.Save()
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "改手机号未注册，在弹窗中调用接口2"})
	}else {
		session:=sessions.Default(c)
		session.Set("user_id",user_id)
		_ = session.Save()
		c.Redirect(302, "https://www.zhihu.com")//重定向301，改成了临时重定向302
	}
}

func Registe(c *gin.Context){ //注册
	session :=sessions.Default(c)
	user_id := session.Get("username").(string)
	username := c.PostForm("name")
	password := c.PostForm("password")
	if UserSignup(username,password,user_id){
		//c.JSON(200, gin.H{"status": http.StatusOK, "message": "注册成功"})
		c.Redirect(302, "https://www.zhihu.com")
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}
}

func Login(c *gin.Context) { //登录
	user_id := c.PostForm("user_id")
	password := c.PostForm("password")
	if UserEmpty(user_id){c.JSON(200, gin.H{"status": http.StatusOK, "message": "请输入手机号或邮箱"})
	} else if UserSignin(user_id,password) {
		session:=sessions.Default(c)
		session.Set("user_id",user_id)
		_ = session.Save()
		//c.SetCookie("user_id", user_id, 10, "/", "localhost", false, true)
		//第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长；第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；第七个表示 cookie 是否支持HttpOnly属性。
		//c.JSON(200, gin.H{"status": http.StatusOK, "message": "登录成功"})
		c.Redirect(http.StatusMovedPermanently, "https://www.zhihu.com")
	} else {
		c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "账号或密码错误"})
	}
}