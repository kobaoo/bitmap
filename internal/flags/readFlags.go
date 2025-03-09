package flags

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Command     string
	Filename    string
	NewFileName string
	MirrorType  []string // Slice to store multiple mirroring options
	FilterType  []string
	RotateType  []string
	CropParams  []string
	Help        bool
}

func ReadFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.Command, "command", "", "Command to execute (header, copy, apply)")
	flag.StringVar(&cfg.Filename, "filename", "", "Source file to process")
	flag.StringVar(&cfg.NewFileName, "newfilename", "", "Output file for processed image")
	flag.Var((*stringSliceFlag)(&cfg.MirrorType), "mirror", "Mirroring options (horizontal, vertical). Can be specified multiple times.")
	flag.Var((*stringSliceFlag)(&cfg.FilterType), "filter", "Filter to apply on bitmap")
	flag.Var((*stringSliceFlag)(&cfg.RotateType), "rotate", "Rotation options (right, left, 90, -90, 180). Can be specified multiple times.")
	flag.Var((*stringSliceFlag)(&cfg.CropParams), "crop", "Crop image: xOffset-yOffset-width-height")
	flag.BoolVar(&cfg.Help, "help", false, "Prints usage information")

	// Default usage (for ./bitmap)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  bitmap <command> [arguments]\n\n")
		fmt.Fprintf(os.Stderr, "The commands are:\n")
		fmt.Fprintf(os.Stderr, "  header    prints bitmap file header information\n")
		fmt.Fprintf(os.Stderr, "  copy      copies the image to a new file\n")
		fmt.Fprintf(os.Stderr, "  apply     applies processing to the image and saves it to the file\n\n")
		fmt.Fprintf(os.Stderr, "Use 'bitmap <command> --help' for more information about a command.\n")
	}

	flag.Parse()

	// If --help is passed, show command-specific usage
	if cfg.Help {
		switch cfg.Command {
		case "header":
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "  bitmap header <source_file>\n\n")
			fmt.Fprintf(os.Stderr, "Description:\n")
			fmt.Fprintf(os.Stderr, "  Prints bitmap file header information\n")
			os.Exit(0)
		case "copy":
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "  bitmap copy <source_file> <output_file>\n\n")
			fmt.Fprintf(os.Stderr, "Description:\n")
			fmt.Fprintf(os.Stderr, "  Copies the image to a new file\n")
			os.Exit(0)
		case "apply":
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fmt.Fprintf(os.Stderr, "  bitmap apply [options] <source_file> <output_file>\n\n")
			fmt.Fprintf(os.Stderr, "The options are:\n")
			fmt.Fprintf(os.Stderr, "  --mirror        mirroring options (horizontal, vertical). Can be specified multiple times.\n")
			fmt.Fprintf(os.Stderr, "  --filter        filter to apply on bitmap\n")
			os.Exit(0)
		default:
			flag.Usage()
			os.Exit(0)
		}
	}

	// Validate required flags
	if cfg.Command == "" || cfg.Filename == "" {
		fmt.Fprintln(os.Stderr, "Error: -command and -filename are required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate command-specific flags
	switch cfg.Command {
	case "header":
		// No additional flags required
	case "copy", "apply":
		if cfg.NewFileName == "" {
			fmt.Fprintln(os.Stderr, "Error: -newfilename is required for this command")
			flag.Usage()
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid command")
		flag.Usage()
		os.Exit(1)
	}

	return cfg
}

// stringSliceFlag is a custom flag type to support multiple values for the same flag.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
