/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/12/8
   Description :
-------------------------------------------------
*/

package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var defaultHttpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// 设置http请求客户端
func SetHttpClient(client *http.Client) {
	defaultHttpClient = client
}

// 请求并将结果写入结构体result中
func Request(req *http.Request, result interface{}) error {
	resp, err := defaultHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %s", err)
	}
	defer resp.Body.Close()

	err = jsoniter.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("解析返回结果失败: %s", err)
	}
	return nil
}
