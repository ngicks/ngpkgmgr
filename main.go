package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/ioiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"golang.org/x/sync/errgroup"
)

var (
	dir   = flag.String("dir", "", "")
	v     = flag.Bool("v", false, "")
	f     = flag.Bool("f", false, "force option: ignores errors")
	n     = flag.String("new", "", "creates command sets for given name")
	debug = flag.Bool("debug", false, "debug")
)

type namedCommandSet struct {
	Name string
	Set  commandSet
}

type commandSet struct {
	Ver         []string `json:"ver,omitzero"`
	CheckLatest []string `json:"checklatest,omitzero"`
	Install     []string `json:"install,omitzero"`
	Update      []string `json:"update,omitzero"`
	After       []string `json:"after,omitzero"`
}

type command string

const (
	commandVer         command = "ver"
	commandChecklatest command = "checklatest"
	commandInstall     command = "install"
	commandUpdate      command = "update"
)

var cmds = []command{commandVer, commandChecklatest, commandInstall, commandUpdate}

func (c commandSet) Select(kind command) []string {
	switch kind {
	default:
		panic(fmt.Errorf("unknown command: %q", kind))
	case commandVer:
		return c.Ver
	case commandChecklatest:
		return c.CheckLatest
	case commandInstall:
		return c.Install
	case commandUpdate:
		return c.Update
	}
}

type commandExecutor struct {
	dir        string
	commandSet namedCommandSet
	stdin      io.Reader
	stdout     io.Writer
	stderr     io.Writer
}

func newCommandExecutor(
	dir string,
	commandSet namedCommandSet,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) *commandExecutor {
	return &commandExecutor{
		dir:        dir,
		commandSet: commandSet,
		stdin:      stdin,
		stdout:     stdout,
		stderr:     stderr,
	}
}

func (e commandExecutor) Exec(
	ctx context.Context,
	kind command,
	ver string,
	verbose bool,
) (string, error) {
	args := e.commandSet.Set.Select(kind)
	if len(args) > 0 {
		dict := dictReplacer{
			"${VER}":  ver,
			"${OS}":   runtime.GOOS,
			"${ARCH}": runtime.GOARCH,
		}
		args = slices.Collect(dict.Map(slices.Values(args)))
	} else {
		for _, suf := range []string{"", ".sh", ".exe", ".bat", ".ps1"} {
			name := filepath.Join(e.dir, e.commandSet.Name, string(kind)+suf)
			_, err := os.Stat(name)
			if err == nil {
				args = append(slices.Clip(args), name)
				break
			}
		}
		if len(args) == 0 {
			return "", fmt.Errorf("command not found")
		}
	}

	cmd := exec.CommandContext(ctx, args[0])
	if len(args) > 1 {
		cmd.Args = args
	}

	cmd.Stdin = e.stdin

	buf := new(bytes.Buffer)
	if !verbose {
		cmd.Stdout = buf
	} else {
		cmd.Stdout = io.MultiWriter(buf, e.stdout)
	}
	cmd.Stderr = e.stderr

	cmd.Env = append(os.Environ(), "OS="+runtime.GOOS, "ARCH="+runtime.GOARCH)
	if ver != "" {
		cmd.Env = append(cmd.Env, "VER="+ver)
	}

	err := cmd.Run()
	return buf.String(), err
}

