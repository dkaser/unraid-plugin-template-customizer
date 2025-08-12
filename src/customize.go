package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
)

func writePluginJSON(pluginOptions PluginOptions) error {
	var config PluginConfig

	config.Name = pluginOptions.PluginName
	config.PackageName = pluginOptions.PackageName
	config.Author = pluginOptions.AuthorName
	config.Support = pluginOptions.SupportURL
	config.Launch = pluginOptions.LaunchPath
	config.Icon = pluginOptions.IconName

	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling plugin.json: %v\n", err)
		return err
	}

	if err := os.WriteFile("plugin/plugin.json", updatedData, 0644); err != nil {
		fmt.Printf("Error writing plugin.json: %v\n", err)
		return err
	}
	fmt.Println("plugin.json updated successfully.")
	return nil
}

func updateSlackDesc(pluginOptions PluginOptions) error {
	// Generate ruler line
	rulerLine := "|-----handy-ruler------------------------------------------------------|"

	// Get repository URL if available
	repoURL := getRepoURL()

	lines := []string{
		// Add ruler line: spaces instead of package name, then the ruler line
		fmt.Sprintf("%s%s", strings.Repeat(" ", len(pluginOptions.PackageName)), rulerLine),
	}

	// Add package name line
	lines = append(lines, fmt.Sprintf("%s: %s", pluginOptions.PackageName, pluginOptions.PackageName))
	lines = append(lines, fmt.Sprintf("%s:", pluginOptions.PackageName))

	// Use go-wordwrap for description
	wrapped := wordwrap.WrapString(pluginOptions.Description, 70)
	descLines := strings.Split(wrapped, "\n")

	// Add description lines (max 6 lines)
	for i, line := range descLines {
		if i < 6 {
			lines = append(lines, fmt.Sprintf("%s: %s", pluginOptions.PackageName, line))
		}
	}

	// Fill remaining description lines
	for len(lines) < 10 {
		lines = append(lines, fmt.Sprintf("%s:", pluginOptions.PackageName))
	}

	// Add repository URL
	if repoURL != "" {
		lines = append(lines, fmt.Sprintf("%s: Repository: %s", pluginOptions.PackageName, repoURL))
	} else {
		lines = append(lines, fmt.Sprintf("%s:", pluginOptions.PackageName))
	}

	// Final empty line
	lines = append(lines, fmt.Sprintf("%s:", pluginOptions.PackageName))

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile("src/install/slack-desc", []byte(content), 0644); err != nil {
		fmt.Printf("Error writing slack-desc: %v\n", err)
		return err
	}
	fmt.Println("slack-desc updated successfully.")
	return nil
}

func updatePluginREADME(path string, pluginOptions PluginOptions) error {
	content := fmt.Sprintf("** %s **\n\n%s\n", pluginOptions.PluginName, pluginOptions.Description)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing plugin README: %v\n", err)
		return err
	}
	fmt.Println("Plugin README updated successfully.")
	return nil
}

func downloadLicense(url string, pluginOptions PluginOptions) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to download license: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read license content: %v\n", err)
		return err
	}

	// Replace placeholders
	licenseText := string(content)
	currentYear := time.Now().Year()
	licenseText = strings.ReplaceAll(licenseText, "{{ year }}", fmt.Sprintf("%d", currentYear))
	licenseText = strings.ReplaceAll(licenseText, "{{ organization }}", pluginOptions.AuthorName)
	licenseText = strings.ReplaceAll(licenseText, "{{ project }}", pluginOptions.PluginName)

	if err := os.WriteFile("LICENSE", []byte(licenseText), 0644); err != nil {
		fmt.Printf("Error writing LICENSE file: %v\n", err)
		return err
	}

	fmt.Println("License downloaded successfully.")

	// Remove LICENSE-unraid-plugin-template if it exists
	if _, err := os.Stat("LICENSE-unraid-plugin-template"); err == nil {
		if err := os.Remove("LICENSE-unraid-plugin-template"); err != nil {
			fmt.Printf("Warning: could not remove LICENSE-unraid-plugin-template: %v\n", err)
		}
	}

	return nil
}

func configureDiagnosticsJSON(pluginOptions PluginOptions) error {
	diagnosticsPath := fmt.Sprintf("src/usr/local/emhttp/plugins/%s/diagnostics.json", pluginOptions.PluginName)
	diagnosticsExamplePath := fmt.Sprintf("src/usr/local/emhttp/plugins/%s/diagnostics.json.example", pluginOptions.PluginName)

	cfg := DiagnosticsConfig{
		Title:             pluginOptions.DisplayName,
		Files:             []string{fmt.Sprintf("/var/log/%s*", pluginOptions.PluginName)},
		SystemDiagnostics: true,
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(diagnosticsPath, data, 0644)
	if err != nil {
		return err
	}

	if _, err := os.Stat(diagnosticsExamplePath); err == nil {
		if err := os.Remove(diagnosticsExamplePath); err != nil {
			fmt.Printf("Warning: could not remove diagnostics.json.example: %v\n", err)
		}
	}

	fmt.Println("diagnostics.json configured.")
	return nil
}
