package cmd

import "fmt"

var (
	version = "development"
	commit  = "none"
	date    = "unknown"
)

func PrintVersion() {
	fmt.Printf("Ecwid Image Downloader: %s (commit: %s) (date: %s)\n", version, commit, date)
}
