package main

import (
	"github.com/rancher/kontainer-driver-metadata/pkg/images"
	"github.com/rancher/kontainer-driver-metadata/pkg/rke"
)

func main() {
	// add drivers init here
	rke.GenerateData()
	images.GenerateRegSyncFile()
}
