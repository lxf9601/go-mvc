package routes

import (
	"../../http"
	"../../web"
)

// 绑定路由映射表
func Bind() *http.Router {
	router := new(http.Router)
	router.Init()
	var home interface{} = new(web.HomeController)

	router.AnyAuth("/", &home, "Index")
	router.Get("/login", &home, "Login")
	router.Get("/json", &home, "JsonResponse")
	router.Post("/register", &home, "Register")
	return router
}
