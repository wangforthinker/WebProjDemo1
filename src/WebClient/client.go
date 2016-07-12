package main

import (
	"io/ioutil"
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
)

type loginInfo struct {
	Username string `json:username`
	Passwd string `json:passwd`
}

type updatePasswdInfo struct {
	Username string `json:username`
	Passwd string `json:passwd`
	NewPasswd string `json:newPasswd`
}

//处理网络request and response
func handle_net_request(data interface{}, url string)(err error, return_type string, returnbody []byte) {
	body,err := json.Marshal(&data)
	if err != nil {
		fmt.Println("json encode error:",err.Error())
		return err,"",nil
	}
	
	rq, err:= http.NewRequest("post", url, bytes.NewReader(body))
	if err != nil {
		fmt.Println("request error:",err.Error())
		return err,"",nil
	}
	
	client := &http.Client{}
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp,err := client.Do(rq)
	
	if err != nil {
		fmt.Println("client connect error:", err.Error())
		return err,"",nil
	}
	
	defer resp.Body.Close()
	
	returntype := resp.Header.Get("returntype")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error:",err.Error())
		return err,"",nil
	}
	
	return nil,returntype,body
}

//url username passwd
func handle_register(args []string) {
	if len(args) < 3 {
		fmt.Println("args error")
		return
	}
	
	var info loginInfo
	info.Username = args[1]
	info.Passwd = args[2]
	
	_,returntype,body := handle_net_request(info, args[0])
	
	fmt.Println(returntype)

	fmt.Println(string(body))
}

//同register 一样的处理过程
func handle_login(args []string) {
	handle_register(args)
}

//args:url username passwd newpasswd
func handle_update_passwd(args []string) {
	if len(args) < 4 {
		fmt.Println("args error")
		return
	}
	
	info := updatePasswdInfo {
		Username : args[1],
		Passwd : args[2],
		NewPasswd : args[3],
	}
	
	_,returntype,body := handle_net_request(info, args[0])
	fmt.Println(returntype)
	fmt.Println(string(body))
	
}

//同register 一样的处理过程
func handle_delete_user(args []string) {
	handle_register(args)
}

//./client methods_name ....
func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("args error")
		return	
	}
	
	new_args := args[2:]
	
	switch(args[1]) {
		case "reg" :
			handle_register(new_args)
		case "login" :
			handle_login(new_args)
		case "password":
			handle_update_passwd(new_args)
		case "deleteuser" :
			handle_delete_user(new_args)
	}
	
	
}