package ipfsmgr

import (
	"log"
	"testing"
	// "os"
	// "path/filepath"
	// "sync"
	// config "github.com/ipfs/go-ipfs-config"
	// icore "github.com/ipfs/interface-go-ipfs-core"
	// ma "github.com/multiformats/go-multiaddr"
	// "github.com/ipfs/go-ipfs/core"
	// "github.com/ipfs/go-ipfs/core/coreapi"
	// "github.com/ipfs/go-ipfs/core/node/libp2p"
	// "github.com/ipfs/go-ipfs/plugin/loader"
	// "github.com/ipfs/go-ipfs/repo/fsrepo"
	// "github.com/libp2p/go-libp2p-core/peer"
)

var mgr *IpfsMgr

func TestMain(m *testing.M) {
	log.SetPrefix("INFO: ")
	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)
	log.Println("begin test")
	mgr = NewIpfsMgr("")
	m.Run()
	log.Println("end test")
}
func TestGetIpfsFile(t *testing.T) {
	log.Println("begin TestGetIpfsFile")
	cidPath := "/ipfs/QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps"
	err := mgr.GetIpfsFile(cidPath, "./tmp/ipfs.pdf")
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("got file back from IPFS %s  and wrote it to ./tmp/ipfs.pdf \n", cidPath)
	}

}
func TestGetIpfsDir(t *testing.T) {
	log.Println("begin TestGetIpfsDir")
	cidPath := "/ipfs/QmdQdu1fkaAUokmkfpWrmPHK78F9Eo9K2nnuWuizUjmhyn"
	err := mgr.GetIpfsDir(cidPath, "./tmp")
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("got file back from IPFS %s  and wrote it to ./tmp \n", cidPath)
	}

}
func TestAddIpfsFile(t *testing.T) {
	log.Println("begin TestAddIpfsFile")
	file := "ipfsmgr_test.go"
	cidFile, err := mgr.AddIpfsFile(file)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("upload ipfs cidFile:" + cidFile.String())
	}
}

func TestAddIpfsDir(t *testing.T) {
	log.Println("begin TestAddIpfsDir")
	path := "./"
	cidFile, err := mgr.AddIpfsDir(path)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("upload ipfs cidFile Dir:" + cidFile.String())
	}

}
