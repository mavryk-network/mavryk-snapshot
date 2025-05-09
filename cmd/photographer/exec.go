package main

import (
	"bytes"
	// "encoding/json"
	"fmt"
	"log"
	"os/exec"
	// "strconv"
	"strings"

	"github.com/mavryk-network/mavryk-snapshot/pkg/snapshot"
)

type SnapshotExec struct {
	snapshotsPath       string
	mavkitClientBinPath string
	mavkitNodeBinPath   string
	mavrykPath          string
	mavrykVolume        string
	mavrykConfig        string
}

type SnapshotHeader struct {
	Version   int    `json:"version"`
	ChainName string `json:"chain_name"`
	Mode      string `json:"mode"`
	BlockHash string `json:"block_hash"`
	Level     int    `json:"level"`
	Timestamp string `json:"timestamp"`
}

type Snapshot struct {
	SnapshotHeader SnapshotHeader `json:"snapshot_header"`
}

func NewSnapshotExec(snapshotsPath, mavkitClientPath, mavkitNodePath, mavrykPath string, mavrykVolume string, mavrykConfig string) *SnapshotExec {
	return &SnapshotExec{snapshotsPath, mavkitClientPath, mavkitNodePath, mavrykPath, mavrykVolume, mavrykConfig}
}

func (s *SnapshotExec) CreateSnapshot(historyMode snapshot.HistoryModeType) {
	log.Println("Creating snapshot.")
	script := "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && " + s.mavkitNodeBinPath + " snapshot export --block head~10 --data-dir " + s.mavrykPath + "/data --config-file " + s.mavrykConfig

	if historyMode == snapshot.ROLLING {
		script = script + " --rolling"
	}

	// if historyMode == snapshot.ARCHIVE {
	// 	script = "wget -qO-  http://127.0.0.1:8732/chains/main/blocks/head/header | sed -E 's/.*\"hash\":\"?([^,\"]*)\"?.*/\\1/'"
	// 	block_hash, _ := s.execScript(script)
	// 	script = "wget -qO-  http://127.0.0.1:8732/chains/main/blocks/head/header | sed -E 's/.*\"level\":\"?([^,\"]*)\"?.*/\\1/'"
	// 	level, _ := s.execScript(script)
	// 	script = "sed -n 's/.*\"chain_name\": \"\\([^\"]*\\)\".*/\\1/p' " + s.mavrykConfig
	// 	chain_name, _ := s.execScript(script)
	// 	filename := s.snapshotsPath + "/" + chain_name.String() + "-" + block_hash.String() + "-" + level.String() + ".archive"
	// 	cleaned_filename := strings.ReplaceAll(filename, "\n", "")
	// 	script = "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && tar cvf - . --exclude='data/*.json' --exclude='client' --exclude='data/lock' -C " + s.mavrykPath + " | lz4 > " + cleaned_filename
	// }

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

// func (s *SnapshotExec) GetArchiveTarballHeaderOutput(filepath string) string {
// 	log.Printf("Getting tarball header output for file: %q. \n", filepath)
// 	script := "wget -qO-  http://127.0.0.1:8732/chains/main/blocks/head/header | sed -E 's/.*\"hash\":\"?([^,\"]*)\"?.*/\\1/'"
// 	block_hash, _ := s.execScript(script)
// 	cleaned_block_hash := strings.ReplaceAll(block_hash.String(), "\n", "")
// 	script = "wget -qO-  http://127.0.0.1:8732/chains/main/blocks/head/header | sed -E 's/.*\"level\":\"?([^,\"]*)\"?.*/\\1/'"
// 	level_string, _ := s.execScript(script)
// 	cleaned_level_string := strings.ReplaceAll(level_string.String(), "\n", "")
// 	level, err := strconv.Atoi(cleaned_level_string)
// 	if err != nil {
// 		panic(err)
// 	}
// 	script = "wget -qO-  http://127.0.0.1:8732/chains/main/blocks/head/header | sed -E 's/.*\"timestamp\":\"?([^,\"]*)\"?.*/\\1/'"
// 	timestamp, _ := s.execScript(script)
// 	cleaned_timestamp := strings.ReplaceAll(timestamp.String(), "\n", "")
// 	script = "sed -n 's/.*\"chain_name\": \"\\([^\"]*\\)\".*/\\1/p' " + s.mavrykConfig
// 	chain_name, _ := s.execScript(script)
// 	cleaned_chain_name := strings.ReplaceAll(chain_name.String(), "\n", "")

// 	// Create an instance of Snapshot
// 	snapshot := Snapshot{
// 		SnapshotHeader: SnapshotHeader{
// 			Version:   7,
// 			ChainName: cleaned_chain_name,
// 			Mode:      "archive",
// 			BlockHash: cleaned_block_hash,
// 			Level:     level,
// 			Timestamp: cleaned_timestamp,
// 		},
// 	}

// 	// Marshal the struct into JSON
// 	jsonData, err := json.Marshal(snapshot)
// 	if err != nil {
// 		log.Fatalf("Error marshaling to JSON: %v", err)
// 	}

// 	// Convert to string and print
// 	jsonString := string(jsonData)
// 	log.Printf("Tarball header output: %q. \n", jsonString)

// 	return jsonString
// }

func (s *SnapshotExec) DeleteLocalSnapshots() {
	log.Println("Deleting local snapshots.")
	script := "rm -rf " + s.snapshotsPath + "/*"
	_, _ = s.execScript(script)
}

func (s *SnapshotExec) execScript(script string) (bytes.Buffer, bytes.Buffer) {
	log.Printf("Executing script: %q. \n", script)
	cmd := exec.Command("sh", "-c", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	if stdout.Len() > 0 {
		log.Printf("stdout: \n%s\n", stdout.String())
	}
	if stderr.Len() > 0 {
		log.Printf("stderr: \n%s\n", stderr.String())
	}

	return stdout, stderr
}
