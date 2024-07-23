# Vhosting or Virtual Hosting

## What is Virtual Hosting?

Virtual hosting is a method for hosting multiple domain names (with separate handling of each name) on a single server (or pool of servers).  This allows one server to share its resources, such as memory and processor cycles, without requiring all services provided to use the same host name. The term virtual hosting is usually used in reference to web servers but the principles do carry over to other Internet services.

One widely used application is shared web hosting. The price for shared web hosting is lower than for a dedicated web server because many customers can be hosted on a single server. It is also very common for a single entity to want to use multiple names on the same machine so that the names can reflect services offered rather than where those services happen to be hosted.

## How vhosting is used with the grow-with-stlgo app

By default the application will start with 3 vhosts configured:

1. localhost - the same as the admin site
2. grow-with-stlgo.localdev.org - strips out the user manipulation part of the website
3. grow-with-stlgo-admin.localdev.org - includes user manipulation screens

## How to use these vhosts

In order to use the vhosts you will need to modify your hosts file to "spoof" your DNS lookup

### Linux

Modify the /etc/hosts file to include these lines

```bash
127.0.0.1 grow-with-stlgo-admin.localdev.org
127.0.0.1 grow-with-stlgo.localdev.org
```

### Windows

Modify the C:\Windows\System32\drivers\etc\hosts file to include these lines

```cmd
127.0.0.1 grow-with-stlgo-admin.localdev.org
127.0.0.1 grow-with-stlgo.localdev.org
```

### User modifications

In order for a user to use these vhosts they must be authorized to do so.  The individual user's vhosts array would need to include the allowable vhost.  Example:

```json
"AUser": {
    "active": true,
    "admin": false,
    "authentication": {
        "id": "user",
        "password": "obf::UEVKdumpyNrU0fuaavSGhFml_7ZpYhXrdmqmh2j_u6kbB_7eW51RLC95_jdiyR-wLk7iDb5dHTWttfmeKcqmDd-3isJOQ6f-P550Qwwo3x4L-Q18jrfDew=="
    },
    "vhosts": [
        "localhost",
        "grow-with-stlgo.localdev.org",
        "grow-with-stlgo-admin.localdev.org"
    ]
}
```

Once started with the new version of the grow-with-stlgo including vhosts you should be able to browse to the individual vhosts:

- https://localhost:10443
- https://grow-with-stlgo.localdev.org:10443
- https://grow-with-stlgo-admin.localdev.org:10443
