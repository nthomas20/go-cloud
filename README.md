![free logo from freelogodesign.org](./img/56928cea69e34e27b2eb76d4eabf81a1.png "go-cloud logo")

# Introduction
go-cloud aims to provide cloud functionality using established open protocols. Inspiration for this project started with webdav protocol and discovering [this gist](https://gist.github.com/darcyliu/336f4b0dd573cda2f5df339a74db0446) -- thanks [darcyliu](https://github.com/darcyliu/)!

Currently webdav functionality is provided via google's [webdav module](https://pkg.go.dev/golang.org/x/net/webdav), in the future a native module may be created #dreams

go-cloud does not strive to provide TLS support. There are already 1,000,000,000 (slightly exaggerated) other ways of providing this. go-cloud should be run behind a reverse proxy (such as haproxy) where TLS will terminate and forward requests on to go-cloud. You should use [certbot](https://certbot.eff.org/).

# Installation

## Requirements

At this point, this project has only been tested on Ubuntu-based 21.04 distributions

* go version 1.16+

## Build

After cloning the repository, execute the following to build

```bash
$ make build
```

# Execution
The goal with this project is to not require editing any configuration file by hand. Additionally, restarting the service should be required only in a small handful of configuration changes (e.g. port change)

When launching the application, it will fork a process and begin executing in the background. Configuration changes will be refreshed every 60 seconds.

## Starting the service

Assuming executing in the root directory of the repository after build

```bash
$ bin/go-cloud start
```

## Stopping the service

Assuming executing in the root directory of the repository after build

```bash
$ bin/go-cloud stop
```

## CLI Configuration

Account configurations may be made via CLI commands in real-time and will be refreshed every 60 seconds.

Replace `${...}` with their appropriate values

The configuration file is stored in `config.yaml`. If this file exists in the current directory it will be utilized, otherwise the file will be located at `~/.go-cloud/config.yaml`

### Adding an account

```bash
$ bin/go-cloud account add --username ${username} --email ${email} --directory ${webdav_directory}
```

### Deleting an account

```bash
$ bin/go-cloud account delete --username ${username}
```

### Adding a password to an account

```bash
$ bin/go-cloud password add --username ${username} --password ${password} --description ${description}
```

### Deleting a password from an account

```bash
$ bin/go-cloud password delete --username ${username} --password ${password}
```



# webdav

The webdav protocol provides file and directory services. The endpoint can be easily mounted in most operating systems. The default port `8080` is used at this time. Currently there is no configuration option to change this port (it's coming!) so you'll need to edit the `config.yaml` file to change it.

There are 2 webdav endpoints supported with no difference between them. If we assume `localhost`, the first mount point is:
* `localhost:8080/webdav/`

The second mount point reflects nextcloud's webdav interface at:
* `localhost:8080/remote.php/dav/files/${USERNAME}/`

In the case of nextcloud endpoint compatibility `${USERNAME}` must match the logged in user.