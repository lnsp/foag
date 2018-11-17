package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const goDockerfile = `
FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
`

const swiftDockerfile = `
FROM swift:latest

WORKDIR /app
COPY . .
RUN swiftc -o /app/main /app/main.swift

CMD ["/app/main"]
`

const jsDockerfile = `
FROM node:8-alpine

WORKDIR /app
COPY . .

CMD ["node", "/app/main.go"]
`

const cDockerfile = `
FROM gcc:4.9

WORKDIR /app
COPY . .
RUN gcc -o /app/main /app/main.c

CMD ["/app/main"]
`

const TriggerPrefix = "/trigger/"
const DeploymentTimeout = 30 * time.Second
const DeploymentImage = "foog-%16.16s"

var (
	dockerfiles = map[string][]byte{
		"go":    []byte(goDockerfile),
		"swift": []byte(swiftDockerfile),
		"js":    []byte(jsDockerfile),
		"c":     []byte(cDockerfile),
	}
)

type Registry struct {
	Items map[string]*Deployment
}

func NewRegistry() *Registry {
	return &Registry{
		Items: make(map[string]*Deployment),
	}
}

func (registry *Registry) Deploy(lang string, source []byte) (*Deployment, error) {
	deployment := NewDeployment(source, lang)
	if err := deployment.Build(); err != nil {
		return nil, fmt.Errorf("failed to build: %v", err)
	}
	registry.Items[deployment.ID] = deployment
	return deployment, nil
}

func (registry *Registry) Run(id string, stdin io.Reader, stdout, stderr io.Writer) error {
	deployment, ok := registry.Items[id]
	if !ok {
		return fmt.Errorf("deployment not found")
	}
	err := deployment.Run(stdin, stdout, stderr)
	if err != nil {
		return err
	}
	return nil
}

type Deployment struct {
	ID       string
	Source   []byte `json:"-"`
	Image    string `json:"-"`
	Date     time.Time
	Language string
	URL      string
	Ready    bool
}

func NewDeployment(source []byte, lang string) *Deployment {
	hasher := sha256.New()
	hasher.Write(source)
	sum := hasher.Sum(nil)
	id := hex.EncodeToString(sum)
	return &Deployment{
		Source:   source,
		ID:       id,
		Image:    fmt.Sprintf(DeploymentImage, id),
		Language: lang,
		Date:     time.Now(),
		URL:      TriggerPrefix + id,
		Ready:    false,
	}
}

func (d *Deployment) Run(stdin io.Reader, stdout, stderr io.Writer) error {
	if !d.Ready {
		return fmt.Errorf("failed to start: function not ready")
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeploymentTimeout)
	cmd := exec.CommandContext(ctx, "docker", "run", d.Image)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	cancel()
	if err != nil {
		return fmt.Errorf("failed to run: %v", err)
	}
	return nil
}

func (d *Deployment) Build() error {
	// Allocate tmp folder
	dir, err := ioutil.TempDir("", d.ID)
	if err != nil {
		return fmt.Errorf("failed to allocate dir: %v", err)
	}
	// Copy source code there
	err = ioutil.WriteFile(filepath.Join(dir, "main."+d.Language), d.Source, 0644)
	if err != nil {
		return fmt.Errorf("failed to write source: %v", err)
	}
	// Copy dockerfile there
	err = ioutil.WriteFile(filepath.Join(dir, "Dockerfile"), dockerfiles[d.Language], 0644)
	if err != nil {
		return fmt.Errorf("failed to write dockerfile: %v", err)
	}
	// Run docker build
	imageName := fmt.Sprintf("foog-%16.16s", d.ID)
	cmd := exec.Command("docker", "build", "-t", imageName, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		fmt.Printf("====== STARTING BUILD %s ======\n", imageName)
		err = cmd.Run()
		if err != nil {
			fmt.Println("====== BUILD FAILED ======")
		} else {
			fmt.Println("====== BUILD SUCCESSFUL ======")
			d.Ready = true
		}
	}()
	return nil
}

type Server struct {
	http.ServeMux
	Registry *Registry
}

func NewServer() *Server {
	srv := &Server{
		ServeMux: *http.NewServeMux(),
		Registry: NewRegistry(),
	}
	srv.HandleFunc(TriggerPrefix, srv.trigger)
	srv.HandleFunc("/deploy", srv.deploy)
	srv.HandleFunc("/list", srv.list)
	return srv
}

func (srv *Server) setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (srv *Server) deploy(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	encoder := json.NewEncoder(w)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encoder.Encode(struct {
			Error   bool
			Message string
		}{true, "can not read body"})
		return
	}
	deployment, err := srv.Registry.Deploy(r.URL.Query().Get("lang"), data)
	if err != nil {
		encoder.Encode(struct {
			Error   bool
			Message string
		}{true, err.Error()})
		return
	}
	encoder.Encode(deployment)
}

func (srv *Server) trigger(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	encoder := json.NewEncoder(w)
	id := filepath.Base(r.URL.Path)
	err := srv.Registry.Run(id, r.Body, w, w)
	if err != nil {
		encoder.Encode(struct {
			Error   bool
			Message string
		}{true, err.Error()})
		return
	}
}

func (srv *Server) list(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	encoder := json.NewEncoder(w)
	deployments := []*Deployment{}
	for _, d := range srv.Registry.Items {
		deployments = append(deployments, d)
	}
	encoder.Encode(deployments)
}

func main() {
	if err := http.ListenAndServe(":8080", NewServer()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
