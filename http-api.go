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
	"go/ast"
	"net/http"
	"strings"
)

var api API

type API []*apiCallHTTP

// Len is part of sort.Interface.
func (s API) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s API) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface.
func (s API) Less(i, j int) bool {
	return s[i].Title < s[j].Title
}

// represents HTTP API
type apiCallHTTP struct {
	Title   string
	Methods []string
	Path    string
	Params  map[string]map[string]*param // kind - name - param
	Headers http.Header
	Context *ast.File
	Desc    []string
}

type param struct {
	Name     string
	Kind     string
	Type     string
	Required bool
	Desc     string
}

func newApiCallHTTP(context *ast.File) *apiCallHTTP {
	call := new(apiCallHTTP)
	call.Context = context
	call.Params = make(map[string]map[string]*param)
	call.Headers = make(http.Header)
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
func (c *apiCallHTTP) parseArg(kind, s string) (string, error) {
	par := new(param)
	par.Kind = kind
	par.Required = true            // TODO parse REQUIRED
	parts := strings.Split(s, ":") // separate name-type from description
	varparts := strings.Split(parts[0], " ")
	par.Name = varparts[0]
	if len(varparts) > 1 {
		par.Type = varparts[1]
	} else {
		par.Type = "string"
	}
	if len(parts) > 1 {
		par.Desc = strings.Join(parts[1:], ":")
	}
	if par.Name == "" {
		return "", errors.New("empty par definition")
	}
	if _, ok := c.Params[kind]; !ok {
		c.Params[kind] = make(map[string]*param)
	}
	c.Params[kind][strings.ToLower(par.Name)] = par
	return par.Name, nil
}

// parse form value description
func (c *apiCallHTTP) parseHeader(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) < 2 {
		return errors.New("parsing error for header")
	}
	if strings.TrimSpace(parts[1]) == "" {
		return errors.New("empty value for header")
	}
	c.Headers.Add(parts[0], parts[1])
	return nil
}
