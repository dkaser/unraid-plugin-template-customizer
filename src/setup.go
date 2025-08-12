package main

import (
	"fmt"
	"os"
)

func runSetup() error {
	// Check if plugin.json exists
	if _, err := os.Stat("plugin/plugin.json"); os.IsNotExist(err) {
		return fmt.Errorf("plugin.json not found. Make sure you're running this from the plugin template root directory")
	}

	// If oldDir does not exist, fail
	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin folder not found. Make sure you're running this from the plugin template root directory")
	}

	pluginOptions, err := surveyConfig()
	if err != nil {
		return fmt.Errorf("error during survey configuration: %v", err)
	}

	// Update plugin.json
	if err := writePluginJSON(pluginOptions); err != nil {
		return fmt.Errorf("error updating plugin.json: %v", err)
	}

	// Update slack-desc
	if err := updateSlackDesc(pluginOptions); err != nil {
		return fmt.Errorf("error updating slack-desc: %v", err)
	}

	newDir := fmt.Sprintf("src/usr/local/emhttp/plugins/%s", pluginOptions.PluginName)

	if err := renamePluginDir(oldDir, newDir); err != nil {
		return err
	}

	if err := updateReadmeIfExists(newDir, pluginOptions); err != nil {
		return err
	}

	if err := handleLicense(pluginOptions); err != nil {
		return err
	}

	if pluginOptions.SetupDiagnostics {
		if err := configureDiagnosticsJSON(pluginOptions); err != nil {
			fmt.Printf("Error configuring diagnostics.json: %v\n", err)
		}
	}

	printSummary(newDir)
	return nil
}

func renamePluginDir(oldDir, newDir string) error {
	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("error renaming plugin directory: %v", err)
	}
	fmt.Println("Plugin directory renamed.")
	return nil
}

func updateReadmeIfExists(newDir string, pluginOptions PluginOptions) error {
	readmePath := fmt.Sprintf("%s/README.md", newDir)
	if _, err := os.Stat(readmePath); err == nil {
		if err := updatePluginREADME(readmePath, pluginOptions); err != nil {
			return fmt.Errorf("error updating plugin README: %v", err)
		}
	}
	fmt.Println("README.md updated successfully.")
	return nil
}

func handleLicense(pluginOptions PluginOptions) error {
	if pluginOptions.LicenseChoice != len(licenses)-1 {
		license := licenses[pluginOptions.LicenseChoice]
		if license.URL != "" {
			if err := downloadLicense(license.URL, pluginOptions); err != nil {
				return fmt.Errorf("error downloading license: %v", err)
			}
		}
	}
	return nil
}

func printSummary(newDir string) {
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("   Setup Complete!")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("1. Review and customize the files in %s/\n", newDir)
	fmt.Println("2. Add your plugin's functionality")
	fmt.Println("3. Test your plugin on your Unraid server")
	fmt.Println("4. Create a release on GitHub when ready")
	fmt.Println()
}
