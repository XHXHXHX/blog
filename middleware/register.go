package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	Handle(gin *gin.Context)
}

var MiddlewareMap map[string] gin.HandlerFunc

func init () {
	MiddlewareMap = make(map[string] gin.HandlerFunc)
	MiddlewareMap["auth"] = Middleware(new(Auth)).Handle
	MiddlewareMap["validator"] = Middleware(new(Validator)).Handle
}