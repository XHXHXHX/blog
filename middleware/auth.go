package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Auth struct {

}

func (this *Auth) Handle(gin *gin.Context) {

	//token := gin.Query("token")
	//
	//if token != "12345" {
	//	fmt.Println("Token Error")
	//	gin.Abort()
	//	return
	//}
	fmt.Println("Auth Continue")
	gin.Next()
}
