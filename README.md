# ENVP

![release](https://github.com/sunggun-yu/envp/actions/workflows/release.yaml/badge.svg)

ENVP is cli wrapper that sets environment variables by profile based configuration when you execute the command line.

## Use Cases I want to solve

These are my actual daily use cases. 😉

- I need to run some command via proxy

  ```bash
  export https_proxy=http://some-proxy:3128
  export no_proxy=127.0.0.1,localhost,some-domain
  kubectl get pods
  curl -IL https://some-website-behind-firewall-and-proxy

  # and sometimes I need to unset proxy to run some other command
  # or execute it in new terminal
  unset https_proxy
  gcloud auth login
  ```

- My team has multiple k8s cluster on different environments. but, should through dedicate proxy per cluster to do kubectl

  ```bash
  # run kubectl for cluster in infrastructure A
  https_proxy=http://some-internal-proxy-A:3128 kubectl get pods

  # run kubectl for cluster in infrastructure B
  https_proxy=http://some-internal-proxy-B:443 kubectl get pods
  ```

- Typing(or Copy and Paste) `http_proxy` part for each env is so annoying. so, created `alias`!

  ```bash
  alias 'ak=https_proxy=http://some-internal-proxy-A:3128 kubectl'
  alias 'gka=https_proxy=http://some-internal-proxy-B-A:443 kubectl'
  alias 'gkb=https_proxy=http://some-internal-proxy-B-B:443 kubectl'
  ...
  alias 'gkz=https_proxy=http://some-internal-proxy-B-Z:443 kubectl'
  ```

  it may works. but I want to simply do `kubectl` (or `k`) 😮‍💨

  also, it doesn't scale for other k8s related commands.
  e.g. `k9s`, `skaffold`, `helm`, and so on.

- My team has multiple VPN servers for the infrastructures. switching VPN back and forth in the local is annoying. so I run docker container that connects to the VPN server and proxying request by Squid Proxy. it works great on browser with FoxyProxy. but still need to set environment variable in my terminal.

- I have multiple servers to run docker remotely

  ```bash
  DOCKER_HOST=ssh://user@workstation-1 docker ps -a
  DOCKER_HOST=ssh://user@workstation-2 docker ps -a
  ```

and more cases like `VAULT_ADDR`, `ARGO_SERVER`, and so on and so on.

## Installation

```bash
brew install sunggun-yu/tap/envp
```

## Usage

> note: command line must start after double dash `--`.

```txt
ENVP is cli wrapper that sets environment variables by profile when you execute the command line

Usage:
  envp profile-name [flags] -- [command line to execute, e.g. kubectl]
  envp [command]

Examples:

  # run command with selected environment variable profile.
  # (example is assuming HTTPS_PROXY is set in the profile)
  envp use profile
  envp -- kubectl cluster-info
  envp -- kubectl get pods
  
  # specify env var profile to use
  envp profile-name -- kubectl get namespaces
  

Available Commands:
  add         Add environment variable profile
  completion  Generate the autocompletion script for the specified shell
  delete      Delete environment variable profile
  edit        Edit environment variable profile
  help        Help about any command
  list        List all environment variable profiles
  show        Print the environment variables of profile
  use         Set default environment variable profile
  version     Print the version of envp

Flags:
  -h, --help   help for envp

Use "envp [command] --help" for more information about a command.
```

### Add profile

```bash
envp add my-proxy \
  -d "profile description" \
  -e HTTPS_PROXY=http://some-proxy:3128 \
  -e "NO_PROXY=127.0.0.1,localhost" \
  -e "DOCKER_HOST=ssh://myuser@some-server"
```

- the value format of environment variable `--env`/`-e` is `NAME=VALUE`. other format will be ignored.
- you can add multiple environment variables by repeating `--env`/`-e`.
- added profile will be set to default profile if default is not set.

### Set default profile

```bash
envp use <profile-name>
```

### Run command line

run command with default profile:

```bash
# set command after double dash --
envp -- kubectl get pods
envp -- kubectl exec -it vault-test-app -- sh
envp -- k9s
envp -- vault login
envp -- docker ps -a
```

run command with specific profile:

```bash
# specify the profile to use. --profile / -p
envp <profile-name> -- kubectl get pods
envp a-a -- k9s
envp a-b -- curl -IL https://some-host
envp g-a -- curl -IL https://some-host
envp g-b -- kubectx g-b && kubectl get pods
envp my-lab -- docker ps -a
envp my-vault-1 -- vault login
```

### List profiles

```bash
envp list
envp ls
```

result:

```txt
  a-profile
* my-lab-1
  test
  vpn-a
  vpn-b
```

- default profile will be marked with `*`

### Show all environment variables of profile

print out default profile's env vars:

```bash
envp show 

ENV_VAR_1=ENV_VAL_1
ENV_VAR_2=ENV_VAL_2
ENV_VAR_3=ENV_VAL_3
```

show env vars of specific profile:

```bash
envp show some-profile

ENV_VAR_1=ENV_VAL_1
ENV_VAR_2=ENV_VAL_2
ENV_VAR_3=ENV_VAL_3
```

show with export option `--export`, `-e`

```bash
envp show --export

# you can export env vars of profile with following command
# eval $(envp show --export)
# eval $(envp show profile-name --export)

export ENV_VAR_1=ENV_VAL_1
export ENV_VAR_2=ENV_VAL_2
export ENV_VAR_3=ENV_VAL_3
```

so that, user can export env vars with `eval $(envp show --export)` command

```bash
eval $(envp show --export)
```

### Edit profile

```bash
envp edit my-proxy \
  -d "updated profile desc" \
  -e "NO_PROXY=127.0.0.1,localhost"
```

- value will be updated for existing env name
- removing env from profile will be added later. please update config file directly for now

### Delete profile

```bash
envp delete profile
envp del another-profile
```

### Nested profile

nested profile is possible natually thanks to `viper`.
you can simply divide group and profile by `.` in the profile name.

```bash
envp add group.profile
envp use group.profile
envp group.profile -- ls -la
envp delete group.profile
envp delete group
```

- if you delete parent profile with delete command, it will also delete all child profiles.

## Config file

config file will be created at `$HOME/.config/envp/config.yaml` initially. and all profiles and environment variables will be stored in this file.

the file format is `YAML` and followed k8s pod env format to avoid map key uncapitalize issue from go/yaml unmarshal.

```yaml
default: vpn-a
profiles:
  <profile-name>:
    env:
      - name: env-name
        value: env-value
  my-lab-1:
    desc: my lab cluster 1
    env:
      - name: DOCKER_HOST
        value: ssh://user@workstation-1
      - name: VAULT_ADDR
        value: https://vault.mylab-1
      - name: ARGO_SERVER
        value: https://argocd.mylab-1
  vpn: # profile group
    pa:
      desc: squid proxy with vpn A
      env:
        - name: HTTP_PROXY
          value: http://192.168.3.3:3128
        - name: NO_PROXY
          value: localhost,127.0.0.1,something
    pb:
      desc: squid proxy with vpn b
      env:
        - name: HTTP_PROXY
          value: http://192.168.3.3:3228
        - name: NO_PROXY
          value: localhost,127.0.0.1,something
```
