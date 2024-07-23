# How to compile apache httpd from source

## What is apache httpd?

The Apache HTTP Server Project is an effort to develop and maintain an open-source HTTP server for modern operating systems including UNIX and Windows. The goal of this project is to provide a secure, efficient and extensible server that provides HTTP services in sync with the current HTTP standards.

The Apache HTTP Server ("httpd") was launched in 1995 and it has been the most popular web server on the Internet since April 1996. It has celebrated its 25th birthday as a project in February 2020.

The Apache HTTP Server is a project of The Apache Software Foundation.

## Why compile from source?

A couple of reasons.

1. Control, you are not tied to system level libraries and can "containerize" your install.
2. Configuration, you can add and remove things as needed.
3. You need the version the system distro doesn't have.

## The steps to compile the web server

### Zlib

https://www.zlib.net/ used for compression:

```bash
wget https://zlib.net/zlib-1.3.1.tar.gz
tar -zxf zlib-1.3.1.tar.gz
cd zlib-1.3.1
./configure --64 --prefix=/opt/app/web/zlib-1.3.1
make
make install
cd /opt/app/web
ln -s zlib-1.3.1 zlib
```

### OpenSSL

https://www.openssl.org/ used for transport layer security

```bash
wget https://openssl.org/source/openssl-3.0.14.tar.gz
tar -zxf openssl-3.0.14.tar.gz
cd openssl-3.0.14
./config --prefix=/opt/app/web/openssl-3.0.14 shared -fPIC
make
make test
make install
cd /opt/app/web
ln -s openssl-3.0.14 openssl
```

### Perl Compatible Regular Expressions (PCRE)

https://github.com/PCRE2Project/pcre2, mostly because this is what I saw my mentor do.

```bash
wget https://github.com/PCRE2Project/pcre2/releases/download/pcre2-10.44/pcre2-10.44.tar.gz
tar -zxf pcre2-10.44.tar.gz
cd pcre2-10.44
./configure --prefix=/opt/app/web/pcre2-10.44
make
make install
cd /opt/app/web
ln -s pcre2-10.44 pcre
```

### Apache Portable Runtime (APR)

https://apr.apache.org/ used to normalize behavior across platforms

```bash
wget https://dlcdn.apache.org/apr/apr-1.7.4.tar.gz
tar -zxf apr-1.7.4.tar.gz
cd apr-1.7.4
./configure --prefix=/opt/app/web/apr-1.7.4 --enable-threads
make
make install
cd /opt/app/web
ln -s apr-1.7.4 apr
```

### Apache Portable Runtime Utils (APR)

https://apr.apache.org/ used to normalize behavior across platforms

```bash
wget https://dlcdn.apache.org/apr/apr-util-1.6.3.tar.gz
tar -zxf apr-util-1.6.3.tar.gz
cd apr-util-1.6.3
./configure --prefix=/opt/app/web/apr-util-1.6.3 --with-apr=/opt/app/web/apr
make
make install
cd /opt/app/web
ln -s apr-util-1.6.3 apr-util
```

### And now apache HTTPD

https://httpd.apache.org/ the web server

```bash
wget https://dlcdn.apache.org/httpd/httpd-2.4.62.tar.gz
tar -zxf httpd-2.4.62.tar.gz
cd httpd-2.4.62
env PCRE_CONFIG=/opt/app/web/pcre/bin/pcre2-config \
./configure --prefix=/opt/app/web/httpd-2.4.62 \
        --enable-most \
        --enable-mods-shared=all \
        --enable-so \
        --with-z=/opt/app/web/zlib \
        --with-ssl=/opt/app/web/openssl \
        --with-apr=/opt/app/web/apr \
        --with-apr-util=/opt/app/web/apr-util \
        --with-pcre=/opt/app/web/pcre
make
make install
cd /opt/app/web
ln -s httpd-2.4.62 httpd
```

## Configuration

### Enable modules

Edit /opt/app/web/httpd/conf/httpd.conf and uncomment the LoadModule definitions until you see a list similar to this:

