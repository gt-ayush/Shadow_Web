package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Client request structure
type ClientRegistrationRequest struct {
	Domain      string `json:"domain"`
	Tld         string `json:"tld"`
	OwnerOrg    string `json:"owner_org"`
	OwnerEmail  string `json:"owner_email"`
	Nameservers []struct {
		Hostname string `json:"hostname"`
		IPv4     string `json:"ipv4"`
	} `json:"nameservers"`
}

// Registry payload structure
type RegistryPayload struct {
	Domain      string `json:"domain"`
	Tld         string `json:"tld"`
	OwnerOrg    string `json:"owner_organization"`
	OwnerEmail  string `json:"owner_email"`
	Registrar   string `json:"registrar"`
	APIKey      string `json:"api_key"`
	Nameservers []struct {
		Hostname string `json:"hostname"`
		IPv4     string `json:"ipv4"`
	} `json:"nameservers"`
}

func main() {
	http.HandleFunc("/api/v1/register", handleClientRegistration)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Registrar API listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleClientRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clientReq ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&clientReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Basic Registrar-Side Normalization
	clientReq.Domain = strings.ToLower(strings.TrimSpace(clientReq.Domain))
	clientReq.Tld = strings.ToLower(strings.TrimSpace(clientReq.Tld))

	// 2. Construct Payload for the Registry
	// In a real environment, API keys are loaded via secure secrets/env vars.
	payload := RegistryPayload{
		Domain:      clientReq.Domain,
		Tld:         clientReq.Tld,
		OwnerOrg:    clientReq.OwnerOrg,
		OwnerEmail:  clientReq.OwnerEmail,
		Registrar:   "ShadowRegistrarPrimary",
		APIKey:      "sec_key_alpha_99887766", // From your 002_seed_data.sql
		Nameservers: clientReq.Nameservers,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 3. Send Request to Registry API (10.89.0.12 on dns-net)
	registryURL := "http://10.89.0.12:8080/domains/register"
	resp, err := http.Post(registryURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to contact registry: %v", err)
		http.Error(w, "Registry unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// 4. Relay Registry Response to Client
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
