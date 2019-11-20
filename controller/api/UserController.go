package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserController struct {

}

func (this *UserController) Login(gin *gin.Context) {
	fmt.Println(gin.Request.Header)
	fmt.Println("I am Detail in UserController")
}
