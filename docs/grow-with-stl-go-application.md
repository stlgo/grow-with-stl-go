# Grow with STL Go

The Grow with STL Go is a sample application to show what can be done with basic Go and some html / javascript.  It demonstrates the following things:

- How to produce a functional product with go
- How to create a configuration for your program and have it correctly populated on first run without externalities.
- A basic web based UI
- Interactions with the backend over both REST and WebSocket interfaces.
- An embedded database that's auto populated with everything you need out of the gate.
- Role based user administration.
- Strong password enforcement, along with proper storage of hashed passwords which are also encrypted in memory.
- Different presentation for user interaction based on role.
- Cryptographically secure endpoints with a self signed certificate auto generated if one isn't present on start.
- Sensitive data becomes obfuscated when written to disk.
- The embedded database is encrypted so random access is not possible without the key.

## Access the UI

Start the application for the first time with:

```bash
$ bin/grow-with-stl-go --loglevel 6
```

Notice the WARN output in the log lines that say what the default passwords are for both user and admin:

```bash
[stl-go] 2024/02/28 12:11:20 stl-go/grow-with-stl-go/pkg/configs/apiuser.go:54: [WARN] Password generated for user 'admin', password 81be7ce107fdee7b483702625ea33358bfb7b1fc468acae537614dd738cee0b - DO NOT USE THIS FOR PRODUCTION
[stl-go] 2024/02/28 12:11:20 stl-go/grow-with-stl-go/pkg/configs/apiuser.go:54: [WARN] Password generated for user 'user', password bef62cbfa5450c6a86f942b5f96ac5c7ca30d3832a500b1d10b38fb4b5c57bfa - DO NOT USE THIS FOR PRODUCTION
```

You can use these ids to log into the web page <https://localhost:10443/>

## Accessing the REST APIs with cURl

### Token request

First thing that must happen is we need to request an authentication token for a specific user:

```bash
 $  curl -i -k https://localhost:10443/REST/v1.0.0/token  -d'{"id":"user", "password":"bef62cbfa5450c6a86f942b5f96ac5c7ca30d3832a500b1d10b38fb4b5c57bfa"}'
```

Output

```bash
HTTP/2 201
content-type: text/plain; charset=utf-8
content-length: 260
date: Wed, 28 Feb 2024 18:13:44 GMT

{
    "sessionID": "42b0cda3-caa7-4578-8784-399db9eb408b",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW88"
}
```

We will need 2 things from this request on subsequent requests.

1. The token will be used for the Authentication header.
2. The sessionID will be used for the sessionID header.

### Get seed inventory with cURL

Command, notice that the session id and token are taken from the result of the above command

- URL: https://localhost:10443/REST/v1.0.0/seeds/getInventory
- Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW8
- Session ID: 42b0cda3-caa7-4578-8784-399db9eb408b

```bash
$  curl -i -k -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW88" -H "sessionID: 42b0cda3-caa7-4578-8784-399db9eb408b" https://localhost:10443/REST/v1.0.0/seeds/getInventory
```

Output

