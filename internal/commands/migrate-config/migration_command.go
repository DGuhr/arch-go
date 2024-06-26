package migrate_config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fdaines/arch-go/api/configuration"

	"github.com/fdaines/arch-go/internal/commands"
	"github.com/fdaines/arch-go/internal/utils/timer"
	"github.com/fdaines/arch-go/internal/validators"
	"gopkg.in/yaml.v2"
)

type migrateConfigCommand struct {
	commands.BaseCommand
	path string
}

func NewCommand(output io.Writer, path string) migrateConfigCommand {
	return migrateConfigCommand{
		commands.BaseCommand{Output: output},
		path,
	}
}

func (dc migrateConfigCommand) Run() {
	var exitCode int
	timer.ExecuteWithTimer(func() {
		configuration := MigrateConfigurationCommand(dc.Output, dc.path)
		if configuration == nil {
			exitCode = 1
		}
	})
	os.Exit(exitCode)
}

func MigrateConfigurationCommand(out io.Writer, path string) *configuration.Config {
	filename := filepath.Join(path, "arch-go.yml")
	conf, err := configuration.LoadConfig(filename)
	if err == nil && conf != nil {
		err2 := validators.ValidateConfiguration(conf)
		if err2 != nil {
			fmt.Fprintf(out, "Invalid Configuration: %+v\n", err2)
			return nil
		}
		fmt.Fprintln(out, "File is already compatible with version 1")
		return conf
	}

	if err == nil {
		fmt.Fprintln(out, "File is already compatible with version 1")
		return conf
	}
	deprecatedConfiguration, err := configuration.LoadDeprecatedConfig(filename)
	if err != nil {
		fmt.Fprintf(out, "Error: %+v\n", err)
		return nil
	}

	fmt.Fprintf(out, "Migrating deprecated configuration to current schema.\n")
	conf = migrateRules(deprecatedConfiguration)
	yamlData, err := yaml.Marshal(&deprecatedConfiguration)
	if err != nil {
		fmt.Fprintf(out, "Error while Marshaling. %+v\n", err)
		return nil
	}
	ok := writeConfiguration(yamlData, "old-arch-go.yml")
	if ok {
		fmt.Fprintf(out, "Deprecated configuration backup at: old-arch-go.yml\n")
	}
	yamlData, err = yaml.Marshal(&conf)
	if err != nil {
		fmt.Fprintf(out, "Error while Marshaling. %+v\n", err)
	}
	ok = writeConfiguration(yamlData, "arch-go.yml")
	if ok {
		fmt.Fprintf(out, "Configuration saved at: arch-go.yml\n")
	}
	return conf
}

func writeConfiguration(data []byte, filename string) bool {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}
	return true
}

func migrateRules(deprecatedConfig *configuration.DeprecatedConfig) *configuration.Config {
	return &configuration.Config{
		Version:           1,
		DependenciesRules: migrateDependencyRules(deprecatedConfig.DependenciesRules),
		ContentRules:      deprecatedConfig.ContentRules,
		FunctionsRules:    deprecatedConfig.FunctionsRules,
		NamingRules:       deprecatedConfig.NamingRules,
	}
}

func migrateDependencyRules(rules []*configuration.DeprecatedDependenciesRule) []*configuration.DependenciesRule {
	var dependencyRules []*configuration.DependenciesRule
	for _, r := range rules {
		dependencyRules = append(dependencyRules, &configuration.DependenciesRule{
			Package:             r.Package,
			ShouldOnlyDependsOn: resolveAllowedDependencies(r),
			ShouldNotDependsOn:  resolveRestrictedDependencies(r),
		})
	}

	return dependencyRules
}

func resolveAllowedDependencies(r *configuration.DeprecatedDependenciesRule) *configuration.Dependencies {
	if len(r.ShouldOnlyDependsOn)+len(r.ShouldOnlyDependsOnExternal) > 0 {
		return &configuration.Dependencies{
			Internal: r.ShouldOnlyDependsOn,
			External: r.ShouldOnlyDependsOnExternal,
		}
	}
	return nil
}

func resolveRestrictedDependencies(r *configuration.DeprecatedDependenciesRule) *configuration.Dependencies {
	if len(r.ShouldNotDependsOn)+len(r.ShouldNotDependsOnExternal) > 0 {
		return &configuration.Dependencies{
			Internal: r.ShouldNotDependsOn,
			External: r.ShouldNotDependsOnExternal,
		}
	}
	return nil
}
