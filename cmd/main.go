package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/loicalleyne/blobmunge"
)

// Reads input data and a base64-encoded bloblang mapping from Stdin
// (must be separated by '}[)(]') and prints the result to Stdout.
// In the event of errors will output to StdErr and exit with a non-0 code.
//
// Example
// input:
// {"id":1234,"dev":"12345ert"}[^^]cm9vdC5pZCA9IHRoaXMuaWQKcm9vdC5kZXZpY2UgPSB0aGlzLmRldg==\n
// outputs:
// {"device":"12345ert","id":1234}

const (
	_                    = iota
	InputScanError       = iota
	MissingSeparator     = iota
	ProblematicInputData = iota
	EmptySourceData      = iota
	EmptyMungeMapping    = iota
	MappingDecodingError = iota
	MappingParseErro     = iota
	MungeApplyError      = iota
)

func main() {
	l := log.New(os.Stderr, "", log.LstdFlags)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		l.Printf("input scan error %v\n", err)
		os.Exit(InputScanError)
	}
	if !strings.Contains(text, "[^^]") {
		l.Printf("missing separator\n")
		os.Exit(MissingSeparator)
	}
	lines := strings.Split(text, "[^^]")
	if len(lines) != 2 {
		l.Printf("problematic input data\n")
		os.Exit(ProblematicInputData)
	}
	if len(lines[0]) == 0 {
		l.Printf("empty source data\n")
		os.Exit(EmptySourceData)
	}
	if len(lines[1]) == 0 {
		l.Printf("empty munge mapping\n")
		os.Exit(EmptyMungeMapping)
	}
	decoded, err := base64.StdEncoding.DecodeString(lines[1])
	if err != nil {
		l.Printf("error decoding munge mapping %v\n", err)
		os.Exit(MappingDecodingError)
	}
	b, err := blobmunge.New(string(decoded))
	if err != nil {
		l.Printf("error parsing munge mapping %v\n", err)
		os.Exit(MappingParseErro)
	}
	munged, err := b.ApplyBloblangMapping(lines[0])
	if err != nil {
		l.Printf("error applying mapping %v\n", err)
		os.Exit(MungeApplyError)
	}
	fmt.Println(string(munged))
}
