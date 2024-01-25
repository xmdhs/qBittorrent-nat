package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/xmddhs/qBittorrent-nat/qbittorrent"
	"github.com/xmdhs/natupnp/natmap"
)

var (
	configPath string
	stun       string
)

func init() {
	flag.StringVar(&configPath, "c", "config.json", "")
	flag.StringVar(&stun, "s", "stunserver.stunprotocol.org:3478", "stun")
	flag.Parse()
}

func main() {
	c := config{}
	b := lo.Must(os.ReadFile(configPath))
	lo.Must0(json.Unmarshal(b, &c))

	ctx := context.Background()

	for {
		func() {
			defer time.Sleep(200 * time.Millisecond)
			err := openPort(ctx, func(pub, pri netip.AddrPort) {
				if pub.Port() == pri.Port() {
					log.Println("公网，无需转发")
					return
				}
				err := Forward(pub.Port(), pri.Port())
				if err != nil {
					log.Println(err)
					return
				}
				q, err := qbittorrent.Login(ctx, c.Root, http.Client{Timeout: 10 * time.Second}, c.UserName, c.PassWord)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("change qBittorrent port to: %v", pub.Port())
				err = q.ChangePort(ctx, pub.Port())
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("%v --> %v --> 127.0.0.1:%v", pub, pri, pub.Port())
			})
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}

}

func openPort(ctx context.Context, finish func(pub, pri netip.AddrPort)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	la, err := natmap.GetLocalAddr()
	if err != nil {
		return fmt.Errorf("openPort: %w", err)
	}

	laddr := netip.MustParseAddrPort(la.String())

	errCh := make(chan error, 1)

	m, maddr, err := natmap.NatMap(ctx, stun, laddr, func(err error) {
		cancel()
		select {
		case errCh <- err:
		default:
		}
	})
	if err != nil {
		return fmt.Errorf("openPort: %w", err)
	}
	defer m.Close()

	finish(maddr, laddr)

	err = <-errCh
	if err != nil {
		return fmt.Errorf("openPort: %w", err)
	}
	return nil

}
