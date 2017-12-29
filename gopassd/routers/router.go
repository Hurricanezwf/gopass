package router

import (
	. "github.com/Hurricanezwf/gopass/gopassd/controllers"
	"net/http"
)

func init() {
	http.HandleFunc("/pass/list", ListPass)
	http.HandleFunc("/pass/add", AddPass)
	http.HandleFunc("/pass/del", DelPass)
	http.HandleFunc("/pass/get", GetPass)
	http.HandleFunc("/pass/modify", ModifyPass)

	// TODO: password sync
}
