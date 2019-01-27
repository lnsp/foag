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
	"strings"
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
RUN swiftc -o /app/main /app/*.swift

CMD ["/app/main"]
`

const jsDockerfile = `
FROM node:8-alpine

WORKDIR /app
COPY . .

CMD ["node", "/app/main.js"]
`

const cDockerfile = `
FROM gcc:4.9

WORKDIR /app
COPY . .
RUN gcc -o /app/main /app/*.c

CMD ["/app/main"]
`

const AliasPrefix = "@"
const BindPrefix = "/bind/"
const TriggerPrefix = "/trigger/"
const DeploymentTimeout = 30 * time.Second
const DeploymentImage = "foag-%16.16s"

var (
	dockerfiles = map[string][]byte{
		"go":    []byte(goDockerfile),
		"swift": []byte(swiftDockerfile),
		"js":    []byte(jsDockerfile),
		"c":     []byte(cDockerfile),
	}
)

type Registry struct {
	Items   map[string]*Deployment
	Aliases map[string]*Alias
}

func NewRegistry() *Registry {
	return &Registry{
		Items:   make(map[string]*Deployment),
		Aliases: make(map[string]*Alias),
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
	deployment := registry.Resolve(id)
	if deployment == nil {
		return fmt.Errorf("deployment not found")
	}
	err := deployment.Run(stdin, stdout, stderr)
	if err != nil {
		return err
	}
	return nil
}

func (registry *Registry) Resolve(identifier string) *Deployment {
	if strings.HasPrefix(identifier, AliasPrefix) {
		id, ok := registry.Aliases[strings.TrimPrefix(identifier, AliasPrefix)]
		if !ok {
			return nil
		}
		identifier = id.For
	}
	return registry.Items[identifier]
}

func (registry *Registry) Find(identifier string) *Deployment {
	var (
		item         *Deployment
		longestMatch int
	)
	for _, d := range registry.Items {
		if strings.HasPrefix(d.ID, identifier) && len(identifier) > longestMatch {
			item = d
			longestMatch = len(identifier)
		}
	}
	return item
}

func (registry *Registry) Bind(alias string, deployment *Deployment) *Alias {
	item := &Alias{
		Name: alias,
		For:  deployment.ID,
		URL:  TriggerPrefix + AliasPrefix + alias,
		Date: time.Now(),
	}
	registry.Aliases[alias] = item
	return item
}

type Deployment struct {
	ID        string
	Source    []byte `json:"-"`
	Image     string
	Date      time.Time
	Language  string
	URL       string
	Ready     bool
	BuildLogs string `json:"-"`
}

type Alias struct {
	Name string
	For  string
	URL  string
	Date time.Time
}

func NewDeployment(source []byte, lang string) *Deployment {
	hasher := sha256.New()
	hasher.Write(source)
	hasher.Write([]byte(lang))
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
	imageName := fmt.Sprintf("foag-%16.16s", d.ID)
	cmd := exec.Command("docker", "build", "-t", imageName, dir)
	go func() {
		tmpfile, err := ioutil.TempFile("", "foagd-*.log")
		if err != nil {
			log.Printf("Failed to hookup build logs: %v\n", err)
			return
		}
		d.BuildLogs = tmpfile.Name()
		cmd.Stdout = tmpfile
		cmd.Stderr = tmpfile
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to build %s (%s): %v\n", d.Image, d.Language, err)
		} else {
			log.Printf("Built image %s (%s)\n", d.Image, d.Language)
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
	srv.HandleFunc(BindPrefix, srv.bind)
	srv.HandleFunc("/describe/", srv.describe)
	srv.HandleFunc("/logs/", srv.logs)
	srv.HandleFunc("/deploy", srv.deploy)
	srv.HandleFunc("/listAlias", srv.listAlias)
	srv.HandleFunc("/list", srv.list)

	return srv
}

func (srv *Server) setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (srv *Server) error(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(struct {
		Error   bool
		Message string
	}{
		Error:   true,
		Message: msg,
	})
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "error (%d): %s\n", status, msg)
		return
	}
	log.Printf("Error while serving request: %s\n", msg)
	w.Header().Add("Content-Type", "application/json")
}

func (srv *Server) deploy(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	encoder := json.NewEncoder(w)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.error(w, http.StatusBadRequest, "Failed to read body")
		return
	}
	deployment, err := srv.Registry.Deploy(r.URL.Query().Get("lang"), data)
	if err != nil {
		srv.error(w, http.StatusInternalServerError, err.Error())
		return
	}
	encoder.Encode(deployment)
}

func (srv *Server) logs(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	id := filepath.Base(r.URL.Path)
	deployment := srv.Registry.Resolve(id)
	if deployment == nil {
		deployment = srv.Registry.Find(id)
	}
	if deployment == nil {
		srv.error(w, http.StatusNotFound, "Deployment not found")
		return
	}
	logfile, err := os.Open(deployment.BuildLogs)
	if err != nil {
		srv.error(w, http.StatusNotFound, "BuildLogs not found")
		return
	}
	defer logfile.Close()
	io.Copy(w, logfile)
}

func (srv *Server) describe(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	id := filepath.Base(r.URL.Path)
	deployment := srv.Registry.Resolve(id)
	if deployment == nil {
		deployment = srv.Registry.Find(id)
	}
	if deployment == nil {
		srv.error(w, http.StatusNotFound, "Deployment not found")
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(deployment)
}

func (srv *Server) bind(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	id := filepath.Base(r.URL.Path)
	alias := r.URL.Query().Get("to")
	deployment := srv.Registry.Resolve(id)
	if deployment == nil {
		srv.error(w, http.StatusNotFound, "Deployment does not exist")
		return
	}
	item := srv.Registry.Bind(alias, deployment)
	encoder := json.NewEncoder(w)
	encoder.Encode(item)
}

func (srv *Server) trigger(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	id := filepath.Base(r.URL.Path)
	err := srv.Registry.Run(id, r.Body, w, w)
	if err != nil {
		srv.error(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (srv *Server) listAlias(w http.ResponseWriter, r *http.Request) {
	srv.setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	encoder := json.NewEncoder(w)
	aliases := []*Alias{}
	for _, a := range srv.Registry.Aliases {
		aliases = append(aliases, a)
	}
	encoder.Encode(aliases)
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
