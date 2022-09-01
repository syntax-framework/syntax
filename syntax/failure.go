package syntax

import (
	"bytes"
	"log"
	"strings"
)

const failedToStartTitle = "application failed to start"

func fail(title string, description string, action string) {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	buf.WriteString(strings.Repeat("*", len(title)) + "\n")
	buf.WriteString(strings.ToUpper(title) + "\n")
	buf.WriteString(strings.Repeat("*", len(title)) + "\n\n")
	buf.WriteString("Description: " + description + "\n")
	buf.WriteString("Action:      " + action)
	log.Fatal(buf.String())
}

func failToStart(description string, action string) {
	fail(failedToStartTitle, description, action)
}
