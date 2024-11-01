package healthcheck

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// getChecksum calculates the SHA-256 checksum of a file
func getChecksum(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Calculate SHA-256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// checkFileChecksum checks if a single file's checksum matches the expected value
func checkFileChecksum(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := getChecksum(filePath)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if actualChecksum != expectedChecksum {
		return true, nil
	}
	return false, nil
}

func VerifyFiles(rootDir string, checkFileNameList string, checkFileSumList string) error {
	if len(checkFileNameList) == 0 {
		return nil
	}
	fileNames := strings.Split(checkFileNameList, ",")
	checkSums := strings.Split(checkFileSumList, ",")
	// Check if the number of files matches the number of checksums
	if len(fileNames) != len(checkSums) {
		return fmt.Errorf("error: The number of files and checksums do not match.")
	}
	// Check each file's checksum
	alteredFiles := []string{}
	for i, fileName := range fileNames {
		expectedChecksum := strings.TrimSpace(checkSums[i])
		filePath := filepath.Join(rootDir, strings.TrimSpace(fileName))
		altered, err := checkFileChecksum(filePath, expectedChecksum)
		if err != nil {
			return err
		}
		if altered {
			alteredFiles = append(alteredFiles, filePath)
		}
	}
	if len(alteredFiles) > 0 {
		return fmt.Errorf("The following files have been altered: `%s`.\n"+
			"Please revert your changes or contact the teaching team "+
			"if you have a valid reason for adjusting them.",
			strings.Join(alteredFiles, "`, `"))
	}
	return nil
}
