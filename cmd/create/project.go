package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/swxctx/gutil"
	"github.com/swxctx/malatd/cmd/info"
	"github.com/swxctx/xlog"
)

type (
	// Project project Information
	Project struct {
		*tplInfo
		codeFiles    map[string]string
		Name         string
		ImprotPrefix string
	}
	Model struct {
		*structType
		ModelStyle       string
		PrimaryFields    []*field
		UniqueFields     []*field
		Fields           []*field
		IsDefaultPrimary bool
		Doc              string
		Name             string
		SnakeName        string
		LowerFirstName   string
		LowerFirstLetter string
	}
)

// NewProject new project.
func NewProject(src []byte) *Project {
	p := new(Project)
	p.tplInfo = newTplInfo(src).Parse()
	p.Name = info.ProjName()
	p.ImprotPrefix = info.ProjPath()
	p.codeFiles = make(map[string]string)
	for k, v := range tplFiles {
		p.codeFiles[k] = v
	}
	for k := range p.codeFiles {
		p.fillFile(k)
	}
	return p
}

func (p *Project) fillFile(k string) {
	v, ok := p.codeFiles[k]
	if !ok {
		return
	}
	v = strings.Replace(v, "${import_prefix}", p.ImprotPrefix, -1)
	switch k {
	case "main.go", "config.go", "config/config.yaml":
		p.codeFiles[k] = v
	case "logic/tmp_code.gen.go":
		p.codeFiles[k] = "// Code generated by 'malatd gen' command.\n" +
			"// The temporary code used to ensure successful compilation!\n" +
			"// When the project is completed, it should be removed!\n\n" + v
	default:
		p.codeFiles[k] = "// Code generated by 'malatd gen' command.\n// DO NOT EDIT!\n\n" + v
	}
}

func mustMkdirAll(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		xlog.Fatalf("[malatd] %v", err)
	}
}

func hasGenSuffix(name string) bool {
	switch name {
	case "README.md", ".gitignore", "main.go", "config.go", "args/const.go",
		"args/var.go", "args/type.go", "api/handler.go", "api/router.go", "config/config.yaml":
		return false
	default:
		return true
	}
}

func (p *Project) Generator() {
	p.gen()
	// make all directorys
	mustMkdirAll("args")
	mustMkdirAll("api")
	mustMkdirAll("config")
	mustMkdirAll("logic")
	// write files
	for k, v := range p.codeFiles {
		if gutil.FileIsExist(k) && !hasGenSuffix(k) {
			continue
		}
		realName := info.ProjPath() + "/" + k
		f, err := os.OpenFile(k, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			xlog.Fatalf("[malatd] create %s error: %v", realName, err)
		}
		if k != "config/config.yaml" {
			b := formatSource(gutil.StringToBytes(v))
			f.Write(b)
		} else {
			f.Write(gutil.StringToBytes(v))
		}
		f.Close()
		fmt.Printf("generate %s\n", realName)
	}
}

// generate all codes
func (p *Project) gen() {
	p.genMainFile()
	p.genConstFile()
	p.genTypeFile()
	p.genRouterFile()
	p.genHandlerFile()
	p.genLogicFile()
}

func (p *Project) genAndWriteReadmeFile(rootGroup string) {
	f, err := os.OpenFile("./README.md", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		xlog.Fatalf("[malatd] create README.md error: %v", err)
	}
	f.WriteString(p.genReadme(rootGroup))
	f.Close()
	fmt.Printf("generate %s\n", info.ProjPath()+"/README.md")
}

func commentToHtml(txt string) string {
	return strings.TrimLeft(strings.Replace(txt, "// ", "<br>", -1), "<br>")
}

func (p *Project) genReadme(rootGroup string) string {
	var text string
	text += commentToHtml(p.tplInfo.doc)
	text += "\n"
	text += "## API Desc\n\n"
	for _, h := range p.tplInfo.HandlerList() {
		text += fmt.Sprintf("### %s\n\n%s\n\n", h.fullName, p.handlerDesc(h, rootGroup))
	}
	r := strings.Replace(__readme__, "${PROJ_NAME}", info.ProjName(), -1)
	r = strings.Replace(r, "${readme}", text, 1)
	return r
}

