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

## Modern Interface (`/v2`)

Enhanced dashboard featuring:

### Status Overview
- **At-a-glance metrics** - Active alerts, affected servers, silenced count, monitored sources
- **Color-coded cards** - Instant visual health assessment
- **Real-time sync indicator** - Shows last update time with pulse animation

### Enhanced Data Table
- **Sortable columns** - Click any header to sort ascending/descending
- **Visual status indicators** - Pulsing dots for active alerts
- **Service tags** - Color-coded alert badges
- **Compact design** - High information density

### Search & Filtering
- **Inclusive search** - Type keywords to filter matches
- **Exclusive search** - Prefix with `!` to exclude (e.g., `!guardian`)
- **Multi-field search** - Searches UID, Account, Hostname, Location, Prom Node, Time, Services

### Keyboard Shortcuts (v2 only)
| Key | Action |
|-----|--------|
| `/` | Focus search box |
| `Esc` | Blur search box |
| `Alt+1` | Switch to Active tab |
| `Alt+2` | Switch to Silenced tab |

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

The server starts on port `:8000` by default.

```
http://localhost:8000     # Legacy interface
http://localhost:8000/v2  # Modern interface
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
│       ├── header.gohtml
│       ├── server.gohtml
│       └── footer.gohtml
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
    Port:           ":8000",
    UpdateInterval: 5 * time.Minute,
    RequestTimeout: 10 * time.Second,
    AlertManagers:  []string{...},
}
```

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

---

## Migration Path

1. Deploy with dual endpoints
2. Users can test `/v2` while `/` remains available
3. Gather feedback on the new interface
4. When ready, swap the templates or update routing

---

## Credits

Original Author: Desmond McDermitt  
Enhanced Version: 2024