const (
	pinnedVersionsFileName = ".pin.json"
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfgDir := *dir

	if cfgDir == "" {
		userCfgDir, err := os.UserConfigDir()
		if err != nil {
			panic(fmt.Errorf("getting os.UserConfigDir: %w", err))
		}
		cfgDir = filepath.Join(userCfgDir, "ngpkgmgr")
	}

	if *n != "" {
		f, err := os.OpenFile(filepath.Join(cfgDir, *n+".json"), os.O_RDWR|os.O_CREATE|os.O_EXCL, fs.ModePerm)
		switch {
		default:
			panic(err)
		case errors.Is(err, fs.ErrExist):
		case err == nil:
			enc := json.NewEncoder(f)
			enc.SetIndent("", "    ")
			err := enc.Encode(commandSet{
				Ver:         []string{},
				Install:     []string{},
				CheckLatest: []string{},
				Update:      []string{},
				After:       []string{},
			})
			_ = f.Close()
			if err != nil {
				panic(err)
			}
		}
		err = os.Mkdir(filepath.Join(cfgDir, *n), fs.ModePerm)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			panic(err)
		}
		for _, c := range cmds {
			scriptName := filepath.Join(cfgDir, *n, string(c))
			switch runtime.GOOS {
			case "windows":
				scriptName += ".ps1"
			default:
				scriptName += ".sh"
			}
			f, err := os.OpenFile(scriptName, os.O_RDWR|os.O_CREATE|os.O_EXCL, fs.ModePerm)
			switch {
			default:
				panic(err)
			case errors.Is(err, fs.ErrExist):
			case err == nil:
				_, err := fmt.Fprintf(f, "#!%s\n", cmp.Or(os.Getenv("SHELL"), "/bin/bash"))
				_ = f.Close()
				if err != nil {
					panic(err)
				}
			}
		}
		return
	}

	var tgt, cmd string
	args := flag.Args()
	switch len(args) {
	case 2:
		tgt, cmd = args[0], args[1]
	case 1:
		cmd = args[0]
	default:
		panic(fmt.Errorf("wrong args length: want 2 or 1, got %d", len(args)))
	}

	if !slices.Contains(cmds, command(cmd)) {
		panic(fmt.Errorf("unknown command: must be one of %v", cmds))
	}

	pinnedVersions := map[string]string{}
	pinFile, err := os.Open(filepath.Join(cfgDir, pinnedVersionsFileName))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			panic(err)
		}
	} else {
		err = json.NewDecoder(pinFile).Decode(&pinnedVersions)
		_ = pinFile.Close()
		if err != nil {
			panic(err)
		}
	}

	for k, v := range pinnedVersions {
		if k != strings.TrimSpace(k) || v != strings.TrimSpace(v) {
			panic(fmt.Errorf("pinned version %q has space prefix and/or suffix in name or version", k))
		}
	}

	var sets []namedCommandSet
	if tgt != "" {
		f, err := os.Open(filepath.Join(cfgDir, tgt+".json"))
		if err == nil {
			var set commandSet
			err = json.NewDecoder(f).Decode(&set)
			_ = f.Close()
			if err != nil {
				panic(err)
			}
			sets = append(sets, namedCommandSet{Name: tgt, Set: set})
		} else if !errors.Is(err, fs.ErrNotExist) {
			panic(err)
		} else {
			s, err := os.Stat(filepath.Join(cfgDir, tgt))
			if err != nil {
				panic(err)
			}
			if !s.IsDir() {
				panic(fmt.Errorf("file %[1]q.json or directory %[1]q must exist", tgt))
			}
			sets = append(sets, namedCommandSet{Name: tgt})
		}
	} else {
		dir, err := os.Open(cfgDir)
		if err != nil {
			panic(err)
		}

		sets, err = hiter.TryAppendSeq(
			sets[:0],
			xiter.Map2(
				func(fi fs.FileInfo, err error) (namedCommandSet, error) {
					switch {
					default:
						return namedCommandSet{}, err
					case fi.Mode().IsRegular() && strings.HasSuffix(fi.Name(), ".json"):
						f, err := os.Open(filepath.Join(cfgDir, fi.Name()))
						if err != nil {
							return namedCommandSet{}, err
						}
						var set commandSet
						err = json.NewDecoder(f).Decode(&set)
						_ = f.Close()
						if err != nil {
							return namedCommandSet{}, err
						}
						return namedCommandSet{Name: strings.TrimSuffix(fi.Name(), ".json"), Set: set}, nil
					case fi.IsDir():
						// directory should contain scripts.
						return namedCommandSet{Name: fi.Name()}, nil
					}
				},
				xiter.Filter2(
					func(fi fs.FileInfo, err error) bool {
						switch {
						default:
							return false
						case err != nil,
							fi.Mode().IsRegular() && strings.HasSuffix(fi.Name(), ".json") && fi.Name() != pinnedVersionsFileName,
							fi.IsDir():
							return true
						}
					},
					ioiter.Readdir(dir),
				),
			),
		)
		_ = dir.Close()
		if err != nil {
			panic(err)
		}
		slices.SortFunc(
			sets,
			func(i, j namedCommandSet) int {
				if c := cmp.Compare(i.Name, j.Name); c != 0 {
					return c
				}
				switch {
				case reflect.ValueOf(i.Set).IsZero():
					// x > y
					return +1
				case reflect.ValueOf(j.Set).IsZero():
					return -1
				default:
					return 0
				}
			},
		)
		// may contain both .json and directory
		sets = slices.CompactFunc(sets, func(i, j namedCommandSet) bool { return i.Name == j.Name })
		sets = topologicalSort(sets)
	}

	if *debug {
		for _, s := range sets {
			fmt.Printf("name = %s, after = %v\n", s.Name, s.Set.After)
		}
		return
	}

	currentVersions := map[string]string{}
	latestVersions := map[string]string{}

	iter := func() iter.Seq[*commandExecutor] {
		return func(yield func(*commandExecutor) bool) {
			for _, set := range sets {
				executor := newCommandExecutor(cfgDir, set, os.Stdin, os.Stdout, os.Stderr)
				if !yield(executor) {
					return
				}
			}
		}
	}

	type targetedExecutor struct {
		tgt      string
		executor *commandExecutor
	}
	var updates []targetedExecutor
	checkVersions := func() {
		gr, gCtx := errgroup.WithContext(ctx)
		gr.SetLimit(5)
		var mu1, mu2 sync.Mutex
		for executor := range iter() {
			gr.Go(func() error {
				out, err := executor.Exec(gCtx, commandVer, "", *v)
				if err != nil || len(out) == 0 {
					if err == nil {
						err = fmt.Errorf("empty output")
					}
					err := fmt.Errorf("ver %q: %w", executor.commandSet.Name, err)
					return err
				}
				mu1.Lock()
				currentVersions[executor.commandSet.Name] = strings.TrimSpace(out)
				mu1.Unlock()
				return nil
			})
			gr.Go(func() error {
				out, err := executor.Exec(gCtx, commandChecklatest, "", *v)
				if err != nil || len(out) == 0 {
					if err == nil {
						err = fmt.Errorf("empty output")
					}
					err = fmt.Errorf("checklatest %q: %w", executor.commandSet.Name, err)
					return err
				}
				mu2.Lock()
				latestVersions[executor.commandSet.Name] = strings.TrimSpace(out)
				mu2.Unlock()
				return nil
			})
		}
		err := gr.Wait()
		if err != nil {
			panic(err)
		}

		for executor := range iter() {
			name := executor.commandSet.Name
			tgt := cmp.Or(pinnedVersions[name], latestVersions[name])
			fmt.Printf("%q: %s -> %s", name, currentVersions[name], tgt)
			if pinnedVersions[name] != "" {
				fmt.Printf("(pinned)")
			}
			if currentVersions[name] == tgt {
				fmt.Printf(": no update\n")
				continue
			}
			updates = append(updates, targetedExecutor{tgt: tgt, executor: executor})
			fmt.Printf("\n")
		}
	}

	switch command(cmd) {
	case commandInstall:
		for executor := range iter() {
			fmt.Printf("installing %q...\n", executor.commandSet.Name)
			out, err := executor.Exec(ctx, commandVer, "", false)
			if err == nil && len(out) > 0 {
				fmt.Printf("Skipping %q: seems already installed at version %s\n", executor.commandSet.Name, strings.TrimSpace(out))
				continue
			}

			out, err = executor.Exec(ctx, commandChecklatest, "", false)
			ver := strings.TrimSpace(out)
			if err != nil {
				ver = ""
				fmt.Printf("fetching latest version failed with err %v\nNow trying with no version specified\n", err)
			}

			_, err = executor.Exec(ctx, commandInstall, cmp.Or(pinnedVersions[executor.commandSet.Name], ver), *v)
			if err != nil {
				err := fmt.Errorf("install %q: %w", executor.commandSet.Name, err)
				if !*f {
					panic(err)
				}
				fmt.Printf("warn: failed: %v\n", err)
			} else {
				fmt.Printf("installing %q done!\n", executor.commandSet.Name)
			}
		}
	case commandVer:
		for executor := range iter() {
			out, err := executor.Exec(ctx, commandVer, "", false)
			if err != nil || len(out) == 0 {
				if err == nil {
					err = fmt.Errorf("empty output")
				}
				err := fmt.Errorf("ver %q: %w", executor.commandSet.Name, err)
				if !*f {
					panic(err)
				}
				fmt.Printf("warn: failed: %v\n", err)
			}
			currentVersions[executor.commandSet.Name] = strings.TrimSpace(out)
		}
		fmt.Printf("%s\n", must(json.MarshalIndent(currentVersions, "", "    ")))
	case commandChecklatest:
		checkVersions()
	case commandUpdate:
		checkVersions()
		for _, t := range updates {
			fmt.Printf("updating %q...\n", t.executor.commandSet.Name)
			_, err := t.executor.Exec(ctx, commandUpdate, t.tgt, *v)
			if err != nil {
				panic(fmt.Errorf("updating %q: %w", t.executor.commandSet.Name, err))
			}
			fmt.Printf("updated %q!\n", t.executor.commandSet.Name)
		}
	}
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

func topologicalSort(s []namedCommandSet) []namedCommandSet {
	type node struct {
		after []*node
		val   namedCommandSet
	}

	nodes := make([]*node, len(s))
	for i, e := range s {
		nodes[i] = &node{val: e}
	}
	for i, n := range nodes {
		for j, nn := range nodes {
			if i == j {
				continue
			}
			if slices.Contains(n.val.Set.After, nn.val.Name) {
				n.after = append(n.after, nn)
			}
		}
	}

	sorted := make([]namedCommandSet, 0, len(s))

	visited := make(map[*node]bool, len(s))
	var visit func(n *node, visited map[*node]bool)
	visit = func(n *node, visited map[*node]bool) {
		if visited[n] {
			return
		}
		visited[n] = true
		for _, nn := range n.after {
			visit(nn, visited)
		}
		sorted = append(sorted, n.val)
	}
	for _, n := range nodes {
		visit(n, visited)
	}

	return sorted
}
