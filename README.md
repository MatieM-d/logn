# LOGN 🔐

A fast, secure, open-source CLI password manager written in Go.

## Features

- 🔑 **AES-256-GCM** encryption for all stored passwords
- 🛡️ **Argon2id** key derivation from master password
- 📋 **Auto-clear clipboard** after 10 seconds
- ⚡ **Single binary** — no dependencies required on target machine
- 💾 **Local storage** — your passwords never leave your machine

## Installation

### Build from source

```bash
git clone https://github.com/MatieM-d/logn.git
cd logn
go build -o logn.exe .
```

### Requirements

- Go 1.22+

## Usage

### Initialize vault

```bash
logn init
```

### Add a password

```bash
logn add github
```

### Get a password

```bash
logn get github
# Password is copied to clipboard and cleared after 10 seconds
```

### List all services

```bash
logn list
```

### Delete a password

```bash
logn delete github
```

### Generate a password

```bash
logn generate
```

## Security

| Feature | Details |
|---|---|
| Encryption | AES-256-GCM |
| Key derivation | Argon2id (64MB, 2 iterations, 4 threads) |
| Random generation | crypto/rand |
| Storage | Local encrypted `.vault` file |
| Clipboard | Auto-cleared after 10 seconds |
| Master password | Never stored, derived key only |

## Project Structure

```
logn/
├── internal/
│   ├── models.go       # Data structures
│   ├── crypto.go       # AES-256-GCM + Argon2id
│   ├── storage.go      # Read/write .vault file
│   ├── vault.go        # Business logic
│   ├── generator.go    # Password generator
│   └── clipboard.go    # Clipboard + auto-clear
├── main.go             # CLI interface
├── go.mod
└── README.md
```

## Important

- Never commit your `.vault` file to git
- Keep your master password safe — there is no recovery option
- The `.vault` file is stored in the directory where you run the program

## License

MIT License — feel free to use, modify and distribute.
