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
	case "build":
		makeReST(os.Stdout, build(path))
	case "serve":
	}
}

// do `build` command
func build(path string) API {
	var (
		call *apiCallHTTP
		api  API
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
						// TODO выбрать по одному варианту имен, выкинуть лишние
						case "url params:", "path params:", "url args:", "path args:":
							pathArg = true
							queryArg = false
							formArg = false
							headerArg = false
							continue
						case "query params:", "query args:":
							pathArg = false
							queryArg = true
							formArg = false
							headerArg = false
							continue
						case "form params:", "form values:":
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
					description = append(description, line)
				}
				call.Desc = description
			}
		}
	}
	return api
}
