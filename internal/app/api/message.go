package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	Log "github.com/gowechat/example/pkg/core"
	"github.com/gowechat/example/pkg/util"
	"io/ioutil"
	"net/http"
	"time"
)

// WXMsg 微信文本消息结构体
type WXMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string

	PicUrl       string
	MediaId      string

	Label        string

	Event        string

	EventKey     string
	MsgId        int64
}

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var Msg WXMsg
	err := c.ShouldBindXML(&Msg)
	if err != nil {
		Log.Logrus().Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	Log.Logrus().Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", Msg.MsgType, Msg.Content)

	// 对接收的消息进行被动回复
	WXMsgReply(c, Msg)
}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string

	Image struct {
		MediaId  string
	}

	Label        string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName      xml.Name `xml:"xml"`
}

const TURING_CHAT_ROBOT_API = "http://www.tuling123.com/openapi/api"
const TRING_CHAT_ROBOT_KEY = "a90fdac0979a333ec36ebc25f11ee1c9"

type chatMessage struct {
	Code int
	Text string
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, Msg WXMsg) {
	//params := make(map[string]interface{})
	//params["userid"] = toUser
	//params["key"] = TRING_CHAT_ROBOT_KEY
	//params["info"] = msg
	//paramsBytes, _ := json.Marshal(params)
	if Msg.MsgType == "event" {
		if Msg.Event == "subscribe" {
			repMsg := WXRepMsg{
				ToUserName:   Msg.FromUserName,
				FromUserName: Msg.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "text",
				Content:      "欢迎加入大连民族大学服务号！！！",
			}
			replyMsg, err := xml.Marshal(repMsg)
			if err != nil {
				Log.Logrus().Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
				return
			}
			_, _ = c.Writer.Write(replyMsg)
		} else if Msg.Event == "CLICK" && Msg.EventKey == "school_bus" {
			repMsg := WXRepMsg{
				ToUserName:   Msg.FromUserName,
				FromUserName: Msg.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "image",
				Image: struct{ MediaId string }{MediaId: "DHY_L992Gr5f-9tvr9czkGaE3OH_Ae6nviNIAvi3zjE"},
			}
			replyMsg, err := xml.Marshal(repMsg)
			if err != nil {
				Log.Logrus().Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
				return
			}
			_, _ = c.Writer.Write(replyMsg)
		}
	} else if Msg.MsgType == "image" {
		repMsg := WXRepMsg{
			ToUserName:   Msg.FromUserName,
			FromUserName: Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "image",
			Image: struct{ MediaId string }{MediaId: Msg.MediaId},
		}

		replyMsg, err := xml.Marshal(repMsg)
		if err != nil {
			Log.Logrus().Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
			return
		}
		_, _ = c.Writer.Write(replyMsg)
	} else if Msg.MsgType == "location" {
		repMsg := WXRepMsg{
			ToUserName:   Msg.FromUserName,
			FromUserName: Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:        Msg.Label,
		}
		replyMsg, err := xml.Marshal(repMsg)
		if err != nil {
			Log.Logrus().Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
			return
		}
		_, _ = c.Writer.Write(replyMsg)
	} else {
		url := fmt.Sprintf(TURING_CHAT_ROBOT_API + "?userid=%s&key=%s&info=%s", Msg.FromUserName, TRING_CHAT_ROBOT_KEY, Msg.Content)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("向图灵机器人建立请求失败", err)
			Log.Logrus().Errorf("向图灵机器人建立请求失败, err=%+v", err)
			util.RenderError(c, err)
			return
		} else {
			fmt.Println("向图灵机器人建立请求成功")
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取内容失败", err)
			return
		}
		chatMsg := chatMessage{}
		err = json.Unmarshal(body, &chatMsg)
		if err != nil {
			fmt.Println("解析json失败")
			return
		}
		repMsg := WXRepMsg{
			ToUserName:   Msg.FromUserName,
			FromUserName: Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      chatMsg.Text,
		}
		replyMsg, err := xml.Marshal(repMsg)
		if err != nil {
			Log.Logrus().Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
			return
		}
		_, _ = c.Writer.Write(replyMsg)
	}
}