/*
Auth: Desmond McDermitt (Enhanced version)
Desc:
	Prometheos - Alert Manager Dashboard
	- Displays alerts from multiple Prometheus Alert Managers
	- Search with include/exclude functionality (use '!' prefix to exclude)
	- Alert silencing capability
	- Auto-refresh every 5 minutes

Endpoints:
	- /    : Legacy interface (original design)
	- /v2  : Modern interface (new dashboard design)

Improvements in this version:
	- Dual interface support (legacy + modern)
	- Added mutex for thread-safe access to shared data
	- Replaced deprecated ioutil with io
	- Added graceful shutdown
	- Improved error handling (non-fatal for individual alert manager failures)
	- Added configuration options
	- Better logging
*/
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Configuration
var config = struct {
	Port           string
	UpdateInterval time.Duration
	RequestTimeout time.Duration
	AlertManagers  []string
}{
	Port:           ":8001",
	UpdateInterval: 5 * time.Minute,
	RequestTimeout: 10 * time.Second,
	AlertManagers: []string{
		"c01.b3.alertmanager.pro.mon.liquidweb.com",
		"c02.b3.alertmanager.pro.mon.liquidweb.com",
		"n01.b2.alertmanager.pro.mon.liquidweb.com",
		"n01.b3.alertmanager.pro.mon.liquidweb.com",
		"n01.b4.alertmanager.pro.mon.liquidweb.com",
		"n01.b5.alertmanager.pro.mon.liquidweb.com",
		"n02.b2.alertmanager.pro.mon.liquidweb.com",
		"n02.b3.alertmanager.pro.mon.liquidweb.com",
		"n02.b4.alertmanager.pro.mon.liquidweb.com",
		"n02.b5.alertmanager.pro.mon.liquidweb.com",
		"n03.b3.alertmanager.pro.mon.liquidweb.com",
		"n04.b3.alertmanager.pro.mon.liquidweb.com",
		"n05.b3.alertmanager.pro.mon.liquidweb.com",
		"n06.b3.alertmanager.pro.mon.liquidweb.com",
		"n07.b3.alertmanager.pro.mon.liquidweb.com",
		"n08.b3.alertmanager.pro.mon.liquidweb.com",
	},
}

// Template sets for each version
var (
	tplV1 *template.Template // Legacy interface
	tplV2 *template.Template // Modern interface
)

// ServerInfo represents a single server's alert information
type serverInfo struct {
	Account  string
	Hostname string
	Location string
	Prom     string
	Time     string
	Service  []string
	UID      string
	Silenced bool
}

// AlertStore provides thread-safe access to alert data
type AlertStore struct {
	mu       sync.RWMutex
	servers  []serverInfo
	lastSync time.Time
}

func NewAlertStore() *AlertStore {
	return &AlertStore{
		servers: make([]serverInfo, 0),
	}
}

func (s *AlertStore) Get() []serverInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent race conditions
	result := make([]serverInfo, len(s.servers))
	copy(result, s.servers)
	return result
}

func (s *AlertStore) Set(servers []serverInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers = servers
	s.lastSync = time.Now()
}

func (s *AlertStore) LastSync() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSync
}

func sortServerList(list map[string]serverInfo) []serverInfo {
	type item struct {
		Key  string
		Time string
	}

	var keys []item
	for k := range list {
		temp := item{
			Key:  k,
			Time: list[k].Time,
		}
		keys = append(keys, temp)
	}

	sort.Slice(keys, func(i int, j int) bool {
		if keys[i].Time == keys[j].Time {
			return keys[i].Key < keys[j].Key
		}
		return keys[i].Time < keys[j].Time
	})

	var sortedServers []serverInfo
	for k := range keys {
		sortedServers = append(sortedServers, list[keys[k].Key])
	}

	return sortedServers
}

