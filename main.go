package main

import (
	"flag"

	"temi/pkg/http"
	"temi/pkg/terminal"
)

var debugUrl = flag.String("url", "", "url for debug/vars, e.g. http://foo.com/debug/vars")

func main() {
	flag.Parse()
	if len(*debugUrl) == 0 {
		flag.Usage()
		return
	}

	terminal.Run(http.NewMemStatsLoader(*debugUrl))
}
