package problem

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
	"where-dream-continues/know-about/sqlconn"
)

func FindMessageByPid(pid,w,f int) []Message { //时间顺序
	db:=sqlconn.Conn()
	var rows *sql.Rows
	var err error
    if f == 0 {rows, err = db.Query("select id,message,praise,stamp,user_id,time from answer where pid=? and now=?  order by praise desc", pid,strconv.Itoa(w))
    }else {rows, err = db.Query("select id,message,praise,stamp,user_id,time from answer where pid=? and now=? order by id desc ", pid, strconv.Itoa(w))}
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var messageSlice []Message
	for rows.Next() {
		var messages Message
		err := rows.Scan(&id, &messages.Message,&messages.Praise,&messages.Stamp, &messages.UserId,&messages.time)
		if  err != nil {
			log.Fatal(err)
		}
		child:= FindMessageByPid(id,w,f)
		messages.ChildMessage = &child
		messageSlice = append(messageSlice, messages)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return messageSlice
}

func FindAnswer(w,f int) ([]Answer,int) {
	db:=sqlconn.Conn()
	var rows *sql.Rows
	var err error
	if f == 0{ rows, err = db.Query("select id,message,praise,stamp,user_id,time from answer where pid= -1 and now=? order by praise desc", w)
	}else if f == 2{ rows, err = db.Query("select id,message,praise,stamp,user_id,time from answer where pid =-1 and now=? order by praise desc limit 3", w)
	}else{rows, err = db.Query("select id,message,praise,stamp,user_id,time from answer where pid= -1 and now=? order by id desc", w)}
	if err != nil {
		fmt.Println("1:",err.Error())
		log.Fatal(err)
	}
	var id int
	var count int
	var answerSlice []Answer
	for rows.Next() {
		var answers Answer
		err := rows.Scan(&id,&answers.Message, &answers.Praise,&answers.Stamp,&answers.UserId,&answers.time)
		if err != nil {
			log.Fatal(err)
		}
		res:=FindMessageByPid(id,w,f)
		var ress []Answer
		ress= MToA(res,ress)
		answers.Comment=cap(ress)
		count++
		answerSlice = append(answerSlice,answers)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return answerSlice,count
}

func MToA(messageSlice []Message,answerSlice []Answer) []Answer {
	for _, messages := range messageSlice {
		message := *messages.ChildMessage
		var answers Answer
		answers.Message,answers.Stamp,answers.Praise,answers.UserId,answers.time=messages.Message,messages.Stamp,messages.Praise,messages.UserId,messages.time
		answerSlice = append(answerSlice,answers)
		if messages.ChildMessage != nil {
			answerSlice=MToA(message,answerSlice)
		}
	}
	return answerSlice
}

func Qsort(arr []Answer, first, last int) {
	flag := first
	left := first
	right := last
	if first >= last {
		return
	}
	for first < last {
		for first < last {
			if arr[last].Praise >= arr[flag].Praise {
				last--
				continue
			}
			arr[last], arr[flag] = arr[flag], arr[last]
			flag = last
			break
		}
		for first < last {
			if arr[first].Praise <= arr[flag].Praise {
				first++
				continue
			}
			arr[first], arr[flag] = arr[flag], arr[first]
			flag = first
			break
		}
	}
	Qsort(arr, left, flag-1)
	Qsort(arr, flag+1, right)
}

func MaxThree(arr []Answer)[]Answer {
	var answerSlice []Answer
	for i:=0; i < 3;i++ {
		var answers Answer
        answers=arr[i]
		answerSlice = append(answerSlice,answers)
	}
	return answerSlice
}

func FindQuestion(w int) Question {
	db:=sqlconn.Conn()
	var questions Question
	err := db.QueryRow(`select message,detail,label,follow,browse from question where id = ?`,w).Scan(&questions.Message,&questions.Detail,&questions.Label,&questions.Follow,&questions.Browse)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	res:=FindMessageByPid(0,w,1)
	var ress []Answer
	ress= MToA(res,ress)
	questions.Comment=cap(ress)
	return questions
}

func MessageInsert(pid int,id int,user_id int,now int,message string) bool {
	db:=sqlconn.Conn()
	timee:=time.Now().Format("2006-01-02 15:04:05")
	stmt,err := db.Prepare("insert into answer(pid,id,user_id,message,now,praise,stamp,time) values (?,?,?,?,?,0,0,?)")
	if err != nil{
		fmt.Println("2:",err.Error())
		log.Fatal(err)
		return false
	}
	defer stmt.Close()

	ret ,err := stmt.Exec(pid,id,user_id,message,now,timee)
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
	res, err := db.Exec("delete from answer where message = ?&& user_id=?",message,user_id)
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func UpdateMesage(message,id string,f int)bool{
	db:=sqlconn.Conn()
	var res sql.Result
	var err error
	if f == 2{res, err = db.Exec("update answer set stamp = stamp +1 where message = ? and user_id=?",message,id)
	}else if f == 0{res, err = db.Exec("update question set follow = follow +1 where message = ? and user_id=?",message,id)
	}else if f == 3{res, err = db.Exec("update question set browse = browse +1 where message = ? and user_id=?",message,id)
	}else if f == 4{res, err = db.Exec("update test set browse = browse +1 where user_id=?",id)
	}else{res, err = db.Exec("update answer set praise = praise +1 where message = ? and user_id=?",message,id)}
	if err != nil{
		log.Fatal(err)
		return false
	}
	if rowsAffected, err := res.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}