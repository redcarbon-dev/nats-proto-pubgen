package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/redcarbon-dev/nats-proto-pubgen/pkg/gen"
)

type renderFile struct {
	path string
	data string
}

func main() {

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	files := iterate(currentDirectory)

	base, err := gen.RenderBaseFile()
	must(err)

	files = append(files, renderFile{
		path: "./pubs/model.go",
		data: base,
	})

	for _, file := range files {
		mustRenderFile(file.path, file.data)
	}
}

func iterate(path string) (files []renderFile) {

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if strings.HasSuffix(info.Name(), ".proto") {
			log.Printf("path: %v, File Name: %s\n", path, info.Name())
			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			res, ok, err := gen.RenderReader(f)
			if err != nil {
				panic(err)
			}
			if ok {
				files = append(files, renderFile{
					path: "./pubs/" + info.Name() + ".pub.go",
					data: res,
				})
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustRenderFile(path string, content string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}
