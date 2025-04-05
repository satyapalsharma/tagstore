# TagStore

**TagStore** is a command-line tool for tagging and searching files on your filesystem. Organize files with arbitrary tags and quickly locate them using boolean tag queries. All data is stored locally in a SQLite database.

---

## Features

- **Tag files** — Associate one or more tags with any file on your filesystem.
- **Remove tags** — Detach tags from files when they are no longer needed.
- **Boolean search** — Find files using AND, OR, and NOT logic in a single query string.
- **List associations** — View all tags for a given file, or all files for a given tag.
- **CSV import/export** — Bulk-import tag associations from a CSV file and export the current database to CSV for backup or migration.
- **Local SQLite storage** — All data lives in a single SQLite file on your machine; no network, no cloud, no accounts.
- **Configurable logging** — Control verbosity with the `TAGSTORE_LOG_LEVEL` environment variable.
- **Graceful shutdown** — Handles OS interrupt and termination signals cleanly.

---

## Project Structure

```
tagstore/
├── main.go                    # Entry point; signal handling and panic recovery
├── go.mod                     # Go module definition
├── README.md                  # This file
├── .env.example               # Environment variable template
├── LICENSE                    # MIT License
│
├── cmd/
│   └── tag/
│       └── main.go            # Binary entry point
│
├── internal/
│   └── commands/
│       ├── add.go             # Legacy add-tag command (internal)
│       ├── remove.go          # Legacy remove-tag command (internal)
│       ├── list_files.go      # Legacy list-files command (internal)
│       └── import.go          # CSV import helper (internal)
│
├── src/
│   ├── cli/                   # Cobra command definitions
│   │   ├── root.go            # Root command; wires config, DB, services
│   │   ├── add.go             # `tagstore add <file> <tag>...`
│   │   ├── remove.go          # `tagstore remove <file> <tag>...`
│   │   ├── list.go            # `tagstore list <file|tag> <name>`
│   │   ├── search.go          # `tagstore search [--and|--or] <tag>...`
│   │   ├── init.go            # `tagstore init`
│   │   ├── import.go          # `tagstore import <input.csv>`
│   │   └── export.go          # `tagstore export <output.csv>`
│   │
│   ├── config/
│   │   └── config.go          # Configuration loaded from environment variables
│   │
│   ├── lib/
│   │   ├── database.go        # SQLite connection wrapper; table creation
│   │   ├── logger.go          # Structured logger with configurable levels
│   │   └── errors.go          # Typed application errors
│   │
│   ├── models/
│   │   ├── file.go            # File model and validation
│   │   ├── tag.go             # Tag model and validation
│   │   ├── file_tag.go        # FileTag (join) model
│   │   └── errors.go          # Model-level error definitions
│   │
│   ├── repository/            # Data-access layer (CRUD + search)
│   │   ├── tag_repository.go
│   │   ├── file_repository.go
│   │   └── file_tag_repository.go
│   │
│   ├── service/               # Business-logic layer
│   │   ├── tag_service.go
│   │   └── search_service.go
│   │
│   ├── types/
│   │   ├── tag.go             # Domain types (Tag, File, FileTag, TagQuery)
│   │   └── interfaces.go      # Repository interfaces
│   │
│   └── utils/
│       ├── validation.go      # Path and tag-name validators
│       └── parser.go          # Boolean query parser
```

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23+ |
| CLI framework | [Cobra](https://github.com/spf13/cobra) v1.10 |
| Database | SQLite (via `go-sqlite3`) |
| Logging | Standard library `log` with custom levels |

---

## Installation

### Prerequisites

- Go 1.23 or later
- A C compiler (CGO is required for `go-sqlite3`)

### Build from source

```bash
# Clone the repository
 git clone https://github.com/your-username/tagstore.git
 cd tagstore

# Download dependencies
 go mod download

# Build the binary
 go build -o tagstore ./cmd/tag

# (Optional) Install to $GOPATH/bin
 go install ./cmd/tag
```

After installation, the `tagstore` command will be available in your terminal.

---

## Configuration

TagStore reads two environment variables:

| Variable | Description | Default |
|----------|-------------|--------|
| `TAGSTORE_DB` | Path to the SQLite database file | `~/.tagstore/tagstore.db` |
| `TAGSTORE_LOG_LEVEL` | Logging verbosity (`debug`, `info`, `warn`, `error`) | `info` |

You can copy `.env.example` to `.env` and export the variables before running TagStore, or set them inline:

```bash
export TAGSTORE_DB=/custom/path/to/tagstore.db
export TAGSTORE_LOG_LEVEL=debug
```

---

## Usage

### Initialize the database

```bash
tagstore init
```

This creates the SQLite database file and all required tables. The database is created automatically on first run, so `init` is optional.

### Add tags to a file

```bash
# Add a single tag
tagstore add /path/to/document.pdf work

# Add multiple tags at once
tagstore add /path/to/document.pdf work urgent 2024
```

### Remove tags from a file

```bash
tagstore remove /path/to/document.pdf urgent
```

### List tags for a file

```bash
tagstore list file /path/to/document.pdf
```

Output:
```
work
urgent
2024
```

### List files for a tag

```bash
tagstore list tag work
```

Output:
```
/path/to/document.pdf
/path/to/report.xlsx
```

### Search files by tags

```bash
# OR search (default) — files with ANY of the tags
tagstore search work personal

# AND search — files with ALL of the tags
tagstore search --and work urgent

# Mixed boolean query using the search service directly
# (space = AND, | = OR, - = NOT)
tagstore search "work -urgent | personal"
```

### Import tags from CSV

```bash
tagstore import tags.csv
```

The CSV file must have a header row with `file_path` and `tag_name` columns:

```csv
file_path,tag_name
/path/to/document.pdf,work
/path/to/document.pdf,urgent
/path/to/report.xlsx,work
```

### Export tags to CSV

```bash
tagstore export backup.csv
```

---

## Boolean Query Syntax

When using the search service (or the `search` command), queries support three operators:

| Operator | Meaning | Example |
|----------|---------|---------|
| Space (implicit AND) | All adjacent tags must be present | `work urgent` |
| `|` (pipe) | OR — any of the grouped tags may be present | `work \| personal` |
| `-` (hyphen prefix) | NOT — the tag must NOT be present | `-urgent` |

**Examples:**

```
work urgent               # files tagged with BOTH work AND urgent
work | personal            # files tagged with work OR personal
work -urgent               # files tagged with work but NOT urgent
work urgent -archived      # files tagged with work AND urgent, but NOT archived
```

---

## API Documentation

TagStore is organized in three layers:

### Repository Layer (`src/repository/`)

| Interface | Methods |
|-----------|---------|
| `TagRepository` | `Create`, `GetByName`, `GetByID`, `List`, `Delete` |
| `FileRepository` | `Create`, `GetByPath`, `GetByID`, `List`, `Delete` |
| `FileTagRepository` | `AddTag`, `RemoveTag`, `GetTagsForFile`, `GetFilesForTag`, `HasTag`, `SearchFiles`, `SearchFilesByTags` |

### Service Layer (`src/service/`)

| Service | Methods |
|---------|---------|
| `TagService` | `AddTagToFile`, `RemoveTagFromFile`, `ListTagsForFile`, `ListFilesForTag` |
| `SearchService` | `SearchFiles`, `SearchFilesByTags`, `ListAllTags` |

### CLI Layer (`src/cli/`)

All commands are built with Cobra. Run `tagstore --help` or `tagstore <command> --help` for inline usage information.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
