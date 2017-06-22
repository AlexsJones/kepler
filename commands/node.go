package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//PackageJSON structure of package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Author          string            `json:"author"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func recursePackages(p *PackageJSON, callback func(key string, value string)) error {

	for key, value := range p.Dependencies {

		callback(key, value)
	}
	return nil
}

//HasPackage searches for packages references
func HasPackage(subPath string, filename string, target string) (bool, error) {

	filepath := path.Join(subPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, err
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, err
	}
	var packagejson PackageJSON
	json.Unmarshal(b, &packagejson)

	var wasFound = false
	recursePackages(&packagejson, func(key string, value string) {
		if strings.Contains(key, target) || strings.Contains(value, target) {
			wasFound = true
			return
		}
	})
	return wasFound, nil
}

//FixLinks ...
func FixLinks(subPath string, filename string, prefix string, target string, shouldDelete bool) error {

	filepath := path.Join(subPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	var packagejson PackageJSON
	json.Unmarshal(b, &packagejson)

	//processing
	recursePackages(&packagejson, func(key string, value string) {
		if strings.Contains(value, target) {
			if shouldDelete {
				delete(packagejson.Dependencies, key)
			} else {
				spli := strings.Split(value, "/")
				subspli := strings.Split(spli[len(spli[1:])], ".")
				foundEntry := subspli[0]
				foundEntry = strings.TrimSuffix(foundEntry, "\"")
				syntax := "file:%s%s"
				value := fmt.Sprintf(syntax, prefix, foundEntry)
				packagejson.Dependencies[key] = value
			}
		}
	})
	o, err := json.MarshalIndent(packagejson, "", "    ")
	if err != nil {
		return err
	}
	o = append(o, []byte("\n")...)

	return ioutil.WriteFile(filepath, o, 0644)

}