```json
{
    "Herb": {
        "category": "Herb",
        "items": {
            "133a8d28-4b22-432e-af68-1a90f521067c": {
                "id": "133a8d28-4b22-432e-af68-1a90f521067c",
                "category": "Herb",
                "genus": "Anethum",
                "species": "graveolens",
                "cultivar": "Ella",
                "commonName": "Dill Weed",
                "description": "Ella is a dwarf dill bred for container and hydroponic growing",
                "hybrid": false,
                "price": 2.67,
                "perPacketCount": 20,
                "packets": 0,
                "image": "/images/herbs/ella_dill.jpg"
            },
            "7a027c43-3bf9-4dca-bfb8-fd70ccfede70": {
                "id": "7a027c43-3bf9-4dca-bfb8-fd70ccfede70",
                "category": "Herb",
                "genus": "Ocimum",
                "species": "basilicum",
                "cultivar": "Genovese",
                "commonName": "Basil",
                "description": "Genovese basil was first bred in the Northwest coastal port of Genoa, gateway to the Italian Riviera.",
                "hybrid": false,
                "price": 4.28,
                "perPacketCount": 100,
                "packets": 94,
                "image": "/images/herbs/genovese_basil.jpg"
            },
            "9b37afb6-e191-4a37-9528-7ca57320dcfe": {
                "id": "9b37afb6-e191-4a37-9528-7ca57320dcfe",
                "category": "Herb",
                "genus": "Allium",
                "species": "schoenoprasum",
                "cultivar": "Polyvert",
                "commonName": "Chive",
                "description": "Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.",
                "hybrid": false,
                "price": 2.18,
                "perPacketCount": 100,
                "packets": 100,
                "image": "/images/herbs/polyvert_chive.jpg"
            },
            "e57ed367-de96-468c-a428-d59c85b4fac8": {
                "id": "e57ed367-de96-468c-a428-d59c85b4fac8",
                "category": "Herb",
                "genus": "Origanum",
                "species": "vulgare",
                "cultivar": "Greek",
                "commonName": "Oregano",
                "description": "Strong oregano aroma and flavor; great for pizza and Italian cooking. Characteristic dark green leaves with white flowers.",
                "hybrid": false,
                "price": 3.95,
                "perPacketCount": 50,
                "packets": 100,
                "image": "/images/herbs/greek_oregano.jpg"
            }
        }
    },
    "Onion": {
        "category": "Onion",
        "items": {
            "0885e242-aa1e-453d-8700-b8dcb4dc270f": {
                "id": "0885e242-aa1e-453d-8700-b8dcb4dc270f",
                "category": "Onion",
                "genus": "Allium",
                "species": "cepa",
                "cultivar": "Patterson",
                "commonName": "Yellow",
                "description": "Patterson’ is a keeper—the longest-storing onion you can find. Straw-colored, globe-shaped bulbs with sweet, mildly pungent yellow flesh",
                "hybrid": false,
                "price": 4.28,
                "perPacketCount": 100,
                "packets": 100,
                "image": "/images/onions/patterson.jpg"
            },
            "14df4848-53db-4272-935a-a60a436a0b30": {
                "id": "14df4848-53db-4272-935a-a60a436a0b30",
                "category": "Onion",
                "genus": "Allium",
                "species": "cepa",
                "cultivar": "Red Wing",
                "commonName": "Red",
                "description": "Uniform, large onions with deep red color. Thick skin, very hard bulbs for long storage. Consistent internal color.",
                "hybrid": true,
                "price": 3.95,
                "perPacketCount": 50,
                "packets": 100,
                "image": "/images/onions/red_wing.jpg"
            },
            "80c3db24-5253-498b-9628-07fbe2b43cf4": {
                "id": "80c3db24-5253-498b-9628-07fbe2b43cf4",
                "category": "Onion",
                "genus": "Allium",
                "species": "cepa",
                "cultivar": "Walla Walla",
                "commonName": "Sweet",
                "description": "Juicy, sweet, regional favorite. In the Northwest,  very large, flattened, ultra-mild onions",
                "hybrid": false,
                "price": 2.18,
                "perPacketCount": 100,
                "packets": 100,
                "image": "/images/onions/walla_walla.jpg"
            },
            "d867d75f-20f3-40c7-884e-699f4fb9e8b1": {
                "id": "d867d75f-20f3-40c7-884e-699f4fb9e8b1",
                "category": "Onion",
                "genus": "Allium",
                "species": "cepa",
                "cultivar": "Ailsa Craig",
                "commonName": "Yellow",
                "description": "Long day. Very well-known globe-shaped heirloom onion that reaches a really huge size—5 lbs is rather common",
                "hybrid": false,
                "price": 2.67,
                "perPacketCount": 20,
                "packets": 100,
                "image": "/images/onions/ailsa_craig.jpg"
            }
        }
    },
    "Pepper": {
        "category": "Pepper",
        "items": {
            "119df5ff-6766-4646-8f3a-01d2ff3002e8": {
                "id": "119df5ff-6766-4646-8f3a-01d2ff3002e8",
                "category": "Pepper",
                "genus": "Capsicum",
                "species": "annuum",
                "cultivar": "Tampiqueno",
                "commonName": "Serrano",
                "description": "This serrano variety comes from the mountains of the Hidalgo and Puebla states of Mexico.  This pepper is 2-3 times hotter than jalapenos",
                "hybrid": false,
                "price": 6.45,
                "perPacketCount": 15,
                "packets": 100,
                "image": "/images/peppers/serrano.jpg"
            },
            "6f5cdedc-3421-4f9e-a23f-ecb2eabd36bd": {
                "id": "6f5cdedc-3421-4f9e-a23f-ecb2eabd36bd",
                "category": "Pepper",
                "genus": "Capsicum",
                "species": "annuum",
                "cultivar": "Zapotec",
                "commonName": "Jalapeno",
                "description": "This jalapeno variety from Oaxaca, Mexico which is a more flavorful, gourmet jalapeño.",
                "hybrid": false,
                "price": 3.78,
                "perPacketCount": 12,
                "packets": 100,
                "image": "/images/peppers/jalapeno.jpg"
            },
            "8c1633cc-faf9-4781-9ae7-e901ba26d122": {
                "id": "8c1633cc-faf9-4781-9ae7-e901ba26d122",
                "category": "Pepper",
                "genus": "Capsicum",
                "species": "annuum",
                "commonName": "Poblano",
                "description": "The poblano is a mild chili pepper originating in the state of Puebla, Mexico. Dried, it is called ancho or chile ancho",
                "hybrid": false,
                "price": 2.96,
                "perPacketCount": 25,
                "packets": 100,
                "image": "/images/peppers/poblano.jpg"
            },
            "8e0c4a5b-4305-4421-a298-d7631218546a": {
                "id": "8e0c4a5b-4305-4421-a298-d7631218546a",
                "category": "Pepper",
                "genus": "Capsicum",
                "species": "annuum",
                "cultivar": "Ozark Giant",
                "commonName": "Bell",
                "description": "Green bell peppers are bell peppers that have been harvested early. Red bell peppers have been allowed to ripen longer.",
                "hybrid": false,
                "price": 4.58,
                "perPacketCount": 50,
                "packets": 100,
                "image": "/images/peppers/green_bell.jpg"
            }
        }
    },
    "Tomato": {
        "category": "Tomato",
        "items": {
            "78410bc7-a1c1-430b-b081-305fbb6438fc": {
                "id": "78410bc7-a1c1-430b-b081-305fbb6438fc",
                "category": "Tomato",
                "genus": "Solanum",
                "species": "lycopersicum",
                "cultivar": "Plum Regal",
                "commonName": "Tomato",
                "description": "Medium-size plants with good leaf cover produce high yields of blocky, 4 oz. plum tomatoes. Fruits have a deep red color with good flavor. Determinate.",
                "hybrid": true,
                "price": 2.96,
                "perPacketCount": 25,
                "packets": 100,
                "image": "/images/tomatoes/determinate/plum_regal.jpg"
            },
            "84507f06-ec5b-4341-81f7-b38cadd5fdcb": {
                "id": "84507f06-ec5b-4341-81f7-b38cadd5fdcb",
                "category": "Tomato",
                "genus": "Solanum",
                "species": "lycopersicum",
                "cultivar": "Carbon",
                "commonName": "Tomato",
                "description": "Indeterminate heirloom. Resists cracking and cat-facing better than other large, black heirlooms. Blocky-round, 10-14 oz. fruit with dark olive shoulders.",
                "hybrid": false,
                "price": 3.78,
                "perPacketCount": 12,
                "packets": 100,
                "image": "/images/tomatoes/indeterminate/carbon.jpg"
            },
            "d78063f7-2110-4339-8fbe-0f9cb1d0ea4d": {
                "id": "d78063f7-2110-4339-8fbe-0f9cb1d0ea4d",
                "category": "Tomato",
                "genus": "Solanum",
                "species": "lycopersicum",
                "cultivar": "Galahad",
                "commonName": "Tomato",
                "description": "Delicious early determinate beefsteak.",
                "hybrid": true,
                "price": 4.58,
                "perPacketCount": 50,
                "packets": 100,
                "image": "/images/tomatoes/determinate/galahad.jpg"
            },
            "e78245b8-859f-48e5-a08e-fdbb1a74418f": {
                "id": "e78245b8-859f-48e5-a08e-fdbb1a74418f",
                "category": "Tomato",
                "genus": "Solanum",
                "species": "lycopersicum",
                "cultivar": "San Marzano",
                "commonName": "Tomato",
                "description": "San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate",
                "hybrid": false,
                "price": 6.45,
                "perPacketCount": 15,
                "packets": 100,
                "image": "/images/tomatoes/indeterminate/san_marzano.jpg"
            }
        }
    }
    }
```

