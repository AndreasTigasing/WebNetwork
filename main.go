package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Title         string `json:"title"`
    Port          int    `json:"port"`
    SyslogEnabled bool   `json:"syslog_enabled"`
    SyslogIP      string `json:"syslog_ip"`
    SyslogPort    int    `json:"syslog_port"`
}

type MacVendor struct {
	MacPrefix   string `json:"Mac Prefix"`
	VendorName  string `json:"Vendor Name"`
	IsPrivate   bool   `json:"Private"`
	BlockType   string `json:"Block Type"`
	LastUpdated string `json:"Last Update"`
}

var tpl *template.Template
var config Config

func main() {
    // Load configuration from GlobalConfig.json
	loadConfig("GlobalConfig.json")
	loadMacVendors("MacVendors.json")

    // Parse templates
    tpl = template.Must(template.ParseGlob("pages/*.html"))

    // HTTP handlers
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
    http.HandleFunc("/", logHandler(homeHandler))
    http.HandleFunc("/devices", logHandler(devicesHandler))
    http.HandleFunc("/sites", logHandler(sitesHandler))
    http.HandleFunc("/mac-lookup", logHandler(mac_lookupHandler))

    // Start the server
	addr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Server started at localhost%s\n", addr)
	http.ListenAndServe(addr, nil)
}

// logHandler is a middleware function that logs the remote address and accessed URL
func logHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Remote Address: %s, Accessed URL: %s\n", r.RemoteAddr, r.URL.Path)
        next.ServeHTTP(w, r)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "index.html", config)
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "devices.html", config)
}

func sitesHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "sites.html", config)
}

func mac_lookupHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "mac-lookup.html", config)
}

func loadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}
}

func loadMacVendors(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening MacVendors.json file: %v", err)
	}
	defer file.Close()

	var vendors []MacVendor
	if err := json.NewDecoder(file).Decode(&vendors); err != nil {
		log.Fatalf("Error decoding MacVendors.json: %v", err)
	}
}