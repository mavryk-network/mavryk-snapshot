package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/mavryk-network/mavryk-snapshot/pkg/snapshot"
)

type SnapshotExec struct {
	snapshotsPath     string
	mavkitNodeBinPath string
	mavrykPath        string
	mavrykConfig      string
}

func NewSnapshotExec(snapshotsPath, mavkitNodePath, mavrykPath string, mavrykConfig string) *SnapshotExec {
	return &SnapshotExec{snapshotsPath, mavkitNodePath, mavrykPath, mavrykConfig}
}

func (s *SnapshotExec) CreateSnapshot(historyMode snapshot.HistoryModeType) {
	log.Println("Creating snapshot.")
	script := "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && " + s.mavkitNodeBinPath + " snapshot export --block head~10 --data-dir " + s.mavrykPath + "/data --config-file " + s.mavrykConfig

	if historyMode == snapshot.ROLLING {
		script = script + " --rolling"
	}

	if historyMode == snapshot.ARCHIVE {
		script = script + " --archive"
	}

	_, _ = s.execScript(script)
}

func (s *SnapshotExec) GetSnapshotName(historyMode snapshot.HistoryModeType) (string, error) {
	log.Println("Getting snapshot names.")
	script := "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && /bin/ls -1a"
	stdout, _ := s.execScript(script)

	snapshotfilenames := strings.Split(stdout.String(), "\n")
	log.Printf("All files found: %v \n", snapshotfilenames)

	for _, filename := range snapshotfilenames {
		if strings.Contains(filename, string(historyMode)) {
			log.Printf("Snapshot file found is: %q. \n", filename)
			return filename, nil
		}
	}

	return "", fmt.Errorf("Snapshot file not found.")
}

func (s *SnapshotExec) GetSnapshotHeaderOutput(filepath string) string {
	log.Printf("Getting snapshot header output for file: %q. \n", filepath)
	script := s.mavkitNodeBinPath + " snapshot info --json " + s.snapshotsPath + "/" + filepath
	stdout, _ := s.execScript(script)
	log.Printf("Snapshot header output: %q. \n", stdout.String())
	return stdout.String()
}

func (s *SnapshotExec) DeleteLocalSnapshots() {
	log.Println("Deleting local snapshots.")
	script := "rm -rf " + s.snapshotsPath + "/*"
	_, _ = s.execScript(script)
}

func (s *SnapshotExec) execScript(script string) (bytes.Buffer, bytes.Buffer) {
	log.Printf("Executing script: %q. \n", script)

	// Set a timeout (e.g., 60 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()

	// Prepare the command with context
	cmd := exec.CommandContext(ctx, "sh", "-c", script)

	// Buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Check if the context was canceled (timeout reached)
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Timeout reached for script: %q, but execution will continue. \n", script)
	} else if err != nil {
		log.Fatalf("%v \n", err)
	}

	// Log stdout and stderr if they contain output
	if stdout.Len() > 0 {
		log.Printf("stdout: \n%s\n", stdout.String())
	}
	if stderr.Len() > 0 {
		log.Printf("stderr: \n%s\n", stderr.String())
	}

	return stdout, stderr
}
