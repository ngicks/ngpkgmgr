package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/ngicks/go-common/exver"
	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	even = flag.Bool("even", false, "pick even")
	odd  = flag.Bool("odd", false, "pick odd")
)

func main() {
	flag.Parse()

	versions := []exver.Version{}
	err := json.NewDecoder(os.Stdin).Decode(&versions)
	if err != nil {
		panic(err)
	}

	if len(versions) == 0 {
		panic(fmt.Errorf("input has zero element"))
	}

	slices.SortFunc(versions, func(i, j exver.Version) int { return -i.Compare(j) })

	found, idx := hiter.FindFunc(
		func(v exver.Version) bool {
			switch {
			case *even:
				return v.Core().Major()%2 == 0
			case *odd:
				return v.Core().Major()%2 == 1
			default:
				return true
			}
		},
		slices.Values(versions),
	)
	if idx < 0 {
		panic("not found")
	}

	fmt.Printf("%s\n", found)
}
