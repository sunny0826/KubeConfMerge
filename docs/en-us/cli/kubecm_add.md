## kubecm add

Add KubeConfig to $HOME/.kube/config

### Synopsis

Add KubeConfig to $HOME/.kube/config

```
kubecm add [flags]
```

>Note: If `-c` is set and **more than one** context is added to the kubeconfig file, the following will occur:
>- If `--context-name` is set, the context will be generated as `<context-name-0>`, `<context-name-1>` ...
>- If `--context-name` is not set, it will be generated as `<file-name-{hash}>` where `{hash}` is the MD5 hash of the file name.

### Examples

```

# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Merge test.yaml with $HOME/.kube/config and rename context name
kubecm add -cf test.yaml --context-name test
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -

```

### Options

```
      --context-name string   override context name when add kubeconfig context
  -c, --cover         Overwrite local kubeconfig files
  -f, --file string   Path to merge kubeconfig files
  
  -h, --help          help for add
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

