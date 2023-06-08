package utiliies

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rancher/kontainer-driver-metadata/pkg/rke"
	"github.com/rancher/rke/types/kdm"
)

const (
	RKE  = "RKE"
	RKE2 = "RKE2"
	K3S  = "K3S"
)

// FromLocalFile loads and returns the data object from the given path to a data.json file.
// It returns an empty data and an error when something goes wrong.
func FromLocalFile() (kdm.Data, error) {
	return FromFile(rke.DataFilePath)
}

// FromFile loads and returns the data object from the given path to a data.json file.
// It returns an empty data and an error when something goes wrong.
func FromFile(path string) (kdm.Data, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return kdm.Data{}, err
	}
	return kdm.FromData(file)
}

// FromURL downloads and returns the data object from the given URL.
// It returns an empty data and an error when something goes wrong.
func FromURL(url string) (kdm.Data, error) {
	resp, err := http.Get(url)
	if err != nil {
		return kdm.Data{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return kdm.Data{}, err
	}
	return kdm.FromData(data)
}

// DownloadFromURL download the file from the provided url, and returns its content as a string.
// It returns an empty string and an error when something goes wrong.
func DownloadFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error code %d and body %s", resp.StatusCode, body)
	}
	if err != nil {
		return "", err
	}

	return string(body), nil
}
