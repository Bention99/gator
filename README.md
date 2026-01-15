# gator

`gator` is a command-line RSS feed aggregator written in Go.
It allows users to register, manage RSS feeds, follow and unfollow feeds, and browse aggregated posts stored in a PostgreSQL database.

---

## Prerequisites

To run `gator`, you need the following installed on your system:

### 1. Go
- Version 1.20 or newer recommended
- Installation instructions: https://go.dev/doc/install

Verify installation:
```
go version
```

### 2. PostgreSQL
- A running PostgreSQL instance is required
- Installation instructions: https://www.postgresql.org/download/

Verify installation:
```
psql --version
```

You will also need a database and connection string (e.g. `postgres://user:password@localhost:5432/gator`).

---

## Installation

Install the `gator` CLI using `go install`:

```
go install github.com/<your-username>/gator@latest
```

After installation, ensure that `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

Verify installation:
```
gator
```

---

## Configuration

`gator` uses a JSON configuration file stored in your home directory.

### Config file location
```
~/.gatorconfig.json
```

### Example configuration

```
{
  "db_url": "postgres://user:password@localhost:5432/gator",
  "current_user_name": ""
}
```

- `db_url`: PostgreSQL connection string
- `current_user_name`: Automatically set when you log in

---

## Running the Program

Once installed and configured, you can run `gator` directly from the terminal:

```
gator <command> [arguments]
```

---

## Available Commands

### Authentication & User Management
- `login <name>` – Log in as an existing user
- `register <name>` – Register a new user
- `reset` – Reset application state
- `users` – List all registered users

### Feed Management
- `addfeed <name> <url>` – Add a new RSS feed
- `feeds` – List all available feeds

### Following Feeds
- `follow <url>` – Follow a feed
- `following` – List followed feeds
- `unfollow <url>` – Unfollow a feed

### Aggregation & Browsing
- `agg <timeframe>` – Fetch and aggregate posts
- `browse [limit]` – Browse aggregated posts

Example:
```
gator browse 10
```

---

## Notes

- This project uses PostgreSQL as its persistence layer.
- Database migrations must be applied before first use.
- RSS feeds may vary in date formats; unsupported formats are safely ignored.

---

## License

This project is provided for educational purposes.