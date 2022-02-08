package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
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
		fmt.Println("could not read cookie file:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.outputPath, os.ModePerm); err != nil {
		fmt.Println("could not verify output directory:", err)
		os.Exit(1)
	}

	// common http stuff
	client := &http.Client{Timeout: cfg.timeout}

	// payslips for each month
	payslipFn := requestFnPayslip(cfg, cookie)
	for cur := cfg.from.Time; cfg.to.After(cur); cur = cur.AddDate(0, 1, 0) {
		id := cur.Format(hilanDateFmtYYYYMM) // request/file id
		fmt.Printf(">> %s... ", id)
		req, err := payslipFn(cur)
		if err != nil {
			fmt.Printf("× [failed to create request: %v]\n", err)
			continue
		}
		data, err := doRequest(client, req)
		if err != nil {
			fmt.Printf("× [failed to fetch payslip: %v]\n", err)
			continue
		}
		file := filepath.Join(cfg.outputPath, fmt.Sprintf("payslip_%s.pdf", id))
		if err := ioutil.WriteFile(file, data, os.ModePerm); err != nil {
			fmt.Printf("× [failed to write payslip: %v]\n", err)
			continue
		}
		fmt.Printf("✓ [~%dK]\n", len(data)/1e3)
	}

	// 106 forms for each year
	f106Fn := requestFn106(cfg, cookie)
	for cur := cfg.from.Time; cfg.to.After(cur); cur = cur.AddDate(1, 0, 0) {
		id := cur.Format(hilanDateFmtYYYY) // request/file id
		fmt.Printf(">> 106-%s...", id)
		req, err := f106Fn(cur)
		if err != nil {
			fmt.Printf("× [failed to create request: %v]\n", err)
			continue
		}
		data, err := doRequest(client, req)
		if err != nil {
			fmt.Printf("× [failed to fetch 106 form: %v]\n", err)
			continue
		}
		file := filepath.Join(cfg.outputPath, fmt.Sprintf("form106_%s.pdf", id))
		if err := ioutil.WriteFile(file, data, os.ModePerm); err != nil {
			fmt.Printf("× [failed to write 106 form: %v]\n", err)
			continue
		}
		fmt.Printf("✓ [~%dK]\n", len(data)/1e3)
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
