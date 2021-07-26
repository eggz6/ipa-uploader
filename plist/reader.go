package plist

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"howett.net/plist"
)

type PList map[string]interface{}

func (p PList) GetString(key string) (string, bool) {
	val, ok := p[key]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", val), true
}

func ReadFromFile(path string) (PList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	res := make(PList)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	_, err = plist.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ReadFrom(r io.Reader) (PList, error) {
	res := make(PList)

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	_, err = plist.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
