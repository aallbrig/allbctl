---
title: "allbctl_completion_zsh"
description: "auto-generated reference"
---

## allbctl completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(allbctl completion zsh)

To load completions for every new session, execute once:

#### Linux:

	allbctl completion zsh > "${fpath[1]}/_allbctl"

#### macOS:

	allbctl completion zsh > $(brew --prefix)/share/zsh/site-functions/_allbctl

You will need to start a new shell for this setup to take effect.


```
allbctl completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.allbctl.yaml)
      --debug           Enable debug telemetry (structured logs, traces, and metrics to stderr)
```

### SEE ALSO

* [allbctl completion](allbctl_completion.md)	 - Generate the autocompletion script for the specified shell
