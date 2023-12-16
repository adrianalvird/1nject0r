package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

var cspPolicies = []string{
	"default-src 'self'",
	"script-src 'self' example.com",
	"style-src 'self' fonts.googleapis.com",
	"font-src 'self' data: fonts.gstatic.com",
	"img-src 'self' data:",
	// Add more CSP policies as needed
}

const banner = `
  ____________ 
 / ___/ __/ _ \
/ /___\ \/ ___/
\___/___/_/                                       
	adrianalvird [ 1.0.0 ]
`

func signalHandler(c chan os.Signal) {
	for {
		sig := <-c
		fmt.Printf("\nExecution stopped by signal: %v\n", sig)
		os.Exit(0)
	}
}

func injectCSPHeader(html string, cspPolicy string) (string, error) {
	// Find the </head> tag index
	headIndex := strings.Index(html, "</head>")
	if headIndex == -1 {
		return "", fmt.Errorf("Error finding </head> tag in HTML")
	}

	// Create a new CSP meta tag and inject it into the HTML
	cspTag := fmt.Sprintf(`<meta http-equiv="Content-Security-Policy" content="%s">`, cspPolicy)
	modifiedHTML := html[:headIndex] + cspTag + html[headIndex:]

	return modifiedHTML, nil
}

func performCSPInjection(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("[-] Error requesting %s: %s\n", url, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("[-] Failed to fetch %s, status code: %d\n", url, response.StatusCode)
		return
	}

	// Read the HTML response
	htmlBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("[-] Error reading HTML response from %s: %s\n", url, err)
		return
	}
	originalHTML := string(htmlBytes)

	for _, cspPolicy := range cspPolicies {
		modifiedHTML, err := injectCSPHeader(originalHTML, cspPolicy)
		if err != nil {
			fmt.Printf("[-] %s\n", err)
			continue
		}

		// Send the modified HTML back to the server
		client := http.Client{}
		req, err := http.NewRequest("POST", url, strings.NewReader(modifiedHTML))
		if err != nil {
			fmt.Printf("[-] Error creating request for %s: %s\n", url, err)
			continue
		}
		req.Header.Set("Content-Type", "text/html; charset=UTF-8")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[-] Error sending modified HTML to %s: %s\n", url, err)
			continue
		}
		defer resp.Body.Close()

		fmt.Printf("[+] CSP Injection successful at %s with policy: %s\n", url, cspPolicy)
	}
}

func main() {
	fmt.Println(banner)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go signalHandler(c)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		targetURL := strings.TrimSpace(scanner.Text())
		performCSPInjection(targetURL)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Error reading from stdin: %s\n", err)
		os.Exit(1)
	}
}

