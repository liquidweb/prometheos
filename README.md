# Prometheos Alert Dashboard

A high-performance alert aggregation dashboard for Prometheus Alert Managers with dual interface support.

## Endpoints

| Endpoint | Interface | Description |
|----------|-----------|-------------|
| `/` | Legacy (v1) | Original interface - familiar layout for existing users |
| `/v2` | Modern (v2) | New dashboard with enhanced UX and features |

Both interfaces share the same data and silencing state.

---

## Legacy Interface (`/`)

The original Prometheos interface with:
- Simple table layout
- Tab-based active/silenced views
- Basic search functionality
- Light gray theme

---

## Modern Interface (`/v2`)

Enhanced dashboard featuring comprehensive monitoring, analytics, and incident management capabilities.

### Theme & Display

#### Light/Dark Mode
- Toggle between light and dark themes
- Press `T` or click the theme toggle button
- Preference saved to localStorage

#### View Density
- **Compact** - Maximum information density, minimal padding
- **Comfortable** - Balanced view (default)
- **Spacious** - Relaxed layout with more whitespace
- Press `D` to cycle through densities

### Status Overview

Primary metrics cards always visible at the top:

| Card | Description |
|------|-------------|
| **Active Alerts** | Total active alerts with trend indicator |
| **Affected Servers** | Unique servers with active alerts |
| **Silenced** | Currently silenced alert count |
| **Health Score** | 0-100 score based on alert severity, velocity, and incidents |
| **Incidents** | Detected incident groups |
| **Sources** | Monitored Prometheus nodes |

### Analytics Sections

Collapsible sections (click header to expand/collapse). All sections auto-collapse by default on first visit.

#### Distribution Analytics
- **By Location** - Geographic distribution of alerts
- **By Account** - Alert count per account
- **Top Services** - Most frequently alerting services

#### Performance Metrics
- **Avg. Alert Age** - Mean time alerts have been open
- **Oldest Alert** - Longest-running active alert
- **Alert Velocity** - Alerts per hour rate
- **Recurring Systems** - Servers with repeated alerts
- **Recurring Accounts** - Accounts with multiple alerts

#### Alert Velocity Chart
- Visual timeline of alert frequency
- Adapts bucket size based on selected time range
- Time range options: 1H, 6H, 24H, 7D

### Incident Management

Automatically detects and groups alerts that represent significant multi-system outages. Configurable thresholds filter out noise to focus on real incidents.

#### Features
- **Configurable Time Window** - 2min, 5min, 10min, 15min, 30min, or 1 hour
- **Minimum Alert Threshold** - Require 2, 3, 4, 5, or 10+ alerts (default: 3)
- **Minimum Host Threshold** - Require 1, 2, 3, or 5+ unique hosts (default: 2)
- **Auto-detection** - Clusters alerts by start time proximity
- **Smart Naming** - Auto-generated names based on scope (multi-location, multi-account, major)
- **Manual Creation** - Select alerts and create custom incident groups
- **Incident Naming** - Click to rename incidents with meaningful titles
- **Status Tracking** - Investigating, Identified, Monitoring, Resolved
- **Notes** - Add timestamped notes to incidents
- **Rescan** - Re-detect incidents while preserving saved ones

#### Incident Criteria
To form an incident, alerts must meet BOTH thresholds:
- **Min Alerts** - At least N alerts in the time window (default: 3+)
- **Min Hosts** - At least N unique hostnames affected (default: 2+)

This prevents single-server alerts or minor issues from creating incident noise.

#### Saved Incidents
Incidents are preserved during rescan if they have:
- Custom name (changed from auto-generated default)
- Status changed from "Investigating"
- Notes attached

### Column Configuration

Customize which columns appear in the alerts table.

#### Features
- **Drag & Drop Reordering** - Arrange columns in preferred order
- **Toggle Visibility** - Show/hide individual columns
- **Column Profiles** - Quick presets (Default, Compact, Detailed, Custom)
- **API Field Explorer** - Discover additional AlertManager fields
- **Live Preview** - See changes before applying

#### API Explorer
Scans available fields from Prometheus AlertManager:
- **Labels** - alertname, severity, instance, job, env, team, etc.
- **Annotations** - summary, description, runbook_url, dashboard_url
- **Core Fields** - Standard table columns

Press `C` to open column configuration.

### Search & Filtering

#### Basic Search
- Type keywords to filter matching alerts
- Searches across all visible columns

#### Exclusive Search
- Prefix with `!` to exclude matches
- Example: `!guardian` hides guardian alerts

#### Advanced Filters Panel
Press `F` to toggle the filters panel:
- **Location** - Filter by datacenter/region
- **Account** - Filter by account name
- **Service** - Filter by service type
- **Prom Node** - Filter by Prometheus source

#### Filter Presets
- Save current filter combinations as named presets
- Quick-apply saved presets from dropdown
- Delete presets you no longer need

### Export Options

Enhanced export menu with multiple formats:

#### Current View
- **CSV** - Spreadsheet-compatible format
- **JSON** - Structured data format

#### With Annotations
- **CSV with Notes** - Includes any notes you've added
- **JSON with Notes** - Full data with annotations

#### Incidents
- **Incidents CSV** - Export incident groups
- **Incidents JSON** - Full incident data with notes

