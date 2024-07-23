# HAProxy

## What is HAProxy?

HAProxy is a free, very fast and reliable reverse-proxy offering high availability, load balancing, and proxying for TCP and HTTP-based applications. It is particularly suited for very high traffic web sites and powers a significant portion of the world's most visited ones. Over the years it has become the de-facto standard opensource load balancer, is now shipped with most mainstream Linux distributions, and is often deployed by default in cloud platforms.

https://www.haproxy.org/

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
