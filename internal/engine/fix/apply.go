package fix

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ApplyPlan(plan Plan, opts ApplyOptions) error {
	for _, f := range plan.Files {
		target := absPath(opts.Root, f.Path)
		currentBytes, err := os.ReadFile(target)
		if err != nil {
			return fmt.Errorf("read %s: %w", f.Path, err)
		}
		current := string(currentBytes)
		if checksum(current) != f.OriginalChecksum {
			return fmt.Errorf("file changed since plan creation: %s", f.Path)
		}

		if opts.BackupDir != "" {
			backupPath := filepath.Join(opts.BackupDir, f.Path)
			if err := copyFile(target, backupPath); err != nil {
				return fmt.Errorf("backup %s: %w", f.Path, err)
			}
		}

		if err := writeFileAtomic(target, []byte(f.UpdatedContent)); err != nil {
			return fmt.Errorf("apply %s: %w", f.Path, err)
		}
	}
	return nil
}

func writeFileAtomic(target string, data []byte) error {
	st, err := os.Stat(target)
	if err != nil {
		return err
	}

	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, ".docflow-fix-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tmpName, st.Mode().Perm()); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpName, target); err != nil {
		cleanup()
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
