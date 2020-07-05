package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//config 設定値
var config = RestoreJSON("config.json")

//Config 設定
type Config struct {
	Port    string `json:"port"`
	Redmine struct {
		Host   string `json:"host"`
		APIKey string `json:"api_key"`
	} `json:"redmine"`
	Backlog struct {
		Host   string `json:"host"`
		APIKey string `json:"api_key"`
	} `json:"backlog"`
	AppScript struct {
		URL            string `json:"url"`
		OAuthTokenPath string `json:"oauth_token_path"`
	} `json:"app_script"`
	Gochikuru struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	} `json:"gochikuru"`
}

//SaveJSON JSONファイルに保存する
func SaveJSON(jsonStruct interface{}, filePath string) error {
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer fp.Close()

	e := json.NewEncoder(fp)
	e.SetIndent("", "  ")
	if err := e.Encode(jsonStruct); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//RestoreJSON JSONから生成
func RestoreJSON(jsonpath string) Config {
	var jsonstruct Config
	file, err := os.Open(jsonpath)
	if err != nil {
		msg := fmt.Sprintf("%s Open error : %v", jsonpath, err)
		fmt.Println(msg)
		return jsonstruct
	}
	defer file.Close()
	d := json.NewDecoder(file)
	d.DisallowUnknownFields() // エラーの場合 json: unknown field "JSONのフィールド名"

	if err := d.Decode(&jsonstruct); err != nil && err != io.EOF {
		msg := fmt.Sprintf("%s Decode error : %v", jsonpath, err)
		fmt.Println(msg)
		return jsonstruct
	}
	return jsonstruct
}

func init() {
	// conf := Config{}
	// SaveJSON(conf, "config.json")
}
