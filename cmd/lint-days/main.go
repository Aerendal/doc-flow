package main

import (
	"flag"
	"os"

	"docflow/internal/dev/worklog/lintdays"
)

func main() {
	root := flag.String("root", "worklog/days", "root directory containing markdown files")
	flag.Parse()
	os.Exit(lintdays.Run(*root, os.Stdout, os.Stderr))
}
