# 🛡️ Database Backup Utility

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Active-success?style=for-the-badge)

A robust, command-line interface (CLI) utility written in **Go** for backing up and restoring various databases. It supports multiple database management systems, local and cloud storage options, compression, and notifications.

> 📚 **Project Reference**: This project is based on the [Database Backup Utility](https://roadmap.sh/projects/database-backup-utility) challenge from [roadmap.sh](https://roadmap.sh).

---

## 🚀 Features

*   **Multi-Database Support**:
    *   🐬 **MySQL** (`mysqldump`)
    *   🐘 **PostgreSQL** (`pg_dump`)
    *   🍃 **MongoDB** (`mongodump`)
    *   🗄️ **SQLite** (File copy)
*   **Flexible Storage**:
    *   📂 **Local Filesystem**
    *   ☁️ **AWS S3**
    *   ☁️ **Google Cloud Storage (GCS)**
*   **Advanced Capabilities**:
    *   📦 **Compression**: Automatic Gzip compression (`.gz`) to save space.
    *   🔔 **Notifications**: Real-time Slack notifications for backup success/failure.
    *   📝 **Logging**: Comprehensive activity logging.

---

## 🛠️ Installation

### Prerequisites
*   **Go** (v1.21 or higher)
*   Database CLI tools installed on the host machine (e.g., `mysqldump`, `pg_dump`, `mongodump`).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/adityapawar327/Database-Backup-Utility.git
cd Database-Backup-Utility

# Build the binary
go build -o backup-tool ./cmd/backup-tool
```

---

## ⚙️ Configuration

Initialize a default configuration file:

```bash
./backup-tool init
```

This creates a `db_backup_config.yaml` file. Edit it with your credentials:

```yaml
storage:
  type: s3 # Options: local, s3, gcs
  path: ./backups # For local storage
  bucket: my-backup-bucket # For S3/GCS
  region: us-east-1 # For S3

notifications:
  slack_webhook: "https://hooks.slack.com/services/..."

databases:
  my_mysql_db:
    type: mysql
    host: localhost
    port: 3306
    user: root
    password: password
    database: my_app_db

  my_postgres_db:
    type: postgres
    host: localhost
    port: 5432
    user: postgres
    password: password
    database: analytics_db
```

---

## 📖 Usage

### 1. Backup a Database

Run a backup for a specific database defined in your config:

```bash
./backup-tool backup my_mysql_db
```

**What happens?**
1.  Connects to the database.
2.  Creates a dump/backup file.
3.  Compresses the file (`.gz`).
4.  Uploads it to the configured storage (Local, S3, or GCS).
5.  Sends a Slack notification.

### 2. Restore a Database

Restore a database from a backup file (local path or cloud object key):

```bash
./backup-tool restore backups/temp_my_mysql_db_20231123.sql.gz my_mysql_db
```

**What happens?**
1.  Downloads the backup file from storage (if remote).
2.  Decompresses the file.
3.  Restores the data into the specified database.

### 3. List Backups

List available backup files in the configured storage:

```bash
./backup-tool list
```

---

## 📂 Project Structure

```
.
├── cmd/
│   └── backup-tool/
│       └── main.go          # CLI Entry Point
├── pkg/
│   ├── core/                # Interfaces (Database, Storage)
│   ├── databases/           # DB Adapters (MySQL, Postgres, etc.)
│   ├── storage/             # Storage Adapters (Local, S3, GCS)
│   └── utils/               # Utilities (Logger, Compressor, Notifier)
├── db_backup_config.yaml    # Configuration File
├── go.mod                   # Go Modules
└── README.md                # Documentation
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
