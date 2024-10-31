// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package stdfs

import "os"

func IsFileExists(path string) (bool, error) {
	sb, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !sb.IsDir() && sb.Mode().IsRegular(), nil
}
