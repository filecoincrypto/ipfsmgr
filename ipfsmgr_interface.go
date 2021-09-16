// package ipfsmgr manager ipfs file, implemented uploading and downloading
// files and directories,connecting peers.
// Author : Andy Zhou <ablozhou@gmail.com>
// Date   : 2021.9.15

package ipfsmgr

import (
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// type IIpfsMgr is an interface of type IpfsMgr
type IIpfsMgr interface {
	// AddIpfsFile upload local file to IPFS
	AddIpfsFile(inputPathFile string) (cidFile icorepath.Path, err error)
	// AddIpfsDir upload local directory to IPFS
	AddIpfsDir(inputPath string) (cidPath icorepath.Path, err error)
	// GetIpfsFile download IPFS cid file to outputPathFile
	GetIpfsFile(cidPath string, outputPathFile string) (err error)
	// GetIpfsDir download IPFS cid path to outputPath
	GetIpfsDir(cidPath string, outputPath string) (err error)
}
