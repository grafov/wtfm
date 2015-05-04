package main

import (
	"fmt"
	"io"
	"strings"
)

func makeReST(out io.Writer, api API) {
	out.Write([]byte("==============\nHTTP API list\n==============\n\n")) // TODO remove nails
	for _, call := range api {
		out.Write([]byte(fmt.Sprintf("\n\n%s %s\n=====%s\n\n", call.Methods, call.Path, strings.Repeat("=", len(call.Path)))))
		out.Write([]byte("Defined in a package: "))
		out.Write([]byte(call.Context.Name.String()))
		if len(call.PathParams) > 0 {
			out.Write([]byte("\n\nParameters in a path:\n"))
			for _, param := range call.PathParams {
				out.Write([]byte(fmt.Sprintf("\n:``%s`` %s:\n%s\n", param.Name, param.Type, param.Desc)))
			}
		}
		if len(call.QueryParams) > 0 {
			out.Write([]byte("\n\nQuery arguments:\n"))
			for _, param := range call.QueryParams {
				out.Write([]byte(fmt.Sprintf("\n:``%s`` %s:\n%s\n", param.Name, param.Type, param.Desc)))
			}
		}
		if len(call.FormParams) > 0 {
			out.Write([]byte("\n\nForm values:\n"))
			for _, param := range call.FormParams {
				out.Write([]byte(fmt.Sprintf("\n:``%s`` %s:\n%s\n", param.Name, param.Type, param.Desc)))
			}
		}
	}
}
