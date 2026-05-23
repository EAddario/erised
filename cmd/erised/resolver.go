package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// NewLocalFileResolver creates a FileResolver that searches for a file by walking
// the directory tree under searchPath on every request.
func NewLocalFileResolver(searchPath string) FileResolver {
	if searchPath == "" {
		return nil
	}

	return func(filename string) ([]byte, error) {
		var fileData []byte
		var walkErr error = ErrFileNotFound
		stopWalk := errors.New("stop walk")

		_ = filepath.WalkDir(searchPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				log.Error().Msg("Invalid path: " + path)
				log.Debug().Msg(fmt.Sprintf("Error: %v", err))
				walkErr = ErrInvalidPath
				return stopWalk
			}

			if !entry.IsDir() && filepath.Base(path) == filename {
				if ct, err := os.ReadFile(path); err != nil {
					log.Error().Msg("Unable to open the file: " + path)
					log.Debug().Msg(fmt.Sprintf("Error: %v", err))
					walkErr = ErrFileAccess
					return stopWalk
				} else {
					log.Info().Msg(fmt.Sprintf("Reading file %v", path))
					fileData = ct
					walkErr = nil
					return stopWalk
				}
			}

			log.Debug().Msg("File " + filename + " not found in " + path)
			return nil
		})

		if walkErr != nil {
			return nil, walkErr
		}
		return fileData, nil
	}
}
