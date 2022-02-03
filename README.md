# Proxy Wrapper

> working on it

This is command line wrapper utility to execute command with multiple proxy servers to simplify following use case.

- I need to run some command with proxy:

  ```bash
  export https_proxy=http://some-internal-proxy:3128
  export no_proxy=127.0.0.1,localhost,some-external
  kubectl get pods
  curl -IL https://some-website-behind-firewall-and-proxy

  # need to cut proxy to run some other command
  unset https_proxy
  gcloud auth login
  # or execute it in new terminal
  ```

- My organization has multi-cloud and each environment has it's own proxy:

  ```bash
  # run kubectl for cluster in AWS
  https_proxy=http://some-internal-proxy-A:3128 kubectl get pods

  # run kubectl for cluster in GCP Project A
  https_proxy=http://some-internal-proxy-G-A:443 kubectl get pods

  # run k9s for cluster in GCP Project B
  https_proxy=http://some-internal-proxy-G-B:443 k9s

  # run some curl for connectivity testing
  https_proxy=http://some-internal-proxy-G-C:443 curl -IL https://some-internal-host
  ```

- Typing/Copy and Paste `http_proxy` part for each env is annoying. so I decided to create `alias`:

  ```bash
  alias 'ak=https_proxy=http://some-internal-proxy-A:3128 kubectl'
  alias 'gka=https_proxy=http://some-internal-proxy-G-A:443 kubectl'
  alias 'gkb=https_proxy=http://some-internal-proxy-G-B:443 kubectl'
  ...
  alias 'gkz=https_proxy=http://some-internal-proxy-G-Z:443 kubectl'
  ```

  but, sometimes I forgot which one is for what and hard to naming the alias as it grows. 😮‍💨

- My company has multiple VPN servers for the infrastructures. switching VPN back and forth in the local is annoying. so I run some docker container to connect the VPN server and proxying request with Squid Proxy. it works great on browser with FoxyProxy. but still need to set environment variable for my terminal.

These are my actual daily use cases. 😉

## Installation

> working on it

```bash
brew install sunggun-yu/tap/prw
```

## Usage

This utility is simplifying setting the proxy environment variable portion by reading proxy server config for the corresponding profile from the configuration file.

```bash
# select which profile to use
prw use <some-proxy-profile-name>
# run command along with prw command
# type command after --
prw -- kubectl get pods
prw -- k9s
```

```bash
# specify the profile to use. --profile / -p
prw -p <some-proxy-profile-name> -- kubectl get pods
prw -p a-a -- k9s
prw -p a-b -- curl -IL https://some-host
prw -p g-a -- curl -IL https://some-host
prw -p g-b -- kubectx g-b && kubectl get pods
```

## Config file

config file will be created at `$HOME/.config/prw/config.yaml` if it is not existing when you run the `pwr`. also it will be updated by sub commands.
but, you can update the config file directly if you need bulk update.

```yaml
default: vpn-a
profiles:
  <profile-name>:
    desc: <description of profile>
    host: "http://<ip>:<port>"
    noproxy: <comma separated string e.g. 127.0.0.1,localhost>
  vpn-a:
    desc: squid proxy with vpn A
    host: http://192.168.3.3:3128
    noproxy: localhost,127.0.0.1,something
  vpn-b:
    desc: squid proxy with vpn B
    host: http://192.168.3.3:3228
```