func (p *Project) handlerDesc(h *handler, rootGroupArg string) string {
	rootGroup := gutil.FieldSnakeString(p.Name)
	if len(rootGroupArg) > 0 {
		rootGroup = rootGroupArg
	}
	uri := path.Join("/", rootGroup, h.uri)
	var text string
	text += commentToHtml(h.doc) + "\n"

	text += fmt.Sprintf("- URI: `%s`\n", uri)

	var method = func(name string, txt string) {
		fields, _ := p.tplInfo.lookupTypeFields(name)
		if len(fields) == 0 {
			// query
			text += fmt.Sprintf("- %s: `GET/POST`\n", txt)
		} else {
			jsonStr := p.fieldsJson(fields)
			if len(jsonStr) <= 2 || jsonStr == "{}" {
				// query
				text += fmt.Sprintf("- %s: `GET/POST`\n", txt)
			} else {
				// body
				text += fmt.Sprintf("- %s: `POST`\n", txt)
			}
		}
	}

	method(h.arg, "METHOD")

	queryParam := "- QUERY:\n"
	for _, param := range h.queryParams {
		doc := param.doc
		if len(doc) == 0 {
			doc = param.comment
		}
		doc = strings.TrimSpace(strings.Replace(doc, "\n//", "", -1))
		if len(doc) > 0 {
			doc = "\t" + doc
		}
		queryParam += fmt.Sprintf("\t- `%s={%s}`%s\n", param.queryName, param.Typ, doc)
	}
	text += queryParam

	var fn = func(name string, txt string) {
		fields, _ := p.tplInfo.lookupTypeFields(name)
		if len(fields) == 0 {
			text += fmt.Sprintf("- %s:\n", txt)
		} else {
			text += fmt.Sprintf("- %s:\n", txt)
			jsonStr := p.fieldsJson(fields)
			var dst bytes.Buffer
			json.Indent(&dst, []byte(jsonStr), "\t", "\t")
			jsonStr = p.replaceCommentJson(dst.String())
			text += fmt.Sprintf("\n\t```js\n\t%s\n\t```\n\n", jsonStr)
		}
	}

	fn(h.arg, "BODY")
	fn(h.result, "RESPONSE")

	return text
}

var ptrStringRegexp = regexp.MustCompile(`(\$\d+)":.*[,\n]{1}`)

func (p *Project) replaceCommentJson(s string) string {
	a := ptrStringRegexp.FindAllStringSubmatch(s, -1)
	for _, ss := range a {
		sub := strings.Replace(ss[0], ss[1], "", 1)
		ptr, _ := strconv.Atoi(ss[1][1:])
		f := (*field)(unsafe.Pointer(uintptr(ptr)))
		doc := f.doc
		if len(doc) == 0 {
			doc = f.comment
		}
		doc = strings.TrimSpace(strings.Replace(doc, "\n//", "", -1))
		if sub[len(sub)-1] == ',' {
			s = strings.Replace(s, ss[0], sub+"\t// {"+f.Typ+"} "+doc, 1)
		} else {
			s = strings.Replace(s, ss[0], sub[:len(sub)-1]+"\t// {"+f.Typ+"} "+doc+"\n", 1)
		}
	}
	return s
}

func (p *Project) fieldsJson(fs []*field) string {
	if len(fs) == 0 {
		return ""
	}
	var text string
	text += "{"
	for _, f := range fs {
		if f.isQuery {
			continue
		}
		fieldName := f.ModelName
		if len(fieldName) == 0 {
			fieldName = gutil.FieldSnakeString(f.Name)
		}
		t := strings.Replace(f.Typ, "*", "", -1)
		var isSlice bool
		if strings.HasPrefix(t, "[]") {
			if t == "[]byte" {
				t = "string"
			} else {
				t = strings.TrimPrefix(t, "[]")
				isSlice = true
			}
		}
		v, ok := baseTypeToJsonValue(t)
		if ok {
			if isSlice {
				text += fmt.Sprintf(`"%s$%d":[%s],`, fieldName, uintptr(unsafe.Pointer(f)), v)
			} else {
				text += fmt.Sprintf(`"%s$%d":%s,`, fieldName, uintptr(unsafe.Pointer(f)), v)
			}
			continue
		}
		if ffs, ok := p.tplInfo.lookupTypeFields(t); ok {
			if isSlice {
				text += fmt.Sprintf(`"%s":[%s],`, fieldName, p.fieldsJson(ffs))
			} else {
				text += fmt.Sprintf(`"%s":%s,`, fieldName, p.fieldsJson(ffs))
			}
			continue
		}
	}
	text = strings.TrimRight(text, ",") + "}"
	return text
}

