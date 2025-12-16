# GNET

**Disclaimer**

GNET is an open-source command-and-control (C2) framework intended **solely** for ethical research, defensive security testing, and educational development.

The developers of GNET **do not** condone or take responsibility for misuse. You are **required** to comply with all applicable local, state, and federal laws. Misuse may result in severe criminal penalties.

---

## Requirements

* Go **1.25.4** or newer
* Windows or Linux

---

## Installation

### 1. Install Go

Choose the instructions for your operating system.

**Ubuntu / Debian**

```bash
sudo apt update && sudo snap install go --classic
```

**Kali Linux**

```bash
sudo apt update && sudo apt install golang -y
```

**RHEL / Rocky / AlmaLinux / CentOS Stream 8/9**

```bash
sudo dnf check-update && sudo dnf install golang -y
```

**Alpine Linux**

```bash
sudo apk update && sudo apk add go
```

**Windows / macOS / Linux (official installer)**

```
https://go.dev/doc/install
```

Verify installation:

```bash
go version
```

---

## Configuration

### 2. Server Configuration

1. Open `server.go` in the project root.
2. Locate the server IP variable.
3. Replace the value with the public IP or hostname of your VPS or host machine.

### 3. Bot Configuration

1. Navigate to:

   ```
   ./bins/compiles/bot.go
   ```
2. Locate the IP variable used by the bot.
3. Set it to the **same IP or hostname** configured in `server.go`.

The server and bot **must** point to the same address to communicate correctly.

---

## Building and Running

### 4. Build the Server

From the project root:

```bash
go build server.go
```

### 5. Run the Server

**Linux**

```bash
chmod +x server
./server &
```

**Windows**

```bash
server
```

The server will now listen for incoming bot connections.

---

## User Management

### Adding Accounts

1. Open the following file:

   ```
   ./data/users.json
   ```
2. Follow the existing JSON structure.
3. Add a new username and password entry.

Example structure:

```json
{
  "user": "username",
  "pass": "password"
}
```

Save the file after making changes. New accounts take effect immediately on the next authentication attempt.

---

## Notes

* This project is intended for controlled environments only.
* Do **not** expose the server to the public internet without proper safeguards.
* Always obtain explicit authorization before testing any system.

---

## License

Open source. Use responsibly.
