package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	appName = "TEA sh!t ::"
)

var (
	src, dist *string
	rebuild   *bool
)

type abiStruct struct {
	ABI interface{} `json:"abi"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "truffle-export-abi",
		Run: proc,
	}

	src = rootCmd.Flags().String("src", "./json", "Source contract json path")
	dist = rootCmd.Flags().String("dist", "./abi", "ABI file dist path")
	rebuild = rootCmd.Flags().Bool("rebuild", false, "Rebuild contracts")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func proc(c *cobra.Command, args []string) {
	fmt.Println(appName, "Working on", *src)
	if *rebuild {
		buildsPath := *src + "/build"
		fmt.Println(appName, "Remove old files", buildsPath)
		err := os.RemoveAll(buildsPath)
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
	}

	fmt.Println(appName, "Export contracts ABI")
	os.Mkdir(*dist, 0755)
	var abi abiStruct
	filepath.Walk(*src, func(innerPath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".json") {
			fmt.Println(appName, "Open", innerPath)
			f, _ := os.Open(innerPath)
			b, _ := ioutil.ReadAll(f)
			json.Unmarshal(b, &abi)
			names := strings.Split(info.Name(), ".")
			b, _ = json.Marshal(&abi.ABI)
			abiName := *dist + "/" + names[0] + ".abi"
			f, _ = os.Create(abiName)
			f.Write(b)
			fmt.Println(appName, "WriteTo", abiName)
		}
		return nil
	})
}
