package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator-10.0.1"
)

type HomeController struct {

}

type param struct {
	Word string `validate:"required"`
}

func (ctr *HomeController) Detail(c *gin.Context) {

	validate := validator.New()
	err := validate.Struct(&param{c.Query("word")})

	if err != nil {
		c.SecureJSON(402, gin.H{
			"msg": "param error",
		})
		return
	}

	c.SecureJSON(200, gin.H{
		"res": "success",
	})
	return
}