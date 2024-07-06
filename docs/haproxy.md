# HAProxy

## What is HAProxy?

## Basic config for the grow-with-stlgo sample application

```bash
global
    log 127.0.0.1:514 local0
    chroot /var/lib/haproxy
    stats socket /run/haproxy/admin.sock mode 660 level admin
    stats timeout 30s
    user haproxy
    group haproxy
    daemon

defaults
    log global
    mode http
    option httplog
    timeout client 10s
    timeout connect 5s
    timeout server 10s
    timeout http-request 10s

frontend grow-with-stl-go
    mode http
    bind :80
    bind :443 ssl crt /home/aschiefe/grow-with-stl-go/etc/ssl.pem

    # redirect http -> https we should only use the encrypted channel
    http-request redirect scheme https unless { ssl_fc }

    default_backend servers

frontend stats
    mode http
    bind *:8404
    stats enable
    stats uri /stats
    stats refresh 10s
    stats admin if LOCALHOST
    stats auth admin:admin

backend servers
    option httpchk
    server server1 127.0.0.1:10443 ssl verify none
```
