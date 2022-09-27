package pkg

import (
	"io"
	"os"
	"regexp"
)

var re = regexp.MustCompile(`^file:\/{0,2}`)

// ReadFile reads the given file into the byte array
func ReadFile(path string) ([]byte, error) {
	// Yes, for now this is just a wrapper around os.ReadFile
	p := re.ReplaceAllString(path, ``)
	return os.ReadFile(p)
}

func WriteFile(path string, data []byte) error {
	// Open a new file for writing only
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0600,
	)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func AppendToFile(source string, to string) error {
	// Open a new file for writing only
	file, err := os.OpenFile(
		source,
		os.O_RDONLY,
		0600,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// open the other file...
	file2, err := os.OpenFile(
		to,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0600,
	)

	if err != nil {
		return err
	}
	defer file2.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}
		file2.Write(buf[:n])
	}

	return err
}

func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
