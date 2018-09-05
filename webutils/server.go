package webutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerFunc struct {
	ApiType    int // 给每个API设置类型,可能每个类型API需要做不同处理,比如有的API需要做签名校验
	Method     string
	HandlePath string
	Handler    gin.HandlerFunc
	Hooks      []gin.HandlerFunc
}

var (
	allControllers []*ControllerFunc
	allMiddleware  []gin.HandlerFunc
)

func AddMiddleware(middleware gin.HandlerFunc) {
	allMiddleware = append(allMiddleware, middleware)
}

func AddController(ctl *ControllerFunc) {
	allControllers = append(allControllers, ctl)
}

func (c *ControllerFunc) AddHook(f gin.HandlerFunc) {
	c.Hooks = append(c.Hooks, f)
}

func ServerRun(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 添加所有的中间件处理函数
	for _, middle := range allMiddleware {
		engine.Use(middle)
	}
	for _, ctl := range allControllers {
		// 添加业务处理函数
		ctl.Hooks = append(ctl.Hooks, ctl.Handler)
		switch ctl.Method {
		// 根据方法添加URL处理函数
		case http.MethodGet:
			engine.GET(ctl.HandlePath, ctl.Hooks...)
		case http.MethodPost:
			engine.POST(ctl.HandlePath, ctl.Hooks...)
		case http.MethodHead:
			engine.HEAD(ctl.HandlePath, ctl.Hooks...)
		case http.MethodPut:
			engine.PUT(ctl.HandlePath, ctl.Hooks...)
		}
	}
	return engine.Run(addr)
}
