package create

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"strings"

	"github.com/swxctx/gutil"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/malatd/cmd/create/structtag"
	"github.com/swxctx/xlog"
)

const (
	// API_TPL_ROUTER name of the interface used to register the pull route in the template
	API_TPL_ROUTER = "__API_TPL__"
)

type (
	tplInfo struct {
		src               []byte
		fileSet           *token.FileSet
		astFile           *ast.File
		doc               string
		apiRouter         *router
		realStructTypes   []*structType
		realStructTypeMap map[string]*structType
		aliasTypes        []*aliasType
		typeImports       []string
	}
	router struct {
		doc      string
		name     string
		handlers []*handler
		children []*router
	}
	handler struct {
		uri         string
		queryParams []*field
		group       *router
		doc         string
		name        string
		fullName    string
		arg         string
		result      string
	}
	structType struct {
		doc              string
		name             string
		fields           []*field
		primaryFields    []*field
		uniqueFields     []*field
		isDefaultPrimary bool
		node             *ast.StructType
	}
	field struct {
		Name      string
		ModelName string
		Typ       string
		isQuery   bool
		queryName string
		anonymous bool
		tag       string
		doc       string
		comment   string
	}
	aliasType struct {
		doc         string
		name        string
		text        string
		rawTypeName string
		rawStruct   *structType
	}
)

func newTplInfo(tplBytes []byte) *tplInfo {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", tplBytes, parser.ParseComments)
	if err != nil {
		xlog.Fatalf("[malatd] %v", err)
	}
	return &tplInfo{
		src:               tplBytes,
		fileSet:           fset,
		astFile:           file,
		doc:               addSlash(file.Doc.Text()),
		apiRouter:         &router{},
		realStructTypeMap: make(map[string]*structType),
	}
}

func (t *tplInfo) Parse() *tplInfo {
	t.parseImports()
	ok := t.hasType(emptyStructType.name)
	if ok {
		xlog.Fatalf("[malatd] Keep structure name cannot be used: %s", emptyStructType.name)
	}
	t.aliasTypes = append(t.aliasTypes, emptyStructType)
	t.collectStructs()
	t.collectAliasTypes()
	t.collectIfaces()
	return t
}

func (t *tplInfo) TypeImportString() string {
	return strings.Join(t.typeImports, "\n")
}

func (t *tplInfo) TypesString() string {
	var a string
	for _, s := range t.aliasTypes {
		a += s.String()
	}
	for _, s := range t.realStructTypes {
		a += s.String()
	}
	return a
}

func (t *tplInfo) ApiHandlerString(ctnFn func(*handler) string) string {
	if ctnFn == nil {
		ctnFn = func(*handler) string { return "return nil,nil" }
	}
	return t.apiRouter.handlerString(ctnFn)
}

func (t *tplInfo) ApiHandlerList() []*handler {
	return t.apiRouter.handlerList()
}

func (t *tplInfo) HandlerList() []*handler {
	return append(t.apiRouter.handlerList())
}

func (t *tplInfo) RouterString(groupName string) string {
	var text string
	text += "\n// APIs...\n"
	text += "{\n"
	text += t.apiRouter.routerString(groupName, "", "")
	text += "}\n"
	return text
}

func (t *tplInfo) getCodeBlock(i interface{}) string {
	var dst bytes.Buffer
	err := format.Node(&dst, t.fileSet, i)
	if err != nil {
		xlog.Fatalf("[malatd] %v", err)
	}
	return dst.String()
}

func (t *tplInfo) parseImports() {
	const codec = `"github.com/swxctx/malatd/binding"`
	t.typeImports = append(t.typeImports, codec)
	for _, imp := range t.astFile.Imports {
		s := t.getCodeBlock(imp)
		if s != codec {
			t.typeImports = append(t.typeImports, s)
		}
	}
}

