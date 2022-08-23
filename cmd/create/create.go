package create

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/swxctx/gutil"
	"github.com/swxctx/malatd/cmd/create/tpl"
	"github.com/swxctx/malatd/cmd/info"
	"github.com/swxctx/xlog"
)

// MalatdTpl template file name
const MalatdTpl = "__malatd__tpl__.go"

// MalatdGenLock the file is used to markup generated project
const MalatdGenLock = "__malatd__gen__.lock"

// CreateProject creates a project.
func CreateProject() {
	xlog.Infof("Generating project: %s", info.ProjPath())

	os.MkdirAll(info.AbsPath(), os.FileMode(0755))
	err := os.Chdir(info.AbsPath())
	if err != nil {
		xlog.Fatalf("[malatd] Jump working directory failed: %v", err)
	}

	// creates base files
	if !gutil.FileIsExist(MalatdGenLock) {
		tpl.Create()
	}

	// read temptale file
	b, err := ioutil.ReadFile(MalatdTpl)
	if err != nil {
		b = []byte(strings.Replace(__tpl__, "__PROJ_NAME__", info.ProjName(), -1))
	}

	// new project code
	proj := NewProject(b)
	proj.Generator()

	// write template file
	f, err := os.OpenFile(MalatdTpl, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		xlog.Fatalf("[malatd] Create files error: %v", err)
	}
	defer f.Close()
	f.Write(formatSource(b))

	tpl.RestoreAsset("./", MalatdGenLock)

	xlog.Infof("Completed code generation!")
}
