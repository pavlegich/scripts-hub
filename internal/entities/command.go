// Package entities contains objects for the application.
package entities

// Command contains data for commands.
type Command struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Script string `json:"script"`
	Output string `json:"output"`
}
