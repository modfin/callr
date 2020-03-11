# Callr

Callr is a small server that, through the use of Twillo, can call an on-call list when an incident occurs.

Simply put, call and wake them up if something goes down.


## Testing, on local

* `ngrok http 8080`
* Add Twillio credentials and set up your environment variables to the docker compose file
* `docker-compose up`