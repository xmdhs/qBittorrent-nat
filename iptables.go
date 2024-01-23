package main

import (
	"fmt"
	"strconv"

	"github.com/coreos/go-iptables/iptables"
)

func Forward(publicPort, privatePort uint16) error {
	iptables4, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return fmt.Errorf("Forward: %w", err)
	}
	iptables6, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		return fmt.Errorf("Forward: %w", err)
	}
	err = iptables4.ClearChain("nat", "qbittorrent")
	if err != nil {
		return fmt.Errorf("Forward: %w", err)
	}

	privatePortStr := strconv.Itoa(int(privatePort))
	publicPortStr := strconv.Itoa(int(publicPort))

	for _, v := range []string{"tcp", "udp"} {
		err = iptables4.AppendUnique("nat", "qbittorrent", "-p", v, "--dport", privatePortStr, "-j", "REDIRECT", "--to-port", publicPortStr)
		if err != nil {
			return fmt.Errorf("Forward: %w", err)
		}
	}

	err = iptables4.DeleteIfExists("nat", "PREROUTING", "-j", "qbittorrent")
	if err != nil {
		return fmt.Errorf("Forward: %w", err)
	}

	err = iptables4.InsertUnique("nat", "PREROUTING", 1, "-j", "qbittorrent")
	if err != nil {
		return fmt.Errorf("Forward: %w", err)
	}

	for _, iptables := range []*iptables.IPTables{iptables4, iptables6} {
		err = iptables.ClearChain("filter", "qbittorrentfilter")
		if err != nil {
			return fmt.Errorf("Forward: %w", err)
		}

		for _, v := range []string{"tcp", "udp"} {
			err = iptables.AppendUnique("filter", "qbittorrentfilter", "-p", v, "--dport", publicPortStr, "-j", "ACCEPT")
			if err != nil {
				return fmt.Errorf("Forward: %w", err)
			}
		}

		err = iptables.DeleteIfExists("filter", "INPUT", "-j", "qbittorrentfilter")
		if err != nil {
			return fmt.Errorf("Forward: %w", err)
		}

		err = iptables.InsertUnique("filter", "INPUT", 1, "-j", "qbittorrentfilter")
		if err != nil {
			return fmt.Errorf("Forward: %w", err)
		}
	}

	return nil
}
