package code

import (
	kinit "dgo/work/base/initialize"
	"fmt"
)

// 配置参数
var IS_TEST_SERVER int = 0

func init() {
	var err error
	IS_TEST_SERVER, err = kinit.Conf.GetInt("server.is_test")
	if err != nil {
		fmt.Println("conf parse fail:", err)
	}

}
