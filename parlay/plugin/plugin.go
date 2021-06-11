package parlayplugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"plunder-app/plunder/parlay/parlaytypes"

	log "github.com/sirupsen/logrus"
)

// The pluginCache contains a map of action->plugin
var pluginCache map[string]string

func init() {
	// Initialise the map
	pluginCache = make(map[string]string)
}

// Find plugins returns an array of all .plugin files
func findPlugins(pluginDir string) ([]string, error) {
	var plugins []string
	// This function will look for all files in a specified directory (defaults to PWD/plugin)
	filepath.Walk(pluginDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			if filepath.Ext(path) == ".plugin" {
				absPath, _ := filepath.Abs(path)

				plugins = append(plugins, absPath)
			}
		}
		return nil
	})
	return plugins, nil
}

func findFunctionInPlugin(pluginPath, functionName string) (plugin.Symbol, error) {

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Debugf("%v", err)
		return nil, fmt.Errorf("Unable to open Plugin [%s]", pluginPath)

	}

	symbol, err := plug.Lookup(functionName)
	if err != nil {
		log.Debugf("%v", err)
		return nil, fmt.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
	}

	return symbol, nil
}

func init() {

	pluginList, err := findPlugins("./plugin")
	if err != nil {
		log.Errorf("%v", err)
	} else {
		log.Debugf("Found [%d] plugins", len(pluginList))
		for x := range pluginList {
			symbol, err := findFunctionInPlugin(pluginList[x], "ParlayActionList")
			if err != nil {
				log.Errorf("%v", err)
				continue
			}

			pluginExec, ok := symbol.(func() []string)
			if !ok {
				log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
				continue
			}

			actions := pluginExec()

			for z := range actions {
				// This will give us a mapping of "action" => plugin
				pluginCache[actions[z]] = pluginList[x]
			}
		}
	}
}

//ListPlugins -
func ListPlugins() {

	pluginList, err := findPlugins("./plugin")
	if err != nil {
		log.Errorf("%v", err)
	} else {
		log.Debugf("Found [%d] plugins", len(pluginList))
		for x := range pluginList {
			symbol, err := findFunctionInPlugin(pluginList[x], "ParlayPluginInfo")
			if err != nil {
				log.Errorf("%v", err)
				continue
			}

			pluginExec, ok := symbol.(func() string)
			if !ok {
				log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
				continue
			}
			sanitizedPath := filepath.Base(pluginList[x])
			fmt.Printf("%s\t%s\n", sanitizedPath, pluginExec())
		}
	}
}

//ListPluginActions -
func ListPluginActions(pluginPath string) {

	symbol, err := findFunctionInPlugin(pluginPath, "ParlayActionList")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	pluginExec, ok := symbol.(func() []string)
	if !ok {
		log.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
		return
	}

	actions := pluginExec()

	symbol, err = findFunctionInPlugin(pluginPath, "ParlayActionDetails")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	pluginExec, ok = symbol.(func() []string)
	if !ok {
		log.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
		return
	}

	descriptions := pluginExec()

	if len(actions) != len(descriptions) {
		log.Warnf("Not all actions have descriptions, contact your plugin provider to have this fixed")
	}

	for x := range actions {
		fmt.Printf("%s\t%s\n", actions[x], descriptions[x])
	}
}

//UsagePlugin returns the usage of a plugin function
func UsagePlugin(pluginPath, action string) {

	symbol, err := findFunctionInPlugin(pluginPath, "ParlayUsage")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	pluginExec, ok := symbol.(func(string) (json.RawMessage, error))
	if !ok {
		log.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
		return
	}
	result, err := pluginExec(action)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	a := parlaytypes.Action{
		Name:       fmt.Sprintf("Example name for action [%s]", action),
		ActionType: action,
		Plugin:     result,
	}
	b, _ := json.MarshalIndent(a, "", "\t")
	fmt.Printf("%s\n", b)
}

// ExecuteAction uses the cache to find an action/plugin mapping
func ExecuteAction(action, host string, raw json.RawMessage) ([]parlaytypes.Action, error) {
	if pluginCache[action] == "" {
		// No KeyMap meaning that the action doesn't map to a plugin
		return nil, fmt.Errorf("Action [%s] does not exist or has no plugin associated with it", action)
	}
	return ExecuteActionInPlugin(pluginCache[action], action, host, raw)
}

// ExecuteActionInPlugin specifies the plugin and action directly
func ExecuteActionInPlugin(pluginPath, action, host string, raw json.RawMessage) ([]parlaytypes.Action, error) {

	// Check a function with the name ParlayExec exists
	symbol, err := findFunctionInPlugin(pluginPath, "ParlayExec")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	log.Debugf("Attempting plugin [%s]", action)
	// Check the function has the correct parameters
	pluginExec, ok := symbol.(func(string, string, json.RawMessage) ([]parlaytypes.Action, error))
	if !ok {
		return nil, fmt.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
	}

	// Pass the action type and the interface to the plugin
	return pluginExec(action, host, raw)
}
