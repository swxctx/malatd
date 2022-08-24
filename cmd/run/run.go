package run

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/swxctx/malatd/cmd/create"
	"github.com/swxctx/malatd/cmd/info"
	"github.com/swxctx/malatd/cmd/run/fsnotify"
	"github.com/swxctx/xlog"
)

// RunProject runs project.
func RunProject() {
	if err := os.Chdir(info.AbsPath()); err != nil {
		xlog.Fatalf("[malatd] Jump working directory failed: %v", err)
	}

	go rewatch()

	xlog.Infof("[malatd] Initializing watcher...")

	select {}
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		xlog.Errorf("[malatd] Fail to open file[ %s ]", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		xlog.Errorf("[malatd] Fail to get file information[ %s ]", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// checkTMPFile returns true if the event was for TMP files.
func checkTMPFile(name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".tmp") {
		return true
	}
	return false
}

func rewatch() {
	rerun()

	time.Sleep(time.Second * 2)

	var eventTime = make(map[string]int64)
	filepath.Walk("./", func(retpath string, f os.FileInfo, err error) error {
		if err != nil || f.IsDir() || checkTMPFile(retpath) || !checkIfWatchExt(retpath) {
			return nil
		}
		eventTime[retpath] = f.ModTime().Unix()
		return nil
	})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		xlog.Errorf("[malatd] Fail to create new watcher[ %v ]", err)
		os.Exit(2)
	}
	for _, path := range readAppDirectories("./") {
		xlog.Infof("[malatd] Watching directory[ %s ]", path)
		err = watcher.Watch(path)
		if err != nil {
			xlog.Errorf("[malatd] Fail to watch curpathectory[ %s ]", err)
			os.Exit(2)
		}
	}

	for {
		select {
		case e := <-watcher.Event:
			isbuild := true

			// Skip TMP files for Sublime Text.
			if checkTMPFile(e.Name) {
				continue
			}
			if !checkIfWatchExt(e.Name) {
				continue
			}

			mt := getFileModTime(e.Name)
			if t := eventTime[e.Name]; mt == t {
				isbuild = false
			}
			eventTime[e.Name] = mt

			if isbuild {
				xlog.Infof("%s", e.String())
				watcher.Close()
				if strings.HasSuffix(e.Name, create.MalatdTpl) {
					create.CreateProject()
				}
				go rewatch()
				return
			}
		case err := <-watcher.Error:
			xlog.Warnf("[malatd] %s", err.Error()) // No need to exit here
		}
	}
}

var (
	state        sync.Mutex
	cmd          *exec.Cmd
	isFirstStart = true
)

func rerun() {
	state.Lock()
	defer state.Unlock()
	xlog.Infof("[malatd] Start build...")
	buildCom := exec.Command("go", "build", "-o", info.FileName())
	buildCom.Env = []string{"GOPATH=" + info.Gopath()}
	for _, env := range os.Environ() {
		if strings.HasPrefix(strings.TrimSpace(env), "GOPATH=") {
			continue
		}
		buildCom.Env = append(buildCom.Env, env)
	}
	buildCom.Stdout = os.Stdout
	buildCom.Stderr = os.Stderr
	err := buildCom.Run()
	if err != nil {
		xlog.Errorf("[malatd] ============== Build failed ===================")
		return
	}
	xlog.Infof("[malatd] Build was successful")

	var start string
	if isFirstStart {
		isFirstStart = false
		xlog.Infof("[malatd] Starting app: %s", info.ProjName())
		start = "Start"
	} else {
		xlog.Infof("[malatd] Restarting app: %s", info.ProjName())
		defer func() {
			if e := recover(); e != nil {
				xlog.Infof("[malatd] Kill.recover -> %v", e)
			}
		}()
		if cmd != nil && cmd.Process != nil {
			err := cmd.Process.Kill()
			cmd.Process.Release()
			if err != nil {
				xlog.Infof("[malatd] Kill -> %v", err)
			}
		}
		start = "Restart"
	}

	go func() {
		cmd = exec.Command("./" + info.FileName())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Start(); err != nil {
			xlog.Errorf("[malatd] Fail to start app[ %s ]", err)
			return
		}
		xlog.Infof("[malatd] %s was successful", start)
		cmd.Wait()
		xlog.Infof("[malatd] Old process was stopped")
	}()
}

func readAppDirectories(dir string) (paths []string) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		xlog.Fatalf("[malatd] read project directorys failed: %v", err)
		return
	}
	useDirectory := false
	for _, fileInfo := range fileInfos {
		if checkIfNotWatch(fileInfo.Name()) {
			continue
		}
		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			paths = append(paths, readAppDirectories(path.Join(dir, fileInfo.Name()))...)
			continue
		}

		if useDirectory == true {
			continue
		}

		if checkIfWatchExt(fileInfo.Name()) {
			paths = append(paths, dir)
			useDirectory = true
		}
	}
	return
}

var notWatch = []string{".md"}

func checkIfNotWatch(name string) bool {
	if name[0] == '_' {
		return true
	}
	for _, s := range notWatch {
		if name == s {
			return true
		}
	}
	return false
}

var watchExts = []string{".go", ".ini", ".yaml", ".toml", ".xml", ".json"}

// checkIfWatchExt returns true if the name HasSuffix <watch_ext>.
func checkIfWatchExt(name string) bool {
	for _, s := range watchExts {
		if strings.HasSuffix(strings.ToLower(name), s) {
			return true
		}
	}
	return false
}
