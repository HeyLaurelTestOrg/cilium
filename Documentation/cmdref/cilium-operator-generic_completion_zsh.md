<!-- This file was autogenerated via cilium-operator --cmdref, do not edit manually-->

## cilium-operator-generic completion zsh

generate the autocompletion script for zsh

### Synopsis


Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for every new session, execute once:
# Linux:
$ cilium-operator-generic completion zsh > "${fpath[1]}/_cilium-operator-generic"
# macOS:
$ cilium-operator-generic completion zsh > /usr/local/share/zsh/site-functions/_cilium-operator-generic

You will need to start a new shell for this setup to take effect.


```
cilium-operator-generic completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [cilium-operator-generic completion](cilium-operator-generic_completion.html)	 - generate the autocompletion script for the specified shell