### Alert Details Drawer

Click any alert row to open a detailed slide-out panel:

- **Full Information** - UID, account, hostname, location, Prom node, start time
- **Duration** - Calculated time since alert started
- **Services** - All services displayed as tags
- **Labels** - Additional data attributes shown
- **Quick Actions**:
  - Silence/Unsilence the alert
  - Copy as formatted text
  - Copy as JSON
- **Navigation** - Use `←` `→` arrow keys to browse through alerts
- **Position Indicator** - Shows "1 of N" for current position

### Bulk Actions

Select multiple alerts using checkboxes to perform operations on all at once:

| Action | Description |
|--------|-------------|
| **Silence** | Silence all selected alerts |
| **Unsilence** | Remove silence from all selected |
| **Create Incident** | Group selected alerts into an incident |
| **Copy** | Copy all selected to clipboard |

The bulk actions bar appears at the bottom of the screen when alerts are selected.

### New Alerts Indicator

Visual notification when new alerts appear:

- **Banner** appears at top when new alerts are detected
- **Count badge** shows number of new alerts
- **Click to scroll** - Jumps to first new alert
- **Green highlight** on new alert rows
- **Persisted state** - Tracks seen alerts across page refreshes

### Shareable URLs

Dashboard state is preserved in the URL for sharing:

- **Share Button** - Copies current view URL to clipboard
- **Encoded State** - Tab, search query, filters, time range
- **Bookmarkable** - Save filtered views as browser bookmarks
- **Back/Forward** - Browser navigation works with dashboard state

Example URL: `/v2?tab=silenced&q=nginx&location=dc1&range=1h`

### Keyboard Shortcuts

Press `?` to view all shortcuts, or `H` for the full help guide.

| Key | Action |
|-----|--------|
| `/` | Focus search box |
| `Esc` | Close modals / Blur inputs / Close drawer |
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
| `←` | Previous alert (in drawer) |
| `→` | Next alert (in drawer) |

### Auto-Refresh

- Configurable refresh intervals: 30s, 1m, 2m, 5m, or Off
- Visual countdown indicator
- Manual refresh with `R` key or refresh button

### Tooltips

Hover over status cards and metrics for explanatory tooltips describing each value.

---

## Alert Managers Monitored

```
c01.b3.alertmanager.pro.mon.liquidweb.com
c02.b3.alertmanager.pro.mon.liquidweb.com
n01.b2.alertmanager.pro.mon.liquidweb.com
n01.b3.alertmanager.pro.mon.liquidweb.com
n01.b4.alertmanager.pro.mon.liquidweb.com
n01.b5.alertmanager.pro.mon.liquidweb.com
n02.b2.alertmanager.pro.mon.liquidweb.com
n02.b3.alertmanager.pro.mon.liquidweb.com
n02.b4.alertmanager.pro.mon.liquidweb.com
n02.b5.alertmanager.pro.mon.liquidweb.com
n03.b3.alertmanager.pro.mon.liquidweb.com
n04.b3.alertmanager.pro.mon.liquidweb.com
n05.b3.alertmanager.pro.mon.liquidweb.com
n06.b3.alertmanager.pro.mon.liquidweb.com
n07.b3.alertmanager.pro.mon.liquidweb.com
n08.b3.alertmanager.pro.mon.liquidweb.com
```

---

## Installation

```bash
# Clone or copy files
git clone <repository>
cd prometheos

# Build
go build -o prometheos main.go

# Run
./prometheos
```

The server starts on port `:8001` by default.

```
http://localhost:8001     # Legacy interface
http://localhost:8001/v2  # Modern interface
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
├── silencedAlerts             # Persistent silenced hosts (shared)
└── README.md
```

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

## Browser Storage

The v2 interface stores user preferences in localStorage:

| Key | Description |
|-----|-------------|
| `prometheos-theme` | Light/dark mode preference |
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

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Legacy dashboard |
| POST | `/` | Legacy dashboard with search/actions |
| GET | `/v2` | Modern dashboard |
| POST | `/v2` | Modern dashboard with search/actions |

### POST Parameters (both versions)

| Parameter | Description |
|-----------|-------------|
| `search-box` | Search query (prefix with `!` to exclude) |
| `silence-alert` | Hostname to silence |
| `remove-silence` | Hostname to unsilence |

---

## Technical Details

### Thread Safety
- Mutex-protected `AlertStore` prevents race conditions
- Safe concurrent access from multiple HTTP requests

### Error Handling
- Individual alert manager failures don't crash the server
- Graceful degradation with partial data

### Performance
- Parallel fetching from all alert managers
- 5-minute update interval (configurable)
- Request timeout prevents hanging connections

### Graceful Shutdown
- Handles SIGINT/SIGTERM signals
- Completes in-flight requests before stopping

### Client-Side Features (v2)
- All preferences stored in localStorage
- No server-side state for UI preferences
- Incident data persists across page refreshes
- Column configurations survive browser restarts

---

## Migration Path

1. Deploy with dual endpoints
2. Users can test `/v2` while `/` remains available
3. Gather feedback on the new interface
4. When ready, swap the templates or update routing

---

## Credits

Original Author: Desmond McDermitt
Enhanced Version: 2024-2025
