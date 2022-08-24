package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

// CreateDoc creates a project doc.
func CreateDoc(rootGroup string) {
	xlog.Infof("Generating README.md: %s, rootGroup: %s", info.ProjPath(), rootGroup)

	os.MkdirAll(info.AbsPath(), os.FileMode(0755))
	err := os.Chdir(info.AbsPath())
	if err != nil {
		xlog.Fatalf("[micro] Jump working directory failed: %v", err)
	}

	// read temptale file
	b, err := ioutil.ReadFile(MalatdTpl)
	if err != nil {
		b = []byte(strings.Replace(__tpl__, "__PROJ_NAME__", info.ProjName(), -1))
	}

	// new project code
	proj := NewProject(b)
	proj.gen()
	proj.genAndWriteReadmeFile(rootGroup)

	// write template file
	f, err := os.OpenFile(MalatdTpl, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		xlog.Fatalf("[micro] Create files error: %v", err)
	}
	defer f.Close()
	f.Write(formatSource(b))

	xlog.Infof("Completed README.md generation by api!")

	// gen err code info
	appendErrorInfo()

	xlog.Infof("Completed README.md generation!")
}

// appendErrorInfo append err code to README.md
func appendErrorInfo() {
	ctn := `
## Error List

|Code|Message|
|------|------|
`
	var appendRow = func(code string, msg string) {
		ctn += fmt.Sprintf("|%s|%s|\n", code, msg)
	}

	b, err := ioutil.ReadFile("rerrs/rerrs.go")
	if err != nil {
		xlog.Errorf("[micro] Append error list error: %v", err)
	}

	re := regexp.MustCompile(`\(([1-9][0-9]*), ".*", `)
	a := re.FindAllStringSubmatch(string(b), -1)
	for _, row := range a {
		if len(row) < 1 {
			xlog.Errorf("[micro] err rows invalid")
			continue
		}
		rows := strings.Split(row[0], ",")
		if len(rows) != 3 {
			continue
		}
		appendRow(row[1], rows[1])
	}

	f, err := os.OpenFile("README.md", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		xlog.Errorf("[micro] Append error list error: %v", err)
	}
	defer f.Close()
	xlog.Infof("ctn-> %s", ctn)
	f.WriteString(ctn)
	xlog.Infof("Appended error list to README.md!")
}