// collectStructs collects and maps structType nodes to their positions
func (t *tplInfo) collectStructs() {
	collectStructs := func(n ast.Node) bool {
		decl, ok := n.(ast.Decl)
		if !ok {
			return true
		}
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return true
		}
		var groupDoc string
		if len(genDecl.Specs) == 1 {
			groupDoc = genDecl.Doc.Text()
		}
		for _, spec := range genDecl.Specs {
			var e ast.Expr
			var structName string
			var doc = groupDoc

			switch x := spec.(type) {
			case *ast.TypeSpec:
				if x.Type == nil {
					continue
				}
				structName = x.Name.Name
				e = x.Type
				if s := x.Doc.Text(); s != "" {
					doc = x.Doc.Text()
				}
			}

			x, ok := e.(*ast.StructType)
			if !ok {
				continue
			}

			if len(x.Fields.List) == 0 {
				if gutil.IsExportedName(structName) {
					a := &aliasType{
						doc:  addSlash(doc),
						name: structName,
						text: fmt.Sprintf("%s = binding.Empty", structName),
					}
					a.rawTypeName = a.text[strings.LastIndex(strings.TrimSpace(strings.Split(a.text, "//")[0]), " ")+1:]
					if a.doc == "" {
						a.doc = fmt.Sprintf("// %s alias of type %s\n", a.name, a.rawTypeName)
					}
					t.aliasTypes = append(t.aliasTypes, a)
				}
				continue
			}

			t.realStructTypes = append(
				t.realStructTypes,
				structType{
					name: structName,
					doc:  addSlash(doc),
					node: x,
				}.init(t),
			)
		}
		return true
	}
	ast.Inspect(t.astFile, collectStructs)
	t.sortStructs()
}

func (t *tplInfo) lookupAliasType(name string) (*aliasType, bool) {
	for _, a := range t.aliasTypes {
		if a.name == name {
			return a, true
		}
	}
	return nil, false
}

func (t *tplInfo) lookupStructType(name string) (*structType, bool) {
	s, ok := t.realStructTypeMap[name]
	return s, ok
}

func (t *tplInfo) hasType(name string) bool {
	_, ok := t.lookupStructType(name)
	if !ok {
		_, ok = t.lookupAliasType(name)
	}
	return ok
}

func (t *tplInfo) lookupTypeFields(name string) ([]*field, bool) {
	s, ok := t.lookupStructType(name)
	if ok {
		return s.fields, true
	}
	a, ok := t.lookupAliasType(name)
	if ok {
		if a.rawStruct != nil {
			return a.rawStruct.fields, true
		}
	}
	return nil, false
}

func (t *tplInfo) collectAliasTypes() {
	collectAliasTypes := func(n ast.Node) bool {
		decl, ok := n.(ast.Decl)
		if !ok {
			return true
		}
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return true
		}
		var groupDoc string
		if len(genDecl.Specs) == 1 {
			groupDoc = genDecl.Doc.Text()
		}
		for _, spec := range genDecl.Specs {
			var aliasName string
			var doc = groupDoc
			x, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if x.Type == nil {
				continue
			}
			switch x.Type.(type) {
			case *ast.Ident, *ast.SelectorExpr:
			default:
				continue
			}

			aliasName = x.Name.Name

			if s := x.Doc.Text(); s != "" {
				doc = x.Doc.Text()
			}
			a := &aliasType{
				name: aliasName,
				doc:  addSlash(doc),
				text: t.getCodeBlock(spec),
			}
			txtArr := strings.Split(a.text, "\n")
			a.text = txtArr[len(txtArr)-1]
			a.rawTypeName = a.text[strings.LastIndex(strings.TrimSpace(strings.Split(a.text, "//")[0]), " ")+1:]
			a.rawStruct = t.realStructTypeMap[strings.TrimLeft(a.rawTypeName, "*")]
			if a.doc == "" {
				a.doc = fmt.Sprintf("// %s alias of type %s\n", a.name, a.rawTypeName)
			}
			t.aliasTypes = append(
				t.aliasTypes,
				a,
			)
		}
		return true
	}
	ast.Inspect(t.astFile, collectAliasTypes)
}

