//go:build ignore

package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	b, err := exec.Command("go", "env", "GOROOT").Output()
	catch(err)
	goroot := string(bytes.TrimSpace(b))

	zipname := filepath.Join(goroot, "lib", "time", "zoneinfo.zip")
	zipf, err := os.Open(zipname)
	catch(err)
	zipfi, err := zipf.Stat()
	catch(err)

	f, err := os.Create(os.Getenv("GOFILE"))
	catch(err)

	fmt.Fprintln(f, `// Code generated by "mktzlist.go"`)
	fmt.Fprintln(f, "//go:generate go run mktzlist.go")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "package %s\n", os.Getenv("GOPACKAGE"))
	fmt.Fprintln(f)
	fmt.Fprintln(f, "var TimeZones = []string{")

	zr, err := zip.NewReader(zipf, zipfi.Size())
	catch(err)
	err = fs.WalkDir(zr, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Fprintf(f, "\t%q,\n", path)
		}
		return nil
	})
	catch(err)

	fmt.Fprintln(f, "}")
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
