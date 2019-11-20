package router

import (
	"fmt"
	"blog/middleware"
	"github.com/gin-gonic/gin"
	"reflect"
)

type registerRoute interface {
	Run() []*RouteGroup
}

var routeModels map[string] []*RouteGroup

type Route struct {
	Route string
	RequestMethod string
	Controller interface{}
	Method string
}

type RouteGroup struct {
	Prefix string
	MiddleWare [] gin.HandlerFunc
	Routes []*Route
}

func init() {
	routeModels = make(map[string] []*RouteGroup)
	routeModels["api"] = registerRoute(new(api)).Run()
}

func middlewareTest(gin *gin.Context) {

	fmt.Println("I am Middleware")

	gin.Next()
}

func InitRouter(router *gin.Engine) {
	for model, arr := range routeModels {
		for _, group := range arr {
			for _, routeInfo := range group.Routes {
				prefix := model + group.Prefix
				registerGinRouter(router, routeInfo, prefix, group.MiddleWare)
			}
		}
	}
}

/*
 * 注册自定义路由到 gin 的路由中
 */
func registerGinRouter(router *gin.Engine, routeInfo *Route, prefix string, handler []gin.HandlerFunc) {
	relativePath := prefix + routeInfo.Route
	cron := reflect.ValueOf(routeInfo.Controller)
	controller_handler := func(context *gin.Context){
		args := []reflect.Value{reflect.ValueOf(context)}
		cron.MethodByName(routeInfo.Method).Call(args)
	}
	handler = append(handler, controller_handler)

	switch routeInfo.RequestMethod {
		case "GET":
			router.GET(relativePath,  handler...)
		case "POST":
			router.POST(relativePath, handler...)
		case "DELETE":
			router.DELETE(relativePath, handler...)
		case "PATCH":
			router.PATCH(relativePath, handler...)
		case "PUT":
			router.PUT(relativePath, handler...)
		case "OPTIONS":
			router.OPTIONS(relativePath, handler...)
		case "HEAD":
			router.HEAD(relativePath, handler...)
		case "ANY":
			router.Any(relativePath, handler...)
		default:
	}
}

func RegisterRoute(request_method, url string, contr interface{}, method string) *RouteGroup {
	return &RouteGroup{
		Routes:[]*Route{
			&Route{
				RequestMethod: request_method,
				Route:url,
				Controller: contr,
				Method: method,
			},
		},
	}
}

func RegisterRouteGet(url string, contr interface{}, method string) *RouteGroup {
	return RegisterRoute("GET", url, contr, method)
}

func RegisterRoutePost(url string, contr interface{}, method string) *RouteGroup {
	return RegisterRoute("POST", url, contr, method)
}

func RegisterRouteAny(url string, contr interface{}, method string) *RouteGroup {
	return RegisterRoute("ANY", url, contr, method)
}

func (this *RouteGroup) RegisterRoute(request_method, url string, contr interface{}, method string) {
	this.Routes = append(this.Routes, &Route{
			RequestMethod:request_method,
			Route:url,
			Controller: contr,
			Method: method,
		},
	)
}

func (this *RouteGroup) RegisterRouteGet(url string, contr interface{}, method string) {
	this.RegisterRoute("GET", url, contr, method)
}

func (this *RouteGroup) RegisterRoutePost(url string, contr interface{}, method string) {
	this.RegisterRoute("POST", url, contr, method)
}

func (this *RouteGroup) RegisterRouteAny(url string, contr interface{}, method string) {
	this.RegisterRoute("ANY", url, contr, method)
}

func RegisterGroup(prefix string, middlewareGroup []string, callback func(group *RouteGroup) *RouteGroup) *RouteGroup {
	var middleware_group [] gin.HandlerFunc
	for _, name := range middlewareGroup {
		if middle, ok := middleware.MiddlewareMap[name]; ok {
			middleware_group = append(middleware_group, middle)
		}
	}
	group := &RouteGroup{
		Prefix:prefix,
		MiddleWare:middleware_group,
	}

	return callback(group)
}