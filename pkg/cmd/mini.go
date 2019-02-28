package cmd // import "minibox.ai/minibox/pkg/cmd"

import (
	"fmt"
	"os"
	"path"

	"github.com/cloudfoundry/jibber_jabber"
	homedir "github.com/mitchellh/go-homedir"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
	"golang.org/x/text/message"
)

var userLocale string

const fallback = "en"

const APP_NAME = "minibox"
const SITE_DNS = "minibox.ai"

func init() {
	cobra.OnInitialize(initConfig)
	flag := rootCmd.PersistentFlags()
	flag.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.minibox.yaml)")
	// flag.Lookup("debug").NoOptDefVal = "true"

	// rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	// rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	site_url := fmt.Sprintf("https://%s", SITE_DNS)
	viper.SetDefault("webserver", site_url)
	viper.SetDefault("apiserver", site_url)

	userLocale, _ = jibber_jabber.DetectIETF()
	viper.SetDefault("userLocale", userLocale)

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(trainCmd)
}

var rootCmd = &cobra.Command{
	Use: "mini",
	SuggestionsMinimumDistance: 2,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Minibox-cli",
	Long:  `All software has versions. This is Minibox's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Minibox cli tools v0.1 -- HEAD")
	},
}

var clientOpt *ClientOption

func Execute() {
	clientOpt = LoadClientConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	cfgFile     string
	projectBase string
	userLicense string
	debug       bool
)

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".minibox" (without extension).
		// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		// dir, err := os.Getwd()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// viper.AddConfigPath(dir)
		viper.AddConfigPath(home)
		os.Mkdir(path.Join(home, ".minibox"), 0755)
		viper.SetConfigName(".minibox/config")
	}

	lang := message.MatchLanguage(userLocale, fallback)

	if err := viper.ReadInConfig(); err != nil {
		p := message.NewPrinter(lang)
		p.Println("Can't read config:", err)
		os.Exit(1)
	}
}
