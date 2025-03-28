package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

const logicalBlocksize diskfs.SectorSize = 2048

func copyFileTreeToISO(sourceRoot string, isoFileSystem *iso9660.FileSystem) error {
	return filepath.WalkDir(sourceRoot, func(sourcePath string, entry fs.DirEntry, err error) error {
		// If an error was encountered while traversing the source file
		// system, return it and stop traversal.
		if err != nil {
			return err
		}

		// All paths on the ISO file system will be relative.
		destinationPath, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return err
		}

		// Process directories.
		if entry.IsDir() {
			return isoFileSystem.Mkdir(destinationPath)
		}

		// Skip irregular files.
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		// Process regular files.
		return func() error {
			// Open the destination file in the ISO file system for writing.
			destination, err := isoFileSystem.OpenFile(destinationPath, os.O_CREATE|os.O_RDWR)
			if err != nil {
				return err
			}
			defer destination.Close()

			// Open the source file for reading.
			source, errorOpeningFile := os.Open(sourcePath)
			if errorOpeningFile != nil {
				return errorOpeningFile
			}
			defer source.Close()

			// Copy the contents of the source file to the ISO file
			_, err = io.Copy(destination, source)
			return err
		}()
	})
}

func buildIsoFromDirectory(sourceRoot string, outputFileName, label string) error {
	// Calculate the total size of the data to be written to the disk image.
	totalSize, err := calculateFileTreeSize(sourceRoot)
	if err != nil {
		return err
	}

	// Prepare a disk image.
	isoDisk, err := diskfs.Create(outputFileName, totalSize, logicalBlocksize)
	if err != nil {
		return err
	}

	// Prepare an ISO file system on the disk image.
	isoFileSystemSpec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: label,
	}
	var isoFileSystem *iso9660.FileSystem
	{
		fsys, err := isoDisk.CreateFilesystem(isoFileSystemSpec)
		if err != nil {
			return err
		}

		var ok bool
		isoFileSystem, ok = fsys.(*iso9660.FileSystem)
		if !ok {
			return fmt.Errorf("the iso9660 file system creation process returned something other than an iso9660 filesystem")
		}
	}

	// Copy all source files and folders to the ISO file system.
	if err := copyFileTreeToISO(sourceRoot, isoFileSystem); err != nil {
		return err
	}

	// Finalize the ISO file system.
	return isoFileSystem.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		DeepDirectories:  false,
		VolumeIdentifier: label,
	})
}
