package internal

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// extractVersionFromCommit extracts version from the latest commit message
func extractVersionFromCommit(currentDir string) (string, error) {
	fmt.Println("Extracting version from commit message...")

	// 优先从最近的提交消息获取版本（用于 CI 环境）
	cmd := exec.Command("git", "log", "-10", "--pretty=format:%s")
	cmd.Dir = currentDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit messages: %v", err)
	}

	commitMessages := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("Recent commit messages: %v\n", commitMessages)

	for _, commitMessage := range commitMessages {
		commitMessage = strings.TrimSpace(commitMessage)
		if commitMessage == "" {
			continue
		}

		// 跳过自动生成的 commit 消息
		if strings.Contains(commitMessage, "Update llpkgstore.json") ||
			strings.Contains(commitMessage, "Merge") ||
			strings.Contains(commitMessage, "chore:") {
			continue
		}

		// 尝试解析版本
		version, err := parseVersionFromCommitMessage(commitMessage)
		if err == nil {
			fmt.Printf("Found version in commit message: %s -> %s\n", commitMessage, version)
			return version, nil
		}
	}

	// 如果提交消息中没有找到版本，尝试从 git tag 获取最新版本
	if version, err := extractVersionFromGitTag(currentDir); err == nil {
		fmt.Printf("No version in commit messages, using git tag: %s\n", version)
		return version, nil
	}

	return "", fmt.Errorf("no version pattern found in recent commit messages or git tags")
}

// extractVersionFromGitTag extracts version from git tags
func extractVersionFromGitTag(currentDir string) (string, error) {
	// 获取最新的 tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = currentDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git tags: %v", err)
	}

	tag := strings.TrimSpace(string(output))
	if tag == "" {
		return "", fmt.Errorf("no git tags found")
	}

	// 验证 tag 格式
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}

	// 简单的版本格式验证
	if !isValidVersionFormat(tag) {
		return "", fmt.Errorf("invalid version format: %s", tag)
	}

	return tag, nil
}

// isValidVersionFormat validates version format (v1.2.3)
func isValidVersionFormat(version string) bool {
	// 移除 v 前缀
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	// 检查是否包含数字和点
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}

	// 检查每个部分是否为数字
	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

// parseVersionFromCommitMessage parses version from commit message
func parseVersionFromCommitMessage(commitMessage string) (string, error) {
	// Define regex patterns for different formats (C++ compatible)
	patterns := []string{
		`Release-as:\s*[^/]+/(v[\d.]+)`, // "Release-as: numpy/v1.26.4" (C++ compatible)
		`Release-as:\s*(v[\d.]+)`,       // "Release-as: v1.26.4"
		`Release:\s*[^/]+/(v[\d.]+)`,    // "Release: numpy/v1.26.4"
		`Release:\s*(v[\d.]+)`,          // "Release: v1.26.4"
		`Version:\s*(v[\d.]+)`,          // "Version: v1.26.4"
		`version:\s*(v[\d.]+)`,          // "version: v1.26.4"
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(commitMessage)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no version pattern found in commit message: %s", commitMessage)
}
