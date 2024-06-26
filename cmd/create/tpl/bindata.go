// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package main generated by go-bindata.
// sources:
// README.md
// __malatd__gen__.lock
// api/handler.go
// api/router.go
// args/const.go
// args/type.go
// args/var.go
// rerrs/rerrs.go
// .gitignore
package tpl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	_ "io/ioutil"
	"os"
	_ "path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _readmeMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8f\x3f\x4f\xc3\x30\x10\xc5\x77\x7f\x8a\x87\x32\x00\x91\xda\xec\xdd\x10\x82\x11\x04\xea\x86\x18\xae\xce\x05\x9b\xda\x39\xeb\x7c\x26\xcd\xb7\x47\x55\x61\xe9\xfa\xfe\xe9\xfd\x3a\x3c\x4a\xce\x34\x8f\xc8\x94\xc8\x46\xe7\x1e\x30\x72\x49\xb2\x66\x9e\x0d\x26\x92\x2a\x26\x51\x7c\x5c\xfc\xcf\xbb\x60\x56\xea\x6e\x18\xbe\xa2\x85\x76\xd8\x7a\xc9\x43\x5d\x4e\xde\x4e\xc3\x25\x71\xff\xb7\x84\xca\xfa\x13\x3d\x63\x52\xca\xbc\x88\x1e\xb7\xce\x75\x1d\x9e\x99\xac\x29\x3b\xb7\xc1\x5b\x8b\xfe\x98\x56\x78\x65\x32\x06\xfd\x37\x8b\xca\x37\x7b\x73\x1b\xbc\xb7\xf9\x4a\xc4\x12\x2d\x20\x88\xc1\x4b\x2e\x31\x91\x45\x99\x9d\xeb\xfb\x97\xd7\xfd\xd3\xae\xef\xb1\x0f\x0c\x5b\x0b\x43\x26\x04\x9a\xc7\xc4\x7a\x5b\x51\xe8\xfc\xc2\x58\x71\x66\x55\xae\x2d\x19\x72\xab\x86\x03\xa3\x9a\x36\x6f\x37\xee\x37\x00\x00\xff\xff\x61\x35\xb0\xd7\x0d\x01\x00\x00")

func readmeMdBytes() ([]byte, error) {
	return bindataRead(
		_readmeMd,
		"README.md",
	)
}

