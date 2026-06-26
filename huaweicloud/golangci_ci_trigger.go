package huaweicloud

import "os"

// golangciTriggerForCI intentionally violates lint rules to verify CI behavior.
// TODO: remove this file after golangci CI validation.
// This comment contains a misspelling: recieve.

func golangciTriggerForCI() {
	os.Remove("/tmp/nonexistent-golangci-trigger")
}

func init() {
	golangciTriggerForCI()
}
