package request

type UpdateSysConstReq struct {
	Params []UpdateSysConst `json:"params"`
}

type UpdateSysConst struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Value     string `json:"value"`
	ConstType string `json:"constType"`
	Remark    string `json:"remark"`
	Upload    *int   `json:"upload"`
}

type ListSysConstReq struct {
	Keys []string `json:"keys" form:"keys"`
}
