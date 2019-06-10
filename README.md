# doh5
A DNS-over-HTTPS enabled SOCKS5 proxy in Go

## Features

- Supports DoH with Cloudflare, Google, and Cloudflare .onion 
- HTTP keepalive to minimize TLS handshakes to DoH provider
- Only support DoH POST requests since GET parameters are more likely to get logged server-side
- High performance, tiny memory footprint

## Usage
```console
$ doh5 --help
Usage of doh5:
  -D [address:]port
    	[address:]port to listen and serve on (default "127.0.0.1:1080")
  -U file
    	unix domain file to listen and serve on
  -q	enable quiet mode
  -r service
    	DoH service to use: cloudflare, google, cloudflare-tor, or none (default "cloudflare")
  -s source
    	source IP to bind to for outgoing connections
```

## Examples

Run a socks proxy on 127.0.0.1 port 1080 with cloudflare DoH (default):<br>
```bash
$ doh5
```
Run a socks proxy on 127.0.0.1 port 9000 with Google DoH:<br>
```bash
$ doh5 -D 9000 -r google
```
Run a socks proxy on 0.0.0.0 port 1080 with no DoH (system resolver):<br>
```bash
$ doh5 -D 0.0.0.0:1080 -r none
```
Running Chrome with socks on the commandline:<br>
```bash
$ google-chrome-stable --incognito --proxy-server=socks5://127.0.0.1:1080
```

In Firefox, go to network settings and set the socks proxy to 127.0.0.1 port 1080.  Make sure network.proxy.socks_remote_dns is set to true.
