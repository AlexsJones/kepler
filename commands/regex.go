package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

//MatchInfile ...
func MatchInfile(spath string, filename string, r *regexp.Regexp, dryRun bool) error {

	fullPath := path.Join(spath, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return err
	}
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	data := []byte(r.ReplaceAllString(string(b), ""))

	if dryRun {
		os.Stdout.Write(data)
	} else {
		err = ioutil.WriteFile(fullPath, data, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Wrote to %s\n", fullPath)
	}
	return nil
}
