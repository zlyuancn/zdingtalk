/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/5/17
   Description :
-------------------------------------------------
*/

package zdingtalk

import (
    "encoding/json"
)

type TextMsg struct {
    Content string `json:"content"`
}

type LinkMsg struct {
    Title      string `json:"title"`
    Text       string `json:"text"`
    MessageURL string `json:"messageUrl"`
    PicURL     string `json:"picUrl,omitempty"`
}

type Markdown struct {
    Title string `json:"title"`
    Text  string `json:"text"`
}

type Button struct {
    Title     string `json:"title"`
    ActionURL string `json:"actionURL"`
}

type ActionCard struct {
    Title          string   `json:"title"`
    Text           string   `json:"text"`
    SingleTitle    string   `json:"singleTitle,omitempty"`
    SingleURL      string   `json:"singleURL,omitempty"`
    Buttons        []Button `json:"btns,omitempty"`
    BtnOrientation string   `json:"btnOrientation,omitempty"` // 0-按钮竖直排列，1-按钮横向排列
}

type FeedLinkMsg struct {
    Title      string `json:"title"`
    MessageURL string `json:"messageURL"`
    PicURL     string `json:"picURL"`
}

type FeedCard struct {
    Links []FeedLinkMsg `json:"links"`
}

type At struct {
    AtMobiles []string `json:"atMobiles,omitempty"`
    IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type Msg struct {
    MsgType    string      `json:"msgtype"`
    Text       *TextMsg    `json:"text,omitempty"`
    Link       *LinkMsg    `json:"link,omitempty"`
    Markdown   *Markdown   `json:"markdown,omitempty"`
    ActionCard *ActionCard `json:"actionCard,omitempty"`
    FeedCard   *FeedCard   `json:"feedCard,omitempty"`
    At         *At         `json:"at,omitempty"`
}

// Text消息
func NewTextMsg(text string) *Msg {
    return &Msg{
        MsgType: "text",
        Text: &TextMsg{
            Content: text,
        },
    }
}

// Link消息
func NewLinkMsg(title, text, msgurl string, picurl ...string) *Msg {
    pic := ""
    if len(picurl) > 0 {
        pic = picurl[0]
    }
    return &Msg{
        MsgType: "link",
        Link: &LinkMsg{
            Title:      title,
            Text:       text,
            MessageURL: msgurl,
            PicURL:     pic,
        },
    }
}

// Markdown消息
func NewMarkdownMsg(title, text string) *Msg {
    return &Msg{
        MsgType: "markdown",
        Markdown: &Markdown{
            Title: title,
            Text:  text,
        },
    }
}

// 整体跳转ActionCard
func NewActionCard(title, text, single_title, single_url string) *Msg {
    return &Msg{
        MsgType: "actionCard",
        ActionCard: &ActionCard{
            Title:       title,
            Text:        text,
            SingleTitle: single_title,
            SingleURL:   single_url,
        },
    }
}

// 独立跳转ActionCard
func NewCustomCard(title, text string, btns ...Button) *Msg {
    return &Msg{
        MsgType: "actionCard",
        ActionCard: &ActionCard{
            Title:   title,
            Text:    text,
            Buttons: append([]Button{}, btns...),
        },
    }
}

// 为独立跳转ActionCard添加按钮
func (m *Msg) AddButton(btns ...Button) *Msg {
    if m.ActionCard != nil {
        m.ActionCard.Buttons = append(m.ActionCard.Buttons, btns...)
    }
    return m
}

// 设置按钮为垂直排列
func (m *Msg) VerticalButton() *Msg {
    if m.ActionCard != nil {
        m.ActionCard.BtnOrientation = "0"
    }
    return m
}

// 设置按钮为水平排列
func (m *Msg) HorizontalButton() *Msg {
    if m.ActionCard != nil {
        m.ActionCard.BtnOrientation = "1"
    }
    return m
}

// FeedCard
func NewFeedCard(links ...FeedLinkMsg) *Msg {
    return &Msg{
        MsgType: "feedCard",
        FeedCard: &FeedCard{
            Links: append([]FeedLinkMsg{}, links...),
        },
    }
}

// 为FeedCard添加Link
func (m *Msg) AddLinks(links ...FeedLinkMsg) *Msg {
    if m.FeedCard != nil {
        m.FeedCard.Links = append(m.FeedCard.Links, links...)
    }
    return m
}

// at指定人, 仅支持 text 和 link 消息. 官方文档说明支持markdown, 但是实际测试并不支持markdown
func (m *Msg) AtMobiles(mobiles ...string) *Msg {
    m.At = &At{
        AtMobiles: append(([]string)(nil), mobiles...),
    }
    return m
}

// at所有人, 仅支持 text 和 link 和 markdown 消息
func (m *Msg) AtAll() *Msg {
    m.At = &At{
        IsAtAll: true,
    }
    return m
}

// 转为body
func (m *Msg) Body() []byte {
    body, _ := json.Marshal(m)
    return body
}
