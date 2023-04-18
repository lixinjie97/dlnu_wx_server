package example

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-laoji/wxbizmsgcrypt"
	"github.com/gowechat/example/config"
	"github.com/gowechat/example/internal/app/api"
	"github.com/gowechat/example/middleWare"
	"github.com/gowechat/example/pkg/core"
	Log "github.com/gowechat/example/pkg/core"
	exampleOffAcount "github.com/gowechat/example/pkg/officialAccount"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"net/http"
)

func init() {
	flag.Parse()
	core.Setup()
	core.Logrus()
}


type QueryParams struct {
	MsgSignature string `form:"msg_signature"`
	TimeStamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
	EchoStr      string `form:"echostr"`
}

const (
	TOKEN  = "cIHgt8ggFwrsecbKFAmgZNz6P"                           //修改成自己配置的token
	AESKEY = "aLRAc7CHKDCIROfcevotERmTc4lwjMCaG09eNWMKGMR"         //修改成自己的aeskey
	CORPID = "wxccf4abb85e4bddfc"                                  //修改成自己的企业ID
)

//Run 程序入口
func Run() error {
	Log.Logrus().Info("start wechat sdk example project")

	cfg := config.GetConfig()
	r := gin.Default()
	r.Use(middleWare.Cors())

	r.LoadHTMLGlob("internal/app/views/*")

	//获取wechat实例
	wc := InitWechat()

	//公众号例子相关操作
	exampleOffAccount := exampleOffAcount.NewExampleOfficialAccount(wc)
	//处理推送消息以及事件
	r.Any("/api/v1/serve", exampleOffAccount.Serve)
	//获取ak

	r.GET("/api/v1/oa/basic/get_access_token", exampleOffAccount.GetAccessToken)
	//获取微信callback IP
	r.GET("/api/v1/oa/basic/get_callback_ip", exampleOffAccount.GetCallbackIP)
	//获取微信API接口 IP
	r.GET("/api/v1/oa/basic/get_api_domain_ip", exampleOffAccount.GetAPIDomainIP)
	//清理接口调用次数
	r.GET("/api/v1/oa/basic/clear_quota", exampleOffAccount.ClearQuota)

	//获取
	//显示首页
	g := r.Group("/dlnu")
	g.StaticFS("/static", http.Dir("internal/app/static"))
	g.GET("/create_menu", exampleOffAccount.CreateMenu)
	g.GET("/wx_server", exampleOffAccount.GetCheckWeixinSign)
	g.POST("/wx_server", api.WXMsgReceive)
	g.GET("/dlnu_view", api.PersonalCenter)
	g.GET("/verify_account", func(c *gin.Context) {
		c.HTML(200, "verify_account.html", gin.H{
		})
	})
	g.GET("/detail", func(c *gin.Context) {
		c.HTML(200, "detail.html", gin.H{
		})
	})
	g.GET("/verify_infomation", api.VerifyInfomation)
	g.POST("/verify", api.Verify)
	//企业微信服务商代开发应用模板callbackurl验证
	g.GET("/callback/customized", func(c *gin.Context) {
		wxbiz := wxbizmsgcrypt.NewWXBizMsgCrypt(TOKEN,
			AESKEY,
			CORPID,
			wxbizmsgcrypt.XmlType)
		var q QueryParams
		if ok := c.Bind(&q); ok == nil {
			echoStr, err := wxbiz.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.EchoStr)
			if err != nil {
				c.JSON(http.StatusNotImplemented, gin.H{"error": err.ErrMsg})
			} else {
				c.Writer.Write(echoStr)
			}
		}
	})
	g.POST("/callback", func(c *gin.Context) {
	})
	return r.Run(cfg.Listen)
}

//Index 显示首页
func Index(c *gin.Context) {
	c.JSON(200, "index")
}

//InitWechat 获取wechat实例
//在这里已经设置了全局cache，则在具体获取公众号/小程序等操作实例之后无需再设置，设置即覆盖
func InitWechat() *wechat.Wechat {
	cfg := config.GetConfig()
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        cfg.Redis.Host,
		Password:    cfg.Redis.Password,
		Database:    cfg.Redis.Database,
		MaxActive:   cfg.Redis.MaxActive,
		MaxIdle:     cfg.Redis.MaxIdle,
		IdleTimeout: cfg.Redis.IdleTimeout,
	}
	redisCache := cache.NewRedis(redisOpts)
	wc.SetCache(redisCache)
	return wc
}
