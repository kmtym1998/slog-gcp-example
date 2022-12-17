package middleware

import "regexp"

func ExtractTraceID(raw []byte) (traceID string, ok bool) {
	if (len(raw)) == 0 {
		return "", false
	}

	// X-Cloud-Trace-Context の値の例: 105445aa7843bc8bf206b12000100000/1;o=1
	// SEE: https://cloud.google.com/trace/docs/setup?hl=ja#force-trace
	matches := regexp.MustCompile(`([a-f\d]+)/([a-f\d]+)`).FindAllSubmatch(raw, -1)
	if len(matches) != 1 {
		return "", false
	}

	sub := matches[0]
	if len(sub) != 3 {
		return "", false
	}

	traceID = string(sub[1])
	ok = true

	return
}
