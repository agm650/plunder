package parlaytypes

import (
	"encoding/json"
)

// TreasureMap - X Marks the spot
// The treasure maps define the automation that will take place on the hosts defined
type TreasureMap struct {
	// An array/list of deployments that will take places as part of this "map"
	Deployments []Deployment `json:"deployments"`
}

// Deployment defines the hosts and the action(s) that should be performed on them
type Deployment struct {
	// Name of the deployment that is taking place i.e. (Install MySQL)
	Name string `json:"name"`
	// An array/list of hosts that these actions should be performed upon
	Hosts []string `json:"hosts"`

	// Parallel allow multiple actions across multiple hosts in parallel
	Parallel         bool `json:"parallel"`
	ParallelSessions int  `json:"parallelSessions"`

	// The actions that should be performed
	Actions []Action `json:"actions"`
}

// Action defines what the instructions that will be executed
type Action struct {
	Name       string `json:"name"`
	ActionType string `json:"type"`
	Timeout    int    `json:"timeout"`

	// File based operations
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	FileMove    bool   `json:"fileMove,omitempty"`

	// Package manager operations
	PkgManager   string `json:"packageManager,omitempty"`
	PkgOperation string `json:"packageOperation,omitempty"`
	Packages     string `json:"packages,omitempty"`

	// Command operations
	Command          string   `json:"command,omitempty"`
	Commands         []string `json:"commands,omitempty"`
	CommandLocal     bool     `json:"commandLocal,omitempty"`
	CommandSaveFile  string   `json:"commandSaveFile,omitempty"`
	CommandSaveAsKey string   `json:"commandSaveAsKey,omitempty"`
	CommandSudo      string   `json:"commandSudo,omitempty"`

	// Piping commands, read in a file and send over stdin, or capture stdout from a local command
	CommandPipeFile string `json:"commandPipeFile,omitempty"`
	CommandPipeCmd  string `json:"commandPipeCmd,omitempty"`

	// Ignore any failures
	IgnoreFailure bool `json:"ignoreFail,omitempty"`

	// Key operations
	KeyFile string `json:"keyFile,omitempty"`
	KeyName string `json:"keyName,omitempty"`

	//Plugin Spec
	Plugin json.RawMessage `json:"plugin,omitempty"`
}
