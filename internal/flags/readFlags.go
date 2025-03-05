package flags

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Command    string
	Filename   string
	MirrorType string
	Filter     string
}

func ReadFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.Command, "command", "", "Command to execute (header, copy)")
	flag.StringVar(&cfg.Filename, "filename", "", "Filename to process")
	flag.StringVar(&cfg.MirrorType, "mirror", "", "Mirroring options (horizontal, vertical)")
	flag.StringVar(&cfg.Filter, "filter", "", "Filter to apply on bitmap")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: go run [options]<filename>\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Validate required flags
	if cfg.Command == "" || cfg.Filename == "" {
		fmt.Fprintln(os.Stderr, "Error: -command and -filename are required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate command
	switch cfg.Command {
	case "header", "copy":
		// Valid commands
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid command")
		flag.Usage()
		os.Exit(1)
	}

	// Validate mirror type
	if cfg.MirrorType != "" {
		switch cfg.MirrorType {
		case "horizontal", "h", "vertical", "v":
			// Valid options
		default:
			fmt.Fprintln(os.Stderr, "Error: invalid mirror type")
			flag.Usage()
			os.Exit(1)
		}
	}

	return cfg
}
