package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var deploymentDetailsTemplate = template.Must(template.New("details").Parse(`ID:      {{ .ID }}
Created: {{ .Date }}
URL:     {{ .URL }}
`))

var foogEndpoint = os.Getenv("FOOG_SERVER")

func errorAndExit(r io.Reader) {
	errBody := struct {
		Error   bool
		Message string
	}{}
	decoder := json.NewDecoder(r)
	decoder.Decode(&errBody)
	if !errBody.Error {
		return
	}
	fmt.Printf("error: %s\n", errBody.Message)
	os.Exit(1)
}

var deployLang = "go"

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Manage builds and images",
}

var buildLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs of image build",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		url := fmt.Sprintf("%s/logs/%s", foogEndpoint, id)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Failed to retrieve build logs")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		io.Copy(os.Stdout, resp.Body)
	},
}

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Use short-names for deployments",
}

var aliasBindCmd = &cobra.Command{
	Use:   "bind [id] [alias]",
	Short: "Bind an alias to a deployment",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, alias := args[0], args[1]
		url := fmt.Sprintf("%s/bind/%s?to=%s", foogEndpoint, id, alias)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println("Failed to bind to alias")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		respBody := struct {
			URL string
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respBody)
		fmt.Printf("Alias bound to %s\n", respBody.URL)
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all bound aliases",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		url := foogEndpoint + "/listAlias"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Failed to list deployments")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		respBody := []struct {
			Name string
			For  string
			URL  string
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respBody)
		fmt.Printf("%-16.16s %-8.8s %s\n", "Name", "For", "URL")
		for _, a := range respBody {
			fmt.Printf("%-16.16s %-8.8s %s\n", a.Name, a.For, a.URL)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployed applications",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		url := foogEndpoint + "/list"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Failed to list deployments")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		respBody := []struct {
			ID    string
			URL   string
			Ready bool
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respBody)

		fmt.Printf("%-8s %-6s %s\n", "ID", "Ready", "URL")
		for _, item := range respBody {
			fmt.Printf("%-8.8s %-6v %s\n", item.ID, item.Ready, item.URL)
		}
	},
}

const describeEndpoint = "%s/describe/%s"

var describeCmd = &cobra.Command{
	Use:   "describe [id | @alias]",
	Short: "View deployment details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf(describeEndpoint, foogEndpoint, args[0]))
		if err != nil {
			fmt.Println("Failed to fetch deployment details")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		respBody := struct {
			Date time.Time
			URL  string
			ID   string
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respBody)
		deploymentDetailsTemplate.Execute(os.Stdout, &respBody)
	},
}

const deployEndpoint = "%s/deploy?lang=%s"

var deployCmd = &cobra.Command{
	Use:   "deploy [file]",
	Short: "Deploy an application to the cloud",
	Args:  cobra.ExactArgs(1),
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

		fileLang := strings.TrimPrefix(filepath.Ext(args[0]), ".")
		if deployLang != "" {
			fileLang = deployLang
		}
		resp, err := http.Post(fmt.Sprintf(deployEndpoint, foogEndpoint, fileLang), "text/plain", file)
		if err != nil {
			fmt.Println("Failed to send deploy request")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errorAndExit(resp.Body)
		}
		respURL := struct {
			URL string
		}{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&respURL)
		fmt.Printf("Your app is now running on %s%s\n", foogEndpoint, respURL.URL)
	},
}

var rootCmd = &cobra.Command{
	Use:   "foog-cli",
	Short: "foog-cli is a serverless toolkit",
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(buildCmd)
	aliasCmd.AddCommand(aliasBindCmd)
	aliasCmd.AddCommand(aliasListCmd)
	buildCmd.AddCommand(buildLogsCmd)
	deployCmd.Flags().StringVarP(&deployLang, "lang", "l", "", "Language container to deploy to")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
