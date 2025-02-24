package main

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	even = flag.Bool("even", false, "pick even")
	odd  = flag.Bool("odd", false, "pick odd")
)

func main() {
	flag.Parse()

	versions := []version{}
	err := json.NewDecoder(os.Stdin).Decode(&versions)
	if err != nil {
		panic(err)
	}

	if len(versions) == 0 {
		panic(fmt.Errorf("input has zero element"))
	}

	slices.SortFunc(versions, func(i, j version) int { return i.Compare(j) })

	found, idx := hiter.FindLastFunc(
		func(v version) bool {
			switch {
			case *even:
				return v.comp[0]%2 == 0
			case *odd:
				return v.comp[0]%2 == 1
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

type version struct {
	leng int
	comp [4]uint
}

func (v *version) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("too short")
	}

	if data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("not a string")
	}

	str := string(data[1 : len(data)-1])

	splitted := strings.Split(str, ".")

	if len(splitted) > 4 {
		return fmt.Errorf("contains too many \".\"")
	}

	var compo [4]uint
	for i, comp := range splitted {
		num, err := strconv.Atoi(comp)
		if err != nil {
			return fmt.Errorf("at %dth: %w", i, err)
		}
		if num < 0 {
			return fmt.Errorf("at %dth: negative num", i)
		}
		if num >= 100000 {
			return fmt.Errorf("at %dth: too large", i)
		}
		compo[i] = uint(num)
	}

	v.leng = len(splitted)
	v.comp = compo

	return nil
}

func (v version) String() string {
	if v.leng == 0 {
		return `0.0.0.0`
	}
	var s strings.Builder
	for i := range v.leng {
		if i > 0 {
			s.WriteByte('.')
		}
		s.WriteString(strconv.Itoa(int(v.comp[i])))
	}
	return s.String()
}

func (v version) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

func (v version) Compare(j version) int {
	for i := range v.comp {
		if c := cmp.Compare(v.comp[i], j.comp[i]); c != 0 {
			return c
		}
	}
	return 0
}
