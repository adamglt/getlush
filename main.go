package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// setup
	log.SetFlags(log.Ldate | log.Ltime)

	// parse flags
	var cfg Config
	cfg.bindFlags(flag.CommandLine)
	flag.Parse()
	if err := cfg.validate(); err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// init i/o
	cookie, err := readCookie(cfg.cookiePath)
	if err != nil {
		fmt.Printf("could not read cookie file: %v", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.outputPath, os.ModePerm); err != nil {
		fmt.Printf("could not verify output directory: %v", err)
		os.Exit(1)
	}

	// common http stuff
	client := &http.Client{Timeout: cfg.timeout}
	requestFn := makeRequestFn(cfg, cookie)

	// start running for each month
	for cur := cfg.from.Time; cfg.to.After(cur); cur = cur.AddDate(0, 1, 0) {
		id := cur.Format(hilanDateFmtYYYYMM) // request/file id
		req, err := requestFn(cur)
		if err != nil {
			log.Printf("ERROR: failed to create request %s: %v", id, err)
			continue
		}
		if cfg.verbose {
			log.Println("DEBUG: getting", req.URL.String())
		}
		data, err := doRequest(client, req)
		if err != nil {
			log.Printf("ERROR: failed to fetch payslip %s: %v", id, err)
			continue
		}
		file := filepath.Join(cfg.outputPath, fmt.Sprintf("payslip_%s.pdf", id))
		if cfg.verbose {
			log.Printf("DEBUG: writing %s [~%dK]", file, len(data)/1e3)
		}
		if err := ioutil.WriteFile(file, data, os.ModePerm); err != nil {
			log.Printf("ERROR: failed to write payslip %s: %v", id, err)
			continue
		}
	}
}

func readCookie(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(data)
	// some browsers export "key: value"
	const copyPrefix = "Cookie: "
	if strings.HasPrefix(s, copyPrefix) {
		s = strings.TrimPrefix(s, copyPrefix)
	}
	return s, nil
}
