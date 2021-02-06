package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/orange-cloudfoundry/cli-plugin-repo-static/web"
	yaml "gopkg.in/yaml.v2"
)

type CLIPR struct {
	RepoIndexPath string `short:"f" long:"filepath" default:"repo-index.yml" description:"Path to repo-index file"`
}

func (cmd *CLIPR) Execute(args []string) error {

	var plugins web.PluginsJson

	b, err := ioutil.ReadFile(cmd.RepoIndexPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &plugins)
	if err != nil {
		return err
	}

	sort.Sort(plugins)

	tmpl, err := template.ParseFiles(filepath.Join("ui", "index.html"))
	if err != nil {
		return err
	}

	indexPage, err := os.Create("./index.html")
	if err != nil {
		return err
	}
	defer indexPage.Close()
	err = tmpl.Execute(indexPage, plugins)
	if err != nil { // should only error if template has syntax errors
		return err
	}

	b, err = json.Marshal(plugins)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./list", b, 0666)
}
