package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// builds a request for a specific month.
type requestFn func(d time.Time) (*http.Request, error)

// common config for all payslip request fns.
func requestFnPayslip(cfg Config, cookie string) requestFn {
	// make payslip base url
	payslipURL := cfg.url
	if !strings.HasSuffix(payslipURL, "/") {
		payslipURL += "/"
	}
	payslipURL += "Hilannetv2/PersonalFile/PdfPaySlip.aspx/"

	return func(d time.Time) (*http.Request, error) {
		// full path looks like:
		// /Hilannetv2/PersonalFile/PdfPaySlip.aspx/PaySlip2020-01.pdf?Date=01/01/2020&UserId=123123123
		fileseg := fmt.Sprintf("PaySlip%s.pdf?Date=%s&UserId=%s%s",
			d.Format(hilanDateFmtYYYYMM), d.Format(hilanDateFmtDDMMYYYY), cfg.org, cfg.emp)
		url := payslipURL + fileseg

		// build req and set cookie as-is
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("new_request: %w", err)
		}
		req.Header.Set("Cookie", cookie)
		return req, nil
	}
}

// common config for all 106 request fns.
func requestFn106(cfg Config, cookie string) requestFn {
	// make payslip base url
	f106URL := cfg.url
	if !strings.HasSuffix(f106URL, "/") {
		f106URL += "/"
	}
	f106URL += "Hilannetv2/PersonalFile/Pdf106.aspx/"

	return func(d time.Time) (*http.Request, error) {
		// full path looks like:
		// /Hilannetv2/PersonalFile/Pdf106.aspx/Pdf106-2020.pdf?UserId=123123123&Year=2020
		yyyy := d.Format(hilanDateFmtYYYY)
		fileseg := fmt.Sprintf("Pdf106-%s.pdf?Year=%s&UserId=%s%s",
			yyyy, yyyy, cfg.org, cfg.emp)
		url := f106URL + fileseg

		// build req and set cookie as-is
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("new_request: %w", err)
		}
		req.Header.Set("Cookie", cookie)
		return req, nil
	}
}

// first few common bytes for all pdf versions.
var pdfPrefix = []byte("%PDF-")

// runs the request with some error checks.
// hilan returns 200 for anything so no status checks.
func doRequest(c *http.Client, req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client_do: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// read entire body, should be <100k
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read_body: %w", err)
	}

	// peek
	if !bytes.HasPrefix(data[:8], pdfPrefix) {
		return nil, fmt.Errorf("response was not a pdf")
	}

	return data, nil
}