```bash
/opt/app/web/httpd/bin$ ./apachectl -t -D DUMP_MODULES
Loaded Modules:
 core_module (static)
 so_module (static)
 http_module (static)
 mpm_event_module (static)
 authn_file_module (shared)
 authn_core_module (shared)
 authz_host_module (shared)
 authz_groupfile_module (shared)
 authz_user_module (shared)
 authz_core_module (shared)
 access_compat_module (shared)
 auth_basic_module (shared)
 reqtimeout_module (shared)
 ext_filter_module (shared)
 include_module (shared)
 filter_module (shared)
 deflate_module (shared)
 mime_module (shared)
 log_config_module (shared)
 env_module (shared)
 headers_module (shared)
 setenvif_module (shared)
 version_module (shared)
 proxy_module (shared)
 proxy_connect_module (shared)
 proxy_ftp_module (shared)
 proxy_http_module (shared)
 proxy_wstunnel_module (shared)
 proxy_ajp_module (shared)
 proxy_balancer_module (shared)
 proxy_express_module (shared)
 slotmem_shm_module (shared)
 slotmem_plain_module (shared)
 ssl_module (shared)
 lbmethod_byrequests_module (shared)
 unixd_module (shared)
 status_module (shared)
 autoindex_module (shared)
 cgid_module (shared)
 dir_module (shared)
 alias_module (shared)
 rewrite_module (shared)
 vhost_alias_module (shared)
```

Add this to the bottom of the httpd.conf file:

```bash
<IfModule deflate_module>
AddOutputFilterByType DEFLATE text/html text/plain text/css text/javascript application/javascript
DeflateCompressionLevel 9
</IfModule>

Include  conf/vhosts.conf
```

Create the conf/vhosts.conf file:

```bash
<IfModule !vhost_alias_module>
  LoadModule vhost_alias_module modules/mod_vhost_alias.so
</IfModule>

Include conf/vhosts/*.conf
```

Create the conf/vhosts directory and then create 3 vhost files.

The default.conf vhost, this will respond to https://localhost:11443:
```bash
<IfModule !headers_module>
  LoadModule headers_module modules/mod_headers.so
</IfModule>
<IfModule !status_module>
  LoadModule status_module modules/mod_status.so
</IfModule>

<VirtualHost *:11080>
  ServerName localhost
  RedirectMatch permanent ^/(.*)$ https://localhost:11443/$1
</VirtualHost>

<VirtualHost *:11443>
  SSLEngine on
  SSLProxyEngine On
  SSLProxyCheckPeerCN on
  SSLProxyCheckPeerExpire on
  SSLCipherSuite ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-C
HACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305
  SSLHonorCipherOrder on

  SSLProtocol -all +TLSv1.3

  SSLCertificateFile "/home/aschiefe/grow-with-stl-go/etc/cert.pem"
  SSLCertificateKeyFile "/home/aschiefe/grow-with-stl-go/etc/key.pem"

  ServerName localhost:11443
  ServerAdmin root@localhost

  CustomLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/default/logs/access.log.%Y.%m.%d 86400" common
  ErrorLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/default/logs/error.log.%Y.%m.%d 86400"

  <IfModule deflate_module>
    AddOutputFilterByType DEFLATE text/html text/plain text/xml
  </IfModule>

  <IfModule rewrite_module>
    RewriteEngine On
    # Disable HTTP TRACE support
    RewriteCond %{REQUEST_METHOD} ^TRACE
    RewriteRule .* - [F]
  </IfModule>

  DocumentRoot "/opt/app/web/vhosts/default/web_root"
  <Directory "/opt/app/web/vhosts/default/web_root">
    Options Indexes FollowSYmLinks Includes
    Order deny,allow
    Require all granted
  </Directory>

  ScriptAlias /cgi-bin/ /opt/app/web/vhosts/default/cgi-bin/
  <Directory "/opt/app/web/vhosts/default/cgi-bin/">
    Options Indexes FollowSYmLinks Includes ExecCGI
    AllowOverride All
    Require all granted
  </Directory>
</VirtualHost>
```

The grow-with-stlgo-admin.localdev.org.conf will respond to https://grow-with-stlgo-admin.localdev.org

