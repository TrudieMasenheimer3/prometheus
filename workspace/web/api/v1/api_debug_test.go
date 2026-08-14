package v1

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDebugPrint(t *testing.T) {
	content, err := os.ReadFile("api.go")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("=== api.go ===")
	fmt.Println(string(content))

	f, err := os.Open("api_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lineNum := 0
	fmt.Println("=== api_test.go ===")
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(line, "limit") || strings.Contains(line, "Limit") || strings.Contains(line, "truncated") {
			fmt.Printf("%d: %s\n", lineNum, line)
		}
	}
	t.Fatal("force failure to see output")
}
