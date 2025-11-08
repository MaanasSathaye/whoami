package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MaanasSathaye/whoami/model"
	"github.com/MaanasSathaye/whoami/parser"
	"github.com/MaanasSathaye/whoami/renderer"
)

func main() {
	var err error
	var maxPages int
	var content []byte
	var resume *model.Resume

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output.pdf> <input.md> [max_pages]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Example: %s resume.pdf resume.md 1\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  max_pages: 1 (compact), 2 (standard), 3+ (spacious). Default: 2\n")
		os.Exit(1)
	}

	outputPath := os.Args[1]
	inputPath := os.Args[2]
	maxPages = 2

	if len(os.Args) >= 4 {
		if maxPages, err = strconv.Atoi(os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: max_pages must be a number\n")
			os.Exit(1)
		}
		if maxPages < 1 {
			fmt.Fprintf(os.Stderr, "Error: max_pages must be at least 1\n")
			os.Exit(1)
		}
	}

	if !strings.HasSuffix(outputPath, ".pdf") {
		outputPath += ".pdf"
	}

	if content, err = os.ReadFile(inputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	if resume, err = parser.ParseMarkdown(content); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown: %v\n", err)
		os.Exit(1)
	}

	if err = renderer.RenderPDF(resume, outputPath, maxPages); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Resume generated successfully: %s (max_pages: %d)\n", outputPath, maxPages)
}
