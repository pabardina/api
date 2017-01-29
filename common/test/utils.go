package test

import "io/ioutil"

func ReadContentFileString(filePath string) string {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(b[:])
}
