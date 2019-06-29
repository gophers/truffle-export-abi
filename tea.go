package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	appName = "TEA sh!t ::"
)

type abiStruct struct {
	ABI interface{} `json:"abi"`
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(appName, "Error", err)
		os.Exit(-1)
	}
	fmt.Println(appName, "Working on", path)

	buildsPath := path + "/build"
	fmt.Println(appName, "Remove old files", buildsPath)
	err = os.RemoveAll(buildsPath)
	if err != nil {
		fmt.Println(appName, "Error", err)
	}

	fmt.Println(appName, "Rebuild contracts ...")
	cmd := exec.Command("truffle", "compile")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(appName, "Error", err)
		os.Exit(-1)
	}
	fmt.Println(appName, "Rebuild contracts done")

	fmt.Println(appName, "Export contracts ABI")
	abiPath := buildsPath + "/abi"
	os.Mkdir(abiPath, 0755)
	var abi abiStruct
	filepath.Walk(buildsPath+"/contracts", func(innerPath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".json") {
			fmt.Println(appName, "Open", innerPath)
			f, _ := os.Open(innerPath)
			b, _ := ioutil.ReadAll(f)
			json.Unmarshal(b, &abi)
			names := strings.Split(info.Name(), ".")
			b, _ = json.Marshal(&abi.ABI)
			abiName := abiPath + "/" + names[0] + ".abi"
			f, _ = os.Create(abiName)
			f.Write(b)
			fmt.Println(appName, "WriteTo", abiName)
		}
		return nil
	})
}
