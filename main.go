package main

/* Copyleft 2015 Alexander I.Grafov aka Axel <grafov@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	var (
		path string = "."
	)
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("command missed")
		flag.Usage()
		os.Exit(1)
	}
	if len(flag.Args()) > 1 {
		path = flag.Arg(1)
	}

	switch flag.Arg(0) {
	case "parse": // analyze sources and make rest/markdown
		parse(path)
		makeReST(os.Stdout)
	case "build": // generate html output
	case "serve": // serve html pages
	}
}

// do `parse` command
func parse(path string) {
	parseDirTree(path)
}

// Recursively parse tree of nested packages
func parseDirTree(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			parsePackage(path)
		}
		return nil
	})
}

// Non recursively parse files of a single package
func parsePackage(path string) API {
	var (
		call *apiCallHTTP
		last string
		err  error
	)
	fset := token.NewFileSet() // positions are relative to fset
	pkgs, err := parser.ParseDir(fset, path, func(f os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, p := range pkgs {
		for _, f := range p.Files {
			for _, c := range f.Comments {
				httpMode := false
				isArg := false
				isHeader := false
				description := []string{}
				call = newApiCallHTTP(f)
				lines := strings.Split(c.Text(), "\n")
				var kind string
				for _, line := range lines {
					sline := strings.TrimSpace(line)
					switch {
					case !httpMode && strings.HasPrefix(sline, "#http"):
						if err := call.parseHttpLine(sline[6:]); err == nil {
							httpMode = true
							api = append(api, call)
							description = append(description, "\n")
							continue
						}
					case isArg:
						switch {
						case strings.HasPrefix(sline, ":"):
							if last, err = call.parseArg(kind, sline[1:]); err == nil {
								continue
							}
							isArg = false
						case sline == "":
							isArg = false
						default:
							call.Params[kind][last].Desc = fmt.Sprintf("%s %s", call.Params[kind][last].Desc, sline)
							continue
						}
					case isHeader:
						switch {
						case strings.Contains(sline, ":"):
							if err = call.parseHeader(sline); err == nil {
								continue
							}
							isHeader = false
						case sline == "":
							isHeader = false
						default:
							continue
						}
					}
					lowsline := strings.ToLower(sline)
					if httpMode && lowsline == "headers:" {
						isArg = false
						isHeader = true
						continue
					}
					if httpMode && strings.HasSuffix(lowsline, "params:") {
						isArg = true
						isHeader = false
						// kind = strings.SplitN(lowsline, " ", 2)[0]
						kind = sline
						continue
					}
					if sline == "" {
						line = "\n"
						isArg = false
						isHeader = false
					}
					description = append(description, line)
				}
				if !sort.StringsAreSorted(call.Methods) {
					sort.Strings(call.Methods)
				}
				call.Desc = description
				if call.Desc == nil {
					call.Desc = []string{"\n"}
				}
				if call.Desc[0] == "\n" {
					call.Title = fmt.Sprintf("%s %s", call.Path, call.Methods)
				} else {
					call.Title = call.Desc[0]
				}
			}
		}
	}
	return api
}
