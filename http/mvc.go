// mvc核心文件
package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"../logc"

	"github.com/valyala/fasthttp"
	"cparrow.com/go-mvc/util"
)

const (
	CONTENT_TYPE_HTML = "text/html; charset=utf-8"
	CONTENT_TYPE_JSON = "application/json;  charset=utf-8"
)

// JSON响应格式
type ApiResponse struct {
	Ret  int         `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// html模板的数据结构
type TplModel map[string]interface{}

type Controller struct {
}

// 返回成功响应
func (this *Controller) Success(data interface{}) *ApiResponse {
	res := ApiResponse{}
	res.Msg = ""
	res.Data = data
	return &res
}

// 返回失败响应
func (this *Controller) Fail(ret int, msg string) *ApiResponse {
	res := ApiResponse{}
	res.Ret = ret
	res.Msg = msg
	return &res
}

// Http上下文
type HttpContext struct {
	RawCtx      *fasthttp.RequestCtx
	contentType string
}

// 获取Http的ContentType
func (this *HttpContext) GetContentType() string {
	return this.contentType
}

// 设置Http的ContentType
func (this *HttpContext) SetContentType(contentType string) {
	this.contentType = contentType
	this.RawCtx.Response.Header.Set("Content-Type", contentType)
}

// 获取表单字段并转化JSON字符串为map
func (this *HttpContext) FormJSON(key string) map[string]interface{} {
	var obj interface{}
	json.Unmarshal(this.RawCtx.FormValue(key), &obj)
	return obj.(map[string]interface{})
}

// 获取表单字段并转为字符串Slice
func (this *HttpContext) FormStringSlice(key string) []string {
	str := this.FormString(key)
	if str != "" {
		return strings.Split(this.FormString(key), ",")
	} else {
		return make([]string, 0, 0)
	}
}

// 判断表单是否存在key
func (this *HttpContext) FormKeyExists(key string) bool {
	if len(this.RawCtx.FormValue(key)) == 0 {
		return false
	} else {
		return true
	}
}

// 获取表单的字符串类型字段
func (this *HttpContext) FormString(key string) string {
	return string(this.RawCtx.FormValue(key))
}

// 获取表单的布尔类型字段
func (this *HttpContext) FormBool(key string) bool {
	i, err := strconv.ParseBool(string(this.RawCtx.FormValue(key)))
	if err != nil {
		return false
	} else {
		return i
	}
}

// 获取表单的无符号整型字段
func (this *HttpContext) FormUint(key string) uint {
	i, err := strconv.ParseUint(string(this.RawCtx.FormValue(key)), 10, 32)
	if err != nil {
		return uint(0)
	} else {
		return uint(i)
	}
}

// 获取表单的整型字段
func (this *HttpContext) FormInt(key string) int {
	i, err := strconv.Atoi(string(this.RawCtx.FormValue(key)))
	if err != nil {
		return 0
	} else {
		return i
	}
}

// mvc路由器
type Router struct {
	routerMap map[string]*RouterLocation
}

// 路由定位
type RouterLocation struct {
	Controller *interface{}
	Handler    string
	Method     string
	IsAuth     bool
}

func (this *Router) Init() {
	this.routerMap = make(map[string]*RouterLocation, 100)
}

func (this *Router) Match(url string) *RouterLocation {
	return this.routerMap[url]
}

func (this *Router) Get(url string, controller *interface{}, handler string) {
	this.routerMap[url] = &RouterLocation{Controller: controller, Handler: handler, Method: http.MethodGet}
}

func (this *Router) Post(url string, controller *interface{}, handler string) {
	this.routerMap[url] = &RouterLocation{Controller: controller, Handler: handler, Method: http.MethodPost}
}

func (this *Router) Any(url string, controller *interface{}, handler string) {
	this.routerMap[url] = &RouterLocation{Controller: controller, Handler: handler, Method: http.MethodPost}
}

func (this *Router) AnyAuth(url string, controller *interface{}, handler string) {
	this.routerMap[url] = &RouterLocation{Controller: controller, Handler: handler, Method: http.MethodPost, IsAuth: true}
}

// 模板格式化日期函数
func ShowTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// http协议处理函数（MVC核心）
func HttpHandler(appPath string, router *Router) func(ctx *fasthttp.RequestCtx) {
	return func (ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				logc.Error(err)
				logc.Error(string(util.PanicTrace(5)))
			}
		}()
		if bytes.HasPrefix(ctx.Path(), []byte("/static")) {
			fs := &fasthttp.FS{
				Root:               appPath + "views",
				IndexNames:         []string{"index.html"},
				GenerateIndexPages: true,
				Compress:           true,
				AcceptByteRange:    true,
			}
			fs.PathRewrite = fasthttp.NewPathSlashesStripper(1)
			fs.NewRequestHandler()(ctx)
			return
		}
		n := router.Match(string(ctx.Path()))
		if n == nil {
			ctx.Write(ctx.Path())
			return
		}
		if n.IsAuth {
			// 可以根据需要进行授权方式的修改
			if len(ctx.Request.Header.Cookie("user_name")) == 0 {
				ctx.Redirect("/login", 200)
				return
			}
		}
		v := reflect.ValueOf(*n.Controller)
		m := v.MethodByName(n.Handler)
		if m.IsNil() {
			// 返回404页面
			ctx.NotFound()
			return
		}
		params := make([]reflect.Value, 1)
		c := new(HttpContext)
		c.RawCtx = ctx
		params[0] = reflect.ValueOf(c)
		vl := m.Call(params)
		if len(vl) > 0 {
			if vl[0].Type().String() != "string" {
				if c.GetContentType() != CONTENT_TYPE_HTML {
					// json响应内容
					c.SetContentType(CONTENT_TYPE_JSON)
					j, _ := json.Marshal(vl[0].Interface())
					if len(j) > 1024 {
						ctx.Response.Header.Add("Content-Encoding", "gzip")
						w := gzip.NewWriter(ctx.Response.BodyWriter())
						defer w.Close()
						w.Write(j)
					} else {
						ctx.Write(j)
						if logc.IsDebug() {
							logc.Debug(ctx.Request.URI().String() + ">>" + string(j))
						}
					}
				}
			} else {
				// http响应内容
				c.SetContentType(CONTENT_TYPE_HTML)
				view, _ := vl[0].Interface().(string)
				if view != "" {
					tplPath := appPath + "views/" + view + ".tpl"
					f, err := os.Open(tplPath)
					if err == nil {
						defer f.Close()
						t := template.New("").Funcs(template.
						FuncMap{"ShowTime": ShowTime})
						t, err = t.ParseGlob(path.Join(appPath+"views/common", "*.tpl"))
						t, err = t.ParseGlob(path.Join(appPath+"views/", "*.tpl"))
						err = t.ExecuteTemplate(ctx.Response.BodyWriter(), view+".tpl", vl[1].Interface())
					}
				}
			}

		}
	}

}
