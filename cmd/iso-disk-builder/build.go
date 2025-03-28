package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildCmd builds an ISO 9660 disk image from a source file system.
type BuildCmd struct {
	SourceDir  string `kong:"source,required,help='Path to the source directory to be copied into the ISO 9660 disk image.'"`
	OutputFile string `kong:"output,required,help='Path to the ISO 9660 disk image to be created. Must have a .iso file name suffix.'"`
	Label      string `kong:"label,help='The ISO 9660 disk image volume label.'"`
}

// Run executes the build command.
func (cmd BuildCmd) Run(ctx context.Context) error {
	// Validate the source directory.
	if cmd.SourceDir == "" {
		return errors.New("a source directory path was not provided")
	}

	sourceInfo, err := os.Stat(cmd.SourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}
	if !sourceInfo.IsDir() {
		return errors.New("the provided source directory path is not a directory")
	}

	// Ensure that the output file path does is valid and not an
	// existing directory.
	{
		outputFileName := filepath.Base(cmd.OutputFile)
		if !strings.HasSuffix(outputFileName, ".iso") {
			return errors.New("the provided output file path must have a \".iso\" file name suffix")
		}
	}

	outputInfo, err := os.Stat(cmd.OutputFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to look for an existing output file: %w", err)
		}
	} else {
		if !outputInfo.Mode().IsRegular() {
			return errors.New("the provided output file path is not a regular file")
		}
	}

	// If a label was not provided, use a default.
	if cmd.Label == "" {
		if sourceDirName := filepath.Base(cmd.SourceDir); sourceDirName != "." {
			cmd.Label = sourceDirName
		} else {
			cmd.Label = "ISO 9660 Disk"
		}
	}

	return buildIsoFromDirectory(cmd.SourceDir, cmd.OutputFile, cmd.Label)
}
