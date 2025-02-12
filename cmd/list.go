package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Project struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

// AppProject represents the ArgoCD AppProject structure
type AppProject struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ArgoCD AppProjects from the configured directory",
	Run: func(cmd *cobra.Command, args []string) {

		// Get list of directories
		appProjectDirs := viper.GetStringSlice("projectPrefixes")
		if len(appProjectDirs) == 0 {
			log.Fatal("No appProjectDirs found in config")
		}

		// Print directories
		for _, dir := range appProjectDirs {
			fmt.Println("AppProject Directory:", dir)
		}

		// Read list of projects
		var projects []Project
		if err := viper.UnmarshalKey("projects", &projects); err != nil {
			log.Fatalf("Error parsing projects: %v", err)
		}

		// Print projects
		//for _, project := range projects {
		//	fmt.Printf("Project: %s, Path: %s\n", project.Name, project.Path)
		//}

		dir := viper.GetString("appProjectDir")
		if dir == "" {
			log.Fatal("appProjectDir not set in config.yaml")
		}

		files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
		if err != nil {
			log.Fatalf("Error reading directory: %v", err)
		}

		if len(files) == 0 {
			fmt.Println("No AppProject YAML files found.")
			return
		}

		for _, file := range files {
			appProject, err := readAppProject(file)
			if err != nil {
				log.Printf("Skipping %s: %v", file, err)
				continue
			}
			fmt.Printf("Found AppProject: %s\n", appProject.Metadata.Name)
		}
	},
}

func readAppProject(filename string) (*AppProject, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var appProject AppProject
	err = yaml.Unmarshal(data, &appProject)
	if err != nil {
		return nil, err
	}

	if appProject.Kind != "AppProject" {
		return nil, fmt.Errorf("not an AppProject resource")
	}

	return &appProject, nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