// checkDefaultFields  CreatedAt UpdatedAt Deleted
func (s *structType) checkDefaultFields() bool {
	return s.getField("CreatedAt") != nil && s.getField("UpdatedAt") != nil && s.getField("Deleted") != nil
}

func (s *structType) getField(fieldName string) *field {
	for _, f := range s.fields {
		if f.Name == fieldName {
			return f
		}
	}
	return nil
}

func (s *structType) setFields(toLeader bool, fields ...*field) {
	for _, f := range fields {
		for i, ff := range s.fields {
			if ff.Name == f.Name {
				s.fields = append(s.fields[:i], s.fields[i+1:]...)
				break
			}
		}
		if toLeader {
			s.fields = append([]*field{f}, s.fields...)
		} else {
			s.fields = append(s.fields, f)
		}
	}
}

func (s *structType) isInvildName() bool {
	return gutil.IsExportedName(s.name)
}

func (s structType) init(t *tplInfo) *structType {
	if !s.isInvildName() {
		xlog.Fatalf("[malatd] Unexported struct name: %s", s.name)
	}
	for _, v := range s.node.Fields.List {
		f := new(field)
		if len(v.Names) > 0 {
			f.Name = v.Names[0].Name
			if !gutil.IsExportedName(f.Name) {
				xlog.Fatalf("[malatd] Unexported field name: %s.%s", s.name, f.Name)
			}
		}
		f.Typ = t.getCodeBlock(v.Type)
		if len(f.Name) == 0 {
			f.anonymous = true
			f.Name = strings.TrimPrefix(f.Typ, "*")
			if !gutil.IsExportedName(f.Name) {
				xlog.Fatalf("[malatd] Unexported anonymous field: %s.%s", s.name, f.Typ)
			}
		}
		f.doc = addSlash(v.Doc.Text())
		f.comment = addSlash(v.Comment.Text())
		if v.Tag != nil {
			f.tag = v.Tag.Value
			f.queryName, f.isQuery = getQueryField(f.tag)
			if len(f.queryName) == 0 {
				f.queryName = gutil.FieldSnakeString(f.Name)
			}
		}
		s.fields = append(s.fields, f)
	}
	s.rangeTags(
		addJsonTag,
	)
	if s.doc == "" {
		s.doc = fmt.Sprintf("// %s comment...\n", s.name)
	}
	return &s
}

var queryRegexp = regexp.MustCompile("<\\s*query\\s*(:[^:>]*)?>")

func getQueryField(tag string) (queryName string, isQuery bool) {
	a := queryRegexp.FindStringSubmatch(tag)
	if len(a) != 2 {
		return
	}
	isQuery = true
	queryName = strings.TrimLeft(a[1], ":")
	queryName = strings.TrimSpace(queryName)
	return
}

func (s *structType) rangeTags(fns ...func(tags *structtag.Tags, f *field, anonymous bool) bool) {
	for _, fn := range fns {
		for _, v := range s.fields {
			logName := v.Name
			if len(logName) == 0 {
				logName = v.Typ
			}
			tags, err := structtag.Parse(strings.TrimSpace(strings.Trim(v.tag, "`")))
			if err != nil {
				xlog.Fatalf("[malatd] %s.%s: %s", s.name, logName, err.Error())
			}
			if !fn(tags, v, len(v.Name) == 0) {
				break
			}
			v.tag = "`" + tags.String() + "`"
		}
	}
}

var addJsonTag = func(tags *structtag.Tags, f *field, anonymous bool) bool {
	tag, _ := tags.Get("json")
	if tag != nil {
		return true
	}
	tags.Set(&structtag.Tag{
		Key:  "json",
		Name: gutil.FieldSnakeString(f.Name),
	})
	return true
}

