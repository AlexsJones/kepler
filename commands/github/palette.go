package github

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlexsJones/kepler/commands/storage"
)

//DeleteIssuePalette destroys the currently saved palette
func deleteIssuePalette() {
	storage.GetInstance().Github.CurrentIssue.Palette = make(map[string]string)
	storage.GetInstance().Save()
}

//ShowCurrentIssuePalette of current working issue
func showCurrentIssuePalette() error {
	for k, v := range storage.GetInstance().Github.CurrentIssue.Palette {
		cmd := exec.Command("git", "branch")
		cmd.Dir = v
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		ar := strings.Split(string(out), " ")
		trimmed := strings.TrimSuffix(string(ar[1]), "\n")
		trimmed = strings.TrimPrefix(trimmed, "*")
		trimmed = strings.TrimSpace(trimmed)
		fmt.Println(fmt.Sprintf("Name: %s Branch: %s Path: %s", k, trimmed, v))
	}
	return nil
}

func deleteRepositoryFromPalette(repo string) error {
	if storage.GetInstance().Github.CurrentIssue == nil {
		return errors.New("there is no working issue set; set with github issue set")
	}
	found := false
	for k := range storage.GetInstance().Github.CurrentIssue.Palette {
		if strings.Compare(k, repo) == 0 {
			found = true
			delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
			storage.GetInstance().Save()
		}
	}
	if found != true {
		return fmt.Errorf("there was no repo matching the name %s in the palette\n", repo)
	}
	return nil
}

func addRespositoryToPalette(repo string) error {
	if storage.GetInstance().Github.CurrentIssue == nil {
		return errors.New("there is no working issue set; set with github issue set")
	}
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		return fmt.Errorf("the named repo %s does not exist as a sub directory of the current working directory", repo)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	p := path.Join(dir, repo)
	storage.GetInstance().Github.CurrentIssue.Palette[repo] = p
	storage.GetInstance().Save()
	return nil
}
