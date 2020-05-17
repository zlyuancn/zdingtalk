
# 朴实无华的dintalk, 支持6种消息类型, 支持secret, 支持失败重试, 两行代码就能推送, github上再也没有比这个更简单的钉钉消息了

# 示例

```go
package main

import "github.com/zlyuancn/zdingtalk"

func main() {
    msg := zdingtalk.NewTextMsg("内容")
    zdingtalk.NewDingTalk("你的access_token").Send(msg)
}
```

# 使用secret

```go
zdingtalk.NewDingTalk("你的access_token").SetSecret("你的secret")
```

# 使用自定义http客户端

```go
client := &http.Client{}
zdingtalk.NewDingTalk("你的access_token").SetHttpClient(client)
```

# 失败重试

```go
retry_num := 2 // 失败最大重试次数
zdingtalk.NewDingTalk("你的access_token").Send(msg, retry_num)
```

# at

> 为什么在msg调用at而不是在send的时候?<br>
> 我们认为是我要将某条消息传达给某人然后发送它, 而不是我要传递某条消息然后给某人.<br>
> 比如你在在聊天的时候, 是写入消息和at的人, 然后发送. 而不是写入消息并发送然后选at的人

```go
msg := zdingtalk.NewTextMsg("内容").AtAll() // at所有人
msg := zdingtalk.NewTextMsg("内容").AtMobiles("159xxx", "137xxx", ...) // at指定人
```

# 不同的消息类型

```go
NewTextMsg(text string)
NewLinkMsg(title, text, msgurl string, picurl ...string)
NewMarkdownMsg(title, text string)
NewActionCard(title, text, single_title, single_url string)
NewCustomCard(title, text string, btns ...Button)
NewFeedCard(links ...FeedLinkMsg)
```