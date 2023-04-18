package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	Log "github.com/gowechat/example/pkg/core"
	"github.com/gowechat/example/pkg/util"
	"net/http"
	"strings"
)

func PushWxMenuCreate(c *gin.Context, accessToken string) error {
	menuStr := `{
		"button": [
			{
				"name": "i 民大", 
				"sub_button": [
					{
						"type": "view", 
						"name": "校园主页", 
						"url": "https://www.dlnu.edu.cn/hhh/index_mobile.htm"
					}, 
					{
						"type": "view", 
						"name": "研究生招生", 
						"url": "https://gd.dlnu.edu.cn/zs/"
					}, 
					{
						"type": "view", 
						"name": "本科招生", 
						"url": "http://zs.dlnu.edu.cn/"
					},
					{
						"type": "view", 
						"name": "就业中心", 
						"url": "http://dlnu.jysd.com/"
					},
					{
						"type": "view", 
						"name": "校园全景", 
						"url": "http://720yun.com/wx/t/9e027jpfxya?plg_nld=1&pano_id=268340&plg_auth=1&plg_uin=1&plg_usr=1&plg_vkey=1&plg_nld=1&plg_dev=1&from=singlemessage&isappinstalled=1"
					}
				]
			}, 
			{
				"name": "人才招聘", 
				"sub_button": [
					{
						"type": "view", 
						"name": "教师招聘", 
						"url": "https://mp.weixin.qq.com/s/BLQsyAzOUCtWuYF4mimk3g"
					}, 
					{
						"type": "view", 
						"name": "高层次人才招聘", 
						"url": "https://mp.weixin.qq.com/s/7LgL5s58Mb0SDvbKEqmarw"
					}
				]
			}, 
			{
				"name": "常用工具", 
				"sub_button": [
					{
						"type": "view", 
						"name": "个人中心", 
						"url": "http://stuhelp.dlnu.edu.cn/dlnu/dlnu_view"
					}, 
					{
						"type": "click", 
						"name": "班车时刻表", 
						"key": "school_bus"
					}
				]
			}
		]
	}`

	menuJsonBytes := []byte(menuStr)

	fmt.Println(accessToken)
	postReq, err := http.NewRequest("POST",
		strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/menu/create", "?access_token=", accessToken}, ""),
		bytes.NewReader(menuJsonBytes))

	if err != nil {
		fmt.Println("向微信发送菜单建立请求失败", err)
		Log.Logrus().Errorf("向微信发送菜单建立请求失败, err=%+v", err)
		util.RenderError(c, err)
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		fmt.Println("client向微信发送菜单建立请求失败", err)
		Log.Logrus().Errorf("client向微信发送菜单建立请求失败, err=%+v", err)
		util.RenderError(c, err)
		return err
	} else {
		fmt.Println("向微信发送菜单建立成功")
	}
	defer resp.Body.Close()

	return nil
}