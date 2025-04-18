package fileIO

import "os"

type Arguments struct {
	InputPath  PathInfo
	OutputPath PathInfo
}

type PathInfo struct {
	Path  string
	Info  os.FileInfo
	Files []string
}
