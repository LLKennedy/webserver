package network

import (
	"path"
	"strings"
)

func getPathNode(pathString string) (head, tail string) {
	pathString = path.Clean("/" + pathString)
	nodeIndex := strings.Index(pathString[1:], "/") + 1
	if nodeIndex <= 0 {
		return pathString[1:], "/"
	}
	return pathString[1:nodeIndex], pathString[nodeIndex:]
}
