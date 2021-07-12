package main

import (
	kruntime "dgo/framework/tools/runtime"
	kinit "dgo/work/base/initialize"
	kroute "dgo/work/base/route"
	kcode "dgo/work/code"
	kcms "dgo/work/control/cms"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"sort"
)

var (
	app *cli.App
)

func init() {
	app = cli.NewApp()
	app.UseShortOptionHandling = true
	//app.Action = geth
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2020 The dgo Authors"
	app.Commands = []cli.Command{
		{
			Name:   "admin",
			Usage:  "run admin",
			Flags:  []cli.Flag{},
			Action: admin,
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func admin(c *cli.Context) error {
	kruntime.MainGetPanic(func() {
		kruntime.Pid("admin.pid")

		kinit.InitConf("")
		kcode.InitConfParam()

		if kcode.IS_TEST_SERVER == 1 {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

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

	return nil
}
