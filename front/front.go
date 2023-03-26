package front

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed build/*
var Build embed.FS

func GetBuildFS() (http.FileSystem, error) {
	fsys, err := fs.Sub(Build, "build")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}
