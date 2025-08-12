package main

type DiagnosticsConfig struct {
	Title             string   `json:"title"`
	Files             []string `json:"files"`
	SystemDiagnostics bool     `json:"system_diagnostics"`
}

type PluginConfig struct {
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	Author      string `json:"author"`
	Support     string `json:"support"`
	Launch      string `json:"launch,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type License struct {
	Name string
	URL  string
}

type PluginOptions struct {
	PluginName       string
	DisplayName      string
	PackageName      string
	AuthorName       string
	SupportURL       string
	LaunchPath       string
	IconName         string
	Description      string
	SetupDiagnostics bool
	LicenseChoice    int
}
