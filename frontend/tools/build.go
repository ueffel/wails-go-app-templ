//go:build ignore

package main

import (
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func CopyFS(srcFS fs.FS, dstDir string, excludes ...string) error {
	excluded := make([]string, 0, len(excludes))
	for _, exclude := range excludes {
		matches, err := fs.Glob(srcFS, exclude)
		if err != nil {
			return err
		}
		excluded = append(excluded, matches...)
	}

	return fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip if excluded
		if slices.Contains(excluded, path) {
			return nil
		}
		fpath, err := filepath.Localize(path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dstDir, fpath)
		if d.IsDir() {
			return os.MkdirAll(newPath, 0777)
		}
		if !d.Type().IsRegular() {
			return &os.PathError{Op: "CopyFS", Path: path, Err: os.ErrInvalid}
		}
		r, err := srcFS.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		info, err := r.Stat()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(newPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666|info.Mode()&0777)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return &os.PathError{Op: "Copy", Path: newPath, Err: err}
		}
		return w.Close()
	})
}

func main() {
	// Compile WASM module
	cmd := exec.Command("go", "build", "-o", "dist/web/app.wasm", "./src")
	cmd.Env = append(os.Environ(), "GOARCH=wasm", "GOOS=js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	must(err)

	// Generate static assets for go-app application
	oldCwd, err := os.Getwd()
	must(err)
	err = os.Chdir("dist")
	cmd = exec.Command("go", "run", "../src")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	must(err)
	err = os.Chdir(oldCwd)
	must(err)

	// Copy frontend to dist overwriting existing files excluding *.go files
	err = CopyFS(os.DirFS("src"), "dist", "components", "components/*", "*.go")
	must(err)
}
