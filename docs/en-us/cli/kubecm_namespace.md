## kubecm namespace

Switch or change namespace interactively

### Synopsis


Switch or change namespace interactively


```
kubecm namespace [flags]
```

![ns](../../static/ns.gif)

### Examples

```

# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system

```

### Options

```
  -h, --help   help for namespace
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

