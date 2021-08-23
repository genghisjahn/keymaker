package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var kidMap map[string]map[string]string

func init() {
	kidMap = make(map[string]map[string]string)
	kidMap["dev"] = make(map[string]string)
	kidMap["qa"] = make(map[string]string)
	kidMap["prod"] = make(map[string]string)

}

// ./makejwt.sh 8544e2f5 proto_json/wizworld.json  500 files/wizworld_qa.rsa jwts/wizworld_qa.jwt

func main() {
	content, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	keydata := []DataItem{}

	err = json.Unmarshal(content, &keydata)
	if err != nil {
		log.Fatal(err)
	}
	dirname := "jwt_" + fmt.Sprintf("%v", time.Now().Unix())
	dirErr := os.Mkdir(dirname, 777)
	if dirErr != nil {
		log.Fatal(dirErr)
	}
	for _, v := range keydata {
		cmdargs := fmt.Sprintf("%v %v %v %v %v", v.KID, v.Proto, v.Expiration, v.SecretPath, dirname+"/"+v.JWTName)
		cmd := exec.Command("/bin/sh", "makejwt.sh", cmdargs)
		cmd.Run()
	}
	time.Sleep(time.Second * 10)

}

type DataItem struct {
	KID        string `json:"kid"`
	Proto      string `json:"proto"`
	SecretPath string `json:"secret_path"`
	Expiration int    `json:"expiration"`
	JWTName    string `json:"jwt_name"`
}
