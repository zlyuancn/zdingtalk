/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/17
   Description :
-------------------------------------------------
*/

package zdingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zlyuancn/zretry"

	"github.com/zlyuancn/zdingtalk/utils"
)

type DingTalk struct {
	token  string
	secret string
}

func NewDingTalk(access_token string) *DingTalk {
	return &DingTalk{
		token: access_token,
	}
}

// 设置secret
func (m *DingTalk) SetSecret(secret string) *DingTalk {
	m.secret = secret
	return m
}

// 发送一条消息
//
// https://ding-doc.dingtalk.com/doc/#/serverapi2/hoy7iv
func (m *DingTalk) Send(msg *Msg, retry_num ...int) error {
	const ApiUrl = "https://oapi.dingtalk.com/robot/send"
	params := url.Values{}
	params.Add("access_token", m.token)
	if m.secret != "" {
		timestamp, sign := m.makeSign()
		params.Add("timestamp", timestamp)
		params.Add("sign", sign)
	}

	req, err := http.NewRequest("POST", ApiUrl+"?"+params.Encode(), bytes.NewReader(msg.Body()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	attempt_count := 1
	if len(retry_num) > 0 && retry_num[0] > 0 {
		attempt_count = 1 + retry_num[0]
	}

	var send_result *SendResult
	err = zretry.DoRetry(func() error {
		result, err := m.send(req)
		if err != nil {
			return err
		}
		send_result = result
		return nil
	}, 0, int32(attempt_count), nil)
	if err != nil {
		return err
	}

	// 绝大部分情况下, 如果返回值解析成功, 但是返回了错误码, 重试也没用, 所以在重试器外部检查返回状态
	if send_result.ErrCode == 0 {
		return nil
	}
	return fmt.Errorf("<%d>: %s", send_result.ErrCode, send_result.ErrMsg)
}

func (m *DingTalk) send(req *http.Request) (*SendResult, error) {
	var result SendResult
	err := utils.Request(req, &result)
	return &result, err
}

func (m *DingTalk) makeSign() (timestamp, sha string) {
	timestamp = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	message := timestamp + "\n" + m.secret
	h := hmac.New(sha256.New, []byte(m.secret))
	h.Write([]byte(message))
	return timestamp, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type SendResult struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
