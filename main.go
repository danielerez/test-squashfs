package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	workDir := os.TempDir()

	_, err := execute(fmt.Sprintf("7z x %s", "rootfs.img"), workDir)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	list, err := execute(fmt.Sprintf("unsquashfs -lc %s", "root.squashfs"), workDir)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	r, err := regexp.Compile(".*nmstatectl")
	if err != nil {
		logrus.Errorf(err.Error())
	}
	match := r.FindString(list)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	binaryPath := strings.TrimPrefix(match, "squashfs-root")
	output, err := execute(fmt.Sprintf("unsquashfs -no-xattrs %s -extract-file %s", "root.squashfs", binaryPath), workDir)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	fmt.Println(output)
}

func execute(command, workDir string) (string, error) {
	var stdoutBytes, stderrBytes bytes.Buffer
	formattedCmd, args := formatCommand(command)
	cmd := exec.Command(formattedCmd, args...)
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	logrus.Debugf("Running cmd: %s %s", formattedCmd, strings.Join(args[:], " "))
	cmd.Dir = workDir
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to execute cmd (%s): %s", cmd, stderrBytes.String())
	}

	return strings.TrimSuffix(stdoutBytes.String(), "\n"), nil
}

func formatCommand(command string) (string, []string) {
	formattedCmd := strings.Split(command, " ")
	return formattedCmd[0], formattedCmd[1:]
}
