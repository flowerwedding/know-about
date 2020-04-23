package home

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime"
	"time"
	"where-dream-continues/know-about/personal"
	"where-dream-continues/know-about/sqlconn"
)

func SelectIdMax(c int)int{
	db:=sqlconn.Conn()
	var id int
	var err error
	if c == 1{err = db.QueryRow("select MAX(id) from answer").Scan(&id)}
	if c == 2{err = db.QueryRow("select MAX(id) from question").Scan(&id)}
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return id
}

func QuestionInsert(label string,id int,user_id string,message string,detail string,browse string,picture string) bool {
	db:=sqlconn.Conn()
	timee:=time.Now().Format("2006-01-02 15:04:05")
	stmt,err := db.Prepare("insert into question(label,id,user_id,message,follow,browse,time,detail,picture) values (?,?,?,?,0,?,?,?,?)")
	if err != nil{
		fmt.Println("2:",err.Error())
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(label,id,user_id,message,browse,timee,detail,picture)
	if err != nil {
		fmt.Println("3:",err.Error())
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func FindQuestion(i int) []Question {
	db:=sqlconn.Conn()
	var rows *sql.Rows
	var err error
	if i == 1{rows, err = db.Query("select message,detail,follow,browse,id from question where browse != -1 order by time desc")
	}else{rows, err = db.Query("select message,detail,follow,browse,id from question where browse != -1 order by follow+browse desc")}
    if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var id int
	var questionSlice []Question
	for rows.Next() {
		var questions Question
		err := rows.Scan(&questions.Message,&questions.Detail, &questions.Follow,&questions.Browse,&id)
		if err != nil {
			log.Fatal(err)
		}
		if i != 0 {child:= FindAnswer(id);questions.Answer = &child}
		questionSlice = append(questionSlice,questions)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return questionSlice
}

func FindAnswer(now int) personal.Answer {
	db:=sqlconn.Conn()
	var user_id string
	var id int
	var answers personal.Answer
	err := db.QueryRow("select  Max(praise) from answer where pid=-1 and now=?  ", now).Scan(&answers.Praise)
	err = db.QueryRow("select id,message,user_id from answer where pid=-1 and now=? and praise=? ", now,answers.Praise).Scan(&id,&answers.Message,&user_id)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	answers.UserId = personal.User_idToName(user_id)
	answers.Comment = personal.FindMessageByPid(id,0,now)
	return answers
}

func DynamicInsert(user_id,message,value string) bool {
	db:=sqlconn.Conn()
	timee:=time.Now().Format("2006-01-02 15:04:05")
	stmt,err := db.Prepare("insert into dynamic(user_id,message,value,time) values (?,?,?,?)")
	if err != nil{
		fmt.Println("2:",err.Error())
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(user_id,message,value,timee)
	if err != nil {
		fmt.Println("3:",err.Error())
		fmt.Println("Failed to insert, err:" ,err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func DeleteMesage(user_id,message string)bool{
	db:=sqlconn.Conn()
	res, err := db.Exec("delete from question where message = ?&& user_id=? ",message,user_id)
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func DynamicDelete(user_id,message string)bool{
	db:=sqlconn.Conn()
	res, err := db.Exec("delete from dynamic where message = ? &&user_id =?",message,user_id)
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func DeleteComment(w int)bool {
	db:=sqlconn.Conn()
	res, err := db.Exec("delete from answer where now =?", w )
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func MessageToId(message string,t int) int {
	db:=sqlconn.Conn()
	var id int
	var err error
	if t == 0 { err = db.QueryRow(`select id from question where message = ?`,message).Scan(&id)
	}else if t == 3{ err = db.QueryRow(`select now from answer where message = ?`,message).Scan(&id)
	}else{ err = db.QueryRow(`select id from answer where message = ?`,message).Scan(&id) }
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return id
}

func Guanzhunil(user_id string) []string {
	db:=sqlconn.Conn()
	rows, err := db.Query(`select user_id from test where user_id != ?`, user_id)
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var user_idSlice []string
	for rows.Next() {
		var user_id string
		err := rows.Scan(&user_id)
		if err != nil {
			log.Fatal(err)
		}
		user_idSlice = append(user_idSlice,user_id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return user_idSlice
}

func Fuzzysearch(message string) []string {
	db:=sqlconn.Conn()
	rows, err := db.Query("select message from question where message like ?", "%" +message + "%")
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