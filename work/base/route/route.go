package route

import (
	kinit "dgo/work/base/initialize"
	kutil "dgo/work/utils"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	kmiddleware "dgo/work/base/middleware"
	"context"
	"html/template"
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type RouteWrapStruct struct {
	Method string
	Path   string
	F      func(*gin.Context)
}
type ControlInterface interface {
	Load() []RouteWrapStruct
}

//-------------------------------------------------------------------------

type RouteStruct struct {
	engine *gin.Engine
	host   string
	port   int
}

func NewRouteStruct(host string, port int) *RouteStruct {
	ts := &RouteStruct{
		host: host,
		port: port,
	}
	ts.engine = gin.New()
	ts.engine.Use( /*gin.Logger(),*/ gin.Recovery())

	ts.engine.Static("/assets", "./view/assets")
	ts.engine.Static("/upload", "./view/upload")

	store := cookie.NewStore([]byte("secret"))
	maxAge, _ := kinit.Conf.GetInt("session.max_age")
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   strings.HasPrefix("SERER HOST", "https"),
		HttpOnly: true,
	})

	sessionName, _ := kinit.Conf.GetString("session.name")
	ts.engine.Use(sessions.Sessions(sessionName, store))

	ts.engine.SetFuncMap(template.FuncMap{
		"htmlSpace": kutil.HtmlSpaceFunc,
	})

	ts.LoadHTML()
	return ts
}
func (ts *RouteStruct) LoadHTML() {
	templates := make([]string, 0)
	twoLevelTemplates, err := filepath.Glob("view/**/*.html")
	if nil != err {
		log.Fatal("load theme templates failed: " + err.Error())
	}

	oneLevelTemplates, err := filepath.Glob("view/*.html")
	if nil != err {
		log.Fatal("load theme templates failed: " + err.Error())
	}

	templates = append(templates, twoLevelTemplates...)
	templates = append(templates, oneLevelTemplates...)
	ts.engine.LoadHTMLFiles(templates...)
}

func (ts *RouteStruct) Load(control ControlInterface) {
	wps := control.Load()
	for _, v := range wps {
		switch v.Method {
		case "GET|POST":
			ts.engine.GET(v.Path, v.F)
			ts.engine.POST(v.Path, v.F)
		case "GET":
			ts.engine.GET(v.Path, v.F)
		case "POST":
			ts.engine.POST(v.Path, v.F)
		case "PUT":
			ts.engine.PUT(v.Path, v.F)
		case "PATCH":
			ts.engine.PATCH(v.Path, v.F)
		case "OPTIONS":
			ts.engine.OPTIONS(v.Path, v.F)
		case "DELETE":
			ts.engine.DELETE(v.Path, v.F)
		default:
			kinit.LogError.Println("not method :", v.Method)
		}
	}
}
func (ts *RouteStruct) SetMode(mode string) {
	gin.SetMode(mode)
}

func (ts *RouteStruct) StartPrometheus(path ...string) {
	routePath := "/metrics"
	if len(path) > 0 {
		routePath = path[0]
	}
	ts.engine.Use(kmiddleware.PromMiddleware(nil))
	ts.engine.GET(routePath, kmiddleware.PromHandler(promhttp.Handler()))
}

func (ts *RouteStruct) SetMiddleware(middleware ...gin.HandlerFunc) {
	ts.engine.Use(middleware...)
}

func (ts *RouteStruct) Run() {
	addr := fmt.Sprintf("%s:%d", ts.host, ts.port)
	server := &http.Server{
		Addr:    addr,
		Handler: ts.engine,
	}
	go func(server *http.Server) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("http listen err: " + err.Error())
		}
	}(server)

	_, _ = fmt.Fprintf(os.Stderr, "[GIN-debug] "+fmt.Sprintf("Listening and serving HTTP on %s\n", addr))
	ts.listenSignal(server)
}

func (ts *RouteStruct) listenSignal(server *http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT) // 2,3,15
	<-ch
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(cxt); err != nil {
		log.Fatal("http shutdown err: " + err.Error())
	}
}

//-------------------------------------------------------------------------

//跨域请求
func MiddlewareCrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
	}
}

// 第二个参数设置:不打印日志的路由
func MiddlewareLoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string

			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			_, _ = fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}
