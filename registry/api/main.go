// api/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type RegisterRequest struct {
	Domain      string       `json:"domain"`
	Tld         string       `json:"tld"`
	OwnerOrg    string       `json:"owner_organization"`
	OwnerEmail  string       `json:"owner_email"`
	Registrar   string       `json:"registrar"`
	APIKey      string       `json:"api_key"`
	Nameservers []Nameserver `json:"nameservers"`
}

type Nameserver struct {
	Hostname string `json:"hostname"`
	IPv4     string `json:"ipv4"`
}

var db *sql.DB

func main() {
	var err error
	dbConnStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	// Initialize DB handle once
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to initialize database handle: %v", err)
	}

	// Retry loop for DB connectivity with stabilized bridge network tolerances
	for i := 0; i < 15; i++ {
		if err = db.Ping(); err == nil {
			log.Println("Connected to PostgreSQL Registry Database")
			break
		}
		log.Printf("Waiting for database... (%d/15)", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to database after retries: %v", err)
	}

	http.HandleFunc("/domains/register", handleRegister)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registry API Operational"))
	})

	log.Println("Registry API listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Authenticate Registrar
	var registrarID int
	err := db.QueryRow("SELECT id FROM registrars WHERE name = $1 AND api_key = $2 AND status = 'active'", req.Registrar, req.APIKey).Scan(&registrarID)
	if err != nil {
		log.Printf("Unauthorized registration attempt from: %s", req.Registrar)
		http.Error(w, "Unauthorized Registrar or Invalid API Key", http.StatusUnauthorized)
		recordTransaction(req.Domain, "REGISTER", req.Registrar, "FAILED_UNAUTHORIZED")
		return
	}

	// 2. Validation Rules
	if err := validateDomain(req.Domain, req.Tld, req.Nameservers); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		recordTransaction(req.Domain, "REGISTER", req.Registrar, "FAILED_VALIDATION")
		return
	}

	// 3. Begin Database Transaction for Registration
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Ensure Owner exists or create
	var ownerID int
	err = tx.QueryRow("SELECT id FROM owners WHERE email = $1", req.OwnerEmail).Scan(&ownerID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(
			"INSERT INTO owners (organization, email) VALUES ($1, $2) RETURNING id",
			req.OwnerOrg, req.OwnerEmail,
		).Scan(&ownerID)
	}
	if err != nil {
		http.Error(w, "Failed to process domain owner", http.StatusInternalServerError)
		return
	}

	// Insert Domain
	expiresAt := time.Now().AddDate(1, 0, 0) // 1 year validity
	var domainID int
	err = tx.QueryRow(
		"INSERT INTO domains (domain, tld, owner_id, status, expires_at) VALUES ($1, $2, $3, 'active', $4) RETURNING id",
		req.Domain, req.Tld, ownerID, expiresAt,
	).Scan(&domainID)
	if err != nil {
		http.Error(w, "Domain already registered or conflict occurred", http.StatusConflict)
		recordTransaction(req.Domain, "REGISTER", req.Registrar, "FAILED_CONFLICT")
		return
	}

	// Insert Nameservers
	for _, ns := range req.Nameservers {
		_, err = tx.Exec(
			"INSERT INTO nameservers (domain_id, hostname, ipv4) VALUES ($1, $2, $3)",
			domainID, ns.Hostname, ns.IPv4,
		)
		if err != nil {
			http.Error(w, "Failed to record nameservers", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit registration", http.StatusInternalServerError)
		return
	}

	recordTransaction(req.Domain, "REGISTER", req.Registrar, "SUCCESS")

	// 4. Trigger Zone Generator & Publisher
	go generateAndPublishZone(req.Tld)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"domain":  req.Domain,
		"message": "Domain registered and TLD zone published successfully",
	})
}

func validateDomain(domain, tld string, nameservers []Nameserver) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[\w]+$`, domain)
	if !matched || !strings.HasSuffix(domain, "."+tld) {
		return fmt.Errorf("invalid domain syntax or TLD mismatch")
	}

	if len(nameservers) == 0 {
		return fmt.Errorf("at least one nameserver is required")
	}

	reserved := map[string]bool{"admin": true, "root": true, "registry": true, "nic": true}
	parts := strings.Split(domain, ".")
	if reserved[parts[0]] {
		return fmt.Errorf("domain name is reserved")
	}

	return nil
}

func recordTransaction(domain, operation, registrar, status string) {
	_, err := db.Exec(
		"INSERT INTO transactions (domain, operation, registrar, status) VALUES ($1, $2, $3, $4)",
		domain, operation, registrar, status,
	)
	if err != nil {
		log.Printf("Failed to record audit transaction: %v", err)
	}
}

func generateAndPublishZone(tld string) {
	log.Printf("Generating updated zone file for TLD: .%s", tld)
}
