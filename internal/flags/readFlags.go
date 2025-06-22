package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Command    string
	Fname      string
	NewFName   string
	MirrorType []string
	FilterType []string
	RotateType []string
	CropParams []string
	Help       bool
}

func ReadFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.Command, "command", "", "Command to execute (header, copy, apply)")
	flag.Var((*stringSliceFlag)(&cfg.MirrorType), "mirror", "Mirroring options (horizontal, vertical). Can be specified multiple times.")
	flag.Var((*stringSliceFlag)(&cfg.FilterType), "filter", "Filter to apply on bitmap")
	flag.Var((*stringSliceFlag)(&cfg.RotateType), "rotate", "Rotation options (right, left, 90, -90, 180). Can be specified multiple times.")
	flag.Var((*stringSliceFlag)(&cfg.CropParams), "crop", "Crop image: xOffset-yOffset-width-height")
	flag.BoolVar(&cfg.Help, "help", false, "Prints usage information")
	flag.BoolVar(&cfg.Help, "helps", false, "Prints usage information")
	flag.BoolVar(&cfg.Help, "h", false, "Prints usage information")

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

	if len(os.Args) < 2 {
		flag.Usage()
		log.Fatal("Error: less arguments than expected.")
	}

	// Extract command type (first argument)
	commandType := os.Args[1]

	// Process the command type
	switch commandType {
	case "apply":
		cfg.Command = "apply"
	case "header":
		cfg.Command = "header"
	case "copy":
		cfg.Command = "copy"
	default:
		log.Fatal("Wrong command, use --help")
	}

	// Remove the first argument (the command) before parsing flags
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)

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

	// This piece of code checks if filename and newfilename was written
	if len(os.Args) < 2 {
		log.Fatal("Error: You have less arguments than expected")
	}
	if cfg.Command != "header" && len(flag.Args()) < 2 {
		log.Fatal("Error: You don't wrote filename or newfilename")
	}
	if cfg.Command != "header" && len(flag.Args()) > 2 {
		log.Fatal("Error: You have more arguments than needed or you should have file names as the last arguments")
	}
	cfg.Fname = flag.Args()[0]
	if len(cfg.Fname) < 5 {
		log.Fatal("Error: too short name of the files")
	} else if cfg.Fname[len(cfg.Fname)-4:] != ".bmp" {
		log.Fatal("Error: not bmp format")
	}

	if cfg.Command != "header" {
		cfg.NewFName = flag.Args()[1]
		if len(cfg.NewFName) < 5 {
			log.Fatal("Error: too short name of the files")
		} else if cfg.NewFName[len(cfg.NewFName)-4:] != ".bmp" {
			log.Fatal("Error: not bmp format")
		}
	}

	// Validate required flags
	if cfg.Command == "" || cfg.Fname == "" {
		fmt.Fprintln(os.Stderr, "Error: -command and -filename are required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate command-specific flags
	switch cfg.Command {
	case "header":
		// No additional flags required
	case "copy", "apply":
		if cfg.NewFName == "" {
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
