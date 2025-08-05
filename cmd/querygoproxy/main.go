package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/ngicks/go-common/exver"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/stringsiter"
)

var (
	onlyLatest = flag.Bool("only-latest", false, "prints only latest version")
	excludePre = flag.Bool("exclude-pre", false, "if set, excludes pre-release version")
)

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

var proxyGolangOrgUrl = must(url.Parse("https://proxy.golang.org"))

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic(fmt.Errorf("accepts only 1 arg"))
	}

	p := flag.Arg(0)

	targetUrl := proxyGolangOrgUrl.JoinPath(p)
	targetUrl = targetUrl.JoinPath("/@v/list")

	client := &http.Client{}

	resp, err := client.Get(targetUrl.String())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("status not ok: %d", resp.StatusCode))
	}

	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	versions, err := hiter.TryCollect(
		hiter.Filter2(func(ver exver.Version, err error) bool {
			return !(*excludePre) || ver.PreRelease() == ""
		},
			hiter.Divide(
				func(s string) (exver.Version, error) {
					return exver.Parse(strings.TrimPrefix(s, "v"))
				},
				stringsiter.SplitFunc(string(bin), -1, nil),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	slices.SortFunc(versions, func(i, j exver.Version) int { return -i.Compare(j) })
	if *onlyLatest && len(versions) > 0 {
		fmt.Println(versions[0].String())
		return
	}
	// concat then print at once.
	// If this command is piped to command line `head -n 1`,
	// the downstream existing early may cause "signal: broken pipe".
	fmt.Println(stringsiter.Join("\n", hiter.Map(exver.Version.String, slices.Values(versions))))
}
