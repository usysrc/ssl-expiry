package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

// NetworkConnector is an interface to abstract network connections.
type NetworkConnector interface {
	Dial(network, address string) (net.Conn, error)
}

// RealNetworkConnector implements NetworkConnector using the real net.Dialer.
type RealNetworkConnector struct{}

func (rnc RealNetworkConnector) Dial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return tls.DialWithDialer(dialer, network, address, tlsConfig)
}

func main() {
	var targetURL string

	// Check if there is any data available from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read target URL from stdin (pipe)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			targetURL = scanner.Text()
		} else {
			fmt.Println("Error reading input.")
			return
		}
	} else if len(os.Args) < 2 {
		fmt.Println("Usage: ssl-expiry <targetURL>")
		return
	} else {
		targetURL = os.Args[1]
	}

	formattedURL := FormatURL(targetURL)
	conn, err := DialNetwork(formattedURL, RealNetworkConnector{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	expiryDate, err := GetCertificateExpiryDate(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if expiryDate.Before(time.Now()) {
		fmt.Println("\033[31mWarning: The certificate has already expired!\033[0m")
	}

	remainingDays := CalculateRemainingDays(expiryDate)

	fmt.Printf("Certificate Expiry Date: %s\n", expiryDate)
	fmt.Printf("Remaining Days: %d\n", remainingDays)
}

// FormatURL formats the given URL to the desired format.
// If the URL starts with "https://", this function removes the prefix,
// trims any trailing "/" character, and appends ":443" to indicate the
// default HTTPS port. The formatted URL is returned as the result.
// for example https://example.com/ is going to return example.com:443
func FormatURL(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(parsedURL)
	}
	if parsedURL.Scheme == "https" {
		targetURL = parsedURL.Host + ":443"
	}
	return targetURL
}

// DialNetwork uses the provided NetworkConnector to establish a network connection.
func DialNetwork(targetURL string, connector NetworkConnector) (net.Conn, error) {
	return connector.Dial("tcp", targetURL)
}

// GetCertificateExpiryDate retrieves the expiry date of the peer certificate.
func GetCertificateExpiryDate(conn net.Conn) (time.Time, error) {
	certChain := conn.(*tls.Conn).ConnectionState().PeerCertificates

	if len(certChain) == 0 {
		return time.Time{}, fmt.Errorf("no certificate found")
	}

	return certChain[0].NotAfter, nil
}

// CalculateRemainingDays calculates the remaining days until the given date.
func CalculateRemainingDays(expiryDate time.Time) int {
	return int(expiryDate.Sub(time.Now()).Hours() / 24)
}
