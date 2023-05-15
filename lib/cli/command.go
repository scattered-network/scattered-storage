package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigMap = helpers.Data

const (
	defaultConfigFile       = "/etc/scattered-storage/rbd-docker-plugin"
	defaultOperationTimeout = 5
)

type Cmd struct {
	ConfigMap    map[string]*ConfigMap
	envPrefix    string
	DebugEnabled bool
	CobraRoot    *cobra.Command
}

// NewCLICommand returns a *Cmd struct that includes a ConfigMap,
// the envPrefix, and a *cobra.Command for the root of the application.
// A local configFile string variable is used to create a "config" flag
// that is used in the initConfiguration PersistentPreRunE step.
func NewCLICommand(
	name, short, long, envPrefix string, parent *cobra.Command, configMap map[string]*ConfigMap,
	run func(cmd *cobra.Command, args []string),
) *Cmd {
	newCmd := &Cmd{
		ConfigMap:    configMap,
		envPrefix:    envPrefix,
		DebugEnabled: false,
		CobraRoot:    nil,
	}

	//nolint:exhaustruct
	cobraCmd := &cobra.Command{
		Use:   name,
		Short: short,
		Long:  long,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return newCmd.initConfiguration(defaultConfigFile)
		},
		Run: run,
	}

	if parent != nil {
		parent.AddCommand(cobraCmd)
	} else {
		cobraCmd.Flags().StringP("config", "c", defaultConfigFile, "config file")
		cobraCmd.Flags().IntP("timeout", "t", defaultOperationTimeout, "timeout for operations (in seconds)")
		newCmd.CobraRoot = cobraCmd
	}

	return newCmd
}

// KeyExists checks the ConfigMap for a given string key.
// Returns false if not found.
func (c *Cmd) KeyExists(key string) bool {
	if _, ok := c.ConfigMap[key]; ok {
		return true
	}

	return false
}

// initConfiguration will attempt to set the path to the configuration file of the application.
// If --config is used then that will have priority, then environment variables, finally
// the default value passed into NewCLICommand() will be used.
func (c *Cmd) initConfiguration(defaultConfigFile string) error {
	c.bindEnvironmentVariables()

	mamba := viper.New()
	mamba.SetConfigType("json")

	if configFile, err := c.CobraRoot.Flags().GetString("config"); err != nil {
		if configFile != "" {
			mamba.SetConfigFile(configFile)
		} else {
			mamba.SetConfigFile(defaultConfigFile)
		}
	}

	if err := mamba.ReadInConfig(); errors.Is(err, os.ErrNotExist) {
		if c.DebugEnabled {
			log.Error().Str("Config", mamba.ConfigFileUsed()).Msg("config file does not exist")
		}
	}

	return nil
}

// bindEnvironmentVariables steps through each flag.
func (c *Cmd) bindEnvironmentVariables() {
	c.CobraRoot.Flags().VisitAll(
		func(flag *pflag.Flag) {
			variableName := strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))
			optionName := fmt.Sprintf("%s_%s", c.envPrefix, variableName)
			optionValue, envFound := os.LookupEnv(optionName)
			if _, found := c.ConfigMap[variableName]; found {
				if envFound && !flag.Changed { // use ENV if found, unless flag is used
					c.ConfigMap[variableName].SetValue(optionValue, flag.Value.Type())
					err := flag.Value.Set(optionValue) // overwrite flag value content with ENV
					if err != nil {
						log.Info().Str("OptionName", optionName).Str("optionValue", optionValue).
							Msgf("Flag %s could not be set to %s", optionName, optionValue)

						return
					}
				} else { // anything else needs value set into config map
					c.ConfigMap[variableName].SetValue(flag.Value, flag.Value.Type())
				}
			}
		},
	)
}
