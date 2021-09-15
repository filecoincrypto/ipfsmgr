//package ipfsmgr manager ipfs file, implement upload and download files,connect peers.
//Author : Andy Zhou <ablozhou@gmail.com>
//Date   : 2021.9.15
package ipfsmgr

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/libp2p/go-libp2p-core/peer"
)

//type IpfsMgr, create repo, collect to peers, add and get ipfs files.
type IpfsMgr struct {
	ctx    context.Context
	cancel context.CancelFunc
	ipfs   icore.CoreAPI
}

// NewIpfsMgr new instance of IpfsMgr
// repoPath could be empty, default is ~/.ipfs
func NewIpfsMgr(repoPath string) *IpfsMgr {
	mgr := new(IpfsMgr)
	ctx, cancel := context.WithCancel(context.Background())
	mgr.ctx = ctx
	mgr.cancel = cancel
	ipfs, err := mgr.CreateNode(repoPath)
	if err != nil {
		panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	}
	mgr.ipfs = ipfs
	return mgr
}

// AddIpfsFile add local file to IPFS
func (mgr *IpfsMgr) AddIpfsFile(inputPathFile string) (cidFile icorepath.Path, err error) {

	defer mgr.cancel()
	file, err := mgr.GetLocalNode(inputPathFile)
	if err != nil {
		err = fmt.Errorf("could not get File: %s", err)
		return nil, err
	}
	cidFile, err = mgr.ipfs.Unixfs().Add(mgr.ctx, file)
	if err != nil {
		err = fmt.Errorf("could not add File: %s", err)
		return nil, err
	}
	return cidFile, err
}

// AddIpfsDir
func (mgr *IpfsMgr) AddIpfsDir(inputPath string) (cidPath icorepath.Path, err error) {
	directory, err := mgr.GetLocalNode(inputPath)
	if err != nil {
		err = fmt.Errorf("could not get File: %s", err)
		return nil, err
	}

	cidDirectory, err := mgr.ipfs.Unixfs().Add(mgr.ctx, directory)
	if err != nil {
		err = fmt.Errorf("could not add Directory: %s", err)
		return nil, err
	}

	log.Printf("added directory to IPFS with CID %s\n", cidDirectory.String())
	return cidDirectory, nil
}

// GetIpfsFile from cidPath string.
func (mgr *IpfsMgr) GetIpfsFile(cidPath string, outputPathFile string) (err error) {
	var cidFile icorepath.Path = icorepath.New(cidPath)
	return mgr.GetIpfsFileFromCid(cidFile, outputPathFile)

}

// GetIpfsFile from cidFile from github.com/ipfs/interface-go-ipfs-core/path.Path
func (mgr *IpfsMgr) GetIpfsFileFromCid(cidFile icorepath.Path, outputPathFile string) (err error) {

	rootNodeFile, err := mgr.ipfs.Unixfs().Get(mgr.ctx, cidFile)
	if err != nil {
		err = fmt.Errorf("could not get file with CID: %s", err)
		return

	}

	err = files.WriteTo(rootNodeFile, outputPathFile)
	if err != nil {
		err = fmt.Errorf("could not write out the fetched CID: %s", err)
		return
	}
	return nil
}

// GetIpfsDir from cidPath string
func (mgr *IpfsMgr) GetIpfsDir(cidPath string, outputPath string) (err error) {
	var cidDirectory icorepath.Path = icorepath.New(cidPath)
	return mgr.GetIpfsFileFromCid(cidDirectory, outputPath)
}

// GetIpfsDirFromCid from cidDirectory from github.com/ipfs/interface-go-ipfs-core/path.Path
func (mgr *IpfsMgr) GetIpfsDirFromCid(cidDirectory icorepath.Path, outputPath string) (err error) {
	rootNodeDirectory, err := mgr.ipfs.Unixfs().Get(mgr.ctx, cidDirectory)
	if err != nil {
		err = fmt.Errorf("could not get file with CID: %s", err)
		return
	}

	err = files.WriteTo(rootNodeDirectory, outputPath)
	if err != nil {
		err = fmt.Errorf("could not write out the fetched CID: %s", err)
		return
	}

	fmt.Printf("Got directory back from IPFS (IPFS path: %s) and wrote it to %s\n", cidDirectory.String(), outputPath)

	return nil
}

//CreateRepo
func (mgr *IpfsMgr) CreateRepo(repoPath string) error {

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return err
	}

	// When creating the repository, you can define custom settings on the repository, such as enabling experimental
	// features (See experimental-features.md) or customizing the gateway endpoint.
	// To do such things, you should modify the variable `cfg`. For example:
	if *flagExp {
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-filestore
		cfg.Experimental.FilestoreEnabled = true
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-urlstore
		cfg.Experimental.UrlstoreEnabled = true
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#directory-sharding--hamt
		cfg.Experimental.ShardingEnabled = true
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#ipfs-p2p
		cfg.Experimental.Libp2pStreamMounting = true
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#p2p-http-proxy
		cfg.Experimental.P2pHttpProxy = true
		// https://github.com/ipfs/go-ipfs/blob/master/docs/experimental-features.md#strategic-providing
		cfg.Experimental.StrategicProviding = true
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return nil
}

// CreateNode Creates an IPFS node and returns its coreAPI
// param repoPath could be empty and set to default ~/.ipfs
// MUST use `ipfs init` to init ~/.ipfs as a repo
func (mgr *IpfsMgr) CreateNode(repoPath string) (api icore.CoreAPI, err error) {
	if repoPath == "" {
		repoPath, err = config.PathRoot()
		if err != nil {
			// shouldn't be possible
			return nil, err
		}

	}
	if err = mgr.CreateRepo(repoPath); err != nil {
		return nil, err
	}

	if err = mgr.SetupPlugins(repoPath); err != nil {
		return nil, err

	}

	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(mgr.ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func (mgr *IpfsMgr) SetupPlugins(pluginsPath string) error {
	// Load any external plugins if available on pluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(pluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

// ConnectToPeers connect peers, there are default bootstrap peers.
// param peers could be nil or empty.
func (mgr *IpfsMgr) ConnectToPeers(peers []string) error {
	var wg sync.WaitGroup

	//process nil peers
	if peers == nil {
		peers = []string{}
	}

	_ = append(bootstrapNodes, peers...)

	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(bootstrapNodes))
	for _, addrStr := range bootstrapNodes {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := mgr.ipfs.Swarm().Connect(mgr.ctx, *peerInfo)
			if err != nil {
				log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
			}
		}(peerInfo)
	}
	wg.Wait()
	return nil
}

func (mgr *IpfsMgr) GetLocalFile(path string) (files.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		return nil, err
	}

	f, err := files.NewReaderPathFile(path, file, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

//GetLocalNode
func (mgr *IpfsMgr) GetLocalNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

/// -------
var bootstrapNodes = []string{
	// IPFS Bootstrapper nodes.
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

	// IPFS Cluster Pinning nodes
	"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",

	// You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
	// "/ip4/127.0.0.1/tcp/4010/p2p/xxx",
	// "/ip4/127.0.0.1/udp/4010/quic/p2p/yyy",
}

/// -------

var flagExp = flag.Bool("experimental", false, "enable experimental features")
