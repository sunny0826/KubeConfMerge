package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// RenameCommand rename cmd struce
type RenameCommand struct {
	BaseCommand
}

// Init RenameCommand
func (rc *RenameCommand) Init() {
	rc.command = &cobra.Command{
		Use:     "rename",
		Short:   "Rename the contexts of kubeconfig",
		Long:    "Rename the contexts of kubeconfig",
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return rc.runRename(cmd, args)
		},
		Example: renameExample(),
	}
}

func (rc *RenameCommand) runRename(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	var kubeItems []Needle
	for key, obj := range config.Contexts {
		if key != config.CurrentContext {
			kubeItems = append(kubeItems, Needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]Needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	// exit option
	kubeItems, err = ExitOption(kubeItems)
	if err != nil {
		return err
	}
	num := SelectUI(kubeItems, "Select The Rename Kube Context")
	kubeName := kubeItems[num].Name
	rename := PromptUI("Rename", kubeName)
	config, err = renameComplet(rename, kubeName, config)
	if err != nil {
		return err
	}
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	return MacNotifier(fmt.Sprintf("Rename [%s] to [%s]\n", kubeName, rename))
}

func renameComplet(rename, kubeName string, config *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	if _, ok := config.Contexts[rename]; ok || rename == kubeName {
		return nil, errors.New("Name: " + rename + " already exists")
	}
	if obj, ok := config.Contexts[kubeName]; ok {
		config.Contexts[rename] = obj
		delete(config.Contexts, kubeName)
		if config.CurrentContext == kubeName {
			config.CurrentContext = rename
		}
	}
	return config, nil
}

func renameExample() string {
	return `
# Renamed the context interactively
kubecm rename
`
}