### Get detail about a specific seed in inventory with cURL

Command, notice that the session id, token, category and seed id are taken from the result of the above command

- URL: https://localhost:10443/REST/v1.0.0/seeds/getDetail/Tomato/e78245b8-859f-48e5-a08e-fdbb1a74418f
- Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW8
- Session ID: 42b0cda3-caa7-4578-8784-399db9eb408b
- Category: Tomato
- Seed ID: e78245b8-859f-48e5-a08e-fdbb1a74418f

```bash
curl -i -k -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW88" -H "sessionID: 42b0cda3-caa7-4578-8784-399db9eb408b" https://localhost:10443/REST/v1.0.0/seeds/getDetail/Tomato/e78245b8-859f-48e5-a08e-fdbb1a74418f
```

Output

```bash
HTTP/2 201
content-type: text/plain; charset=utf-8
content-length: 406
date: Wed, 28 Feb 2024 18:26:35 GMT

{
    "id": "e78245b8-859f-48e5-a08e-fdbb1a74418f",
    "category": "Tomato",
    "genus": "Solanum",
    "species": "lycopersicum",
    "cultivar": "San Marzano",
    "commonName": "Tomato",
    "description": "San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate",
    "hybrid": false,
    "price": 6.45,
    "perPacketCount": 15,
    "packets": 100,
    "image": "/images/tomatoes/indeterminate/san_marzano.jpg"
}
```