func baseTypeToJsonValue(t string) (string, bool) {
	if t == "bool" {
		return "false", true
	} else if t == "string" || t == "[]byte" || t == "time.Time" {
		return `""`, true
	} else if strings.HasPrefix(t, "int") || t == "rune" {
		return "-0", true
	} else if strings.HasPrefix(t, "uint") || t == "byte" {
		return "0", true
	} else if strings.HasPrefix(t, "float") {
		return "-0.000000", true
	}
	return "", false
}

func (p *Project) genMainFile() {
	p.replace("main.go", "${project_gen_time}", time.Now().Format("2006-01-02 15:04:05"))
	p.replace("main.go", "${service_api_prefix}", gutil.FieldSnakeString(p.Name))
	p.replace("config.go", "${service_api_prefix}", gutil.FieldSnakeString(p.Name))
}

func (p *Project) genConstFile() {
	var text string
	p.replaceWithLine("args/const.gen.go", "${const_list}", text)
}

func (p *Project) genTypeFile() {
	p.replaceWithLine("args/type.gen.go", "${import_list}", p.tplInfo.TypeImportString())
	p.replaceWithLine("args/type.gen.go", "${type_define_list}", p.tplInfo.TypesString())
}

func (p *Project) genRouterFile() {
	p.replaceWithLine(
		"api/router.gen.go",
		"${register_router_list}",
		p.tplInfo.RouterString("_group"),
	)
}

func (p *Project) genHandlerFile() {
	if len(p.tplInfo.ApiHandlerList()) > 0 {
		s := p.tplInfo.ApiHandlerString(func(h *handler) string {
			// bind args
			hc := fmt.Sprintf("// bind arg\narg := new(args.%s)\n\tif err := binding.Binder(ctx, arg); err != nil {\n\t\tctx.RenderRerr(td.RerrInternalServer.SetReason(err.Error()))\n\t\treturn\n\t}\n\n", h.arg)
			// api logic
			hc += fmt.Sprintf("// api logic\n\tresult, rerr := logic.%s(ctx, arg)\n\tif rerr != nil {\n\t\tctx.RenderRerr(rerr)\n\t\treturn\n\t}\n\tctx.Render(result)", h.fullName)
			return hc
		})
		p.replaceWithLine("api/handler.gen.go", "${handler_api_define}", s)
	} else {
		delete(p.codeFiles, "api/handler.gen.go")
		os.Remove("api/handler.gen.go")
	}
}

func (p *Project) genLogicFile() {
	var s string
	for _, h := range p.tplInfo.HandlerList() {
		name := h.fullName
		s += fmt.Sprintf(
			"%sfunc %s(ctx *td.Context,arg *args.%s)(*args.%s,*td.Rerror){\nreturn new(args.%s),nil\n}\n\n",
			h.doc, name, h.arg, h.result, h.result,
		)
	}
	p.replaceWithLine("logic/tmp_code.gen.go", "${logic_api_define}", s)
}

func (p *Project) replace(key, placeholder, value string) string {
	a := strings.Replace(p.codeFiles[key], placeholder, value, -1)
	p.codeFiles[key] = a
	return a
}

func (p *Project) replaceWithLine(key, placeholder, value string) string {
	return p.replace(key, placeholder, "\n"+value)
}

func formatSource(src []byte) []byte {
	b, err := format.Source(src)
	if err != nil {
		xlog.Fatalf("[malatd] format error: %v\ncode:\n%s", err, src)
	}
	return b
}
