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
				pathArg := false
				queryArg := false
				formArg := false
				headerArg := false
				description := []string{}
				call = newApiCallHTTP(f)
				lines := strings.Split(c.Text(), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					switch {
					case !httpMode && strings.HasPrefix(line, "#http"):
						if err := call.parseHttpLine(line[6:]); err == nil {
							httpMode = true
							api = append(api, call)
							description = append(description, "\n")
							continue
						}
					case pathArg:
						if strings.HasPrefix(line, ":") {
							if err := call.parsePathArg(line[1:]); err == nil {
								continue
							}
						}
						pathArg = false
					case queryArg:
						if strings.HasPrefix(line, ":") {
							if err := call.parseQueryArg(line[1:]); err == nil {
								continue
							}
						}
						queryArg = false
					case formArg:
						if strings.HasPrefix(line, ":") {
							if err := call.parseFormArg(line[1:]); err == nil {
								continue
							}
						}
						formArg = false
					case headerArg:
					}
					switch {
					case httpMode:
						switch strings.ToLower(line) {
						case "path params:":
							pathArg = true
							queryArg = false
							formArg = false
							headerArg = false
							continue
						case "query params:":
							pathArg = false
							queryArg = true
							formArg = false
							headerArg = false
							continue
						case "form params:":
							pathArg = false
							queryArg = false
							formArg = true
							headerArg = false
							continue
						case "header params:":
							pathArg = false
							queryArg = false
							formArg = false
							headerArg = true
							continue
						}
					}
					if line == "" {
						line = "\n"
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
