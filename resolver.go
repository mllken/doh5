// DoH resolver implementation.
// TODO: test what happens when the url is bad

package main

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

const (
	cloudflareURI    = "https://cloudflare-dns.com/dns-query"
	googleURI        = "https://dns.google.com/experimental"
	cloudflareTorURI = "https://dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion/dns-query"

	TorAddress = "127.0.0.1:9050"
)

func NewResolver(provider string) (*net.Resolver, error) {
	var pURI string
	var dial proxy.Dialer = proxy.Direct

	switch provider {
	case "cloudflare":
		pURI = cloudflareURI
	case "cloudflare-tor":
		var err error
		dial, err = proxy.SOCKS5("tcp", TorAddress, &proxy.Auth{"doh5", "doh5"}, proxy.Direct)
		if err != nil {
			return nil, err
		}
		pURI = cloudflareTorURI
	case "google":
		pURI = googleURI
	case "", "none":
		pURI = ""
	default:
		return nil, errors.New("invalid DoH provider given")
	}

	if pURI == "" {
		log.Println("DoH disabled.  Using system resolver for DNS")
		return nil, nil
	}
	log.Printf("using DoH provider %s (%s)", provider, pURI)

	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// global http client for keepalive
	// can also set socks proxy for cloudflare-tor
	client := &http.Client{
		Transport: &http.Transport{
			Dial:                dial.Dial,
			MaxConnsPerHost:     3,
			MaxIdleConnsPerHost: 32,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	go func() {
		for {
			b := make([]byte, 4096)
			n, addr, err := pc.ReadFrom(b)
			if err != nil {
				log.Println(err)
				continue
			}
			req, err := http.NewRequest("POST", pURI, bytes.NewBuffer(b[:n]))
			if err != nil {
				log.Println(err)
				continue
			}
			req.Header.Set("Accept", "application/dns-udpwireformat")
			req.Header.Set("Content-Type", "application/dns-udpwireformat")

			go handle(addr, pc, client, req)
		}
	}()
	res := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial(network, pc.LocalAddr().String())
		},
	}
	return res, nil
}

// go-routine
func handle(addr net.Addr, pc net.PacketConn, client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = pc.WriteTo(rawResp, addr)
	if err != nil {
		log.Println(err)
		return
	}
}
