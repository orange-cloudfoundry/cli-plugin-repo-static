package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

	var pluginsOrigin web.PluginsJson
	resp, err := http.Get("https://raw.githubusercontent.com/cloudfoundry/cli-plugin-repo/master/repo-index.yml")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &pluginsOrigin)
	if err != nil {
		return err
	}

	var plugins web.PluginsJson
	b, err = ioutil.ReadFile(cmd.RepoIndexPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &plugins)
	if err != nil {
		return err
	}

	for _, plugin := range plugins.Plugins {
		AddOrReplacePlugin(pluginsOrigin, plugin)
	}

	sort.Sort(pluginsOrigin)

	tmpl, err := template.ParseFiles(filepath.Join("ui", "index.html"))
	if err != nil {
		return err
	}

	indexPage, err := os.Create("./index.html")
	if err != nil {
		return err
	}
	defer indexPage.Close()
	err = tmpl.Execute(indexPage, pluginsOrigin)
	if err != nil { // should only error if template has syntax errors
		return err
	}

	b, err = json.Marshal(pluginsOrigin)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./list", b, 0666)
}

func AddOrReplacePlugin(pluginsOrigin web.PluginsJson, plugin web.Plugin) {
	for i, pluginOrigin := range pluginsOrigin.Plugins {
		if pluginOrigin.Name == plugin.Name {
			pluginsOrigin.Plugins[i] = plugin
			return
		}
	}
	pluginsOrigin.Plugins = append(pluginsOrigin.Plugins, plugin)
}
