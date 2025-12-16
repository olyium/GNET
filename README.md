# GNET Official readme.MD

_Disclaimed:_ **GNET is an open source, C2, meant for the soul purpose, of ethical researching, and educational development. The developers of GNET, do not claim your actions, that you commit towards, whilst using GNET. You MUST, follow the laws of your jursidiction, failure to do so can result in life changing crimminal charges.**

## Requirements

- Golang 1.25.4
- Windows/Linux OS

## Setup

1. Install Golang:
**Ubuntu/Debian**: ```bash
sudo apt update && sudo snap install go --classic```
**Kali Linux**: ```bash
sudo apt update && sudo apt install golang -y```
**RHEL / Rocky / AlmaLinux / CentOS Stream 8 / 9**: ```bash
sudo dnf check-update && sudo dnf install golang -y```
**Alpine Linux**: ```bash
sudo apk update && sudo apk add go```
**Windows 10/11 / Linux / Mac**: ```bash
https://go.dev/doc/install```

2. Configure Server:
Locate **server.go** ```./server.go```, then within the variable, there is an IP. Change this IP, to your VPS's IP, or which IP the host is. Then locate the compiles folder ```./bins/compiles```, within here you will find **bot.go**. Open **bot.go** and do the same with the IP located in **bot.go**

3. Launch Server:
Go back to the main directory, and build the server. ```go build server.go```. Once the build is succesfull you want to run (if you're on Linux) ```chmod +x server``` then you want to run the server. (if you're on Linux) ```./server &```, (if you're on Windows) ```server &```

## Adding Accounts

1. Locate ```./data/users.json```. Then you simply, follow the file structure. Add user, and pass.