package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Validator struct {

}

func (this *Validator) Handle (gin *gin.Context) {
	fmt.Println("Validator Continue")
	gin.Next()
}