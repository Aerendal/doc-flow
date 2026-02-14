package lintdays

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Run lints markdown entries in worklog/days and returns process-like exit code:
// 0 = ok, 1 = lint issues found, 2 = usage/runtime error.
func Run(root string, out, errOut io.Writer) int {
	if strings.TrimSpace(root) == "" {
		fmt.Fprintln(errOut, "lint-days: --root cannot be empty")
		return 2
	}

	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(out, "lint-days: root does not exist, skipping: %s\n", root)
			return 0
		}
		fmt.Fprintf(errOut, "lint-days: cannot stat root %q: %v\n", root, err)
		return 2
	}
	if !info.IsDir() {
		fmt.Fprintf(errOut, "lint-days: root is not a directory: %s\n", root)
		return 2
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintf(errOut, "lint-days: cannot read root %q: %v\n", root, err)
		return 2
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.EqualFold(name, "index.md") {
			continue
		}
		if strings.HasSuffix(strings.ToLower(name), ".md") {
			files = append(files, filepath.Join(root, name))
		}
	}
	sort.Strings(files)

	issues := 0
	for _, p := range files {
		ok, line, readErr := hasHeadingLine(p)
		if readErr != nil {
			fmt.Fprintf(errOut, "lint-days: read error %s: %v\n", p, readErr)
			issues++
			continue
		}
		if !ok {
			fmt.Fprintf(errOut, "lint-days: missing H1 heading in %s (first line: %q)\n", p, line)
			issues++
		}
	}

	if issues > 0 {
		fmt.Fprintf(errOut, "lint-days: found %d issue(s)\n", issues)
		return 1
	}

	fmt.Fprintf(out, "lint-days: ok (%d file(s))\n", len(files))
	return 0
}

func hasHeadingLine(path string) (ok bool, firstLine string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		if sc.Err() != nil {
			return false, "", sc.Err()
		}
		return false, "", nil
	}
	first := strings.TrimSpace(sc.Text())
	return strings.HasPrefix(first, "# "), first, nil
}
