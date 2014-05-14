package main

import (
	"bufio"
	"fmt"
	"github.com/docopt/docopt.go"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type CliString struct {
	cli string
}

type SelfPath struct {
	Path string
}

const usageTemplate string = `
Usage: {{.Path}} (a|x|b|d|r) [-0] [--] [-|PATH...]

{{.Path}} is a simple file/dir path parsing tool. 

Commands:
  a         Absolute path
  x         File extension
  b         Base name
  d         Parent Directory
  r         Real path of symlink

Arguments:
PATH        path of file or directory

Options:
  -h        Show this help
  -0        Paths are seperated by NUL charater`

var logger = log.New(os.Stderr, "", 0)

func main() {
	arguments := parseCli()
	var delim byte
	if arguments["-0"].(bool) {
		delim = '\x00'
	} else {
		delim = '\n'
	}
	pathChan := make(chan string, 0)
	defer close(pathChan)
	quitChan := make(chan int, 0)
	defer close(quitChan)
	go func() {
		for {
			select {
			case p := <-pathChan:
				switch {
				case arguments["a"].(bool):
					abs, err := filepath.Abs(p)
					die(err)
					fmt.Println(abs)
				case arguments["b"].(bool):
					fmt.Println(filepath.Base(p))
				case arguments["x"].(bool):
					fmt.Println(filepath.Ext(p))
				case arguments["d"].(bool):
					fmt.Println(filepath.Dir(p))
				case arguments["r"].(bool):
					realp, err := filepath.EvalSymlinks(p)
					die(err)
					fmt.Println(realp)
				}
			case <-quitChan:
				return
			}
		}
	}()
	paths := arguments["PATH"].([]string)
	if arguments["-"].(bool) || len(paths) == 0 {
		stdinReader := bufio.NewReader(os.Stdin)
		for {
			p, e := stdinReader.ReadString(delim)
			if e != nil && p == "" {
				break
			}
			pathChan <- strings.TrimSuffix(p, string(delim))
		}
	}
	for _, p := range paths {
		pathChan <- p
	}
	quitChan <- 1
}

func parseCli() map[string]interface{} {
	selfname := SelfPath{filepath.Base(os.Args[0])}

	clitmpl := template.Must(template.New("usageTemplate").Parse(usageTemplate))

	usage := &CliString{}
	err := clitmpl.Execute(usage, selfname)
	die(err)
	arguments, _ := docopt.Parse(usage.cli, nil, true, "", false)
	return arguments
}

func die(e error) {
	if e != nil {
		logger.Fatal(e)
	}
}

func (s *CliString) Write(p []byte) (int, error) {
	s.cli += string(p)
	return len(p), nil
}

func (s *CliString) String() string {
	return s.cli
}
