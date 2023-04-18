package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gowechat/example/pkg/core"
	Log "github.com/gowechat/example/pkg/core"
)

type student struct {
	XM    string `db:"XM"`
	YWXM  string `db:"YWXM"`
	SFZH  string `db:"SFZH"`
	XB    string `db:"XB"`
	DH    string `db:"DH"`
	EMAIL string `db:"EMAIL"`
	MZMC  string `db:"MZMC"`
	KQMC  string `db:"KQMC"`
	BYZX  string `db:"BYZX"`
	GKZF  *string `db:"GKZF"`
	XSM   string `db:"XSM"`
	ZYM   string `db:"ZYM"`
	BJH   string `db:"BJH"`
	XQM   string `db:"XQM"`
}

func VerifyInfomation(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//如果账号或密码为空不执行下面的内容
	if username == "" || password == "" {
		return
	}

	//增加一步账户校验
	var pwd string
	err := core.Db.QueryRow("select MM from NEWJW.SYS_YHMMB where ZJH = :1", username).Scan(&pwd)
	if err != nil {
		Log.Logrus().Println(err)
		return
	}

	if password != pwd {
		return
	}

	sql := "select NEWJW.XS_XJB.XM, NEWJW.XS_XJB.YWXM, NEWJW.XS_XJB.SFZH, NEWJW.XS_XJB.XB," +
		"NEWJW.XS_GRXXB.DH, NEWJW.XS_GRXXB.EMAIL," +
		"NEWJW.CODE_MZB.MZMC, NEWJW.CODE_KQB.KQMC, NEWJW.XS_XJB.BYZX," +
		"NEWJW.XS_XJB.GKZF, NEWJW.CODE_XSB.XSM, NEWJW.CODE_ZYB.ZYM," +
		"NEWJW.XS_XJB.BJH, NEWJW.CODE_XAQB.XQM " +
		"from NEWJW.XS_XJB, NEWJW.XS_GRXXB, NEWJW.CODE_XSB, NEWJW.CODE_ZYB," +
		"NEWJW.CODE_MZB, NEWJW.CODE_XAQB, NEWJW.CODE_KQB " +
		"where NEWJW.XS_XJB.XH = NEWJW.XS_GRXXB.XH(+) " +
		"and NEWJW.CODE_XSB.XSH(+) = NEWJW.XS_XJB.XSH " +
		"and NEWJW.XS_XJB.ZYH = NEWJW.CODE_ZYB.ZYH(+) " +
		"and NEWJW.CODE_MZB.MZDM(+) = NEWJW.XS_XJB.MZDM " +
		"and NEWJW.CODE_XAQB.XQH(+) = NEWJW.XS_XJB.XQH " +
		"and NEWJW.CODE_KQB.KQH(+) = NEWJW.XS_XJB.KQH " +
		"and NEWJW.XS_XJB.XH = :1"
	var stu student
	err = core.Db.QueryRow(sql, username).Scan(&stu.XM, &stu.YWXM, &stu.SFZH, &stu.XB, &stu.DH, &stu.EMAIL, &stu.MZMC, &stu.KQMC, &stu.BYZX, &stu.GKZF, &stu.XSM, &stu.ZYM, &stu.BJH, &stu.XQM)
	if err != nil {
		Log.Logrus().Println(err)
		return
	}
	c.HTML(200, "verify_infomation.html", gin.H{
		"name": stu.XM,
		"department": stu.XSM,
		"major": stu.ZYM,
		"email": stu.EMAIL,
		"gender": stu.XB,
		"uid": username,
		"class_": stu.BJH,
	})
}
