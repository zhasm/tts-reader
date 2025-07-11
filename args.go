package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	Verbose  bool
	Language string
	Speed    float64 = 0.8
	Content  string
)

func ParseArgs() {
	// Short flags
	flag.BoolVar(&Verbose, "v", false, "verbose mode")
	flag.StringVar(&Language, "l", "fr", "language (fr, jp, pl)")
	flag.Float64Var(&Speed, "s", 0.8, "speed (float)")

	// Support long flags before flag.Parse()
	args := os.Args[1:]
	newArgs := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--verbose":
			Verbose = true
		case "--language":
			if i+1 < len(args) {
				Language = args[i+1]
				i++
			}
		case "--speed":
			if i+1 < len(args) {
				if val, err := strconv.ParseFloat(args[i+1], 64); err == nil {
					Speed = val
					i++
				}
			}
		default:
			newArgs = append(newArgs, args[i])
		}
	}
	os.Args = append([]string{os.Args[0]}, newArgs...) // update os.Args for flag.Parse()
	flag.Parse()
	// Positional argument (content)
	remaining := flag.Args()
	if len(remaining) > 0 {
		Content = remaining[0]
	}

	// Validate language
	if Language != "" && Language != "fr" && Language != "jp" && Language != "pl" {
		fmt.Println("Invalid language. Choose from: fr, jp, pl")
		os.Exit(1)
	}
	if len(Content) == 0 {
		fmt.Println("Usage: [options] <content>")
		fmt.Println("Options:")
		fmt.Println("  -v, --verbose           verbose mode")
		fmt.Println("  -l, --language <lang>   language (fr, jp, pl)")
		fmt.Println("  -s, --speed <float>     speed (default 0.8)")
		os.Exit(1)
	}

	// Print the parsed arguments
	VPrintln("Parsed arguments:")
	VPrintf("  Verbose : %v\n", Verbose)
	VPrintf("  Language: %s\n", Language)
	VPrintf("  Speed   : %v\n", Speed)
	VPrintf("  Content : %s\n", Content)
}
