package main

import (
	"blog/router"
	"github.com/gin-gonic/gin"
)

func main() {

	//routeModel := route.GetRoute()
	app := gin.New()
	router.InitRouter(app)

	app.Run(":8088")
	//
	//for model, groups := range routeModel {
	//	for _, group := range groups {
	//		prefix := model + group.Prefix
	//	}
	//}

}
