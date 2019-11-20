package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HomeController struct {

}

func (this *HomeController) Detail(gin *gin.Context) {
	fmt.Println(gin.Request.Method)
	fmt.Println("I am Detail in HomeController")
}