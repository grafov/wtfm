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
	"fmt"
	"io"
	"sort"
	"strings"
)

var ln = []byte("\n")

func makeReST(out io.Writer) {
	out.Write([]byte(".. autogenerated by wtfm\n\n"))
	sort.Sort(api)
	// mk := make([]string, len(api))
	// i := 0
	// for k, _ := range api {
	// 	mk[i] = k
	// 	i++
	// }
	// sort.Strings(mk)
	for _, call := range api {
		//		call := api[key]
		out.Write([]byte(fmt.Sprintf("\n\n%s\n=%s\n", call.Title, strings.Repeat("=", len(call.Title)))))
		out.Write([]byte(fmt.Sprintf("\n``%s %s``\n", call.Methods, call.Path)))
		if len(call.PathParams) > 0 {
			out.Write([]byte("\nParameters in a path:\n\n"))
			out.Write([]byte("\n.. glossary::\n"))
			pk := make([]string, len(call.PathParams))
			i := 0
			for k, _ := range call.PathParams {
				pk[i] = k
				i++
			}
			sort.Strings(pk)
			for _, key := range pk {
				param := call.PathParams[key]
				out.Write([]byte(fmt.Sprintf("  ``%s`` %s\n    %s\n", param.Name, param.Type, param.Desc)))
			}
		}
		if len(call.QueryParams) > 0 {
			out.Write([]byte("\nParameters of a query string:\n\n"))
			out.Write([]byte("\n.. glossary::\n"))
			pk := make([]string, len(call.QueryParams))
			i := 0
			for k, _ := range call.QueryParams {
				pk[i] = k
				i++
			}
			sort.Strings(pk)
			for _, key := range pk {
				param := call.QueryParams[key]
				out.Write([]byte(fmt.Sprintf("  ``%s`` %s\n    %s\n", param.Name, param.Type, param.Desc)))
			}
		}
		if len(call.FormParams) > 0 {
			out.Write([]byte("\nParameters in a body of a request:\n\n"))
			out.Write([]byte("\n.. glossary::\n"))
			pk := make([]string, len(call.FormParams))
			i := 0
			for k, _ := range call.FormParams {
				pk[i] = k
				i++
			}
			sort.Strings(pk)
			for _, key := range pk {
				param := call.FormParams[key]
				out.Write([]byte(fmt.Sprintf("  ``%s`` %s\n    %s\n", param.Name, param.Type, param.Desc)))
			}
		}
		out.Write(ln)
		if len(call.Desc) > 1 {
			for _, l := range call.Desc[1:] {
				out.Write([]byte(l))
				out.Write(ln)
			}
		}
	}
}
