# Sample app

## curl

Curl token command

```bash
 curl -i -k https://localhost:10443/REST/v1.0.0/token  -d'{"id":"admin", "password":"some password"}'
```

Output

```bash
$ curl -i -k https://localhost:10443/REST/v1.0.0/token  -d'{"id":"admin", "password":"6ebe1e8d6668de2eb541c8cac70ed304f8f98ddcac1460e1c093778e8739e1bf"}'
HTTP/1.1 201 Created
Date: Tue, 30 Jan 2024 02:14:52 GMT
Content-Length: 262
Content-Type: text/plain; charset=utf-8

{"sessionID":"927feb83-19ad-4ef6-bc3b-013ac4cb85e0","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY1ODQ0OTIsInNlc3Npb25JRCI6IjkyN2ZlYjgzLTE5YWQtNGVmNi1iYzNiLTAxM2FjNGNiODVlMCIsInVzZXJuYW1lIjoiYWRtaW4ifQ._9KqHTrFT0j_ibgk-V-3_Mo6viLyxxfShNWLHRBc7hU"}
```
