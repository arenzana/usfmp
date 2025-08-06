package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arenzana/usfmp/internal/formatter"
	"github.com/arenzana/usfmp/pkg/usfm"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFormat string
	outputFile   string
	verbose      bool
	quiet        bool
	strict       bool

	// Version information
	buildVersion = "dev"
	buildCommit  = "unknown"
	buildDate    = "unknown"
	buildBy      = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "usfmp [input-file-or-directory]",
	Short: "A USFM (Unified Standard Format Marker) parser for biblical texts",
	Long: `usfmp is a command-line tool for parsing USFM (Unified Standard Format Marker) files.
It can process single files or entire directories of USFM files and output them in various formats.

USFM is a markup format used for biblical texts. More information: https://docs.usfm.bible/usfm/3.1/index.html`,
	Args: cobra.ExactArgs(1),
	RunE: run,
}

// SetVersionInfo sets version information from build-time variables
func SetVersionInfo(version, commit, date, builtBy string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date
	buildBy = builtBy

	// Set version for cobra
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s, by: %s)",
		buildVersion, buildCommit, buildDate, buildBy)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Output format flag
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "json",
		"Output format: json, txt, tsv, pdf")

	// Output file flag
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "",
		"Output file (default: stdout)")

	// Verbosity flags
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose output")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false,
		"Quiet mode - minimal output")

	// Parsing options
	rootCmd.Flags().BoolVar(&strict, "strict", false,
		"Strict mode - fail on unknown markers")
}

// run is the main command execution function
func run(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Validate flags
	if err := validateFlags(); err != nil {
		return err
	}

	// Configure parser options
	parseOptions := usfm.ParseOptions{
		StrictMode:        strict,
		IncludeFootnotes:  true,
		IncludeReferences: true,
	}

	parser := usfm.NewParser(parseOptions)

	// Check if input is file or directory
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("cannot access input path: %w", err)
	}

	var files []string
	if info.IsDir() {
		// Process directory
		files, err = findUSFMFiles(inputPath)
		if err != nil {
			return fmt.Errorf("error finding USFM files: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no USFM files found in directory: %s", inputPath)
		}

		logInfo("Found %d USFM files", len(files))
	} else {
		// Single file
		files = []string{inputPath}
	}

	// Parse files
	var documents []*usfm.Document
	for _, file := range files {
		logInfo("Parsing file: %s", file)

		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", file, err)
		}

		doc, err := parser.Parse(f, file)
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file %s: %w", file, closeErr)
		}

		if err != nil {
			return fmt.Errorf("error parsing file %s: %w", file, err)
		}

		documents = append(documents, doc)
		logInfo("Successfully parsed %s - %s", doc.ID, doc.MainTitle)
	}

	// Format and output results
	return outputResults(documents)
}

// validateFlags checks that flag combinations are valid
func validateFlags() error {
	if quiet && verbose {
		return fmt.Errorf("cannot use both --quiet and --verbose flags")
	}

	validFormats := []string{"json", "txt", "tsv", "pdf"}
	for _, format := range validFormats {
		if outputFormat == format {
			return nil
		}
	}

	return fmt.Errorf("invalid output format: %s (valid: %s)",
		outputFormat, strings.Join(validFormats, ", "))
}

// findUSFMFiles recursively finds USFM files in a directory
func findUSFMFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isUSFMFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// isUSFMFile checks if a file is likely a USFM file based on extension
func isUSFMFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".sfm" || ext == ".usfm"
}

// outputResults formats and outputs the parsed documents
func outputResults(documents []*usfm.Document) error {
	var output string
	var err error

	switch outputFormat {
	case "json":
		output, err = formatter.FormatJSON(documents)
	case "txt":
		output, err = formatter.FormatText(documents)
	case "tsv":
		output, err = formatter.FormatTSV(documents)
	case "pdf":
		return fmt.Errorf("PDF output not yet implemented")
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	// Write output
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		logInfo("Output written to: %s", outputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

// logInfo prints informational messages unless in quiet mode
func logInfo(format string, args ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, "[INFO] "+format+"\n", args...)
	}
}
