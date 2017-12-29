package main

import (
	"github.com/Hurricanezwf/gopass/gopassd/g"
	"github.com/Hurricanezwf/gopass/gopassd/service"
)

func main() {
	if err := service.Run(); err != nil {
		g.Log.Error("Run service failed, %v", err)
		return
	}

	if err := service.PasswordMgr.Add("test", "zwf", "dfsf"); err != nil {
		g.Log.Error("Add password failed, %v", err)
		return
	}
	g.Log.Info("Add password OK")

}
