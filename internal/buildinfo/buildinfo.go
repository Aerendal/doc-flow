package buildinfo

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func FullVersion() string {
	return fmt.Sprintf("%s (commit=%s date=%s)", Version, Commit, Date)
}
