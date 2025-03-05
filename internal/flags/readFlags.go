package flags

import "flag"

var (
	filename   string
	mirrorType []string
	filter     string
	rotates    []string
	crop       string
)

func ReadFlags() {
	flag.StringVar(&filename, "header", "", "Filename to show header information")
	flag.StringVar(&mirrorType, "mirror", "", "Mirroring options")
	flag.StringVar(&filter, "filter", "", "Filter to apply on bitmap")
	flag.Func() // for slices
}
