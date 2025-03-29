package version

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	leng int
	comp [4]uint
}

func Parse(s string) (Version, error) {
	var v Version
	err := v.UnmarshalText([]byte(s))
	return v, err
}

func (v *Version) UnmarshalText(data []byte) error {
	splitted := strings.Split(string(data), ".") // TODO: use cut for more efficiency

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

func (v *Version) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("too short")
	}

	if data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("not a string")
	}

	return v.UnmarshalText(data[1 : len(data)-1])
}

func (v Version) Nums() []uint {
	ret := make([]uint, v.leng)
	copy(ret, v.comp[:])
	return ret
}

func (v Version) NumsFull() [4]uint {
	return v.comp
}

func (v Version) String() string {
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

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.String() + "\""), nil
}

func (v Version) Compare(j Version) int {
	for i := range 4 {
		if c := cmp.Compare(v.comp[i], j.comp[i]); c != 0 {
			return c
		}
	}
	return 0
}
