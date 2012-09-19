package main

import (
	"fmt"
	"strings"
)

const RedirectTemplate = `<html><head>
<meta http-equiv="refresh" content="0;url=%s">
</head><body></body></html>
`

func WriteRedirect(srcFilename, dstURL string) error {
	n := strings.Count(srcFilename, "/")
	prefix := ""
	for i := 1; i < n; i++ {
		prefix += "../"
	}

	return WriteOutput(
		[]byte(fmt.Sprintf(RedirectTemplate, prefix+dstURL)),
		srcFilename,
	)
}

func WriteRedirectsFor(y, m, d int, t, dstURL string) {
	t = Display2Filename(t)
	prefix := fmt.Sprintf("%s/%s", *outputPath, *blogPath)

	// redirect YYYY-MM-DD.html, YYYY-MM-DD-TTTTT.html, YYYY/MM/DD/index.html
	redirectFiles := []string{
		fmt.Sprintf("%s/%4d-%02d-%02d-%s.%s", prefix, y, m, d, t, *outputExtension),
		fmt.Sprintf("%s/%4d-%02d-%02d.%s", prefix, y, m, d, *outputExtension),
		fmt.Sprintf("%s/%4d/%02d/%02d/index.%s", prefix, y, m, d, *outputExtension),
	}

	// if necessary, redirect YYYY/M/DD/TTTTT.html and YYYY/M/DD/index.html
	if m < 10 {
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%d/%02d/%s.%s", prefix, y, m, d, t, *outputExtension))
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%d/%02d/index.%s", prefix, y, m, d, *outputExtension))
	}

	// if necessary, redirect YYYY/MM/D/TTTTT.html and YYYY/MM/D/index.html
	if d < 10 {
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%02d/%d/%s.%s", prefix, y, m, d, t, *outputExtension))
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%02d/%d/index.%s", prefix, y, m, d, *outputExtension))
	}

	// if necessary, redirect YYYY/M/D/TTTTT.html and YYYY/M/D/index.html
	if m < 10 && d < 10 {
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%d/%d/%s.%s", prefix, y, m, d, t, *outputExtension))
		redirectFiles = append(redirectFiles, fmt.Sprintf("%s/%4d/%d/%d/index.%s", prefix, y, m, d, *outputExtension))
	}

	for _, redirectFile := range redirectFiles {
		if err := WriteRedirect(strings.ToLower(redirectFile), dstURL); err != nil {
			Logf("writing redirect %s: %s", redirectFile, err)
		}
	}
}