func (s *structType) String() string {
	r := fmt.Sprintf("%stype %s struct {\n", s.doc, s.name)
	for _, f := range s.fields {
		if f.anonymous {
			r += fmt.Sprintf("%s  %s  %s  %s", f.doc, f.Typ, f.tag, f.comment)
		} else {
			r += fmt.Sprintf("%s  %s  %s  %s", f.doc+f.Name, f.Typ, f.tag, f.comment)
		}
		if r[len(r)-1] != '\n' {
			r += "\n"
		}
	}
	r += "}\n\n"
	return r
}

func (t *tplInfo) sortStructs() {
	var lastList []*structType
	for _, v := range t.realStructTypes {
		lastList = append(lastList, v)
		t.realStructTypeMap[v.name] = v
	}
	t.realStructTypes = lastList
}

func (t *tplInfo) collectIfaces() {
	var apiIface *ast.InterfaceType
	ast.Inspect(t.astFile, func(n ast.Node) bool {
		var e ast.Expr
		var ifaceName string
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Type == nil {
				return true
			}
			ifaceName = x.Name.Name
			e = x.Type
		}
		x, ok := e.(*ast.InterfaceType)
		if !ok {
			return true
		}
		switch ifaceName {
		case API_TPL_ROUTER:
			apiIface = x
		}
		return true
	})
	t.collectApis(t.apiRouter, apiIface)
}

func (t *tplInfo) collectApis(r *router, i *ast.InterfaceType) bool {
	if i == nil {
		return false
	}
	for _, f := range i.Methods.List {
		switch n := f.Type.(type) {
		case *ast.Ident:
			if n.Obj == nil {
				continue
			}
			x := n.Obj.Decl.(*ast.TypeSpec)
			child := new(router)
			child.name = x.Name.Name
			child.doc = addSlash(f.Doc.Text())
			if len(child.doc) == 0 {
				child.doc = fmt.Sprintf("// %s controller\n", child.name)
			}
			if t.collectApis(child, x.Type.(*ast.InterfaceType)) {
				r.children = append(r.children, child)
			}

		case *ast.FuncType:
			var funcName string
			if len(f.Names) > 0 {
				funcName = f.Names[0].Name
			} else {
				xlog.Fatalf("[malatd] no name of the function: %s", t.getCodeBlock(i))
			}
			h, err := t.getHandler(n)
			if err != nil {
				xlog.Fatalf("[malatd] %s.%s: %s", r.name, funcName, err.Error())
			}
			h.name = funcName
			h.group = r
			h.doc = addSlash(f.Doc.Text())
			if len(h.doc) == 0 {
				h.doc = fmt.Sprintf("// %s handler\n", h.name)
			}
			r.handlers = append(r.handlers, h)
		}
	}
	return true
}

func (t *tplInfo) getHandler(f *ast.FuncType) (*handler, error) {
	if f.Params.NumFields() != 1 {
		return nil, fmt.Errorf("the number of method parameters should be %d", 1)
	}
	var numResults int = 1
	if f.Results.NumFields() != numResults {
		return nil, fmt.Errorf("the number of method results should be %d", numResults)
	}
	h := new(handler)
	var ft = f.Params.List[0].Type
	if se, ok := ft.(*ast.StarExpr); ok {
		ft = se.X
	}
	switch ftype := ft.(type) {
	case *ast.Ident:
		ok := t.hasType(ftype.Name)
		if !ok {
			return nil, fmt.Errorf("the type of method parameter should be clearly defined")
		}
		h.arg = ftype.Name
		if fields, ok := t.lookupTypeFields(ftype.Name); ok {
			for _, v := range fields {
				if v.isQuery {
					h.queryParams = append(h.queryParams, v)
				}
			}
		}
	case *ast.StructType:
		if len(ftype.Fields.List) != 0 {
			return nil, fmt.Errorf("the type of method parameter should be clearly defined")
		}
		h.arg = emptyStructType.name
	}
	// push handler has no result
	if numResults == 0 {
		return h, nil
	}
	// pull handler has result
	ft = f.Results.List[0].Type
	if se, ok := ft.(*ast.StarExpr); ok {
		ft = se.X
	}
	switch ftype := ft.(type) {
	case *ast.Ident:
		ok := t.hasType(ftype.Name)
		if !ok {
			return nil, fmt.Errorf("the type of method result is not struct type")
		}
		h.result = ftype.Name
	case *ast.StructType:
		if len(ftype.Fields.List) != 0 {
			return nil, fmt.Errorf("the type of method result should be pointer")
		}
		h.result = emptyStructType.name
	}
	return h, nil
}

