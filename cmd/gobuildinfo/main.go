package main

import (
	"debug/buildinfo"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	onlyMain = flag.Bool("main", false, "print only main version. v prefix is stripped.")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic(fmt.Errorf("accepts only 1 arg"))
	}

	p := flag.Arg(0)

	if !strings.HasPrefix(p, "./") && !strings.HasPrefix(p, "/") {
		base, err := gobinPath()
		if err != nil {
			panic(err)
		}
		p = filepath.Join(base, path.Base(p))
	}

	info, err := buildinfo.ReadFile(p)
	if err != nil {
		panic(err)
	}

	if *onlyMain {
		fmt.Println(strings.TrimPrefix(info.Main.Version, "v"))
		return
	}

	bin, _ := json.MarshalIndent(info, "", "    ")
	fmt.Printf("%s\n", bin)
}

func gobinPath() (path string, err error) {
	bin, err := exec.Command("go", "env", "GOBIN").Output()
	if err != nil {
		return "", err
	}
	if path = strings.TrimSpace(string(bin)); len(path) > 0 {
		return path, err
	}

	bin, err = exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return "", err
	}

	if path = strings.TrimSpace(string(bin)); len(path) > 0 {
		return filepath.Join(path, "bin"), err
	}

	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, "go", "bin"), nil
}
