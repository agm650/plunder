package main

import (
	"encoding/json"
	"fmt"

	"plunder-app/plunder/parlay/parlaytypes"
)

const pluginInfo = `This plugin is used to managed docker automation`

type image struct {

	// Image details
	ImageNames []string `json:"imageName"`
	ImageFiles []string `json:"imageFile"`

	DockerLocalSudo  bool `json:"localSudo"`
	DockerRemoteSudo bool `json:"remoteSudo"`
}

type tag struct {

	// A list of sources and target tags
	SourceNames []string `json:"sourceNames,omitempty"`
	TargetNames []string `json:"targetNames,omitempty"`

	// These two fields are used to change out a tag (e.g. version number) or the repository itself
	TargetTag  string `json:"imageTag,omitempty"`
	TargetRepo string `json:"imageRepo,omitempty"`
}

// Dummy main function
func main() {}

// ParlayActionList - This should return an array of actions
func ParlayActionList() []string {
	return []string{
		"docker/image",
		"docker/tag"}
}

// ParlayActionDetails - This should return an array of action descriptions
func ParlayActionDetails() []string {
	return []string{
		"This action automates the management of docker images",
		"This action manages the tagging of docker images"}
}

// ParlayPluginInfo - returns information about the plugin
func ParlayPluginInfo() string {
	return pluginInfo
}

// ParlayUsage - Returns the json that matches the specific action
// <- action is a string that defines which action the usage information should be
// <- raw - raw JSON that will be manipulated into a correct struct that matches the action
// -> err is any error that has been generated
func ParlayUsage(action string) (raw json.RawMessage, err error) {

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
	switch action {
	case "docker/image":
		a := image{
			ImageFiles: []string{"./my_image.tar.gz", "./my__other_image.tar.gz"},
			ImageNames: []string{"gcr.io/my_image:latest", "gcr.io/my_other_image:latest"},
		}
		// In order to turn a struct into an map[string]interface we need to turn it into JSON

		return json.Marshal(a)
	case "docker/tag":
		a := tag{
			SourceNames: []string{"gcr.io/my_image:latest"},
			TargetNames: []string{"internal_repo/my_image:1.0"},
		}
		// In order to turn a struct into an map[string]interface we need to turn it into JSON

		return json.Marshal(a)
	default:
		return raw, fmt.Errorf("Action [%s] could not be found", action)
	}
}

// ParlayExec - Parses the action and the data that the action will consume
// <- action a string that details the action to be executed
// <- raw - raw JSON that will be manipulated into a correct struct that matches the action
// -> actions are an array of generated actions that the parser will then execute
// -> err is any error that has been generated
func ParlayExec(action, host string, raw json.RawMessage) (actions []parlaytypes.Action, err error) {

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
	switch action {
	case "docker/image":
		var img image
		// Unmarshall the JSON into the struct
		err = json.Unmarshal(raw, &img)
		return img.generateImageActions(host), err
	case "docker/tag":
		var t tag
		// Unmarshall the JSON into the struct
		err = json.Unmarshal(raw, &t)
		return t.generateTagActions(host)
	default:
		return
	}
}
