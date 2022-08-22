package rerrs

import td "github.com/swxctx/malatd"

var (
	RerrUserNotExists = td.NewRerror(20001, "用户不存在", "")
)
