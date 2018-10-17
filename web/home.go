package web

import (
	"../http"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// WEB控制器
type HomeController struct {
	http.Controller
}

// Index 使用tpl展示的页面，返回view地址
func (this *HomeController) Index(ctx *http.HttpContext) (view string, data interface{}) {
	return "index", nil
}

// 返回json数据
// 注意返回值是ApiResponse
func (this *HomeController) JsonResponse(ctx *http.HttpContext) *http.ApiResponse {
	return this.Success(nil)
}

// Register 注册用户
// 注意返回值是ApiResponse
func (this *HomeController) Register(ctx *http.HttpContext) *http.ApiResponse {
	username := ctx.FormString("user_name")
	password := ctx.FormString("password")
	var err error = nil
	if username == "" || password == "" {
		err = errors.New("用户名或密码不能为空")
	}
	if err != nil {
		return this.Fail(-1, err.Error())
	} else {
		return this.Success(nil)
	}
}

// Login 用户登录
func (this *HomeController) Login(ctx *http.HttpContext) (view string, data interface{}) {
	if ctx.RawCtx.IsGet() {
		return "login", nil
	} else {
		userName := ctx.FormString("user_name")
		password := ctx.FormString("password")
		if userName == "admin" && password == "admin" {
			cookie := &fasthttp.Cookie{}
			cookie.SetKey("user_name")
			cookie.SetValue(userName)
			ctx.RawCtx.Response.Header.SetCookie(cookie)
			ctx.RawCtx.Redirect("/", 200)
			return "", nil
		} else {
			return "login", map[string]string{"err": "账号密码错误"}
		}
	}
}
