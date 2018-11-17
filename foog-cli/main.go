package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var deployLang = "go"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Expected file name")
			return
		}
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("Failed to open file")
			return
		}
		defer file.Close()
		url := os.Getenv("FOOG_SERVER") + "/deploy?lang=" + strings.TrimPrefix(filepath.Ext(args[0]), ".")
		resp, err := http.Post(url, "application/octet-stream", file)
		if err != nil {
			fmt.Println("Failed to send deploy request")
			return
		}
		defer resp.Body.Close()
		respURL := struct {
			URL string
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respURL)
		fmt.Printf("Your app is now running on %s%s\n", os.Getenv("FOOG_SERVER"), respURL.URL)
	},
}

var rootCmd = &cobra.Command{
	Use:   "foog-cli",
	Short: "foog-cli is a serverless toolkit",
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&deployLang, "lang", "l", "go", "Language container to deploy to")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
