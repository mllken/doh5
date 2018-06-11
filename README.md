# doh5 - A DNS-over-HTTPS SOCKS5 proxy written in Go

## Features

- Supports DNS-over-HTTPS with Cloudflare, Google, and Cloudflare .onion 
- HTTP keepalive to minimize TLS handshakes to DoH provider
- Only support POST requests since GET parameters are more likely to get logged server-side
- High performance, tiny memory footprint

## Examples

Run a socks proxy on 127.0.0.1 port 1080 with cloudflare DoH (default):<br>
```./doh5```

Run a socks proxy on 127.0.0.1 port 9000 with Google DoH:<br>
```./doh5 -D 9000 -r google```

Run a socks proxy on 0.0.0.0 port 1080 with no DoH (system resolver):<br>
```./doh5 -D 0.0.0.0:1080 -r none```

Running Chrome with socks on the commandline:<br>

```google-chrome-stable --incognito --proxy-server=socks5://127.0.0.1:1080```

In Firefox, go to network settings and set the socks proxy to 127.0.0.1 port 1080
