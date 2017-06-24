package sftp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"

	"golang.org/x/crypto/ssh"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
	"github.com/pkg/sftp"
)

const (
	FlagFile int = 1 << iota
	FlagDir
)

// Client
type Client struct {
	sftp.Client
}

func loadPrivateKey(privateKeyPath string) (sss ssh.Signer, err error) {
	pemBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return
	}
	sss, err = ssh.ParsePrivateKey(pemBytes)
	return
}

// BuildSSHURL builds up a SSH URL to feed ssh
// @host: hostname
// @port: port
func BuildSSHURL(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

// PrivateKeyAuthMethod returns a ssh.AuthMethod initialized with a private key
// @privateKeyPath: path to the private key to use
func PrivateKeyAuthMethod(privateKeyPath string) ssh.AuthMethod {
	sss, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(sss)

}

// New returns a new SFTP Client
// @host: hostname to connect to
// @port: port on the hostname
// @username: username to login
// @sams: list of ssh.AuthMethod to use
func New(host, port, username string, sams ...ssh.AuthMethod) (*Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: sams,
	}
	client, err := ssh.Dial("tcp", BuildSSHURL(host, port), config)
	if err != nil {
		return nil, err
	}
	tmpSftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	return (*Client)(unsafe.Pointer(tmpSftpClient)), err
}

// ResolveSymlink resolve a path and returns the path of the file pointed if
// it is a symlink
// @path: path to resolve
func (sc *Client) ResolveSymlink(path string) string {
	pointer, err := sc.ReadLink(path)
	if err != nil {
		return path
	}
	return pointer
}

// Walk walks recursively through the SFTP
// @root: root path to start walking through
func (sc *Client) Walk(root string) <-chan fswalker.WalkItem {
	iterChannel := make(chan fswalker.WalkItem)
	dirsToProcess := []string{root}
	go func() {
		for len(dirsToProcess) > 0 {
			dirs, files := []os.FileInfo{}, []os.FileInfo{}
			dirpath := dirsToProcess[len(dirsToProcess)-1]
			dirsToProcess = dirsToProcess[:len(dirsToProcess)-1]
			filesInfo, err := sc.ReadDir(sc.ResolveSymlink(dirpath))
			if err != nil {
				fmt.Printf("Error reading directory (%s): %s\n", err.Error(), dirpath)
			} else {
				for _, fileInfo := range filesInfo {
					switch {
					case fileInfo.Mode().IsDir():
						dirs = append(dirs, fileInfo)
						dirsToProcess = append(dirsToProcess, filepath.Join(dirpath, fileInfo.Name()))
					case fileInfo.Mode().IsRegular():
						files = append(files, fileInfo)
					case fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink:
						sympath := filepath.Join(dirpath, fileInfo.Name())
						pointerFI, err := sc.Stat(sympath)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error reading symlink (%s): %s\n", err.Error(), sympath)
						} else {
							switch {
							case pointerFI.Mode().IsDir():
								dirs = append(dirs, fileInfo)
								dirsToProcess = append(dirsToProcess, sc.Join(dirpath, fileInfo.Name()))
							case pointerFI.Mode().IsRegular():
								files = append(files, fileInfo)
							}
						}
					}
				}
			}
			iterChannel <- fswalker.WalkItem{Dirpath: dirpath,
				Dirs:  dirs,
				Files: files,
				Err:   err}
		}
		close(iterChannel)
	}()
	return iterChannel
}

func matchPatterns(str string, cPatterns []*regexp.Regexp) bool {
	for _, cPattern := range cPatterns {
		if cPattern.MatchString(str) {
			return true
		}
	}
	return false
}

func filter(osfi *[]os.FileInfo, cPatterns []*regexp.Regexp) {
	for i := 0; i < len(*osfi); {
		if matchPatterns((*osfi)[i].Name(), cPatterns) == false {
			*osfi = append((*osfi)[:i], (*osfi)[i+1:]...)
			continue
		}
		i++
	}
}

func (sc *Client) Find(root string, flag int, patterns ...string) (iterChannel chan fswalker.WalkItem) {
	var cPatterns []*regexp.Regexp
	iterChannel = make(chan fswalker.WalkItem)
	for _, pattern := range patterns {
		pattern = strings.TrimLeft(pattern, "^")
		pattern = strings.TrimRight(pattern, "$")
		cPattern, err := regexp.Compile(fmt.Sprintf("^%s$", pattern))
		if err == nil {
			cPatterns = append(cPatterns, cPattern)
		}
	}
	go func() {
		for wi := range sc.Walk(root) {
			if wi.Err != nil {
				iterChannel <- wi
			} else {
				switch flag {
				case FlagFile:
					wi.Dirs = []os.FileInfo{}
					filter(&wi.Files, cPatterns)
				case FlagDir:
					wi.Files = []os.FileInfo{}
					filter(&wi.Dirs, cPatterns)
				case FlagFile | FlagDir:
					filter(&wi.Dirs, cPatterns)
					filter(&wi.Files, cPatterns)
				}
				iterChannel <- wi
			}
		}
		close(iterChannel)
	}()
	return
}
