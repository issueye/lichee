package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func FindDir(dir string, args ...interface{}) []string {
	var suffix string
	if len(args) > 0 {
		suffix = args[0].(string)
		if suffix[0] != '.' {
			suffix = "." + suffix
		}
	}

	dirList := make([]string, 0)
	fileinfo, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	// 遍历这个文件夹
	for _, fi := range fileinfo {
		// 重复输出制表符，模拟层级结构
		// fmt.Println(strings.Repeat("\t", num))
		// 判断是不是目录
		if fi.IsDir() {
			tmpList := FindDir(dir+`\`+fi.Name(), args...)
			dirList = append(dirList, tmpList...)
		} else {
			path := dir + `\` + fi.Name()
			if suffix != "" {
				if filepath.Ext(path) == suffix {
					fmt.Println(`文件：`, path)
					dirList = append(dirList, path)
				}
			} else {
				dirList = append(dirList, path)
			}
		}
	}
	return dirList
}
