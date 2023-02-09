package spotify

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// ErrTokenNotFound is returned when we couldn't find the token at the default or set location
	ErrTokenNotFound = errors.New("token not found. Specify the correct location of the token")
)

type tokenUtils struct {
	location string
}

func (t *tokenUtils) get() ([]byte, error) {
	tokenFileLocation := filepath.Join(t.location, "token")
	_, err := os.Stat(t.location)
	// check if the directory exists
	if os.IsNotExist(err) {
		return nil, ErrTokenNotFound
	}

	b, err := os.ReadFile(tokenFileLocation)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *tokenUtils) create(value []byte) error {
	_, err := os.Stat(t.location)
	// check if the directory exists
	if os.IsNotExist(err) {
		os.Mkdir(t.location, 0644)
	}

	tokenFileLocation := filepath.Join(t.location, "token")
	_, err = os.Stat(tokenFileLocation)
	if os.IsNotExist(err) {
		//create the file and write to the file location
		err = os.WriteFile(
			tokenFileLocation,
			value,
			fs.FileMode(os.O_RDWR|os.O_CREATE),
		)

		if err != nil {
			return err
		}
	} else {
		// the file exists, so we over-write it
		f, err := os.Open(tokenFileLocation)
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write(value)
	}

	return nil
}

func (t *tokenUtils) valid(token []byte) bool {
	return false
}
