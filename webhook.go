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
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/zlyuancn/zretry"
)

const DingTalkApiUrl = "https://oapi.dingtalk.com/robot/send"

type DingTalk struct {
    token  string
    secret string
    client *http.Client
}

func NewDingTalk(access_token string) *DingTalk {
    return &DingTalk{
        token:  access_token,
        client: &http.Client{},
    }
}

// 设置secret
func (m *DingTalk) SetSecret(secret string) *DingTalk {
    m.secret = secret
    return m
}

// 设置自定义http客户端
func (m *DingTalk) SetHttpClient(client *http.Client) *DingTalk {
    m.client = client
    return m
}

// 发送一条消息
func (m *DingTalk) Send(msg *Msg, retry_num ...int) error {
    params := url.Values{}
    params.Add("access_token", m.token)
    if m.secret != "" {
        timestamp, sign := m.makeSign()
        params.Add("timestamp", timestamp)
        params.Add("sign", sign)
    }

    req, err := http.NewRequest("POST", DingTalkApiUrl+"?"+params.Encode(), bytes.NewReader(msg.Body()))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    attempt_count := 1
    if len(retry_num) > 0 && retry_num[0] > 0 {
        attempt_count = 1 + retry_num[0]
    }
    return zretry.DoRetry(func() error {
        return m.send(req)
    }, 0, int32(attempt_count), nil)
}

func (m *DingTalk) send(req *http.Request) error {
    resp, err := m.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    result := new(SendResult)
    err = json.NewDecoder(resp.Body).Decode(result)
    if err != nil {
        return err
    }

    if result.ErrCode == 0 {
        return nil
    }

    return fmt.Errorf("<%d>: %s", result.ErrCode, result.ErrMsg)
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
