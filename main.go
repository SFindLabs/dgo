package main

import (
	kruntime "dgo/framework/tools/runtime"
	kinit "dgo/work/base/initialize"
	kroute "dgo/work/base/route"
	kcms "dgo/work/control/cms"
	"flag"
	"runtime"
)

func main() {
	types := flag.String("t", "", "启动类型，admin：启动程序")
	flag.Parse()

	if *types == "admin" {
		kruntime.MainGetPanic(func() {
			kruntime.Pid("admin.pid")

			if runtime.GOOS != "windows" {
				kinit.InitLog("admin")
			}

			kinit.InitMysql()
			kinit.InitRedis()

			host, _ := kinit.Conf.GetString("server.host")
			port, _ := kinit.Conf.GetInt("server.port")
			r := kroute.NewRouteStruct(host, port)
			//r.SetMiddleware(kroute.MiddlewareCrossDomain())
			r.SetMiddleware(kroute.MiddlewareLoggerWithWriter(kinit.LogError))

			//开启prometheus监控
			r.StartPrometheus()

			r.Load(kcms.NewCms())
			r.Load(kcms.NewKeys())

			r.Run()
		})
	}
}
