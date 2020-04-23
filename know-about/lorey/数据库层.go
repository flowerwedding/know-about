package lorey

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"where-dream-continues/know-about/sqlconn"
)

func UserSignup(username string, passwd string,user_id string) bool {//注册是否成功
	db:=sqlconn.Conn()
	stmt,err := db.Prepare("insert into test(name,password,user_id,gender,browse,background,headportrait,introduce,address,industry,occupation,education,synopsis) values (?,?,?,?,0,0,0,0,0,0,0,0,0)")
	if err != nil{
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(username, passwd,user_id,"未知")
	if err != nil {
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func UserSignin(user_id string,password string) bool {//密码是否正确
	db:=sqlconn.Conn()
	var passwd string
	err := db.QueryRow(`select password from test where user_id = ?`,user_id).Scan(&passwd)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	//fmt.Println("password:"+password,"passwd"+passwd)
    if  passwd==password{ return true}
	return false
}

func UserExist(username string) bool {//用户是否存在
	db:=sqlconn.Conn()
	var passwd string
	err := db.QueryRow(`select password from test where user_id = ?`,username).Scan(&passwd)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
    if passwd==""{return true}
	return false
}

func UserEmpty(username string)bool { //func 输入是否为空
    if username == "" {return true}
    return false
}

/*
进入知乎
c.Redirect(http.StatusMovedPermanently, "https://www.zhihu.com")
 */