### Purchase a specific seed in inventory with cURL

Command, notice that the session id, token, category and seed id are taken from the result of the above command

- URL: https://localhost:10443/REST/v1.0.0/seeds/getDetail/Tomato/e78245b8-859f-48e5-a08e-fdbb1a74418f
- Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW8
- Session ID: 42b0cda3-caa7-4578-8784-399db9eb408b
- Category: Tomato
- Seed ID: e78245b8-859f-48e5-a08e-fdbb1a74418f

```bash
curl -i -k -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkxNDc2MjQsInNlc3Npb25JRCI6IjQyYjBjZGEzLWNhYTctNDU3OC04Nzg0LTM5OWRiOWViNDA4YiIsInVzZXJuYW1lIjoidXNlciJ9.a1cC0qYPt-kJl9-vqoN_MnUN2eYZZVung3CRPA0iW88" -H "sessionID: 42b0cda3-caa7-4578-8784-399db9eb408b" -d '{"category": "Tomato", "id": "e78245b8-859f-48e5-a08e-fdbb1a74418f", "quantity": "1"}' https://localhost:10443/REST/v1.0.0/seeds/purchase
```

Output

```bash
{
    "id": "e78245b8-859f-48e5-a08e-fdbb1a74418f",
    "category": "Tomato",
    "genus": "Solanum",
    "species": "lycopersicum",
    "cultivar": "San Marzano",
    "commonName": "Tomato",
    "description": "San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate",
    "hybrid": false,
    "price": 6.45,
    "perPacketCount": 15,
    "packets": 99,
    "image": "/images/tomatoes/indeterminate/san_marzano.jpg"
}
```

