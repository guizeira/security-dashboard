package handlers

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LookupRequest struct {
	Target string `json:"target"`
}

func Index(c *fiber.Ctx) error {
	return c.SendFile("./templates/index.html")
}

// ======================================================
// DNS
// ======================================================

func LookupDNS(c *fiber.Ctx) error {

	req := new(LookupRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	ips, _ := net.LookupIP(req.Target)
	mxRecords, _ := net.LookupMX(req.Target)
	nsRecords, _ := net.LookupNS(req.Target)

	var ipList []string
	var mxList []string
	var nsList []string

	for _, ip := range ips {
		ipList = append(ipList, ip.String())
	}

	for _, mx := range mxRecords {
		mxList = append(mxList, mx.Host)
	}

	for _, ns := range nsRecords {
		nsList = append(nsList, ns.Host)
	}

	return c.JSON(fiber.Map{
		"ips": ipList,
		"mx":  mxList,
		"ns":  nsList,
	})
}

// ======================================================
// WHOIS (CORRIGIDO)
// ======================================================

func LookupWhois(c *fiber.Ctx) error {

	req := new(LookupRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	ips, _ := net.LookupIP(req.Target)

	whoisInfo := fiber.Map{
		"owner":   "N/A",
		"country": "N/A",
		"as":      "N/A",
		"city":    "N/A",
		"region":  "N/A",
		"lat":     0,
		"lon":     0,
	}

	if len(ips) == 0 {
		return c.JSON(whoisInfo)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"http://ip-api.com/json/%s?fields=status,message,org,country,as,city,regionName,lat,lon",
		ips[0].String(),
	)

	resp, err := client.Get(url)
	if err != nil {
		return c.JSON(whoisInfo)
	}
	defer resp.Body.Close()

	var result map[string]interface{}

	if json.NewDecoder(resp.Body).Decode(&result) == nil {

		if status, ok := result["status"].(string); ok && status != "success" {
			return c.JSON(whoisInfo)
		}

		if v, ok := result["org"].(string); ok {
			whoisInfo["owner"] = v
		}

		if v, ok := result["country"].(string); ok {
			whoisInfo["country"] = v
		}

		if v, ok := result["as"].(string); ok {
			whoisInfo["as"] = v
		}

		if v, ok := result["city"].(string); ok {
			whoisInfo["city"] = v
		}

		if v, ok := result["regionName"].(string); ok {
			whoisInfo["region"] = v
		}

		if v, ok := result["lat"].(float64); ok {
			whoisInfo["lat"] = v
		}

		if v, ok := result["lon"].(float64); ok {
			whoisInfo["lon"] = v
		}
	}

	return c.JSON(whoisInfo)
}

// ======================================================
// SSL
// ======================================================

func LookupSSL(c *fiber.Ctx) error {

	req := new(LookupRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	certInfo := fiber.Map{
		"active": false,
	}

	conn, err := tls.DialWithDialer(
		&net.Dialer{
			Timeout: 5 * time.Second,
		},
		"tcp",
		fmt.Sprintf("%s:443", req.Target),
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)

	if err == nil {

		cert := conn.ConnectionState().PeerCertificates[0]

		certInfo = fiber.Map{
			"active": true,
			"issuer": cert.Issuer.CommonName,
			"days_left": int(
				time.Until(cert.NotAfter).Hours() / 24,
			),
		}

		conn.Close()
	}

	return c.JSON(certInfo)
}

// ======================================================
// NMAP STREAM
// ======================================================

func StreamScan(c *fiber.Ctx) error {

	target := strings.TrimSpace(c.Query("target"))

	if target == "" {
		return c.Status(400).SendString("missing target")
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Minute,
		)
		defer cancel()

		cmd := exec.CommandContext(
			ctx,
			"nmap",
			"-Pn",
			"-sV",
			"--version-light",
			"-T4",
			"--stats-every", "1s",
			"--min-rate", "2000",
			"--max-retries", "1",
			"--host-timeout", "2m",
			"--open",
			"-oN", "-",
			"-p",
			"21,22,25,53,80,110,111,123,135,139,143,161,389,443,445,465,587,631,993,995,1025,1433,1521,1723,1883,2049,2082,2083,2086,2087,2095,2096,2181,2375,2376,3000,3001,3306,3389,4443,5432,5601,5672,5900,5985,5986,6379,6443,6667,7001,7070,7443,7777,8000,8008,8080,8081,8088,8443,8888,9000,9090,9091,9200,9300,9443,10000,11211,15672,27017",
			target,
		)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: failed stdout\n\n")
			w.Flush()
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: failed stderr\n\n")
			w.Flush()
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(w, "event: error\ndata: failed start nmap\n\n")
			w.Flush()
			return
		}

		// progress
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.Contains(line, "About") && strings.Contains(line, "%") {

					start := strings.Index(line, "About ")
					end := strings.Index(line, "%")

					if start != -1 && end != -1 && end > start {
						percent := strings.TrimSpace(line[start+6 : end])

						fmt.Fprintf(w, "event: progress\ndata: %s\n\n", percent)
						w.Flush()
					}
				}
			}
		}()

		openPorts := []fiber.Map{}
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if strings.Contains(line, "/tcp") && strings.Contains(line, "open") {

				fields := strings.Fields(line)

				if len(fields) >= 3 {

					portNum := strings.Split(fields[0], "/")[0]
					service := fields[2]

					version := ""
					if len(fields) > 3 {
						version = strings.Join(fields[3:], " ")
					}

					openPorts = append(openPorts, fiber.Map{
						"port":    portNum,
						"status":  "open",
						"service": service,
						"version": version,
					})

					portJson, _ := json.Marshal(fiber.Map{
						"port":    portNum,
						"status":  "open",
						"service": service,
						"version": version,
					})

					fmt.Fprintf(w, "event: port\ndata: %s\n\n", string(portJson))
					w.Flush()
				}
			}
		}

		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(w, "event: error\ndata: scan failed\n\n")
			w.Flush()
			return
		}

		jsonData, _ := json.Marshal(openPorts)

		fmt.Fprintf(w, "event: result\ndata: %s\n\n", string(jsonData))
		fmt.Fprintf(w, "event: progress\ndata: 100\n\n")
		fmt.Fprintf(w, "event: done\ndata: complete\n\n")
		w.Flush()
	})

	return nil
}
