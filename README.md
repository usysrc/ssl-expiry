# SSL Expiry Checker

The SSL Expiry Checker is a simple command-line tool written in Go that checks the expiry date of SSL certificates for a given target URL. It provides information about the remaining days until the SSL certificate expires.

## Installation
Needs golang 1.20+ installed.

To use the ssl-expiry tool, you need to have Go installed on your system. If you haven't installed Go, you can download it from the official Go website: https://golang.org/dl/

Once you have Go installed, you can install the ssl-expiry tool using the following command:

```bash
go install github.com/usysrc/ssl-expiry@latest
```

## Usage

Run the SSL Expiry Checker with a target URL:

```sh
ssl-expiry <targetURL>
```

Replace `<targetURL>` with the URL of the website whose SSL certificate you want to check. The tool will display the certificate expiry date and the remaining days until expiry.

Example:

```sh
ssl-expiry https://www.example.com
```

Or without the schema and with a port:

```
ssl-expiry www.example.com:443
```

## Testing

To run the tests for the project, use the following command:

```sh
go test
