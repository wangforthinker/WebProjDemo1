package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"strconv"
	"os"
)

var sqlproxy SqlProxy

type loginInfo struct {
	Username string `json:username`
	Passwd string `json:passwd`
}

type updatePasswdInfo struct {
	Username string `json:username`
	Passwd string `json:passwd`
	NewPasswd string `json:newPasswd`
}

// response format, header {returntype: int} (1,error 2,register resp 3, login resp ..)
type errorRespInfo struct {
	Msgerror string `json:msgerror`
	Syserror string `json:syserror`
}

//type regRespInfo struct
type regRespInfo struct {
	Message string `json:message`
}

//
type loginRespInfo struct {
	Message string `json:message`
}

//
type updatePasswdRespInfo struct {
	Message string `json:message`
}

func handleJsonDecodeForRq(rq *http.Request, data interface{}) error {
	body,err := ioutil.ReadAll(rq.Body)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(body, &data)
}

func handleError(w http.ResponseWriter, error_msg string,err error) {
	w.Header().Add("returntype","error") 
	var syserror string
	
	if err == nil {
		syserror = ""
	}else {
		syserror = err.Error()
	}
	
	errorinfo := errorRespInfo{
		Msgerror : error_msg,
		Syserror : syserror,
	}
	
	body,e := json.Marshal(errorinfo)
	if e != nil {
		fmt.Println("response info encode error")
	}else {
		w.Write(body)
	}
}

func handleRegister(w http.ResponseWriter, rq *http.Request){
	
	var info loginInfo
	err := handleJsonDecodeForRq(rq, &info)
	if err != nil {
		handleError(w,"decode request error",err)
		return
	}
	
//not lock ...
	b,_ :=sqlproxy.CheckDup("user","username",info.Username)
	if b == true {
		handleError(w,"username has been registered",nil)
		return
	}
	
	msg,err := sqlproxy.InsertData("user", &info)
// not lock	
	if msg != "" || err != nil {
		//fmt.Println("err : ",err.Error())
		handleError(w,msg,err)
		return
	}
	
	w.Header().Set("returntype","register_resp")
	
	reg_resp := regRespInfo {
		Message : "reg ok",
	}
	
	body,err := json.Marshal(reg_resp)
	w.Write(body)
	
}

func handleLogin(w http.ResponseWriter, rq *http.Request){
	
	var info loginInfo
	err := handleJsonDecodeForRq(rq, &info)
	if err != nil {
		handleError(w,"",err)
		return
	}
	
	b,_ := sqlproxy.CheckDup("user", "username",info.Username, "passwd", info.Passwd)
	if b != true {
		handleError(w,"username or passwd error",nil)
		return
	}
	
	w.Header().Set("returntype","login_resp")
	
	login_resp := loginRespInfo {
		Message : "login ok",
	}
	
	body,err := json.Marshal(login_resp)
	w.Write(body)
	
}

//update passwd
func handlePasswd(w http.ResponseWriter, rq *http.Request) {
	var info updatePasswdInfo
	err := handleJsonDecodeForRq(rq, &info)
	if err != nil {
		handleError(w,"",err)
		return
	}
	
	b,_ := sqlproxy.CheckDup("user", "username", info.Username, "passwd", info.Passwd)
	if b != true {
		handleError(w,"username or passwd error",nil)
		return
	}
	
	effrows,err := sqlproxy.UpdateData("username","passwd","user", info.Username, info.NewPasswd)
	if err != nil || effrows == 0{
		handleError(w, "update passwd error",err)
		return
	}
	
	w.Header().Set("returntype","password_resp")
	passwd_resp := updatePasswdRespInfo {
		Message : "update passwd ok",
	}
	
	body,err := json.Marshal(passwd_resp)
	w.Write(body)
	
}

//delete user
func handleDeleteUser(w http.ResponseWriter, rq *http.Request) {
	var info loginInfo
	err := handleJsonDecodeForRq(rq, &info)
	if err != nil {
		handleError(w,"",err)
		return
	}
	
	b,_ := sqlproxy.CheckDup("user", "username", info.Username, "passwd", info.Passwd)
	if b != true {
		handleError(w,"username or passwd error",nil)
		return
	}
	
	effrows,err := sqlproxy.DeleteData("username","user",info.Username)
	if err != nil || effrows == 0{
		handleError(w, "delete user error",err)
		return
	}
	
	w.Header().Set("returntype","delete_user_resp")
	passwd_resp := updatePasswdRespInfo {
		Message : ("delete user " + info.Username + " success"),
	}
	
	body,err := json.Marshal(passwd_resp)
	w.Write(body)
}

// ./server port database_addr database_username database_passwd database_name
func main () {
	args := os.Args
	if len(args) < 6 {
		fmt.Println("args error")
		return
	}
	
	_,err := strconv.Atoi(args[1])
	
	if err != nil {
		fmt.Println("port error:",err.Error())
		return
	}
	
	//err =sqlproxy.Connect("tcp(localhost:3306)/web1","root","****")
	err =sqlproxy.Connect(args[2],args[3],args[4],args[5])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	defer sqlproxy.Close()
	
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/password", handlePasswd)
	http.HandleFunc("/deleteuser", handleDeleteUser)
	
	
	err = http.ListenAndServe(":" + args[1], nil)
	if err != nil {
		fmt.Println("listen error:", err.Error())
	}
}