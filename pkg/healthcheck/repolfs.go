package healthcheck

func RepoLFS(rootDir string) error {
	// cmd := exec.Command("git", "lfs", "fsck", "--pointers")
	// cmd.Dir = rootDir
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return fmt.Errorf(
	// 		"error running `git lfs fsck --pointers`: %w, output: %s",
	// 		err,
	// 		output,
	// 	)
	// }
	return nil
}
