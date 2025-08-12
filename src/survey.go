package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func surveyConfig() (PluginOptions, error) {
	var pluginOptions PluginOptions

	fmt.Println("This program will help you customize your Unraid plugin template.")
	fmt.Println("Please provide the following information:")
	fmt.Println()

	pluginNamePrompt := &survey.Input{Message: "Plugin Name (file name):"}
	_ = survey.AskOne(pluginNamePrompt, &pluginOptions.PluginName, survey.WithValidator(nameValidator))
	pluginOptions.PluginName = strings.TrimSpace(pluginOptions.PluginName)

	_ = survey.AskOne(&survey.Input{Message: "Display Name:"}, &pluginOptions.DisplayName)

	packageNamePrompt := &survey.Input{Message: "Package Name (use unraid-pluginname format):"}
	_ = survey.AskOne(packageNamePrompt, &pluginOptions.PackageName, survey.WithValidator(nameValidator))
	pluginOptions.PackageName = strings.TrimSpace(pluginOptions.PackageName)

	_ = survey.AskOne(&survey.Input{Message: "Author Name:"}, &pluginOptions.AuthorName)

	supportURLPrompt := &survey.Input{Message: "Support URL (forum thread):"}
	_ = survey.AskOne(supportURLPrompt, &pluginOptions.SupportURL, survey.WithValidator(urlValidator))
	pluginOptions.SupportURL = strings.TrimSpace(pluginOptions.SupportURL)
	if pluginOptions.SupportURL == "" {
		pluginOptions.SupportURL = "https://forums.unraid.net/"
	}

	fmt.Println()
	fmt.Println("Optional fields (press Enter to skip):")
	_ = survey.AskOne(&survey.Input{Message: "Launch Path (e.g., Tools/PluginName):"}, &pluginOptions.LaunchPath)
	_ = survey.AskOne(&survey.Input{Message: "Font Awesome Icon (e.g., fa-cog) or leave blank for custom icon:"}, &pluginOptions.IconName)
	fmt.Println()
	fmt.Println("Plugin Description (for installation display):")
	_ = survey.AskOne(&survey.Input{Message: "Enter a brief description of your plugin:"}, &pluginOptions.Description, survey.WithValidator(descriptionValidator))

	_ = survey.AskOne(&survey.Confirm{Message: "Would you like to enable Plugin Diagnostics?", Default: false}, &pluginOptions.SetupDiagnostics)

	// License selection
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("   License Selection")
	fmt.Println("=========================================")
	fmt.Println()
	licenseOptions := make([]string, len(licenses))
	for i, l := range licenses {
		licenseOptions[i] = l.Name
	}
	_ = survey.AskOne(&survey.Select{Message: "Select a license for your plugin:", Options: licenseOptions}, &pluginOptions.LicenseChoice)

	// Show summary before making changes
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("   Configuration Summary")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("Your plugin will be customized with the following settings:")
	fmt.Printf("- Plugin Name: %s\n", pluginOptions.PluginName)
	fmt.Printf("- Display Name: %s\n", pluginOptions.DisplayName)
	fmt.Printf("- Package Name: %s\n", pluginOptions.PackageName)
	fmt.Printf("- Author: %s\n", pluginOptions.AuthorName)
	fmt.Printf("- Support URL: %s\n", pluginOptions.SupportURL)
	if pluginOptions.LaunchPath != "" {
		fmt.Printf("- Launch Path: %s\n", pluginOptions.LaunchPath)
	}
	if pluginOptions.IconName != "" {
		fmt.Printf("- Icon: %s\n", pluginOptions.IconName)
	}
	fmt.Printf("- Description: %s\n", pluginOptions.Description)
	if pluginOptions.LicenseChoice != len(licenses)-1 {
		fmt.Printf("- License: %s\n", licenses[pluginOptions.LicenseChoice].Name)
	}
	fmt.Printf("- Plugin Diagnostics enabled: %v\n", pluginOptions.SetupDiagnostics)
	fmt.Println()

	// Final confirmation
	var proceed bool
	_ = survey.AskOne(&survey.Confirm{Message: "Proceed with these changes?", Default: true}, &proceed)
	if !proceed {
		fmt.Println("Setup cancelled. No changes made.")
		return PluginOptions{}, fmt.Errorf("setup cancelled by user")
	}

	return pluginOptions, nil
}
