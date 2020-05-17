package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkFile(folder, name string, isFolder bool) (string, error) {
	if folder == "" {
		return "", fmt.Errorf("%s is empty", name)
	}

	if !filepath.IsAbs(folder) {
		folder, _ = filepath.Abs(folder)
	}

	st, err := os.Stat(folder)

	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%s does not exist: %s", name, folder)
		}

		return "", fmt.Errorf("unable to stat %s at %q: %w", name, folder, err)
	}

	if isFolder && !st.IsDir() {
		return "", fmt.Errorf("%s must point to a folder in your filesystem, not a file", folder)
	}

	return folder, nil
}