## Accessing the WebSocket APIs

You can use [Postman](https://www.postman.com/downloads/) to create WebSocket requests.  To do this you'll have to go to the [file menu -> new -> WebSocket](https://learning.postman.com/docs/sending-requests/websocket/create-a-websocket-request/)

### Connect WebSocket

Request URL:
wss://localhost:10443/ws/v1.0.0

Click connect, you should see this output:

```json
{
    "route": "websocketclient",
    "type": "initialize",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "timestamp": 1709933441446
}
```

### Authenticate WebSocket

Notice the session id is from the above response

Send message:

```json
{
    "route": "websocketclient",
    "type": "auth",
    "component": "authenticate",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "authentication": {
        "id": "user",
        "password": "dc92056dce84c58d5273bd8fc6dbd359d27e99a9e233e8b63e7e317bd78bb68"
    }
}
```

Output

```json
{
    "route": "websocketclient",
    "type": "auth",
    "component": "authenticate",
    "subComponent": "approved",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "timestamp": 1709933460572,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk5MzcwNjAsInNlc3Npb25JRCI6ImY3YTgwMTMwLWUxNzctNDU4Mi04ZmVmLTFkZGU1NWVmOTFlMyIsInVzZXJuYW1lIjoidXNlciJ9._ZhePlTvCuS5BU64epFCrk-NoKJQCbEWeFd2yANPdkM",
    "isAdmin": false
}
```

### Request seed inventory with WebSocket

Notice the session id and token are from the above response

Send message:

```json
{
    "route": "seeds",
    "type": "getInventory",
    "component": "getInventory",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk5MzcwNjAsInNlc3Npb25JRCI6ImY3YTgwMTMwLWUxNzctNDU4Mi04ZmVmLTFkZGU1NWVmOTFlMyIsInVzZXJuYW1lIjoidXNlciJ9._ZhePlTvCuS5BU64epFCrk-NoKJQCbEWeFd2yANPdkM"
}
```

Output

```json
{
    "route": "seeds",
    "type": "getInventory",
    "component": "getInventory",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "timestamp": 1709933522666,
    "data": {
        "Herb": {
            "category": "Herb",
            "items": {
                "457a054c-487e-4fd5-91e3-3208d351e682": {
                    "id": "457a054c-487e-4fd5-91e3-3208d351e682",
                    "category": "Herb",
                    "genus": "Allium",
                    "species": "schoenoprasum",
                    "cultivar": "Polyvert",
                    "commonName": "Chive",
                    "description": "Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.",
                    "hybrid": false,
                    "price": 2.18,
                    "perPacketCount": 100,
                    "packets": 100,
                    "image": "/images/herbs/polyvert_chive.jpg"
                },
                "a4b54775-4ad9-4671-818d-85f7a4932bd4": {
                    "id": "a4b54775-4ad9-4671-818d-85f7a4932bd4",
                    "category": "Herb",
                    "genus": "Anethum",
                    "species": "graveolens",
                    "cultivar": "Ella",
                    "commonName": "Dill Weed",
                    "description": "Ella is a dwarf dill bred for container and hydroponic growing",
                    "hybrid": false,
                    "price": 2.67,
                    "perPacketCount": 20,
                    "packets": 98,
                    "image": "/images/herbs/ella_dill.jpg"
                },
                "e4857eb8-8bb4-46f2-b660-051d15cf3022": {
                    "id": "e4857eb8-8bb4-46f2-b660-051d15cf3022",
                    "category": "Herb",
                    "genus": "Ocimum",
                    "species": "basilicum",
                    "cultivar": "Genovese",
                    "commonName": "Basil",
                    "description": "Genovese basil was first bred in the Northwest coastal port of Genoa, gateway to the Italian Riviera.",
                    "hybrid": false,
                    "price": 4.28,
                    "perPacketCount": 100,
                    "packets": 100,
                    "image": "/images/herbs/genovese_basil.jpg"
                },
                "f0632288-6519-479c-b345-65998d387aca": {
                    "id": "f0632288-6519-479c-b345-65998d387aca",
                    "category": "Herb",
                    "genus": "Origanum",
                    "species": "vulgare",
                    "cultivar": "Greek",
                    "commonName": "Oregano",
                    "description": "Strong oregano aroma and flavor; great for pizza and Italian cooking. Characteristic dark green leaves with white flowers.",
                    "hybrid": false,
                    "price": 3.95,
                    "perPacketCount": 50,
                    "packets": 100,
                    "image": "/images/herbs/greek_oregano.jpg"
                }
            }
        },
        "Onion": {
            "category": "Onion",
            "items": {
                "194bc187-633f-4fd0-ab49-2331b743e420": {
                    "id": "194bc187-633f-4fd0-ab49-2331b743e420",
                    "category": "Onion",
                    "genus": "Allium",
                    "species": "cepa",
                    "cultivar": "Red Wing",
                    "commonName": "Red",
                    "description": "Uniform, large onions with deep red color. Thick skin, very hard bulbs for long storage. Consistent internal color.",
                    "hybrid": true,
                    "price": 3.95,
                    "perPacketCount": 50,
                    "packets": 100,
                    "image": "/images/onions/red_wing.jpg"
                },
                "3e14e501-a5b9-4d8b-ac05-64c8952ab5e4": {
                    "id": "3e14e501-a5b9-4d8b-ac05-64c8952ab5e4",
                    "category": "Onion",
                    "genus": "Allium",
                    "species": "cepa",
                    "cultivar": "Walla Walla",
                    "commonName": "Sweet",
                    "description": "Juicy, sweet, regional favorite. In the Northwest,  very large, flattened, ultra-mild onions",
                    "hybrid": false,
                    "price": 2.18,
                    "perPacketCount": 100,
                    "packets": 95,
                    "image": "/images/onions/walla_walla.jpg"
                },
                "5a47f887-3fd5-4ead-99ca-ed14bad694af": {
                    "id": "5a47f887-3fd5-4ead-99ca-ed14bad694af",
                    "category": "Onion",
                    "genus": "Allium",
                    "species": "cepa",
                    "cultivar": "Patterson",
                    "commonName": "Yellow",
                    "description": "Patterson’ is a keeper—the longest-storing onion you can find. Straw-colored, globe-shaped bulbs with sweet, mildly pungent yellow flesh",
                    "hybrid": false,
                    "price": 4.28,
                    "perPacketCount": 100,
                    "packets": 100,
                    "image": "/images/onions/patterson.jpg"
                },
                "7bf6183a-972c-4d14-8828-40d2b4742710": {
                    "id": "7bf6183a-972c-4d14-8828-40d2b4742710",
                    "category": "Onion",
                    "genus": "Allium",
                    "species": "cepa",
                    "cultivar": "Ailsa Craig",
                    "commonName": "Yellow",
                    "description": "Long day. Very well-known globe-shaped heirloom onion that reaches a really huge size—5 lbs is rather common",
                    "hybrid": false,
                    "price": 2.67,
                    "perPacketCount": 20,
                    "packets": 100,
                    "image": "/images/onions/ailsa_craig.jpg"
                }
            }
        },
        "Pepper": {
            "category": "Pepper",
            "items": {
                "17a8c6d1-4be8-410c-9e97-ba88e7557d54": {
                    "id": "17a8c6d1-4be8-410c-9e97-ba88e7557d54",
                    "category": "Pepper",
                    "genus": "Capsicum",
                    "species": "annuum",
                    "cultivar": "Ozark Giant",
                    "commonName": "Bell",
                    "description": "Green bell peppers are bell peppers that have been harvested early. Red bell peppers have been allowed to ripen longer.",
                    "hybrid": false,
                    "price": 4.58,
                    "perPacketCount": 50,
                    "packets": 100,
                    "image": "/images/peppers/green_bell.jpg"
                },
                "28c22773-547c-4a28-b803-ba8b4298c117": {
                    "id": "28c22773-547c-4a28-b803-ba8b4298c117",
                    "category": "Pepper",
                    "genus": "Capsicum",
                    "species": "annuum",
                    "cultivar": "Tampiqueno",
                    "commonName": "Serrano",
                    "description": "This serrano variety comes from the mountains of the Hidalgo and Puebla states of Mexico.  This pepper is 2-3 times hotter than jalapenos",
                    "hybrid": false,
                    "price": 6.45,
                    "perPacketCount": 15,
                    "packets": 99,
                    "image": "/images/peppers/serrano.jpg"
                },
                "3e3cd1fb-48a0-41f4-8f3e-30542d8b4daf": {
                    "id": "3e3cd1fb-48a0-41f4-8f3e-30542d8b4daf",
                    "category": "Pepper",
                    "genus": "Capsicum",
                    "species": "annuum",
                    "commonName": "Poblano",
                    "description": "The poblano is a mild chili pepper originating in the state of Puebla, Mexico. Dried, it is called ancho or chile ancho",
                    "hybrid": false,
                    "price": 2.96,
                    "perPacketCount": 25,
                    "packets": 100,
                    "image": "/images/peppers/poblano.jpg"
                },
                "9f9d76b1-eda3-4d98-85d9-df8effb04ea8": {
                    "id": "9f9d76b1-eda3-4d98-85d9-df8effb04ea8",
                    "category": "Pepper",
                    "genus": "Capsicum",
                    "species": "annuum",
                    "cultivar": "Zapotec",
                    "commonName": "Jalapeno",
                    "description": "This jalapeno variety from Oaxaca, Mexico which is a more flavorful, gourmet jalapeño.",
                    "hybrid": false,
                    "price": 3.78,
                    "perPacketCount": 12,
                    "packets": 100,
                    "image": "/images/peppers/jalapeno.jpg"
                }
            }
        },
        "Tomato": {
            "category": "Tomato",
            "items": {
                "18d87133-f086-4c02-b477-1bd1af6f4309": {
                    "id": "18d87133-f086-4c02-b477-1bd1af6f4309",
                    "category": "Tomato",
                    "genus": "Solanum",
                    "species": "lycopersicum",
                    "cultivar": "Plum Regal",
                    "commonName": "Tomato",
                    "description": "Medium-size plants with good leaf cover produce high yields of blocky, 4 oz. plum tomatoes. Fruits have a deep red color with good flavor. Determinate.",
                    "hybrid": true,
                    "price": 2.96,
                    "perPacketCount": 25,
                    "packets": 100,
                    "image": "/images/tomatoes/determinate/plum_regal.jpg"
                },
                "4c5e38f1-edab-440e-99cd-29c557b925cd": {
                    "id": "4c5e38f1-edab-440e-99cd-29c557b925cd",
                    "category": "Tomato",
                    "genus": "Solanum",
                    "species": "lycopersicum",
                    "cultivar": "San Marzano",
                    "commonName": "Tomato",
                    "description": "San Marzano is considered one of the best paste tomatoes of all time, with Old World look and taste.  Indeterminate",
                    "hybrid": false,
                    "price": 6.45,
                    "perPacketCount": 15,
                    "packets": 100,
                    "image": "/images/tomatoes/indeterminate/san_marzano.jpg"
                },
                "8496e153-0a6c-425f-b9c5-8d55dc4b0d3f": {
                    "id": "8496e153-0a6c-425f-b9c5-8d55dc4b0d3f",
                    "category": "Tomato",
                    "genus": "Solanum",
                    "species": "lycopersicum",
                    "cultivar": "Carbon",
                    "commonName": "Tomato",
                    "description": "Indeterminate heirloom. Resists cracking better than other large, black heirlooms. Blocky-round, 10-14 oz. fruit with dark olive shoulders.",
                    "hybrid": false,
                    "price": 3.78,
                    "perPacketCount": 12,
                    "packets": 100,
                    "image": "/images/tomatoes/indeterminate/carbon.jpg"
                },
                "c119da73-5648-435a-a22f-80d09c86098d": {
                    "id": "c119da73-5648-435a-a22f-80d09c86098d",
                    "category": "Tomato",
                    "genus": "Solanum",
                    "species": "lycopersicum",
                    "cultivar": "Galahad",
                    "commonName": "Tomato",
                    "description": "Delicious early determinate beefsteak.",
                    "hybrid": true,
                    "price": 4.58,
                    "perPacketCount": 50,
                    "packets": 100,
                    "image": "/images/tomatoes/determinate/galahad.jpg"
                }
            }
        }
    }
}
```

### Get detail about a specific seed in inventory with Websocket

Notice the session id, token, category and seed id are from the above response

Send message:

```json
{
    "route": "seeds",
    "type": "getDetail",
    "component": "Herb",
    "subComponent": "457a054c-487e-4fd5-91e3-3208d351e682",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk5MzcwNjAsInNlc3Npb25JRCI6ImY3YTgwMTMwLWUxNzctNDU4Mi04ZmVmLTFkZGU1NWVmOTFlMyIsInVzZXJuYW1lIjoidXNlciJ9._ZhePlTvCuS5BU64epFCrk-NoKJQCbEWeFd2yANPdkM"
}
```

Output

```json
{
    "route": "seeds",
    "type": "getDetail",
    "component": "Herb",
    "subComponent": "457a054c-487e-4fd5-91e3-3208d351e682",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "timestamp": 1709933629722,
    "data": {
        "id": "457a054c-487e-4fd5-91e3-3208d351e682",
        "category": "Herb",
        "genus": "Allium",
        "species": "schoenoprasum",
        "cultivar": "Polyvert",
        "commonName": "Chive",
        "description": "Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.",
        "hybrid": false,
        "price": 2.18,
        "perPacketCount": 100,
        "packets": 100,
        "image": "/images/herbs/polyvert_chive.jpg"
    }
}
```

### Purchase a specific seed in inventory with Websocket

Notice the session id, token, category and seed id are from the above response

Send message:

```json
{
    "route": "seeds",
    "type": "purchase",
    "component": "Herb",
    "subComponent": "457a054c-487e-4fd5-91e3-3208d351e682",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk5MzcwNjAsInNlc3Npb25JRCI6ImY3YTgwMTMwLWUxNzctNDU4Mi04ZmVmLTFkZGU1NWVmOTFlMyIsInVzZXJuYW1lIjoidXNlciJ9._ZhePlTvCuS5BU64epFCrk-NoKJQCbEWeFd2yANPdkM",
    "data": {
        "id": "7a027c43-3bf9-4dca-bfb8-fd70ccfede70",
        "quantity": "1"
    }
}
```

Output

```json
{
    "route": "seeds",
    "type": "purchase",
    "component": "Herb",
    "subComponent": "457a054c-487e-4fd5-91e3-3208d351e682",
    "sessionID": "f7a80130-e177-4582-8fef-1dde55ef91e3",
    "timestamp": 1709933729735,
    "data": {
        "id": "457a054c-487e-4fd5-91e3-3208d351e682",
        "category": "Herb",
        "genus": "Allium",
        "species": "schoenoprasum",
        "cultivar": "Polyvert",
        "commonName": "Chive",
        "description": "Suitable for growing in field or containers. Dark green leaves with very good uniformity. USDA Certified Organic.",
        "hybrid": false,
        "price": 2.18,
        "perPacketCount": 100,
        "packets": 99,
        "image": "/images/herbs/polyvert_chive.jpg"
    }
}
```
