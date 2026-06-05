package main

import (
	"github.com/rancher/kontainer-driver-metadata/pkg/data"
	"github.com/rancher/kontainer-driver-metadata/pkg/images"
)

func main() {
	// add drivers init here
	data.GenerateData()
	images.GenerateRegSyncFile()
}
