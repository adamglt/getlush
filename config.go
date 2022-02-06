package main

import (
	"errors"
	"flag"
	"strings"
	"time"
)

type Config struct {
	url string // hilan's base url
	org string // parent org id
	emp string // running employee id

	from hilanDate // fetch start
	to   hilanDate // fetch end

	cookiePath string // path to cookie file
	outputPath string // path to output dir

	timeout time.Duration
}

func (cfg *Config) bindFlags(fs *flag.FlagSet) {
	const (
		traianaURL   = "https://traiana.net.hilan.co.il/"
		traianaOrgID = "9133"
	)
	fs.StringVar(&cfg.url, "url", traianaURL, "Hilan's base URL")
	fs.StringVar(&cfg.org, "org", traianaOrgID, "Parent organization ID")
	fs.StringVar(&cfg.emp, "emp", "", "Employee ID [required]")
	fs.Var(&cfg.from, "from", "First payslip to fetch (YYYY-MM) [required]")
	fs.Var(&cfg.to, "to", "Last payslip to fetch (YYYY-MM) [required]")
	fs.StringVar(&cfg.cookiePath, "cookie", "hilan.cookie", "Path to cookie file")
	fs.StringVar(&cfg.outputPath, "out", "getlush_out/", "Directory path for fetched pdfs")
	fs.DurationVar(&cfg.timeout, "t", 10*time.Second, "Single request timeout")
}

func (cfg Config) validate() error {
	var msgs []string
	if cfg.emp == "" {
		msgs = append(msgs, "employee id cannot be empty")
	}
	if cfg.from.IsZero() || cfg.to.IsZero() {
		msgs = append(msgs, "a date range must be specified using 'from' and 'to'")
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, "\n"))
	}
	return nil
}

const (
	// for some reason we need both
	hilanDateFmtYYYYMM   = "2006-01"
	hilanDateFmtDDMMYYYY = "02/01/2006"
)

// flag parsing for yyyy-mm fmt
type hilanDate struct {
	time.Time
}

func (h hilanDate) String() string {
	return h.Format(hilanDateFmtYYYYMM)
}

func (h *hilanDate) Set(s string) error {
	t, err := time.Parse(hilanDateFmtYYYYMM, s)
	if err != nil {
		return err
	}
	h.Time = t
	return nil
}
