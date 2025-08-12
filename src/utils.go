package main

import (
	"os/exec"
	"strings"
)

func getRepoURL() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Convert SSH or HTTPS git URL to https://github.com/owner/repo format
	url := strings.TrimSpace(string(output))
	url = strings.TrimSuffix(url, ".git")
	if strings.HasPrefix(url, "git@github.com:") {
		// Convert git@github.com:owner/repo to https://github.com/owner/repo
		url = "https://github.com/" + strings.TrimPrefix(url, "git@github.com:")
	} else if strings.HasPrefix(url, "https://github.com/") {
		// Already in correct format
		// do nothing
	} else {
		// Unknown format
		return ""
	}

	return url
}