```bash
<IfModule !headers_module>
  LoadModule headers_module modules/mod_headers.so
</IfModule>
<IfModule !status_module>
  LoadModule status_module modules/mod_status.so
</IfModule>

<VirtualHost *:11080>
  ServerName grow-with-stlgo-admin.localdev.org
  RedirectMatch permanent ^/(.*)$ https://grow-with-stlgo-admin.localdev.org:11443/$1
</VirtualHost>

<VirtualHost *:11443>
  SSLEngine on
  SSLCipherSuite ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-C
HACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305
  SSLHonorCipherOrder on

  SSLProtocol -all +TLSv1.3

  SSLCertificateFile "/home/aschiefe/grow-with-stl-go/etc/cert.pem"
  SSLCertificateKeyFile "/home/aschiefe/grow-with-stl-go/etc/key.pem"

  ServerName grow-with-stlgo-admin.localdev.org:11443
  ServerAdmin root@localhost

  CustomLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/grow-with-stlgo-admin.localdev.org/logs/access.log.%Y.%m.%d 86400" common
  ErrorLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/grow-with-stlgo-admin.localdev.org/logs/error.log.%Y.%m.%d 86400"

  <IfModule deflate_module>
    AddOutputFilterByType DEFLATE text/html text/plain text/xml
  </IfModule>

  <IfModule rewrite_module>
    RewriteEngine On
    # Disable HTTP TRACE support
    RewriteCond %{REQUEST_METHOD} ^TRACE
    RewriteRule .* - [F]
  </IfModule>

  DocumentRoot "/home/aschiefe/grow-with-stl-go/web/grow-with-stlgo-admin"
  <Directory "/home/aschiefe/grow-with-stl-go/web/grow-with-stlgo-admin">
    Options Indexes FollowSYmLinks Includes
    Order deny,allow
    Require all granted
  </Directory>
</VirtualHost>
```

The grow-with-stlgo.localdev.org.conf will respond to https://grow-with-stlgo.localdev.org

```bash
<IfModule !headers_module>
  LoadModule headers_module modules/mod_headers.so
</IfModule>
<IfModule !status_module>
  LoadModule status_module modules/mod_status.so
</IfModule>

<VirtualHost *:11080>
  ServerName grow-with-stlgo.localdev.org
  RedirectMatch permanent ^/(.*)$ https://grow-with-stlgo.localdev.org:11443/$1
</VirtualHost>

<VirtualHost *:11443>
  SSLEngine on
  SSLProxyEngine On
  SSLProxyCheckPeerCN on
  SSLProxyCheckPeerExpire on
  SSLCipherSuite ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-C
HACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305
  SSLHonorCipherOrder on

  SSLProtocol -all +TLSv1.3

  SSLCertificateFile "/home/aschiefe/grow-with-stl-go/etc/cert.pem"
  SSLCertificateKeyFile "/home/aschiefe/grow-with-stl-go/etc/key.pem"

  ServerName grow-with-stlgo.localdev.org:11443
  ServerAdmin root@localhost

  CustomLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/grow-with-stlgo.localdev.org/logs/access.log.%Y.%m.%d 86400" common
  ErrorLog "|/opt/app/web/httpd/bin/rotatelogs -l /opt/app/web/vhosts/grow-with-stlgo.localdev.org/logs/error.log.%Y.%m.%d 86400"

  <IfModule deflate_module>
    AddOutputFilterByType DEFLATE text/html text/plain text/xml
  </IfModule>

  <IfModule rewrite_module>
    RewriteEngine On
    # Disable HTTP TRACE support
    RewriteCond %{REQUEST_METHOD} ^TRACE
    RewriteRule .* - [F]
  </IfModule>

  DocumentRoot "/home/aschiefe/grow-with-stl-go/web/grow-with-stlgo"
  <Directory "/home/aschiefe/grow-with-stl-go/web/grow-with-stlgo">
    Options Indexes FollowSYmLinks Includes
    Order deny,allow
    Require all granted
  </Directory>
</VirtualHost>
```

To start apache:

```bash
/opt/app/web/httpd/bin/apachectl start
```

To stop apache:

```bash
/opt/app/web/httpd/bin/apachectl stop
```
