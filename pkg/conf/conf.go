package conf

import (
	"fmt"
	"os"
	"strings"
)

type LoadErr []error

func (le LoadErr) Error() string {
	msg := []string{}
	for i, err := range le {
		msg = append(msg, fmt.Sprintf("[%d/%d] %s", i+1, len(le), err.Error()))
	}
	return strings.Join(msg, ", ")
}

// Load attempts to get the configuration from the provided loader(s).
// It returns an error if all loaders fail.
func Load(into any, loaders ...Loader) error {
	var loadErr LoadErr
	for _, tryLoad := range loaders {
		err := tryLoad(into)
		if err != nil {
			loadErr = append(loadErr, err)
			continue
		}
		break // success
	}
	if len(loadErr) == len(loaders) {
		return loadErr
	}
	return nil
}

// MustLoad is like Load but panics if all loaders fail.
func MustLoad(into any, loaders ...Loader) {
	if err := Load(into, loaders...); err != nil {
		panic(err)
	}
}

// Loader is used to load configuration from a file or an environment variable for example.
type Loader func(into any) error

func LoadFile(fpath string, decoder Decoder) Loader {
	return func(into any) error {
		fbytes, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		return decoder(fbytes, into)
	}
}

// Decoder is used to decode the data in loaded from a file.
// For example: json.Unmarshal implements this interface.
type Decoder func(raw []byte, into any) error
