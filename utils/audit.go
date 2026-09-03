package utils

import "context"

func AuditDirectory(ctx context.Context, dir string, opts ...ScanOptions) ([]Finding, error) {
	return ScanDirectoryWithBetterleaks(ctx, dir, nil, opts...)
}
