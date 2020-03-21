# Callr

Callr is a small service that, through the use of Twillo, can call an on-call list when an incident occurs.
It is meant to be used in with together with another service that monitor uptime, eg. Statuscake, Pingdom 
or really anything that can perform a HTTP GET or POST when an call needs to be made.
 
Simply put, call and wake someone if things goes down. 

## Running on local

* Create a Twillo account and buy a phone number to be used.  
* Run `ngrok http 8080` in a terminal (Twillo need to be able make requests to Callr)
* Create a docker compose file containing the twillio and ngrok data
* `docker-compose up`

```bash 
$ git clone https://github.com/modfin/callr
$ cd callr

$ echo 'version: "3.0"
services:
  callr:
    build:
      context: ./cmd/callrd
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - "PORT=8080"
      - "TWIL_SID=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
      - "TWIL_TOKEN=YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
      - "TWIL_PHONE=+46123456789"
      - "BASE_URL=https://<NGROK_ID>.ngrok.io"
      - "DATA_PATH=/go/src/callr/store"
      - "BASIC_AUTH_USER=admin"
      - "BASIC_AUTH_PASS=qwerty"
      - "INCIDENT_TOKEN=ABCDEFGHIJ"
      - "INCIDENT_ROTTEN_DURATION=4h"
    volumes:
      - .:/go/src/callr:cached
      - gopkg:/go/pkg:cached

volumes:
  gopkg:
' > docker-compose.yml

$ docker-compose up
callr_1  | ---- API ----
callr_1  | Incident reporting: POST/GET: https://b94a4401.ngrok.io/incident?token=ABCDEFGHIJ
callr_1  |     Incident is rotten after: 4h0m0s
callr_1  | 
callr_1  | ---- GUI ----
callr_1  | Page at: https://b94a4401.ngrok.io
```



## Production

**Docker swarm example**

* Create a small node on a VPS provider such as Digital Ocean or Linode 
* Point a DNS `A` record to the IP of the node
* Install docker on the node
* Run the following (with the correct credentials)

```bash 
mkdir -p /root/callr
mkdir -p /root/callr-tls
echo '
version: "3.2"
services:
  callr:
    image: modfin/callrd:latest
    ports:
      - target: 80
        published: 80
        protocol: tcp
        mode: host
      - target: 443
        published: 443
        protocol: tcp
        mode: host
    environment:
      - "AUTO_TLS=true"
      - "AUTO_TLS_DIR=/callr-tls"
      - "DEBUG=true"
      - "TWIL_SID=XXXXXXXXXXXXXXXXXXXXXXXX"
      - "TWIL_TOKEN=YYYYYYYYYYYYYYYYYYYYYYY"
      - "TWIL_PHONE=+4612345678"
      - "BASE_URL=https://callr.example.com"
      - "DATA_PATH=/callr-data"
      - "BASIC_AUTH_USER=aUser"
      - "BASIC_AUTH_PASS=aPassword"
      - "INCIDENT_TOKEN=ABCDEFGHIJ"
    volumes:
      - /root/callr:/callr-data
      - /root/callr-tls:/callr-tls
' > docker-compose.yml

$ docker stack deploy callr -c ./docker-compose.yml

```

* Go to `https://callr.example.com` and add a on-call list
* Add the post hook `https://callr.example.com/incident?token=ABCDEFGHIJ` to you monitoring service

## TODO
* Testing
* Some clean up 
* Add SQL based DAO
* Rewriting some stuff to be able to run as a cloud function.

 
