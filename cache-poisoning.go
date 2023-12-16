package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
)

const banner = `
   ______     ____        _                
  / ____/    / __ \____  (_)________  ____ 
 / /  ______/ /_/ / __ \/ / ___/ __ \/ __ \
/ /__/_____/ ____/ /_/ / (__  ) /_/ / / / /
\____/    /_/    \____/_/____/\____/_/ /_/ 
                                           
	adrianalvird [ 1.0.0 ]
`

func signalHandler(c chan os.Signal) {
	for {
		sig := <-c
		fmt.Printf("\nExecution stopped by signal: %v\n", sig)
		os.Exit(0)
	}
}

func loadTargetURLs() []string {
	var targetURLs []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		targetURLs = append(targetURLs, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Error reading from stdin: %s\n", err)
		os.Exit(1)
	}
	return targetURLs
}

func performCachePoisoningAttack(url string, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[-] Error creating request for %s: %s\n", url, err)
		return
	}

	// Make an initial request to capture the ETag
	initialResp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] Error capturing ETag for %s: %s\n", url, err)
		return
	}

	// Check if the response contains an ETag header
	if etag, ok := initialResp.Header["ETag"]; ok {
		// Use the captured ETag value in a conditional request
		req.Header.Set("If-None-Match", etag[0])
		conditionalResp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[-] Error making conditional request for %s: %s\n", url, err)
			return
		}

		// Check the conditional response to see if the resource has changed
		if conditionalResp.StatusCode == http.StatusNotModified {
			fmt.Printf("[+] Cache poisoning successful at %s\n", url)
			return
		}
	}

	fmt.Printf("[-] Cache poisoning failed at %s\n", url)
}

func main() {
	fmt.Println(banner)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go signalHandler(c)

	var wg sync.WaitGroup

	client := &http.Client{}

	targetURLs := loadTargetURLs()

	for _, url := range targetURLs {
		wg.Add(1)
		go performCachePoisoningAttack(url, client, &wg)
	}

	wg.Wait()
}

