# LOGN 🔐

A fast, secure, open-source CLI password manager written in Go.

## Features

- 🔑 **AES-256-GCM** encryption for all stored passwords
- 🛡️ **Argon2id** key derivation from master password
- 📋 **Auto-clear clipboard** after 10 seconds
- ⚡ **Single binary** — no dependencies required on target machine
- 💾 **Local storage** — your passwords never leave your machine
- 🎨 **Color-coded interface** — instant visual feedback on password strength
- 🔍 **Search** — quickly find passwords by service name
- 🛡️ **Password checker** — audit all your passwords at once
- 💾 **Backup & restore** — easily back up your vault
- 📤 **Import / Export** — supports Bitwarden, LastPass, KeePass and CSV

## Installation

### Build from source

```bash
git clone https://github.com/MatieM-d/logn.git
cd logn
go build -o logn.exe .
```

### Requirements

- Go 1.22+

### Run from anywhere (Windows)

Add the folder containing `logn.exe` to your system `PATH`:

1. Open **Start** → search **Environment Variables**
2. Click **Environment Variables**
3. Under **System variables** find `Path` → **Edit**
4. Click **New** and add your folder path (e.g. `D:\Projects\logn`)
5. Click **OK** and restart your terminal

## Usage

### Initialize vault

```bash
logn init
```

On first run you will be asked where to store the `.vault` file. You can use the current directory or specify a custom path.

### Add a password

```bash
logn add github
# Enter login, password or leave empty to auto-generate
```

### Get a password

```bash
logn get github
# Password is copied to clipboard and cleared after 10 seconds
```

### Edit a password

```bash
logn edit github
# Leave any field empty to keep the current value
# Type "generate" in the password field to auto-generate a new one
```

### List all passwords

```bash
logn list
# Shows all services with logins and passwords
# Passwords are color-coded: green = strong, red = weak
```

### Search by service name

```bash
logn search git
# Returns all entries matching the query
```

### Check password strength

```bash
# Check a specific password
logn check github

# Check all passwords at once
logn check
```

Password checker verifies:
- Minimum length of 8 characters
- At least one uppercase letter
- At least one digit
- At least one special character (!@#$%...)

### Delete a password

```bash
logn delete github
```

### Generate a password

```bash
logn generate
# Generates a secure 20-character password
```

### Export passwords

```bash
logn export
# Exports all passwords to a CSV file
# Warning: exported file contains passwords in plain text!
```

### Import passwords

```bash
# Import from CSV (duplicates are skipped)
logn import bitwarden_export.csv

# Import and overwrite existing entries
logn import bitwarden_export.csv --overwrite
```

Supported import formats (auto-detected):
- **Bitwarden** — export from Account → Tools → Export
- **LastPass** — export from Account Options → Export
- **KeePass** — export from File → Export → CSV
- **LOGN** — previously exported LOGN CSV

### Backup vault

```bash
logn backup
# Creates a timestamped backup: logn-backup-2024-01-15_10-30-00.vault
```

### Restore from backup

```bash
logn restore D:\Projects\logn\backups\logn-backup-2024-01-15_10-30-00.vault
```

> You can also copy the `.vault` file to another machine and use it there with the same master password.

## Security

| Feature | Details |
|---|---|
| Encryption | AES-256-GCM |
| Key derivation | Argon2id (64MB, 2 iterations, 4 threads) |
| Random generation | crypto/rand |
| Storage | Local encrypted `.vault` file |
| Clipboard | Auto-cleared after 10 seconds |
| Master password | Never stored, derived key only |
| Vault file permissions | 600 (owner read/write only) |

## How it works

```
Master password + Salt → Argon2id → Encryption key
                                         ↓
                               AES-256-GCM encrypt
                                         ↓
                               .vault file on disk
```

Your master password is **never stored anywhere**. Every time you run LOGN it derives the encryption key from your password and the salt stored in the `.vault` file. If you forget your master password, there is no recovery option.

## Project Structure

```
logn/
├── internal/
│   ├── models.go       # Data structures
│   ├── crypto.go       # AES-256-GCM + Argon2id
│   ├── storage.go      # Read/write .vault file
│   ├── vault.go        # Business logic
│   ├── generator.go    # Password generator
│   ├── clipboard.go    # Clipboard + auto-clear
│   ├── colors.go       # Color-coded CLI output
│   ├── export.go       # CSV export
│   └── import.go       # CSV import
├── main.go             # CLI interface
├── go.mod
└── README.md
```

## Important

- **Never commit** your `.vault` file or `config.json` to git
- Keep your **master password safe** — there is no recovery option
- The `.vault` location is set during `logn init` and saved in `config.json`
- Backups are stored in the `backups/` folder and are also excluded from git
- Exported CSV files contain passwords in **plain text** — handle with care

## License

Copyright (C) 2024 MatieM-d

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
