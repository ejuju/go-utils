package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Load attempts to get the configuration from the provided loader(s).
// It returns an error if all loaders fail.
func Load(into any, loaders ...TryLoader) error {
	var loadErr loadErr
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
func MustLoad(into any, loaders ...TryLoader) {
	if err := Load(into, loaders...); err != nil {
		panic(err)
	}
}

type loadErr []error

func (le loadErr) Error() string {
	msg := []string{}
	for i, err := range le {
		msg = append(msg, fmt.Sprintf("[%d/%d] %s", i+1, len(le), err.Error()))
	}
	return strings.Join(msg, ", ")
}

// TryLoader is used to load configuration from a file or an environment variable for example.
type TryLoader func(into any) error

func TryLoadFile(fpath string, decoder Decoder) TryLoader {
	return func(into any) error {
		fbytes, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		return decoder(fbytes, into)
	}
}

func TryLoadString(s string, decoder Decoder) TryLoader {
	return func(into any) error { return decoder([]byte(s), into) }
}

// Decoder is used to decode the data in loaded from a file.
// For example: json.Unmarshal implements this interface.
type Decoder func(raw []byte, into any) error

func JSONDecoder(raw []byte, into any) error { return json.Unmarshal(raw, into) }
