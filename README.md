# Prometheos Alert Dashboard

A high-performance alert aggregation dashboard for Prometheus Alert Managers with dual interface support.

## Quick Start

```bash
# Build
go build -o prometheos main.go

# Run
./prometheos
```

The server starts on port `:8001` by default.

- **Legacy Interface:** http://localhost:8001/
- **Modern Interface:** http://localhost:8001/v2

## Interfaces

| Endpoint | Interface | Description |
|----------|-----------|-------------|
| `/` | Legacy (v1) | Original interface with simple table layout and basic search |
| `/v2` | Modern (v2) | Feature-rich dashboard with analytics, incidents, and advanced filtering |

Both interfaces share the same alert data and silencing state.

---

## Modern Interface Features (`/v2`)

### Status Cards

Real-time metrics displayed at the top of the dashboard:

| Card | Description |
|------|-------------|
| **Active Alerts** | Total alerts requiring attention |
| **Affected Servers** | Unique servers currently alerting |
| **Incidents** | Auto-detected groups of correlated alerts |
| **Silenced** | Currently muted alert count |
| **Alert Managers** | Number of sources being monitored (click for details) |

### Theme & Display

- **Light/Dark Mode** — Press `T` or click the theme toggle
- **View Density** — Press `D` to cycle: Compact, Comfortable, Spacious
- Preferences persist in localStorage

### Search & Filtering

**Basic Search**
- Type keywords to filter alerts across all columns
- Prefix with `!` to exclude matches (e.g., `!guardian` hides guardian alerts)
- Press `/` to focus the search box

**Advanced Filters** (Press `F`)
- Data Center, Section, Row, Rack
- Prom Node
- Time Range (1H, 6H, 24H, 7D)
- Service

**Filter Presets**
- Save filter combinations as named presets
- Quick-apply saved presets from dropdown

### Table Columns

Default columns include:
- UID, Account, Hostname
- DC, Section, Row, Rack, Location
- Prom Node, Start Time, Duration
- Services, Actions

**Column Configuration** (Press `C`)
- Drag & drop to reorder columns
- Toggle column visibility
- Presets: Default, Compact, Detailed, Custom
- API Explorer to discover additional AlertManager fields

### Analytics Sections

Collapsible sections with insights:

- **Distribution** — Alerts by location, account, and top services
- **Performance Metrics** — Average alert age, oldest alert, alert velocity, recurring systems
- **Alert Velocity Chart** — Visual timeline with configurable time ranges

### Incident Management

Automatically detects and groups alerts that may represent correlated infrastructure issues.

**Detection Criteria**
- Configurable time window (2min to 1 hour)
- Minimum alert threshold (2-10+ alerts)
- Minimum host threshold (1-5+ unique hosts)

**Managing Incidents**
- Click incident name to rename
- Set status: Investigating, Identified, Monitoring, Resolved
- Add timestamped notes
- Rescan to re-detect (preserves saved incidents)

**Location-Based Grouping**

Location codes follow the format `B3S2R9K8U27`:
- **B3** — Datacenter 3
- **S2** — Section 2
- **R9** — Row 9
- **K8** — Rack/Cabinet 8
- **U27** — Rack Unit 27

### Alert Details Drawer

Click any alert row to open a slide-out panel with:
- Full details: UID, account, hostname, location, duration
- All services displayed as tags
- Quick actions: Silence/unsilence, copy as text or JSON
- Navigate between alerts with `←` `→` arrow keys

### Bulk Operations

Select multiple alerts using checkboxes:
- **Silence** — Mute all selected alerts
- **Unsilence** — Unmute all selected alerts
- **Create Incident** — Group selected alerts
- **Copy** — Copy selected alerts to clipboard

### Export Options

- **CSV / JSON** — Export current filtered view
- **With Annotations** — Include notes
- **Incidents** — Export incident data

### Shareable URLs

Dashboard state is encoded in the URL:
- Share button copies URL with current tab, search, filters, and time range
- URLs are bookmarkable
- Browser back/forward navigation works

---

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

---

## Configuration

Edit the config struct in `main.go`:

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
    AlertManagers:  []string{...},
}
```

---

## Project Structure

```
prometheos/
├── main.go                    # Go server with dual endpoint support
├── templates/
│   ├── v1/                    # Legacy templates
│   │   ├── header.gohtml
│   │   ├── server.gohtml
│   │   └── footer.gohtml
│   └── v2/                    # Modern templates
│       ├── header.gohtml      # CSS, HTML structure, modals
│       ├── server.gohtml      # Controls, filters, table
│       └── footer.gohtml      # JavaScript managers
├── silencedAlerts             # Persistent silenced hosts (shared between interfaces)
└── README.md
```

---

## Browser Storage

The v2 interface stores preferences in localStorage:

| Key | Description |
|-----|-------------|
| `prometheos-theme` | Light/dark mode |
| `prometheos-density` | View density setting |
| `prometheos-refresh-interval` | Auto-refresh interval |
| `prometheos-collapsed-sections` | Analytics section states |
| `prometheos-analytics-range` | Selected time range |
| `prometheos-incidents` | Saved incident data |
| `prometheos-incident-window` | Incident clustering window |
| `prometheos_column_config` | Column visibility and order |
| `prometheos-filter-presets` | Saved filter presets |
| `prometheos-known-alerts` | Tracks seen alerts for new alert detection |

---

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Legacy dashboard |
| POST | `/` | Legacy dashboard with search/actions |
| GET | `/v2` | Modern dashboard |
| POST | `/v2` | Modern dashboard with search/actions |

### POST Parameters

| Parameter | Description |
|-----------|-------------|
| `search-box` | Search query (prefix with `!` to exclude) |
| `silence-alert` | Hostname to silence |
| `remove-silence` | Hostname to unsilence |

---

## Technical Details

- **Thread Safety** — Mutex-protected AlertStore prevents race conditions
- **Error Handling** — Individual alert manager failures don't crash the server
- **Performance** — Parallel fetching from all alert managers
- **Auto-Refresh** — Configurable intervals (30s, 1m, 2m, 5m, or Off)
- **Graceful Shutdown** — Handles SIGINT/SIGTERM signals

---

## Credits

Original Author: Desmond McDermitt
