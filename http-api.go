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
	"errors"
	"fmt"
	"go/ast"
	"net/http"
	"sort"
	"strings"
)

var api API = make(map[string]*apiCallHTTP, 0)

type API map[string]*apiCallHTTP

func (a API) Set(call *apiCallHTTP) {
	if !sort.StringsAreSorted(call.Methods) {
		sort.Strings(call.Methods)
	}
	a[fmt.Sprintf("%s%s", call.Path, call.Methods)] = call
}

// represents HTTP API
type apiCallHTTP struct {
	Methods      []string
	Path         string
	PathParams   map[string]*param
	QueryParams  map[string]*param
	FormParams   map[string]*param
	HeaderParams map[string]*param
	Headers      []http.Header
	Context      *ast.File
	Desc         []string
}

type param struct {
	Name     string
	Kind     parType
	Type     string
	Required bool
	Desc     string
}

type parType uint8

const (
	UrlParam parType = iota
	QueryParam
	FormParam
	HeaderParam
)

func newApiCallHTTP(context *ast.File) *apiCallHTTP {
	call := new(apiCallHTTP)
	call.Context = context
	call.PathParams = make(map[string]*param)
	call.QueryParams = make(map[string]*param)
	call.FormParams = make(map[string]*param)
	call.HeaderParams = make(map[string]*param)
	return call
}

func (c *apiCallHTTP) parseHttpLine(s string) error {
	parts := strings.Split(s, " ")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if part[0] == '/' {
			c.Path = strings.Join(parts[i:], " ")
			break
		}
		c.Methods = append(c.Methods, part)
	}
	if len(c.Methods) == 0 {
		c.Methods = []string{"GET"}
	}
	if len(c.Path) == 0 {
		return errors.New("empty path for HTTP call")
	}
	return nil
}

// parse path argument description
func (c *apiCallHTTP) parsePathArg(s string) error {
	param := new(param)
	param.Kind = UrlParam
	param.Required = true          // TODO parse REQUIRED
	parts := strings.Split(s, ":") // separate name-type from description
	varparts := strings.Split(parts[0], " ")
	param.Name = varparts[0]
	if len(varparts) > 1 {
		param.Type = varparts[1]
	} else {
		param.Type = "string"
	}
	if len(parts) > 1 {
		param.Desc = strings.Join(parts[1:], ":")
	}
	if param.Name == "" {
		return errors.New("empty path param definition")
	}
	c.PathParams[strings.ToLower(param.Name)] = param
	return nil
}

// parse query argument description
func (c *apiCallHTTP) parseQueryArg(s string) error {
	param := new(param)
	param.Kind = QueryParam
	param.Required = false         // TODO parse REQUIRED
	parts := strings.Split(s, ":") // separate name-type from description
	varparts := strings.Split(parts[0], " ")
	param.Name = varparts[0]
	if len(varparts) > 1 {
		param.Type = varparts[1]
	} else {
		param.Type = "string"
	}
	if len(parts) > 1 {
		param.Desc = strings.Join(parts[1:], ":")
	}
	if param.Name == "" {
		return errors.New("empty query arg definition")
	}
	c.QueryParams[strings.ToLower(param.Name)] = param
	return nil
}

// parse form value description
func (c *apiCallHTTP) parseFormArg(s string) error {
	// TODO parse REQUIRED
	param := new(param)
	param.Kind = FormParam
	parts := strings.Split(s, ":") // separate name-type from description
	varparts := strings.Split(parts[0], " ")
	param.Name = varparts[0]
	if len(varparts) > 1 {
		param.Type = varparts[1]
	} else {
		param.Type = "string"
	}
	if len(parts) > 1 {
		param.Desc = strings.Join(parts[1:], ":")
	}
	if param.Name == "" {
		return errors.New("empty form value definition")
	}
	c.FormParams[strings.ToLower(param.Name)] = param
	return nil
}
