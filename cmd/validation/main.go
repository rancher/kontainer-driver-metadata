package main

import (
	"os"

	"github.com/rancher/kontainer-driver-metadata/pkg/validation"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := validation.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
