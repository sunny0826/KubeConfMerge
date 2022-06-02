package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// NamespaceCommand namespace cmd struct
type NamespaceCommand struct {
	BaseCommand
}

// Init NamespaceCommand
func (nc *NamespaceCommand) Init() {
	nc.command = &cobra.Command{
		Use:   "namespace",
		Short: "Switch or change namespace interactively",
		Long: `
Switch or change namespace interactively
`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"ns"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nc.runNamespace(cmd, args)
		},
		Example: namespaceExample(),
	}
}

func (nc *NamespaceCommand) runNamespace(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	currentContext := config.CurrentContext
	contNs := config.Contexts[currentContext].Namespace
	namespaceList, err := GetNamespaceList(contNs)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		// exit option
		namespaceList = append(namespaceList, Namespaces{Name: "<Exit>", Default: false})
		num := selectNamespace(namespaceList)
		config.Contexts[currentContext].Namespace = namespaceList[num].Name
	} else {
		err := changeNamespace(args, namespaceList, currentContext, config)
		if err != nil {
			return err
		}
	}
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	return nil
}

func changeNamespace(args []string, namespaceList []Namespaces, currentContext string, config *clientcmdapi.Config) error {
	for _, ns := range namespaceList {
		if ns.Name == args[0] {
			config.Contexts[currentContext].Namespace = args[0]
			fmt.Printf("Namespace: 「%s」 is selected.\n", args[0])
			return nil
		}
	}
	return errors.New("Can not find namespace: " + args[0])
}

func selectNamespace(namespaces []Namespaces) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F6A9 {{if .Default}} {{ .Name | red }} * {{else}} {{ .Name | red }} {{end}}",
		Inactive: "{{if .Default}} {{ .Name | cyan }} * {{else}} {{ .Name | cyan }} {{end}}",
		Selected: "\U0001F680" + `{{if ne .Name "<Exit>" }}  Namespace: {{ .Name | green }} is selected.{{end}}`,
	}
	searcher := func(input string, index int) bool {
		pepper := namespaces[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		if input == "q" && name == "<exit>" {
			return true
		}
		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:     "Select Namespace:",
		Items:     namespaces,
		Templates: templates,
		Size:      uiSize,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	if namespaces[i].Name == "<Exit>" {
		fmt.Println("Exited.")
		os.Exit(1)
	}
	return i
}

func namespaceExample() string {
	return `
# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system
`
}
