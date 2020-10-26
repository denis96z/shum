# shum

## What is _shum_?

_shum_ is a UNIX daemon that allows executing shell commands via HTTP API calls

**Attention!**

The tool's primary use case is non-production environment deployment utility.
It is NOT meant to be a production-ready totally-secure tool, as the main goal
is to make its use as easy as possible.

## Build

```
$ make
```

## Install

(TODO: provide packages)

```
# make install
```

## Quick start

1. Modify /etc/shum.yaml to look like this
```
server:
  addr: "127.0.0.1"
  port: 8065
auth:
  clients:
    - client_id: test_id
      client_secret: test_secret
shell:
  bin: "/bin/sh"
  args:
    - "-c"
  commands:
    uname:
      command: "uname -a"
      reveal_output: true
    deploy:
      command: "cd /root/src/docker"
      async: true
```

2. Execute
```
# systemctl restart shum
```

3. Check if _shum_ is healthy
```
# curl 'http://127.0.0.1:8065/status'
```

4. Execute command by calling _shum_ API
```
# curl 'http://127.0.0.1:8065/cmd/uname'
```
