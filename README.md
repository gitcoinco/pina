# Piña 🍍

Piña is a clone of the [Pinata](https://www.pinata.cloud/) API,
**exclusively designed for local development environments** to streamline your development process and minimize
the need for direct usage of Pinata's services during development.

## API

```
get /ipfs/{CID}
post /pinning/pinJSONToIPFS
post /pinning/pinFileToIPFS
```

### Test

`go test`

### Run

`go build && ./pina -port 8000 -public ./public`

### Run in docker

```
make docker-build
make docker-run
```
