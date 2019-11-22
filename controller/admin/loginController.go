package admin

import (
	"blog/library"
	"blog/library/log"
	"blog/model"
	"github.com/sirupsen/logrus"

	//"blog/library/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator-10.0.1"
)

type LoginController struct {

}

type Param struct {
	Passwd string `validate:"required"`
}

const EspeciallyHandsomeId = 1
const LoginTokenExpireHour = 2

func (ctr *LoginController) Login (c *gin.Context) {

	err, param := ctr.Validate(c)

	if err != nil {
		c.SecureJSON(200, gin.H{
			"msg": "param error",
			"err": err,
		})
		return
	}

	passwd := library.Md5Encryption(param.Passwd)
	uid, err := model.FindAdmin(passwd)
	if err != nil {
		c.SecureJSON(200, gin.H{
			"msg": "select error",
			"err": err,
		})
		return
	}
	c.SecureJSON(200, gin.H{
		"msg": "success",
		"uid": uid,
	})

	//Redis, err := redis.GetNewClient()
	//if err != nil {
	//	c.SecureJSON(200, gin.H{
	//		"msg": "system error",
	//	})
	//	return
	//}



}

func (ctr *LoginController) Validate (c *gin.Context) (error, *Param) {
	param := &Param{c.PostForm("passwd"),}
	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		return err, nil
	}

	return nil, param
}