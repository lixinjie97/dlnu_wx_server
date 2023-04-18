package officialAccount

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gowechat/example/config"
	"github.com/gowechat/example/internal/app/api"
	Log "github.com/gowechat/example/pkg/core"
	"github.com/gowechat/example/pkg/util"
	"sort"
)

//GetAccessToken 获取ak
func (ex *ExampleOfficialAccount) GetAccessToken(c *gin.Context) {
	ak, err := ex.officialAccount.GetAccessToken()
	if err != nil {
		Log.Logrus().Errorf("get ak error, err=%+v", err)
		util.RenderError(c, err)
		return
	}
	util.RenderSuccess(c, ak)
}

//CreateMenu 创建菜单
func (ex *ExampleOfficialAccount) CreateMenu(c *gin.Context) {
	ak, err := ex.officialAccount.GetAccessToken()
	if err != nil {
		Log.Logrus().Errorf("get ak error, err=%+v", err)
		util.RenderError(c, err)
		return
	}
	util.RenderSuccess(c, ak)

	err = api.PushWxMenuCreate(c, ak)
	if err != nil {
		Log.Logrus().Errorf("create menu error, err=%+v", err)
		util.RenderError(c, err)
	}
}

//GetCallbackIP 获取微信callback IP地址
func (ex *ExampleOfficialAccount) GetCallbackIP(c *gin.Context) {
	ipList, err := ex.officialAccount.GetBasic().GetCallbackIP()
	if err != nil {
		Log.Logrus().Errorf("GetCallbackIP error, err=%+v", err)
		util.RenderError(c, err)
		return
	}
	util.RenderSuccess(c, ipList)
}

//GetAPIDomainIP 获取微信callback IP地址
func (ex *ExampleOfficialAccount) GetAPIDomainIP(c *gin.Context) {
	ipList, err := ex.officialAccount.GetBasic().GetAPIDomainIP()
	if err != nil {
		Log.Logrus().Errorf("GetAPIDomainIP error, err=%+v", err)
		util.RenderError(c, err)
		return
	}
	util.RenderSuccess(c, ipList)
}

//GetAPIDomainIP  清理接口调用次数
func (ex *ExampleOfficialAccount) ClearQuota(c *gin.Context) {
	err := ex.officialAccount.GetBasic().ClearQuota()
	if err != nil {
		Log.Logrus().Errorf("ClearQuota error, err=%+v", err)
		util.RenderError(c, err)
		return
	}
	util.RenderSuccess(c, "success")
}

func (ex *ExampleOfficialAccount) GetCheckWeixinSign(c *gin.Context) {
	globalCfg := config.GetConfig()
	token := globalCfg.Token//自己填的token
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		c.Writer.Write([]byte(echostr))

	} else {
		Log.Logrus().Println("微信API验证失败")
	}
}