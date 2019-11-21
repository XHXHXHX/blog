package main

import (
	"blog/library/log"
	"blog/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func main() {
	defer MainPanicHandler()

	app := gin.New()
	router.InitRouter(app)

	app.Run(":8088")

}

func MainPanicHandler() {
	if err := recover(); err != nil {
		log.New().WithFields(logrus.Fields{
			"err":err,
		}).Error("main panic: %v", err)

		log.New().Error("main panic: debug stack:", string(debug.Stack()))
	}
	log.Close()
}