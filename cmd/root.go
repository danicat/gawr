/*
Copyright © 2023 Daniela Petruzalek

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"gawr/crawler"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gawr [website]",
	Example: "gawr -f 1 -m 10 https://example.com",
	Short:   "A breadth-first search (BFS) web crawler",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Duration(viper.GetInt("frequency")) * time.Second
		freq := rate.Every(t)

		c, err := crawler.NewCrawler(args[0], freq, 1)
		if err != nil {
			log.Printf("error creating crawler: %s", err)
			return
		}

		c.MaxVisits = viper.GetInt("max-visits")
		c.VisitFn = func(u url.URL, content string) {
			links, err := crawler.ExtractLinks(content)
			if err != nil {
				log.Println(err)
				return
			}

			fmt.Println("Visited: " + u.String())
			fmt.Println("Found:")
			for _, l := range links {
				fmt.Printf("\t%s\n", l.String())
			}
		}

		c.FilterFn = func(u url.URL) bool {
			return strings.HasPrefix(u.String(), args[0])
		}

		err = c.Crawl()
		if err != nil {
			log.Printf("error crawling: %s", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gawr.yaml)")
	rootCmd.PersistentFlags().IntP("max-visits", "m", 0, "Maximum number of links to visit (0 = disabled)")
	rootCmd.PersistentFlags().IntP("frequency", "f", 10, "Frequency in seconds. e.g. 10 means sending one crawling request every 10 seconds.")
	viper.BindPFlag("max-visits", rootCmd.PersistentFlags().Lookup("max-visits"))
	viper.BindPFlag("frequency", rootCmd.PersistentFlags().Lookup("frequency"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gawr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gawr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
