package ipfsmgr

import (
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

type IIpfsMgr interface {
	AddIpfsFile(inputPathFile string) (cidFile icorepath.Path, err error)
	AddIpfsDir(inputPath string) (cidPath icorepath.Path, err error)
	GetIpfsFile(cidPath string, outputPathFile string) (err error)
	GetIpfsDir(cidPath string, outputPath string) (err error)
}
