# learnelastic

First spin up the docker-compose:
```
docker-compose up -d
```

Next run the `main.go`.

Included is a very basic wrapper for the low-level ElasticSearch client, as it doesn't
include the ablilty to decode responses / errors, and requires explcit body drain and
close like Go's low-level HTTP client.
