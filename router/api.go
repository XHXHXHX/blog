package router

import (
	c "blog/controller/api"
	"blog/middleware"
	"github.com/gin-gonic/gin"
)

type api struct {
	GlobalMiddlewre []string
}

/*
 * api 下的全局中间件
 */
func (this *api) globalMiddleware() []string {
	return []string{
		"auth",
	}
}

/*
 * 路由组
 */
func (this *api) Run() []*RouteGroup {
	register := []*RouteGroup{
		RegisterRouteGet("/home/detail", new(c.HomeController), "Detail"),
		RegisterGroup("", []string{"validator"}, func(group *RouteGroup) *RouteGroup {
			group.RegisterRouteGet("/user/login", new(c.UserController), "Login")
			return group
		}),
	}



	return this.setGlobalMiddleware(register)
}

/*
 * 注册全局路由
 * 放到最前面
 */
func (this *api) setGlobalMiddleware(routeGroup []*RouteGroup) []*RouteGroup {
	global_middleware := this.globalMiddleware()
	var middleware_group [] gin.HandlerFunc
	for _, name := range global_middleware {
		if middle, ok := middleware.MiddlewareMap[name]; ok {
			middleware_group = append(middleware_group, middle)
		}
	}
	for _, group := range routeGroup {
		middleware_group_tmp := middleware_group
		middleware_group_tmp = append(middleware_group_tmp, group.MiddleWare...)
		group.MiddleWare = middleware_group_tmp
	}

	return routeGroup
}