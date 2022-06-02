package cmd

/*
Copyright © 2020 Guo Xudong

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	uiSize  int
)

// Cli cmd struct
type Cli struct {
	rootCmd *cobra.Command
}

//NewCli returns the cli instance used to register and execute command
func NewCli() *Cli {
	cli := &Cli{
		rootCmd: &cobra.Command{
			Use:   "kubecm",
			Short: "KubeConfig Manager.",
			Long:  printLogo(),
		},
	}
	cli.rootCmd.SetOut(os.Stdout)
	cli.rootCmd.SetErr(os.Stderr)
	cli.setFlags()
	return cli
}

func (cli *Cli) setFlags() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flags := cli.rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", *kubeconfig, "path of kubeconfig")
	flags.IntVar(&uiSize, "ui-size", 4, "number of list items to show in menu at once")
}

// Run command
func (cli *Cli) Run() error {
	// check and format kubeconfig path
	config, err := CheckAndTransformFilePath(cfgFile)
	if err != nil {
		return err
	}
	err = flag.Set("config", config)
	if err != nil {
		return err
	}
	return cli.rootCmd.Execute()
}
func homeDir() string {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir
	}
	// cross compile support
	if runtime.GOOS == "windows" {
		return homeWindows()
	}
	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return ""
	}
	result := strings.TrimSpace(stdout.String())
	if result == "" {
		fmt.Println("blank output when reading home directory")
		os.Exit(0)
	}

	return result
}

func homeWindows() string {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		fmt.Println("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
		os.Exit(0)
	}

	return home
}
func printLogo() string {
	panel := pterm.DefaultHeader.WithMargin(8).
		WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).Sprint("Manage your kubeconfig more easily.")
	// 	s, _ := pterm.DefaultBigText.WithLetters(
	//	pterm.NewLettersFromStringWithStyle("kube", pterm.NewStyle(pterm.FgLightGreen)),
	//	pterm.NewLettersFromStringWithStyle("cm", pterm.NewStyle(pterm.FgLightBlue))).Srender()
	logo := pterm.FgLightGreen.Sprint(`
██   ██ ██    ██ ██████  ███████  ██████ ███    ███ 
██  ██  ██    ██ ██   ██ ██      ██      ████  ████ 
█████   ██    ██ ██████  █████   ██      ██ ████ ██ 
██  ██  ██    ██ ██   ██ ██      ██      ██  ██  ██ 
██   ██  ██████  ██████  ███████  ██████ ██      ██
`)
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "Tips",
		Style: pterm.NewStyle(pterm.BgBlue, pterm.FgLightWhite),
	}
	url := pterm.Info.Sprintf("Find more information at: %s", pterm.LightMagenta("https://kubecm.cloud"))
	return fmt.Sprintf(`
%s%s
%s
`, panel, logo, url)
}
