package shell

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

//DefaultFailedCode constant
const DefaultFailedCode = 1

// RedirectToConsole Show command execution through the standard output or
// standard error Exit if error
func RedirectToConsole(out []byte, err error) {
	if err != nil {
		if out != nil {
			log.Error().Msg(string(out))
		}
		log.Fatal().Err(err).Msg(err.Error())
	}
	log.Info().Msg(string(out))
}

// RunCommand Run a shell command
func RunCommand(command string, args ...string) error {
	log.Debug().Msgf("Executing command: %s %s",
		command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

// GetFullPath Get Full Path
func GetFullPath(relativePath string) string {
	currentPath, _ := os.Getwd()
	fullPathTemplate := path.Join(currentPath, relativePath)
	return fullPathTemplate
}

// ExistsFileOrDirectory Exists File or Directory
func ExistsFileOrDirectory(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FindFilesByPattern Find Files in a directory by pattern
func FindFilesByPattern(targetDir string, pattern []string) ([]string, error) {

	for _, v := range pattern {
		matches, err := filepath.Glob(path.Join(targetDir, v))

		if err != nil {
			return nil, err
		}

		if len(matches) != 0 {
			return matches, nil
		}
	}

	return nil, nil
}

// RemoveDirectory Remove a directory
func RemoveDirectory(directory string) error {
	d, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(directory, name))
		if err != nil {
			return err
		}
	}

	return os.Remove(directory)
}
