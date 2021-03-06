// +build !windows,!darwin,!linux

package appdir

import (
	"fmt"
	"path/filepath"
)

func SetHomeDir(dir string) {
	// do nothing
}

func general(app string) string {
	return InHomeDir(fmt.Sprintf(".%s", app))
}

func logs(app string) string {
	return filepath.Join(general(app), "logs")
}
