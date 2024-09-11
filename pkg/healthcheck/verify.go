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
func checkFileChecksum(rootDir, fileName, expectedChecksum string) error {
	filePath := filepath.Join(rootDir, strings.TrimSpace(fileName))
	actualChecksum, err := getChecksum(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filePath, err)
	}
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("Checksum for %s failed. Expected %s, but got %s. Please revert your changes or contact the teaching team if you have a valid reason for adjusting them.", filePath, expectedChecksum, actualChecksum)
	}
	return nil
}

func VerifyFiles(rootDir string, checkFileNameList string, checkFileSumList string) error {
	if len(checkFileNameList) == 0 {
		return nil
	}
	fileNames := strings.Split(checkFileNameList, ",")
	checkSums := strings.Split(checkFileSumList, ",")
	// Check each file's checksum
	for i, fileName := range fileNames {
		expectedChecksum := strings.TrimSpace(checkSums[i])
		err := checkFileChecksum(rootDir, fileName, expectedChecksum)
		if err != nil {
			return err
		}
	}
	return nil
}