// indexHandler renders the dashboard page using the specified template set
func indexHandler(w http.ResponseWriter, r *http.Request, txt []serverInfo, tpl *template.Template) {
	var ttlAlerts int
	var htmlOut bytes.Buffer
	tplWriter := io.MultiWriter(&htmlOut)

	// Count total alerts
	for _, v := range txt {
		ttlAlerts += len(v.Service)
	}

	// Get last update time safely
	updateTime := "N/A"
	if len(txt) > 0 {
		updateTime = txt[len(txt)-1].Time
	}

	// Template data for header
	alerts := struct {
		ServerCount int
		AlertCount  int
		UpdateTime  string
	}{
		ServerCount: len(txt) - 1, // Subtract 1 for the timestamp entry
		AlertCount:  ttlAlerts,
		UpdateTime:  updateTime,
	}

	if alerts.ServerCount < 0 {
		alerts.ServerCount = 0
	}

	// Execute templates
	if err := tpl.ExecuteTemplate(tplWriter, "head", alerts); err != nil {
		log.Printf("Error executing head template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(tplWriter, "body", txt); err != nil {
		log.Printf("Error executing body template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(tplWriter, "foot", nil); err != nil {
		log.Printf("Error executing foot template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlOut.Bytes())
}

// jsonToStruct converts JSON from alert managers into serverInfo structs
func jsonToStruct(jsonData [][]map[string]interface{}) map[string]serverInfo {
	serverlist := make(map[string]serverInfo)

	for _, v := range jsonData {
		for _, vv := range v {
			labels, ok := vv["labels"].(map[string]interface{})
			if !ok {
				continue
			}

			uid := fmt.Sprintf("%v", labels["uniq_id"])
			timeStr := fmt.Sprintf("%v", vv["startsAt"])
			timeAr := strings.Split(timeStr, "T")

			var timeStamp string
			if len(timeAr) >= 2 {
				timeStamp = timeAr[0] + " " + strings.Split(timeAr[1], ".")[0]
			} else {
				timeStamp = timeStr
			}

			if _, ok := serverlist[uid]; !ok {
				var tempSrvStruct serverInfo

				if labels["group"] == "guardian" {
					tempSrvStruct = serverInfo{
						Hostname: fmt.Sprintf("%v", labels["hostname"]),
						Account:  fmt.Sprintf("%v", labels["policy_description"]),
						Location: fmt.Sprintf("%v", labels["disksafe_description"]),
						Prom:     fmt.Sprintf("%v", labels["prom_serv"]),
						Service:  []string{fmt.Sprintf("%v", labels["alertname"])},
						Time:     timeStamp,
						UID:      uid,
					}
				} else {
					tempSrvStruct = serverInfo{
						Hostname: fmt.Sprintf("%v", labels["Hostname"]),
						Account:  fmt.Sprintf("%v", labels["Account"]),
						Location: fmt.Sprintf("%v", labels["Location"]),
						Prom:     fmt.Sprintf("%v", labels["prom_serv"]),
						Service:  []string{fmt.Sprintf("%v", labels["alertname"])},
						Time:     timeStamp,
						UID:      uid,
					}
				}

				serverlist[tempSrvStruct.UID] = tempSrvStruct
			} else {
				// Check for duplicate services
				found := false
				alertName := fmt.Sprintf("%v", labels["alertname"])
				for _, k := range serverlist[uid].Service {
					if k == alertName {
						found = true
						break
					}
				}

				if !found {
					temp := serverlist[uid].Service
					temp = append(temp, alertName)

					if key, ok := serverlist[uid]; ok {
						key.Service = temp
						serverlist[uid] = key
					}
				}

				// Update to earliest time
				if key, ok := serverlist[uid]; ok {
					if key.Time > timeStamp {
						key.Time = timeStamp
						serverlist[uid] = key
					}
				}
			}
		}
	}

	return serverlist
}

// getJSON fetches alerts from all configured alert managers
func getJSON() [][]map[string]interface{} {
	var result [][]map[string]interface{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	client := &http.Client{
		Timeout: config.RequestTimeout,
	}

	for _, host := range config.AlertManagers {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			url := fmt.Sprintf("http://%s:9093/api/v2/alerts?silenced=false&inhibited=false&group=uniq_id", host)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("Error creating request for %s: %v", host, err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error fetching from %s: %v", host, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("Non-OK response from %s: %d", host, resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response from %s: %v", host, err)
				return
			}

			var results []map[string]interface{}
			if err := json.Unmarshal(body, &results); err != nil {
				log.Printf("Error parsing JSON from %s: %v", host, err)
				return
			}

			mu.Lock()
			result = append(result, results)
			mu.Unlock()
		}(host)
	}

	wg.Wait()
	return result
}

// Initialize templates for both versions
func init() {
	var err error

	// Load legacy templates (v1)
	tplV1, err = template.ParseGlob("templates/v1/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to load v1 templates: %v", err)
	}
	log.Println("Loaded v1 (legacy) templates")

	// Load modern templates (v2)
	tplV2, err = template.ParseGlob("templates/v2/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to load v2 templates: %v", err)
	}
	log.Println("Loaded v2 (modern) templates")
}

// searchService filters the server list based on search criteria
func searchService(crit string, list []serverInfo) []serverInfo {
	if len(crit) == 0 {
		return list
	}

	crit = strings.TrimSpace(crit)
	include := true

	if len(crit) > 0 && crit[0] == '!' {
		include = false
		crit = crit[1:]
	}

	if len(crit) == 0 {
		return list
	}

	var tempMap []serverInfo
	critLower := strings.ToLower(crit)

	for k := range list {
		// Skip the timestamp entry during search
		if list[k].UID == "TS" {
			continue
		}

		found := strings.Contains(strings.ToLower(list[k].UID), critLower) ||
			strings.Contains(strings.ToLower(list[k].Account), critLower) ||
			strings.Contains(strings.ToLower(list[k].Hostname), critLower) ||
			strings.Contains(strings.ToLower(list[k].Location), critLower) ||
			strings.Contains(strings.ToLower(list[k].Prom), critLower) ||
			strings.Contains(strings.ToLower(list[k].Time), critLower)

		if !found {
			for s := range list[k].Service {
				if strings.Contains(strings.ToLower(list[k].Service[s]), critLower) {
					found = true
					break
				}
			}
		}

		if !include {
			if !found {
				tempMap = append(tempMap, list[k])
			}
		} else if found {
			tempMap = append(tempMap, list[k])
		}
	}

	// Always append the timestamp entry at the end
	if len(list) > 0 && list[len(list)-1].UID == "TS" {
		tempMap = append(tempMap, list[len(list)-1])
	}

	return tempMap
}

// checkUpdate refreshes the alert data from all sources
func checkUpdate(store *AlertStore) {
	log.Println("Fetching alerts from alert managers...")
	start := time.Now()

	jsonData := getJSON()
	servMap := jsonToStruct(jsonData)
	servList := sortServerList(servMap)

	// Add timestamp entry
	var timeStamp serverInfo
	timeString := fmt.Sprintf("%02d:%02d:%02d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	timeStamp.Time = timeString
	timeStamp.UID = "TS"

	servList = append(servList, timeStamp)
	store.Set(servList)

	log.Printf("Alert update complete: %d servers, took %v", len(servList)-1, time.Since(start))
}

// addSilencedAlert adds a hostname to the silenced list
func addSilencedAlert(alert string) error {
	file, err := os.OpenFile("silencedAlerts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open silencedAlerts file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(alert + "\n"); err != nil {
		return fmt.Errorf("failed to write to silencedAlerts file: %w", err)
	}

	log.Printf("Silenced alert for: %s", alert)
	return nil
}

// getSilencedHosts returns the list of silenced hostnames
func getSilencedHosts() ([]string, error) {
	hosts := []string{}

	file, err := os.Open("silencedAlerts")
	if os.IsNotExist(err) {
		return hosts, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open silencedAlerts file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := strings.TrimSpace(scanner.Text())
		if host != "" {
			hosts = append(hosts, host)
		}
	}

	return hosts, scanner.Err()
}

// setSilencedAlerts marks servers as silenced based on the silenced hosts list
func setSilencedAlerts(servers []serverInfo) []serverInfo {
	hosts, err := getSilencedHosts()
	if err != nil {
		log.Printf("Error getting silenced hosts: %v", err)
		return servers
	}

	hostSet := make(map[string]bool)
	for _, h := range hosts {
		hostSet[h] = true
	}

	for i := range servers {
		if hostSet[servers[i].Hostname] {
			servers[i].Silenced = true
		}
	}

	return servers
}

// removeSilenced removes a hostname from the silenced list
func removeSilenced(servers []serverInfo, silenced string) []serverInfo {
	hosts, err := getSilencedHosts()
	if err != nil {
		log.Printf("Error getting silenced hosts: %v", err)
		return servers
	}

	// Filter out the silenced hostname
	newHosts := []string{}
	for _, h := range hosts {
		if h != silenced {
			newHosts = append(newHosts, h)
		}
	}

	// Rewrite the file
	file, err := os.OpenFile("silencedAlerts", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error opening silencedAlerts for writing: %v", err)
		return servers
	}
	defer file.Close()

	for _, v := range newHosts {
		file.WriteString(v + "\n")
	}

	// Update the server list
	for i := range servers {
		if servers[i].Hostname == silenced {
			servers[i].Silenced = false
		}
	}

	log.Printf("Unsilenced alert for: %s", silenced)
	return servers
}

// createHandler creates an HTTP handler for a specific template version
func createHandler(store *AlertStore, tpl *template.Template, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a copy of the server list
		tmp := store.Get()

		// Parse form data
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
		}

		// Handle search
		searchCriteria := ""
		if val, ok := r.Form["search-box"]; ok && len(val) > 0 {
			searchCriteria = val[0]
		}

		// Handle silence request (supports bulk operations)
		if vals, ok := r.Form["silence-alert"]; ok {
			for _, val := range vals {
				if val != "" {
					if err := addSilencedAlert(val); err != nil {
						log.Printf("Error silencing alert %s: %v", val, err)
					}
				}
			}
		}

		// Handle unsilence request (supports bulk operations)
		if vals, ok := r.Form["remove-silence"]; ok {
			for _, val := range vals {
				if val != "" {
					tmp = removeSilenced(tmp, val)
				}
			}
		}

		// Apply search filter
		tmp = searchService(searchCriteria, tmp)

		// Apply silenced status
		tmp = setSilencedAlerts(tmp)

		log.Printf("[%s] Serving %s request from %s", version, r.Method, r.RemoteAddr)
		indexHandler(w, r, tmp, tpl)
	}
}

func main() {
	store := NewAlertStore()

	// Initial data load
	log.Println("Starting Prometheos Alert Dashboard...")
	log.Println("Endpoints:")
	log.Println("  /    - Legacy interface")
	log.Println("  /v2  - Modern interface")
	checkUpdate(store)

	// Set up periodic updates
	updateTimer := time.NewTicker(config.UpdateInterval)
	defer updateTimer.Stop()

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Background update goroutine
	go func() {
		for {
			select {
			case <-quit:
				log.Println("Shutting down update routine...")
				return
			case <-updateTimer.C:
				checkUpdate(store)
			}
		}
	}()

	// HTTP handlers for both versions
	// Legacy interface at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only handle exact root path for legacy
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		createHandler(store, tplV1, "v1")(w, r)
	})

	// Modern interface at /v2
	http.HandleFunc("/v2", createHandler(store, tplV2, "v2"))
	http.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect /v2/ to /v2
		if r.URL.Path == "/v2/" {
			http.Redirect(w, r, "/v2", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})

	// Create server with graceful shutdown
	server := &http.Server{
		Addr:         config.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
