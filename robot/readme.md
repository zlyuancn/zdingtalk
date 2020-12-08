
# 朴实无华的dingtalk机器人, 支持6种消息类型, 支持secret, 支持失败重试, 两行代码就能推送, github上再也没有比这个更简单方便快捷的钉钉消息了

# 支持

- [x] text类型
- [x] link类型
- [x] markdown类型
- [x] 整体跳转ActionCard类型
- [x] 独立跳转ActionCard类型
- [x] FeedCard类型

# 示例

```go
package main

import "github.com/zlyuancn/zdingtalk/robot"

func main() {
    msg := robot.NewTextMsg("内容")
    robot.NewDingTalk("你的access_token").Send(msg)
}
```

# 使用secret

```go
robot.NewDingTalk("你的access_token").SetSecret("你的secret")
```

# 失败重试

```go
retry_num := 2 // 失败最大重试次数
robot.NewDingTalk("你的access_token").Send(msg, retry_num)
```

# at

> 为什么在msg调用at而不是在send的时候?<br>
> 你在聊天的时候, 应该是写入消息和@的人然后发送. 而不是写入消息并发送后再选@的人

```go
msg := robot.NewTextMsg("内容").AtAll() // at所有人
msg := robot.NewTextMsg("内容").AtMobiles("159xxx", "137xxx", ...) // at指定人
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