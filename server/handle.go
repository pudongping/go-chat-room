package server

import (
	"os"
	"path/filepath"
)

var rootDir string

func RegisterHandle() {
	inferRootDir()

	// 广播消息处理
	// go logic.Broadcaster.Start()

	// 	http.HandleFunc("/", homeHandleFunc)
	//	http.HandleFunc("/ws", WebSocketHandleFunc)

}

// 推断出项目根目录
func inferRootDir() {
	// 通过 os.Getwd() 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		// 这里要确保项目根目录下存在 template 目录
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
