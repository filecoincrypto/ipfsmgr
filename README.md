# IPFS Manager 
- Author: Andy Zhou <ablozhou@gmail.com>
- Date: 2021-09-15

## Abstract
This is a golang module to manager IPFS file and directory.

It's a part of IPFS grant service for enterprice.

The main functions of the library is:
- Add files and directories to IPFS
- Get files and directories from IPFS

This module tested on Ubuntu Linux 2020

If you have any problem please contact the author.

## Env
- [golang 1.16+](https://golang.org/doc/install)
- [ipfs v0.9.1+](https://dist.ipfs.io/#go-ipfs)

## Getting started

get ipfsmgr mod

```
go get -u github.com/filecoincrypto/ipfsmgr
```

then import the mod  as normal.
```
import (
  m "github.com/filecoincrypto/ipfsmgr"
)
```

## Build from source
First you must add golang and ipfs CLI.

```
git clone https://github.com/filecoincrypto/ipfsmgr.git

go clean --modcache

go mod tidy

go build -mod=mod

go install
```

## Install ipfs on Linux
```
wget https://dist.ipfs.io/go-ipfs/v0.9.1/go-ipfs_v0.9.1_linux-amd64.tar.gz --no-check-certificate
tar -xvzf go-ipfs_v0.9.1_linux-amd64.tar.gz

> x go-ipfs/install.sh
> x go-ipfs/ipfs

cd go-ipfs
sudo bash install.sh

> Moved ./ipfs to /usr/local/bin
ipfs --version

> ipfs version 0.9.1
```

## Running a test

To run the test, just do:

```
> go test
```

# Trouble shooting
- missing go.sum entry for module providing package ...
  run `go build -mod=mod` will generate go.sum
- go-multiaddr-net@v0.2.0/registry.go:25:17: undefined: manet.NetCodec
  go-multiaddr v0.3.3 and not v0.4.0.
- libp2p/go-libp2p-noise@v0.2.0/handshake.go:209:21: cannot assign error to err in multiple assignment
  go-libp2p-core v0.8.6 and not v0.9.0
- go-ipfs@v0.9.1/core/coreapi/path.go:52:18: undefined: resolver.ResolveOnce 
  replace go-ipfs v0.9.1 to latest v0.10.0-rc1
- panic: failed to spawn ephemeral node: no IPFS repo found in /home/zhh/.ipfs.
  please install go-ipfs and run: 'ipfs init'. 

- failed to spawn ephemeral node: failed to init ephemeral node: unknown datastore type: flatfs
  run `ipfs init` first
- failed to spawn ephemeral node: no version file found, please run 0-to-1 migration tool.
See https://github.com/ipfs/fs-repo-migrations/blob/master/run.md
Sorry for the inconvenience. In the future, these will run automatically.
  not init repo directory correctly. run `ipfs init` first.

- failed to sufficiently increase receive buffer size
  run `sudo sysctl -w net.core.rmem_max=2500000`,
  This command would increase the maximum receive buffer size to roughly 2.5 MB

# Reference
- [install ipfs](https://docs.ipfs.io/install/)
- [download ipfs binary](https://dist.ipfs.io/#go-ipfs)
- [go-ipfs Core API](https://godoc.org/github.com/ipfs/interface-go-ipfs-core)
- [config a Node](https://docs.ipfs.io/how-to/configure-node/)