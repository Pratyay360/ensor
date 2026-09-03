package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/betterleaks/betterleaks/config"
	"github.com/betterleaks/betterleaks/detect"
	"github.com/betterleaks/betterleaks/report"
	"github.com/betterleaks/betterleaks/sources"
)

type Finding struct {
	RuleID      string
	Description string
	File        string
	StartLine   int
	Secret      string
	// ValidationStatus is the liveness result from betterleaks live
	// validation ("" when validation was not run). Values: valid, invalid,
	// revoked, unknown, error, needs_validation.
	ValidationStatus string
	ValidationReason string
}

// ScanOptions controls how a scan is performed.
type ScanOptions struct {
	// Validate enables live secret validation against the corresponding
	// providers (AWS, GitHub, Stripe, ...). Requires network access.
	Validate bool
	// SkipInvalid drops findings the provider confirmed as invalid, which are
	// almost always false positives. Only applied when Validate is set.
	SkipInvalid bool
	// ValidationTimeout is the per-request timeout for validation requests.
	// Zero uses the default of 10s.
	ValidationTimeout time.Duration
	// ValidationWorkers is the number of concurrent validation workers.
	// Zero uses the default of 10.
	ValidationWorkers int
	// ValidationMaxRequestsPerTarget caps validation requests per provider
	// target. Zero means unlimited.
	ValidationMaxRequestsPerTarget int
	// ValidationRequestsPerSecond globally rate-limits validation requests.
	// Zero means unlimited.
	ValidationRequestsPerSecond float64
}

// ScanDirectoryWithBetterleaks recursively scans dir for secrets using the
// betterleaks detector and returns the findings. 
func ScanDirectoryWithBetterleaks(ctx context.Context, dir string, progress func(string), opts ...ScanOptions) ([]Finding, error) {
	var opt ScanOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var detector *detect.Detector
	if opt.Validate {
		valOpts := detect.ValidationOptions{
			Enabled:              true,
			Workers:              opt.ValidationWorkers,
			Timeout:              opt.ValidationTimeout,
			MaxRequestsPerTarget: opt.ValidationMaxRequestsPerTarget,
			RequestsPerSecond:    opt.ValidationRequestsPerSecond,
		}
		if valOpts.Workers <= 0 {
			valOpts.Workers = 10
		}
		if valOpts.Timeout <= 0 {
			valOpts.Timeout = 10 * time.Second
		}
		cfg, err := config.Default()
		if err != nil {
			return nil, fmt.Errorf("creating default config: %w", err)
		}
		detector = detect.NewDetectorContext(ctx, cfg, valOpts)
	} else {
		var err error
		detector, err = detect.NewDetectorDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("creating detector: %w", err)
		}
	}
	detector.MaxDecodeDepth = 10
	detectorSkip := detector.SkipFunc()
	skip := func(attrs map[string]string) bool {
		if detectorSkip != nil && detectorSkip(attrs) {
			return true
		}
		if path := attrs[sources.AttrPath]; path != "" && IgnoreFile(path, dir) {
			return true
		}
		return false
	}

	src := &sources.Files{
		Path:            dir,
		Sema:            detector.Sema,
		ShouldSkip:      skip,
		MaxFileSize:     100 * 1000 * 1000, 
		MaxArchiveDepth: 2,
	}

	var findings []Finding
	for result := range detector.Run(ctx, src) {
		if result.Err != nil {
			continue
		}
		f := result.Finding
		if f.Secret == "" {
			continue
		}
		if opt.Validate && opt.SkipInvalid && f.ValidationStatus == report.ValidationStatusInvalid {
			continue
		}
		if progress != nil {
			progress(f.File)
		}
		findings = append(findings, Finding{
			RuleID:           f.RuleID,
			Description:      f.Description,
			File:             f.File,
			StartLine:        f.StartLine,
			Secret:           f.Secret,
			ValidationStatus: string(f.ValidationStatus),
			ValidationReason: f.ValidationReason,
		})
	}

	return findings, nil
}
