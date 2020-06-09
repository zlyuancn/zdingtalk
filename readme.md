
# 朴实无华的dingtalk, 支持6种消息类型, 支持secret, 支持失败重试, 两行代码就能推送, github上再也没有比这个更简单方便快捷的钉钉消息了

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
> 你在聊天的时候, 应该是写入消息和@的人然后发送. 而不是写入消息并发送后再选@的人

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