---
weight: 4
title: "Databases"
---

# Status DB

Detect and display information about database management systems (DBMS) installed on the system.

## Usage

```bash
# Show all detected databases
allbctl status db

# Show specific database
allbctl status db sqlite3
allbctl status db postgres
allbctl status db mysql

# Show detailed information
allbctl status db --detail
allbctl status db sqlite3 --detail
```

## Output

### Summary Mode (Default)
Shows detected databases with status and version:

```
sqlite3:     installed sqlite3 3.37.2, 3 .db files
postgres:    running PostgreSQL 14.5
mysql:       installed mysql Ver 8.0.36
redis:       running Redis server v=6.0.16
```

### Detailed Mode
Shows comprehensive information for each database:

```
Sqlite3:
----------------------------------------
  Client Binary:  /usr/bin/sqlite3
  Client Version: 3.37.2 2022-01-06 13:25:41
  Database Files: (3 found in ~/src)
    - ~/src/test-project/test.db
    - ~/src/test-project/test.sqlite
    - ~/src/test-project/app.sqlite3
  Environment Variables:
    SQLITE_HISTORY=/home/user/.sqlite_history

Postgres:
----------------------------------------
  Client Binary:  /usr/bin/psql
  Client Version: psql (PostgreSQL) 14.5
  Server Binary:  /usr/lib/postgresql/14/bin/postgres
  Status:         RUNNING
  Environment Variables:
    PGHOST=localhost
    PGPORT=5432
    PGUSER=postgres

Databases detected: sqlite3, postgres
```

### Specific Database
Shows information for one database only:

```bash
allbctl status db sqlite3 --detail
```

If the database is not detected:
```
Database 'mongodb' not detected on this system
```
(exit code 1)

## Supported Databases

1. **SQLite3** - File-based SQL database
   - Finds .db, .sqlite, .sqlite3 files in ~/src
   - Shows file paths in detailed mode

2. **MySQL** - Open-source RDBMS
   - Detects client and server
   - Shows running status

3. **MariaDB** - MySQL fork
   - Detects client and server
   - Shows running status

4. **PostgreSQL** - Advanced open-source RDBMS
   - Detects client and server
   - Shows running status

5. **MongoDB** - NoSQL document database
   - Detects mongosh client and mongod server
   - Shows running status

6. **Redis** - In-memory key-value store
   - Detects redis-cli and redis-server
   - Shows running status

7. **Cassandra** - Distributed NoSQL database
   - Detects cqlsh client and cassandra server
   - Shows running status

8. **Oracle** - Enterprise RDBMS
   - Detects sqlplus client
   - Shows environment variables

9. **SQL Server** - Microsoft RDBMS
   - Detects sqlcmd client
   - Shows environment variables

## Information Displayed

### Summary Mode
- Database name
- Running status (installed/running)
- Client version (first line, truncated)
- Number of database files (file-based DBs only)

### Detailed Mode
- Client binary path
- Client version (full output)
- Server binary path (if applicable)
- Server version (if different from client)
- Server running status (RUNNING / not running)
- Database files with full paths (file-based DBs)
- Environment variables (filtered by database-specific prefixes)
- Summary line listing all detected databases

## Detection Logic

### Client Detection
- Checks if client binary exists in PATH (`which <binary>`)
- Gets version using standard version flags
- Only shows database if client is installed

### Server Detection
- Checks if server binary exists in PATH
- Uses `pgrep` to check if server process is running
- Shows "RUNNING" or "not running" status

### Database Files
- Searches `~/src` directory recursively
- Finds files with known extensions (SQLite: .db, .sqlite, .sqlite3)
- Shows count in summary, full paths in detailed mode

### Environment Variables
Collects variables with known prefixes:
- SQLite: `SQLITE_*`
- MySQL: `MYSQL_*`
- MariaDB: `MYSQL_*`, `MARIADB_*`
- PostgreSQL: `PG*`, `POSTGRES_*`
- MongoDB: `MONGO_*`
- Redis: `REDIS_*`
- Cassandra: `CASSANDRA_*`
- Oracle: `ORACLE_*`, `TNS_*`
- SQL Server: `MSSQL_*`

## Examples

### Development Machine
```bash
$ allbctl status db
sqlite3:     installed sqlite3 3.37.2, 5 .db files
postgres:    running PostgreSQL 14.5
redis:       running Redis server v=6.0.16
```

### Production Server
```bash
$ allbctl status db
postgres:    running PostgreSQL 14.5
mongodb:     running db version v5.0.9
redis:       running Redis server v=6.2.7
```

### SQLite Project Details
```bash
$ allbctl status db sqlite3 --detail
Sqlite3:
----------------------------------------
  Client Binary:  /usr/bin/sqlite3
  Client Version: 3.37.2 2022-01-06 13:25:41
  Database Files: (12 found in ~/src)
    - ~/src/myapp/data/users.db
    - ~/src/myapp/data/products.db
    - ~/src/blog/blog.sqlite3
    ...
```

## Integration

The main `allbctl status` command includes a one-line database summary:

```
Runtimes:  Python (3.12.3), Go (1.25.5)
Databases: sqlite3, postgres (running), redis (running)
```

Shows detected databases with "(running)" indicator for active servers.

## Notes

- Only detects databases with installed client binaries
- File search limited to ~/src directory for performance
- Server status detection requires process running as same user
- Environment variables read from current shell environment
