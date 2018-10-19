package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
)

func main() {
	/*
	json1 := `{ "openid": "OPENID", "session_key": "SESSIONKEY", }`

	json2 := `{ "openid": "OPENID", "session_key": "SESSIONKEY", "unionid": "UNIONID" }`

	json3 := `{ "errcode": 40029, "errmsg": "invalid code" }`
	*/

	/*
	type res struct {
		openid string,
		session_key string,
		errcode int,
		errmsg string,
	}
	*/


	myurl := "http://192.168.1.13:6001/login?identityType=username&identifier=test1&credential=123456"
	resp, err := http.PostForm(myurl, url.Values{})
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	strbody := string(body)

	//json str 转map
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(strbody), &dat); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(dat)
		fmt.Println(dat["code"])
		fmt.Println(dat["data"])
		fmt.Println(dat["msg"])
	}
}
