package ipfsmgr

import (
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

type IpfsMgrInterface interface {
	AddIpfsFile(inputPathFile string) (cidFile icorepath.Path, err error)
	AddIpfsDir(inputPath string) (cidPath icorepath.Path, err error)
	GetIpfsFile(inputPath string) error
	GetIpfsDir(cidPath string, outputPath string) (err error)
}
