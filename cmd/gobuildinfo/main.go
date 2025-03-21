package main

import (
	"debug/buildinfo"
	"encoding/json"
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic(fmt.Errorf("accepts only 1 arg"))
	}

	p := flag.Arg(0)

	info, err := buildinfo.ReadFile(p)
	if err != nil {
		panic(err)
	}

	bin, _ := json.MarshalIndent(info, "", "    ")
	fmt.Printf("%s\n", bin)
}
