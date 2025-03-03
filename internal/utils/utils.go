package utils

import "io"

func ReadFile(r io.Reader) (string, error) {
	line := make([]byte, 0, 1024)

	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		c := buf[0]
		if err != nil {
			if err == io.EOF {
				return string(line), nil
			}
			return "", err
		}
		if n > 0 {
			line = append(line, c)
		}
	}
}
