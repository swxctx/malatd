package args

type LoginArgs struct {
	AppVer   string `query:"app_ver" json:"app_ver"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	AppVer   string `json:"app_ver"`
	Username string `json:"username"`
}
