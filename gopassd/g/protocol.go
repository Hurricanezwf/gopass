package g

type CommonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type AddPassReq struct {
	PassName string `json:"passname"` // 密码的名字
	PassVal  string `json:"passval"`  // 密码的值
}

type AddPassResp struct {
	CommonResp
}
