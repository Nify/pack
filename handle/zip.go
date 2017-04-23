package handle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"archive/zip"

	"io"

	. "github.com/Nify/pack/tools"
)

var wg sync.WaitGroup

func ReadFiles(dir string) []string {
	buf := new(bytes.Buffer)
	buf.Grow(1024)
	wg.Add(1)
	go readSubFile(buf, dir)
	wg.Wait()
	names := strings.Split(string(buf.Bytes()), "\n")
	names = names[0 : len(names)-1]
	return names
}
func readSubFile(buf *bytes.Buffer, name string) {
	defer wg.Done()
	f, err := os.Open(name)
	HandleErr(fmt.Sprintf("open %s ", name), err)
	defer f.Close()
	fi, err := f.Stat()
	HandleErr("get fileinfo : ", err)
	if !fi.IsDir() {
		if _, err := buf.WriteString(name + "\n"); err != nil {
			HandleErr(fmt.Sprintf("write %s to buf ", name), err)
		}
		return
	}
	newfis, err := f.Readdir(-1)
	HandleErr("read "+name+"dir fileinfo : ", err)
	if len(newfis) == 0 {
		if _, err := buf.WriteString(name + "\n"); err != nil {
			HandleErr(fmt.Sprintf("write %s to buf ", name), err)
		}
		return
	}
	for _, newfi := range newfis {
		wg.Add(1)
		go readSubFile(buf, filepath.Join(name, newfi.Name()))
	}
}

// AddAllZip add current working direction's subdirs and subfiels  to zip archive
func AddAllZip(zipname string) {
	dir, err := os.Getwd()
	HandleErr("get work dir  ", err)
	names := ReadFiles(dir)
	if zipname == "" {
		zipname = filepath.Base(dir)
		zipname += ".zip"
	}
	dir += "/"
	zipFile, err := os.Create(zipname)
	HandleErr("init archive zip file ", err)
	defer zipFile.Close()
	w := zip.NewWriter(zipFile)
	defer func() {
		err = w.Close()
		if err != nil {
			HandleErr("close archive zip ", err)
		}
	}()
	for _, name := range names {
		AddFile(w, dir, name)
	}
}

// AddFile add a file to archive
func AddFile(w *zip.Writer, dir, name string) {
	cleanName := strings.TrimPrefix(name, dir)
	f, err := os.Open(name)
	HandleErr("open file ", err)
	defer f.Close()
	fi, err := f.Stat()
	HandleErr("get file fileinfo : ", err)
	header, err := zip.FileInfoHeader(fi)
	HandleErr("get zip fileheader ", err)
	header.Name = cleanName
	if fi.IsDir() {
		header.Name += "/"
		header.Method = zip.Store
		w.CreateHeader(header)
		return
	}
	zw, err := w.CreateHeader(header)
	HandleErr(" create archive writer ", err)
	_, err = io.Copy(zw, f)
	HandleErr("copy file to archive file ", err)
}
