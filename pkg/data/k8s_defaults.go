package data

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

const (
	DataFilePath = "./data/data.json"
)

type Data struct {
	K3S  map[string]interface{} `json:"k3s,omitempty"`
	RKE2 map[string]interface{} `json:"rke2,omitempty"`
}

var (
	DriverData     Data
	MissedTemplate map[string][]string
)

func initData() {
	DriverData = Data{
		K3S:  map[string]interface{}{},
		RKE2: map[string]interface{}{},
	}

	if err := readFile("./channels.yaml", DriverData.K3S); err != nil {
		panic(err)
	}
	if err := readFile("./channels-rke2.yaml", DriverData.RKE2); err != nil {
		panic(err)
	}
}

func readFile(input string, data map[string]interface{}) error {
	bytes, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, &data)
}

func GenerateData() {
	initData()

	if len(MissedTemplate) != 0 {
		logrus.Warnf("found k8s versions without a template")
		for plugin, data := range MissedTemplate {
			logrus.Warnf("no %s template for k8sVersions %v \n", plugin, data)
		}
	}

	// todo: zip file
	strData, _ := json.MarshalIndent(DriverData, "", " ")
	jsonFile, err := os.Create(DataFilePath)
	if err != nil {
		panic(fmt.Errorf("err creating data file %v", err))
	}
	defer jsonFile.Close()
	_, err = jsonFile.Write(strData)
	if err != nil {
		panic(fmt.Errorf("err writing jsonFile %v", err))
	}
	fmt.Println("finished generating data.json")
}

func FromData(b []byte) (Data, error) {
	d := &Data{}
	if err := json.Unmarshal(b, d); err != nil {
		return Data{}, err
	}
	return *d, nil
}
