package utils

import "fmt"

// ComposeSignString combines HTTP method, path, access token, body hash, and timestamp to create a sign string
func ComposeSignString(method, path, accessToken, bodyHash, timestamp string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", method, path, accessToken, bodyHash, timestamp)
}