var emptyStructType = &aliasType{
	doc:  "// EmptyStruct alias of type struct {}\n",
	name: "EmptyStruct",
	text: "EmptyStruct = binding.Empty",
}

func (a *aliasType) String() string {
	return fmt.Sprintf("%stype %s\n", a.doc, a.text)
}

func (r *router) handlerList() []*handler {
	var hs []*handler
	hs = append(hs, r.handlers...)
	for _, child := range r.children {
		hs = append(hs, child.handlerList()...)
	}
	return hs
}

func (r *router) handlerString(ctnFn func(*handler) string) string {
	var text string
	if len(r.handlers) > 0 {
		for _, h := range r.handlers {
			if len(r.name) > 0 {
				text += fmt.Sprintf(
					"%sfunc %sHandle(ctx *td.Context) {\n%s\n}\n\n",
					h.doc, h.name, strings.TrimSpace(ctnFn(h)),
				)
			} else {
				text += fmt.Sprintf(
					"%sfunc %sHandle() {\n%s\n}\n\n",
					h.doc, h.name, strings.TrimSpace(ctnFn(h)),
				)
			}
		}
	}
	for _, child := range r.children {
		text += child.handlerString(ctnFn)
	}
	return text
}

func (r *router) routerString(groupName, fullNamePrefix, uriPrefix string) string {
	var text, subGroupName string
	if len(r.name) > 0 {
		fullNamePrefix = joinName(fullNamePrefix, r.name)
		uriPrefix = path.Join("/", uriPrefix, td.ToUriPath(r.name))
		groupModule := strings.ToLower(r.name)
		text += fmt.Sprintf("\n\n%s := srv.Group(\"%s\")\n", groupModule, uriPrefix)
		if len(r.handlers) > 0 {
			for _, h := range r.handlers {
				h.fullName = joinName(fullNamePrefix, h.name)
				h.uri = path.Join("/", uriPrefix, td.ToUriPath(h.name))
				text += fmt.Sprintf("%s.Post(\"%s\", %sHandle)\n", groupModule, td.ToUriPath(h.name), h.name)
				text += fmt.Sprintf("%s.Get(\"%s\", %sHandle)\n", groupModule, td.ToUriPath(h.name), h.name)
			}
		}
		if len(r.children) > 0 {
			for _, child := range r.children {
				subGroupName = "_" + firstLowerLetter(r.name) + r.name[1:]
				text += child.routerString(subGroupName, fullNamePrefix, uriPrefix)
			}
		}
	} else {
		for _, h := range r.handlers {
			h.fullName = joinName(fullNamePrefix, h.name)
			h.uri = path.Join("/", uriPrefix, td.ToUriPath(h.name))
			text += fmt.Sprintf("srv.Post(\"%s\", %sHandle)", td.ToUriPath(h.name), h.name)
			text += fmt.Sprintf("srv.Get(\"%s\", %sHandle)", td.ToUriPath(h.name), h.name)
		}
		for _, child := range r.children {
			text += child.routerString(groupName, fullNamePrefix, uriPrefix)
		}
	}
	return text
}

func joinName(a, b string) string {
	a = strings.Trim(a, "_")
	b = strings.Trim(b, "_")
	if a == "" {
		return b
	}
	a = strings.ToUpper(a[:1]) + a[1:]
	if b == "" {
		return a
	}
	b = strings.ToUpper(b[:1]) + b[1:]
	return a + "_" + b
}

func firstLowerLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0]))
}

func addSlash(txt string) (comment string) {
	r := bufio.NewReader(strings.NewReader(txt))
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		comment += "// " + string(line) + "\n"
	}
	return
}
