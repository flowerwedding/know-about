package personal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"strconv"
	"where-dream-continues/know-about/sqlconn"
)

func UpdateMesage(user_id,name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis string)bool{
	db:=sqlconn.Conn()
	res, err := db.Exec("update test set name=?,gender=?,background=?,headportrait=?,introduce=?,address=?,industry=?,occupation=?,education=?,synopsis=? where user_id = ? ",name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis,user_id)
	if err != nil{
		fmt.Println(err.Error())
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func See(user_id string)(string,string,string,string,string,string,string,string,string,string){
	db:=sqlconn.Conn()
	var name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis string
	err := db.QueryRow(`select name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis from test where user_id=?`, user_id).Scan(&name,&gender,&background,&headportrait,&introduce,&address,&industry,&occupation,&education,&synopsis)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return  name,gender,background,headportrait,introduce,address,industry,occupation,education,synopsis
}

func FindDynamic(user_id,value string) []string {
	db:=sqlconn.Conn()
    rows, err := db.Query(`select message from dynamic where user_id = ? and value = ? order by time desc`, user_id, value)
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var messageSlice []string
	for rows.Next() {
		var message string
		err := rows.Scan(&message)
		if err != nil {
			log.Fatal(err)
		}
		messageSlice = append(messageSlice,message)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return messageSlice
}

func FinfDynamic(user_id string) ([]string,[]string) {
	db:=sqlconn.Conn()
	rows, err := db.Query(`select message,value from dynamic where user_id = ? order by time desc`, user_id)
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var messageSlice []string
	var valueSlice []string
	for rows.Next() {
		var message string
		var value string
		err := rows.Scan(&message,&value)
		if err != nil {
			log.Fatal(err)
		}
		messageSlice = append(messageSlice,message)
		valueSlice = append(valueSlice,value)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return messageSlice,valueSlice
}

func FindQuestionll(user_id,message string) (int,string,int,string) {
	db:=sqlconn.Conn()
	var detail,picture string
	var follow,id int
	err := db.QueryRow(`select id,detail,follow,picture from question where user_id = ? and message = ?`, user_id, message).Scan(&id,&detail,&follow,&picture)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return id,detail,follow,picture
}

func FindWentill(now int)  string {
	db:=sqlconn.Conn()
	var message string
	err := db.QueryRow(`select message from question where id = ?`, now).Scan(&message)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return message
}

func FindTest(user_id string) (string,string,string,string) {
	db:=sqlconn.Conn()
	var name,headportrait,background,browse string
	err := db.QueryRow(`select name,browse,headportrait,background from test where user_id = ?`, user_id).Scan(&name,&browse,&headportrait,&background)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return name,headportrait,background,browse
}

func FindCount(user_id string) (int,int,int,int,int,int,int,int,int,int,int) {
	db:=sqlconn.Conn()
	rows, err := db.Query("select value from dynamic where user_id = ? order by time desc", user_id)
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	count_huida,count_wenti,count_wenzhang,count_zhuanlan,count_xiangfa,count_woguanzhuderen,count_woguanzhudehuati,count_woguanzhudezhuanlan,count_woguanzhudewenti,count_woguanzhudeshoucangjia,count_shoucang:=0,0,0,0,0,0,0,0,0,0,0
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
        if value =="回答"{count_huida++}else if value == "问题"{count_wenti++}else if value == "文章"{count_wenzhang++}else if value =="专栏"{count_zhuanlan++}else if value=="想法"{count_xiangfa++}else if value =="我关注的人"{count_woguanzhuderen++} else
		if value =="我关注的话题"{count_woguanzhudehuati++}else if value =="我关注的专栏"{count_woguanzhudezhuanlan++}else if value =="我关注的问题"{count_woguanzhudewenti++}else if value=="我关注的收藏夹"{count_woguanzhudeshoucangjia++}else if value=="收藏"{count_shoucang++}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return count_huida,count_wenti,count_wenzhang,count_zhuanlan,count_xiangfa,count_woguanzhuderen,count_woguanzhudehuati,count_woguanzhudezhuanlan,count_woguanzhudewenti,count_woguanzhudeshoucangjia,count_shoucang
}

func FinfCount(id string) (int,[]string) {
	db:=sqlconn.Conn()
	rows, err := db.Query("select value,user_id from dynamic where message = ? order by time desc", id)
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var user_idSlice []string
    count_guanzhuwoderen:=0
	for rows.Next() {
		var value string
		var user_id string
		err := rows.Scan(&value,&user_id)
		if err != nil {
			log.Fatal(err)
		}
		if value == "我关注的人"{count_guanzhuwoderen++;user_idSlice = append(user_idSlice,user_id)}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return count_guanzhuwoderen,user_idSlice
}

func FindMaxAnswerll(now int) Answer {
	db:=sqlconn.Conn()
	var user_id string
	var id int
	var answers Answer
	err := db.QueryRow("select  Max(praise) from answer where pid=-1 and now=?  ", now).Scan(&answers.Praise)
	err = db.QueryRow("select id,message,user_id from answer where pid=-1 and now=? and praise=? ", now,answers.Praise).Scan(&id,&answers.Message,&user_id)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	answers.UserId = User_idToName(user_id)
	answers.Comment = FindMessageByPid(id,0,now)
	return answers
}

func FindAnswerll(message string) (Answer,int) {
	db:=sqlconn.Conn()
	var answers Answer
	var id,now int
	var user_id string
	err := db.QueryRow("select user_id,id,praise,now from answer where pid=-1 and message=? ",message).Scan(&user_id,&id,&answers.Praise,&now)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	answers.UserId = User_idToName(user_id)
	answers.Comment = FindMessageByPid(id,0,now)
	return answers,now
}

func User_idToName(user_id string) string {
	db:=sqlconn.Conn()
	var name string
	err := db.QueryRow(`select name from test where user_id = ?`,user_id).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return name
}

func FindMessageByPid(pid,count,w int) int {//计算递归的评论数
	db:=sqlconn.Conn()
	rows, err := db.Query("select id from answer where pid=? and now=?", pid, strconv.Itoa(w))
	if err != nil {
		log.Fatal(err)
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if  err != nil {
			log.Fatal(err)
		}
		count++
		count= FindMessageByPid(id,count,w)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return count
}