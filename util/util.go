package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
)

func DumpJson(v interface{}, label string) {
	foo, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("failed to dump object to JSON: %s", err.Error())
	}
	fmt.Printf("%s%s\n", label, foo)
}

func DumpObject(v interface{}, label string) {
	fmt.Printf("%s%#v\n", label, v)
}

func SplitPathVar(value string) []string {
	return strings.Split(value, ":")
}

func IsExecutable(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := stat.Mode()
	statSys, ok := stat.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	if mode&0001 != 0 {
		// "other" executable
		return true
	} else if mode&0010 != 0 && int(statSys.Gid) == os.Getgid() {
		// "group" executable
		return true
	} else if mode&0100 != 0 && int(statSys.Uid) == os.Getuid() {
		// "user" executable
		return true
	}
	return false
}

func FindExecutables(name string, paths []string, allowPartial bool) []string {
	ret := []string{}
	for _, path := range paths {
		if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
			continue
		}
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			if allowPartial {
				if f.Name()[:len(name)] == name {
					ret = append(ret, fmt.Sprintf("%s/%s", path, f.Name()))
				}
			} else if f.Name() == name {
				ret = append(ret, fmt.Sprintf("%s/%s", path, f.Name()))
			}
		}
	}
	return ret
}
