// package ipfsmgr manager ipfs file, implemented uploading and downloading
// files and directories,connecting peers.
// Author : Andy Zhou <ablozhou@gmail.com>
// Date   : 2021.9.15

package ipfsmgr

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

var mgr IIpfsMgr
var outputBasePath string

func TestMain(m *testing.M) {
	log.SetPrefix("INFO: ")
	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)
	log.Println("begin test")
	var err error
	outputBasePath, err = ioutil.TempDir("./", "tmp")
	if err != nil {
		panic(fmt.Errorf("could not create output dir (%v)", err))
	}
	fmt.Printf("output folder: %s\n", outputBasePath)
	mgr = NewIpfsMgr("")
	// init twice
	mgr = NewIpfsMgr("")
	m.Run()
	log.Println("end test")
}

func TestAddIpfsFile(t *testing.T) {
	log.Println("begin TestAddIpfsFile")
	file := "testdata/ipfs.svg"
	cidFile, err := mgr.AddIpfsFile(file)
	if err != nil {
		panic(fmt.Errorf("TestAddIpfsFile  %s failed! (%v)", file, err))
	} else {
		log.Println("upload ipfs cidFile:" + cidFile.String())
	}

	// test add twice
	cidFile1, err := mgr.AddIpfsFile(file)
	if err != nil {
		panic(fmt.Errorf("TestAddIpfsFile 2 %s failed! (%v)", file, err))
	} else {
		log.Println("upload ipfs2 cidFile:" + cidFile1.String())
	}
}

func TestAddIpfsDir(t *testing.T) {
	log.Println("begin TestAddIpfsDir")
	path := "testdata"
	cidFile, err := mgr.AddIpfsDir(path)
	if err != nil {
		panic(fmt.Errorf("TestAddIpfsDir  %s failed!(%v)", path, err))
	} else {
		log.Println("upload ipfs Dir cid:" + cidFile.String())
	}

}

func TestGetIpfsFile(t *testing.T) {
	log.Println("begin TestGetIpfsFile")
	cidPath := "/ipfs/QmSUNuGugzSoaeMZuoFSJZyViXXv1YUcjwzRbK4ARppqxx"
	err := mgr.GetIpfsFile(cidPath, outputBasePath+"/ipfs.svg")
	if err != nil {
		panic(fmt.Errorf("TestGetIpfsFile  %s failed!(%v)", cidPath, err))
	} else {
		log.Printf("got file back from IPFS %s  and wrote it to %s/ipfs.svg \n", cidPath, outputBasePath)
	}

}
func TestGetIpfsDir(t *testing.T) {
	log.Println("begin TestGetIpfsDir")
	cidPath := "/ipfs/QmfWkDdAuvVmuF8Kuo5wAj1MtXZNkq7nJ4Xvjpae59AXXX"
	err := mgr.GetIpfsDir(cidPath, outputBasePath+"/ipfsdir/")
	if err != nil {
		panic(fmt.Errorf("TestGetIpfsDir  %s failed!(%v)", cidPath, err))
	} else {
		log.Printf("got file back from IPFS %s  and wrote it to %s/ipfsdir \n", cidPath, outputBasePath)
	}

}