func readmeMd() (*asset, error) {
	bytes, err := readmeMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "README.md", size: 269, mode: os.FileMode(420), modTime: time.Unix(1661234216, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var ___malatd__gen__Lock = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x1c\xca\xb1\x0d\x02\x31\x0c\x05\xd0\x9e\x29\x3e\x37\xc0\x0d\x40\x4f\x0b\xcd\x0d\x90\x28\xb1\x13\x4b\x89\x2d\x61\x03\x61\x7b\x04\xed\xd3\xbb\xdd\x8f\xeb\xe5\x84\xa3\x8b\x23\xcd\x3c\x28\x2a\x1a\x69\x42\xb1\x39\xb3\x56\x98\x8e\x0f\x8a\xbd\xe8\xe1\x60\x19\xe4\x78\x4b\x74\x44\x27\x6c\x7b\x23\xdd\x9b\x6d\xf0\x27\xb3\x2c\x08\xff\xfd\xd7\x40\x4b\x3c\xfc\xfc\x0d\x00\x00\xff\xff\x42\x12\x72\xa5\x60\x00\x00\x00")

func __malatd__gen__LockBytes() ([]byte, error) {
	return bindataRead(
		___malatd__gen__Lock,
		"__malatd__gen__.lock",
	)
}

func __malatd__gen__Lock() (*asset, error) {
	bytes, err := __malatd__gen__LockBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "__malatd__gen__.lock", size: 96, mode: os.FileMode(420), modTime: time.Unix(1661251175, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _apiHandlerGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x48\x4c\xce\x4e\x4c\x4f\x55\x48\x2c\xc8\xe4\x02\x04\x00\x00\xff\xff\x0c\x0c\x0a\x62\x0c\x00\x00\x00")

func apiHandlerGoBytes() ([]byte, error) {
	return bindataRead(
		_apiHandlerGo,
		"api/handler.go",
	)
}

func apiHandlerGo() (*asset, error) {
	bytes, err := apiHandlerGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "api/handler.go", size: 12, mode: os.FileMode(420), modTime: time.Unix(1571309527, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _apiRouterGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x04\xc0\xc1\x0d\xc2\x30\x0c\x05\xd0\x33\x9e\xc2\xea\xa9\x45\xa8\x1d\x83\x0b\x37\x26\x08\x4e\x08\x16\xa4\x8e\x9c\x9f\x52\x09\xb1\x7b\x5f\x0d\xf2\x0e\x39\x71\xa8\x4a\xa4\xa5\x9a\x83\x47\x3a\x21\xf2\x90\x15\xaf\xfe\x98\xc5\xca\xd2\xbe\xbb\x60\x5f\x4a\xf8\x04\xc4\x81\x26\xa2\x67\x5f\x85\xdd\x3a\xd2\xcd\xb2\xca\xd8\x7c\xe3\x33\xe2\x7c\x4f\xbe\x25\xbf\xb0\x9b\xe1\xea\xd6\x2b\x37\xb8\xae\x79\xe2\x1f\xfd\xe9\x08\x00\x00\xff\xff\xe7\x8e\x1f\xf5\x6e\x00\x00\x00")

func apiRouterGoBytes() ([]byte, error) {
	return bindataRead(
		_apiRouterGo,
		"api/router.go",
	)
}

func apiRouterGo() (*asset, error) {
	bytes, err := apiRouterGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "api/router.go", size: 110, mode: os.FileMode(420), modTime: time.Unix(1661332010, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _argsConstGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x48\x4c\xce\x4e\x4c\x4f\x55\x48\x2c\x4a\x2f\xe6\xe2\x4a\xce\xcf\x2b\x2e\x51\xd0\xd0\xe4\x02\x04\x00\x00\xff\xff\xf5\x9a\x10\x2f\x17\x00\x00\x00")

func argsConstGoBytes() ([]byte, error) {
	return bindataRead(
		_argsConstGo,
		"args/const.go",
	)
}

func argsConstGo() (*asset, error) {
	bytes, err := argsConstGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "args/const.go", size: 23, mode: os.FileMode(420), modTime: time.Unix(1571309527, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _argsTypeGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x48\x4c\xce\x4e\x4c\x4f\x55\x48\x2c\x4a\x2f\xe6\xe2\x2a\xa9\x2c\x48\x55\xd0\xd0\xe4\x02\x04\x00\x00\xff\xff\x61\x1b\x80\x25\x16\x00\x00\x00")

func argsTypeGoBytes() ([]byte, error) {
	return bindataRead(
		_argsTypeGo,
		"args/type.go",
	)
}

func argsTypeGo() (*asset, error) {
	bytes, err := argsTypeGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "args/type.go", size: 22, mode: os.FileMode(420), modTime: time.Unix(1571309527, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _argsVarGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\x48\x4c\xce\x4e\x4c\x4f\x55\x48\x2c\x4a\x2f\xe6\xe2\x2a\x4b\x2c\x52\xd0\xd0\xe4\x02\x04\x00\x00\xff\xff\xa5\xca\xdc\xfb\x15\x00\x00\x00")

func argsVarGoBytes() ([]byte, error) {
	return bindataRead(
		_argsVarGo,
		"args/var.go",
	)
}

func argsVarGo() (*asset, error) {
	bytes, err := argsVarGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "args/var.go", size: 21, mode: os.FileMode(420), modTime: time.Unix(1571309527, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _rerrsRerrsGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\xca\x4d\xae\x82\x30\x10\x00\xe0\x35\x73\x8a\x49\x57\x90\x10\x5a\x0e\xf0\x56\x6f\xe5\xc6\x18\x6e\x30\xd2\x09\x36\x52\x8a\xd3\xe1\xe7\xf8\x46\x83\xae\x5c\x7f\xdf\x4c\xfd\x9d\x06\x46\x61\x91\x0c\x10\xe2\x9c\x44\xb1\x84\x42\x3d\x9a\x21\xe8\x6d\xb9\x36\x7d\x8a\x36\x6f\x7b\xaf\xbb\x8d\x34\x92\x7a\x03\x15\xc0\x4a\xf2\x7a\xd6\x62\xc7\x22\xa7\x69\xa5\x31\xf8\x0b\x09\x45\x56\x16\x64\x91\x24\x50\xfc\xb4\x3f\x54\xdf\x9c\x79\xeb\xde\xa9\x6c\x9d\x73\xae\xad\xd1\x1c\x11\xbf\xd3\xd4\x68\xfe\xd3\xa4\x14\xa6\x8c\xe1\x50\xe1\xc7\xc2\x59\x71\xfe\xac\x6c\x2a\xa8\xe0\x19\x00\x00\xff\xff\x8d\xca\x2a\x0c\xca\x00\x00\x00")

func rerrsRerrsGoBytes() ([]byte, error) {
	return bindataRead(
		_rerrsRerrsGo,
		"rerrs/rerrs.go",
	)
}

func rerrsRerrsGo() (*asset, error) {
	bytes, err := rerrsRerrsGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "rerrs/rerrs.go", size: 202, mode: os.FileMode(420), modTime: time.Unix(1661223879, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _Gitignore = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\x4f\x4b\xc4\x40\x0c\xc5\xef\xef\xa3\x04\x0c\x28\x28\x5e\x3d\x78\x11\xfc\x83\x57\x91\xa1\x9d\xa6\x63\xd7\x6d\x13\x67\xa6\x6b\xdd\x65\xfd\xec\x32\x3b\x15\xf6\xf2\x92\xf7\x23\xbc\x24\xc4\x0a\xe2\x06\xc4\x49\xe1\xb4\xdd\xc0\x65\x49\x19\xc4\x6f\xd7\x37\xb7\xbb\xaf\x77\xac\x95\x75\x2e\xd4\x07\xbd\xe4\xa0\xb5\xbb\x62\x0f\xe7\x83\xba\x4e\xfa\x79\xfa\x37\x41\xf3\x8f\x49\x2a\x53\x27\x2f\x8b\x69\xcc\x4c\x35\x79\x6c\x86\xa9\x06\xc8\x22\x55\x7f\x41\xbc\x2e\xb5\xa8\x3d\x88\x63\x13\x41\xbc\x1f\x0c\xc4\x61\x5f\x78\xea\x40\xdc\x8e\x45\x7d\x1f\x0a\xb1\xbc\x80\x78\xab\x01\x34\xe9\xc7\x6c\xeb\x81\xbb\xe4\xb5\x2b\xc1\x43\x27\xa7\xb7\xe6\x76\x3b\x8c\x72\x61\x51\x37\xe2\xf3\x19\xf9\xd6\xf8\x99\xac\xf1\x82\xc3\xe1\xe5\xf5\xf9\xc1\x3d\xdd\x3d\xde\x1f\x8f\xf8\x0b\x00\x00\xff\xff\x67\xf4\xec\xbf\x13\x01\x00\x00")

func GitignoreBytes() ([]byte, error) {
	return bindataRead(
		_Gitignore,
		".gitignore",
	)
}

func Gitignore() (*asset, error) {
	bytes, err := GitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: ".gitignore", size: 275, mode: os.FileMode(420), modTime: time.Unix(1637820424, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"README.md":            readmeMd,
	"__malatd__gen__.lock": __malatd__gen__Lock,
	"api/handler.go":       apiHandlerGo,
	"api/router.go":        apiRouterGo,
	"args/const.go":        argsConstGo,
	"args/type.go":         argsTypeGo,
	"args/var.go":          argsVarGo,
	"rerrs/rerrs.go":       rerrsRerrsGo,
	".gitignore":           Gitignore,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	".gitignore":           &bintree{Gitignore, map[string]*bintree{}},
	"README.md":            &bintree{readmeMd, map[string]*bintree{}},
	"__malatd__gen__.lock": &bintree{__malatd__gen__Lock, map[string]*bintree{}},
	"api": &bintree{nil, map[string]*bintree{
		"handler.go": &bintree{apiHandlerGo, map[string]*bintree{}},
		"router.go":  &bintree{apiRouterGo, map[string]*bintree{}},
	}},
	"args": &bintree{nil, map[string]*bintree{
		"const.go": &bintree{argsConstGo, map[string]*bintree{}},
		"type.go":  &bintree{argsTypeGo, map[string]*bintree{}},
		"var.go":   &bintree{argsVarGo, map[string]*bintree{}},
	}},
	"rerrs": &bintree{nil, map[string]*bintree{
		"rerrs.go": &bintree{rerrsRerrsGo, map[string]*bintree{}},
	}},
}}
