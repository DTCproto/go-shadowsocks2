package main

import (
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DTCproto/go-shadowsocks2/core"
	"github.com/DTCproto/go-shadowsocks2/files/filestream"
)

var config struct {
	Verbose        bool
	ConfigFilePath string
	UDPTimeout     time.Duration
}

var addr struct {
	TcpAddr string
	UdpAddr string
}

func main() {

	flag.StringVar(&config.ConfigFilePath, "config-file-path", "./config.json", "config file path(./config.json)")
	flag.Parse()

	flags := filestream.ConfigFilePath{
		Verbose:    false,
		UDPTimeout: 5 * time.Minute,
	}

	if config.ConfigFilePath != "" {
		err := flags.ParseConfigFile(config.ConfigFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	config.Verbose = flags.Verbose
	config.UDPTimeout = flags.UDPTimeout

	addr.TcpAddr = ":" + flags.ListenPort
	if flags.Server != "" {
		tcpAddrTmp, cipherTmp, passwordTmp, err := parseServerConfig(flags.Server)
		if err != nil {
			log.Fatal(err)
		}
		addr.TcpAddr = tcpAddrTmp
		flags.Cipher = cipherTmp
		flags.Password = passwordTmp
	}
	addr.UdpAddr = addr.TcpAddr

	if flags.Plugin != "" {
		tcpAddrTmp, err := startPlugin(flags.Plugin, flags.PluginOpts, addr.TcpAddr, true)
		if err != nil {
			log.Fatal(err)
		}
		addr.TcpAddr = tcpAddrTmp
	}

	// []byte("")  capacity->32,length->0
	cipher, err := core.PickCipher(flags.Cipher, []byte(""), flags.Password)
	if err != nil {
		log.Fatal(err)
	}
	go udpRemote(addr.UdpAddr, cipher.PacketConn)
	go tcpRemote(addr.TcpAddr, cipher.StreamConn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	killPlugin()
}

// URL方式[不推荐]
func parseServerConfig(serverUrl string) (addr, cipher, password string, err error) {
	if strings.HasPrefix(serverUrl, "ss://") {
		var u *url.URL
		u, err = url.Parse(serverUrl)
		if err != nil {
			return
		}
		addr = u.Host
		if u.User != nil {
			cipher = u.User.Username()
			password, _ = u.User.Password()
		}
	} else {
		err = errors.New("JSON FILE ERROR: Server format is incorrect")
	}
	return
}
