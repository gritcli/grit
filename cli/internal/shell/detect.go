package shell

import (
	"os"
	"path/filepath"
)

// Type is an enumeration of supported shells.
type Type int

const (
	// UnknownType is the default value for an unknown shell.
	UnknownType Type = iota

	// ZshType is the Z shell.
	ZshType

	// BashType is the Bourne Again shell.
	BashType

	// FishType is the Friendly Interactive Shell.
	FishType

	// PowerShellType is the Windows shell.
	PowerShellType
)

// Detect returns the type of shell that is currently running.
func Detect() Type {
	if os.Getenv("PSVersionTable") != "" {
		return PowerShellType
	}

	shell := os.Getenv("SHELL")

	switch filepath.Base(shell) {
	case "zsh":
		return ZshType
	case "bash":
		return BashType
	case "fish":
		return FishType
	}

	return UnknownType
}
