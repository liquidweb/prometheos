# Prometheos

A high-performance alert aggregation dashboard for Prometheus Alert Managers. Prometheos consolidates alerts from multiple Alert Manager instances into a unified interface, providing real-time monitoring, incident management, and alert silencing capabilities.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Interfaces](#interfaces)
  - [Legacy Interface (v1)](#legacy-interface-v1)
  - [Modern Interface (v2)](#modern-interface-v2)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Keyboard Shortcuts](#keyboard-shortcuts)
- [Browser Storage](#browser-storage)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Development](#development)
- [Credits](#credits)

## Features

- **Multi-Source Aggregation** — Fetch alerts from multiple Prometheus Alert Managers in parallel
- **Dual Interface** — Legacy table view and modern feature-rich dashboard
- **Real-Time Updates** — Configurable auto-refresh intervals
- **Alert Silencing** — Mute alerts by hostname with persistent storage
- **Incident Management** — Auto-detect and group correlated alerts
- **Advanced Filtering** — Filter by datacenter, section, row, rack, time range, and service
- **Export Options** — Download alerts as CSV or JSON
- **Shareable URLs** — Encode dashboard state in URLs for bookmarking and sharing
- **Thread-Safe** — Mutex-protected data store prevents race conditions
- **Graceful Shutdown** — Handles SIGINT/SIGTERM signals cleanly

## Prerequisites

- Go 1.18 or later
- Network access to Prometheus Alert Manager instances

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/liquidweb/prometheos.git
cd prometheos

# Build the binary
go build -o prometheos main.go
```

### Verify Installation

```bash
./prometheos --help
```

## Quick Start

1. **Configure Alert Managers** — Edit the `AlertManagers` slice in `main.go` (see [Configuration](#configuration))

2. **Build and Run**
   ```bash
   go build -o prometheos main.go
   ./prometheos
   ```

3. **Access the Dashboard**
   - Legacy Interface: http://localhost:8001/
   - Modern Interface: http://localhost:8001/v2

The server logs startup information and begins fetching alerts immediately:

```
Starting Prometheos Alert Dashboard...
Endpoints:
  /    - Legacy interface
  /v2  - Modern interface
Fetching alerts from alert managers...
Alert update complete: 42 servers, took 1.234s
Server listening on :8001
```

## Interfaces

Prometheos provides two interfaces that share the same underlying alert data and silencing state.

| Endpoint | Interface | Description |
|----------|-----------|-------------|
| `/` | Legacy (v1) | Simple table layout with basic search |
| `/v2` | Modern (v2) | Feature-rich dashboard with analytics and incident management |

### Legacy Interface (v1)

The original interface provides:

- Tabular alert display sorted by start time
- Basic text search across all fields
- Exclude search with `!` prefix
- Per-alert silence/unsilence actions

### Modern Interface (v2)

The modern interface includes all legacy features plus advanced capabilities organized into the following sections.

#### Status Cards

Real-time metrics displayed at the top of the dashboard:

| Card | Description |
|------|-------------|
| **Active Alerts** | Total alerts requiring attention |
| **Affected Servers** | Unique servers currently alerting |
| **Incidents** | Auto-detected groups of correlated alerts |
| **Silenced** | Currently muted alert count |
| **Alert Managers** | Number of sources being monitored (click for details) |

#### Theme and Display

| Feature | Toggle | Description |
|---------|--------|-------------|
| Light/Dark Mode | `T` or theme button | Switch between color themes |
| View Density | `D` | Cycle through Compact, Comfortable, Spacious |

Preferences persist in localStorage across sessions.

#### Search and Filtering

**Basic Search**

- Type keywords to filter alerts across all columns
- Prefix with `!` to exclude matches (e.g., `!guardian` hides guardian alerts)
- Press `/` to focus the search box

**Advanced Filters** (Press `F` to toggle)

| Filter | Description |
|--------|-------------|
| Data Center | Filter by datacenter identifier |
| Section | Filter by datacenter section |
| Row | Filter by row number |
| Rack | Filter by rack/cabinet number |
| Prom Node | Filter by Prometheus node source |
| Time Range | 1H, 6H, 24H, or 7D |
| Service | Filter by service/alert name |

**Filter Presets**

- Save filter combinations as named presets
- Quick-apply saved presets from dropdown

#### Table Columns

Default visible columns:

- UID, Account, Hostname
- DC, Section, Row, Rack, Location
- Prom Node, Start Time, Duration
- Services, Actions

**Column Configuration** (Press `C`)

- Drag and drop to reorder columns
- Toggle column visibility
- Presets: Default, Compact, Detailed, Custom
- API Explorer to discover additional AlertManager fields

#### Analytics Sections

Collapsible sections providing insights into alert patterns:

| Section | Contents |
|---------|----------|
| Distribution | Alerts by location, account, and top services |
| Performance Metrics | Average alert age, oldest alert, alert velocity, recurring systems |
| Alert Velocity Chart | Visual timeline with configurable time ranges |

#### Incident Management

Automatically detects and groups alerts that may represent correlated infrastructure issues.

**Detection Criteria**

| Parameter | Range | Description |
|-----------|-------|-------------|
| Time Window | 2min - 1hour | Alerts within this window are considered related |
| Alert Threshold | 2-10+ | Minimum alerts to form an incident |
| Host Threshold | 1-5+ | Minimum unique hosts to form an incident |

**Managing Incidents**

- Click incident name to rename
- Set status: Investigating, Identified, Monitoring, Resolved
- Add timestamped notes
- Rescan to re-detect (preserves saved incidents)

**Location Code Format**

Location codes follow the format `B3S2R9K8U27`:

| Code | Meaning |
|------|---------|
| B3 | Datacenter 3 |
| S2 | Section 2 |
| R9 | Row 9 |
| K8 | Rack/Cabinet 8 |
| U27 | Rack Unit 27 |

#### Alert Details Drawer

Click any alert row to open a slide-out panel with:

- Full details: UID, account, hostname, location, duration
- All services displayed as tags
- Quick actions: Silence/unsilence, copy as text or JSON
- Navigate between alerts with `←` `→` arrow keys

#### Bulk Operations

Select multiple alerts using checkboxes:

| Action | Description |
|--------|-------------|
| Silence | Mute all selected alerts |
| Unsilence | Unmute all selected alerts |
| Create Incident | Group selected alerts into an incident |
| Copy | Copy selected alerts to clipboard |

#### Export Options

| Format | Description |
|--------|-------------|
| CSV | Export current filtered view as CSV |
| JSON | Export current filtered view as JSON |
| With Annotations | Include notes in export |
| Incidents | Export incident data |

#### Shareable URLs

Dashboard state is encoded in the URL:

- Share button copies URL with current tab, search, filters, and time range
- URLs are bookmarkable
- Browser back/forward navigation works

## Configuration

Edit the configuration struct in `main.go`:

```go
var config = struct {
    Port           string
    UpdateInterval time.Duration
    RequestTimeout time.Duration
    AlertManagers  []string
}{
    Port:           ":8001",
    UpdateInterval: 5 * time.Minute,
    RequestTimeout: 10 * time.Second,
    AlertManagers:  []string{
        "alertmanager1.example.com",
        "alertmanager2.example.com",
        // Add your Alert Manager hosts here
    },
}
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `Port` | `:8001` | HTTP server listen address |
| `UpdateInterval` | `5 * time.Minute` | How often to fetch alerts from sources |
| `RequestTimeout` | `10 * time.Second` | HTTP timeout for Alert Manager requests |
| `AlertManagers` | (see main.go) | List of Alert Manager hostnames |

Alert Manager endpoints are accessed via:
```
http://{hostname}:9093/api/v2/alerts?silenced=false&inhibited=false&group=uniq_id
```

## API Reference

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Legacy dashboard |
| POST | `/` | Legacy dashboard with search/actions |
| GET | `/v2` | Modern dashboard |
| POST | `/v2` | Modern dashboard with search/actions |

### POST Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `search-box` | string | Search query (prefix with `!` to exclude) |
| `silence-alert` | string | Hostname to silence (supports multiple values) |
| `remove-silence` | string | Hostname to unsilence (supports multiple values) |

### Response

Both endpoints return HTML. The response includes:

- Server count and total alert count
- Alert data sorted by start time (oldest first)
- Last update timestamp

## Keyboard Shortcuts

Press `?` for quick reference or `H` for the full help guide.

| Key | Action |
|-----|--------|
| `/` | Focus search box |
| `Esc` | Close modal / blur input / close drawer |
| `Alt+1` | Switch to Active tab |
| `Alt+2` | Switch to Silenced tab |
| `F` | Toggle filters panel |
| `T` | Toggle light/dark theme |
| `D` | Cycle view density |
| `C` | Open column configuration |
| `H` | Open help guide |
| `R` | Refresh data |
| `?` | Show keyboard shortcuts |
| `Ctrl+A` | Select all visible alerts |
| `←` `→` | Navigate alerts in drawer |

## Browser Storage

The v2 interface stores preferences in localStorage:

| Key | Description |
|-----|-------------|
| `prometheos-theme` | Light/dark mode preference |
| `prometheos-density` | View density setting |
| `prometheos-refresh-interval` | Auto-refresh interval |
| `prometheos-collapsed-sections` | Analytics section states |
| `prometheos-analytics-range` | Selected analytics time range |
| `prometheos-incidents` | Saved incident data |
| `prometheos-incident-window` | Incident clustering window |
| `prometheos_column_config` | Column visibility and order |
| `prometheos-filter-presets` | Saved filter presets |
| `prometheos-known-alerts` | Tracks seen alerts for new alert detection |

## Project Structure

```
prometheos/
├── main.go                 # Go server with dual endpoint support
├── silencedAlerts          # Persistent silenced hosts (plain text, one per line)
├── templates/
│   ├── v1/                 # Legacy interface templates
│   │   ├── header.gohtml   # HTML head and header
│   │   ├── server.gohtml   # Alert table body
│   │   └── footer.gohtml   # Footer and scripts
│   └── v2/                 # Modern interface templates
│       ├── header.gohtml   # CSS, HTML structure, modals
│       ├── server.gohtml   # Controls, filters, table
│       └── footer.gohtml   # JavaScript managers
└── README.md
```

## Architecture

### Data Flow

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Alert Manager 1 │     │ Alert Manager 2 │     │ Alert Manager N │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────┬───────┴───────────────────────┘
                         │ Parallel HTTP requests
                         ▼
              ┌─────────────────────┐
              │    Prometheos       │
              │  ┌───────────────┐  │
              │  │  AlertStore   │  │  Thread-safe storage
              │  │  (mutex)      │  │
              │  └───────────────┘  │
              │          │          │
              │  ┌───────┴───────┐  │
              │  │    │    │     │  │
              │  ▼    ▼    ▼     ▼  │
              │  v1   v2  API  Silence│
              └─────────────────────┘
                         │
                         ▼
              ┌─────────────────────┐
              │    Web Browser      │
              │  (HTML Dashboard)   │
              └─────────────────────┘
```

### Key Components

| Component | Description |
|-----------|-------------|
| `AlertStore` | Thread-safe storage with RWMutex for concurrent access |
| `getJSON()` | Parallel fetcher for all Alert Manager sources |
| `jsonToStruct()` | Converts Alert Manager JSON to internal structs |
| `searchService()` | Filters alerts by search criteria |
| `createHandler()` | Factory function for version-specific HTTP handlers |

### Technical Details

| Aspect | Implementation |
|--------|----------------|
| Thread Safety | `sync.RWMutex` protects shared alert data |
| Error Handling | Individual Alert Manager failures logged but don't crash server |
| Performance | Goroutines fetch from all sources in parallel |
| Auto-Refresh | Configurable intervals: 30s, 1m, 2m, 5m, or Off |
| Graceful Shutdown | Handles SIGINT/SIGTERM with 30-second timeout |
| HTTP Timeouts | Read: 15s, Write: 15s, Idle: 60s |

## Development

### Running Locally

```bash
# Build and run
go build -o prometheos main.go && ./prometheos

# Or run directly
go run main.go
```

### Adding Alert Managers

Edit the `AlertManagers` slice in `main.go` and rebuild:

```go
AlertManagers: []string{
    "new-alertmanager.example.com",
    // ...
},
```

### Template Development

Templates use Go's `html/template` package with `.gohtml` extension:

- `header.gohtml` — Defines `head` template (HTML head, CSS, opening tags)
- `server.gohtml` — Defines `body` template (alert table/content)
- `footer.gohtml` — Defines `foot` template (scripts, closing tags)

Changes to templates require a server restart.

### Silenced Alerts Storage

The `silencedAlerts` file stores one hostname per line:

```
server1.example.com
server2.example.com
```

This file is shared between both interfaces and persists across restarts.

## Credits

Original Author: Desmond McDermitt
