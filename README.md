# Prebid Cache

This application stores short-term data for use in Prebid.

It exists to support Video Ads from Prebid.js, as well as prebid-native

## API

### POST /cache

Adds one or more values to the cache. Values can be given as either JSON or XML. A sample request is below.

```json
{
  "puts": [
    {
      "type": "xml",
      "value": "<tag>Your XML content goes here.</tag>"
    },
    {
      "type": "json",
      "value": [1, true, "JSON value of any type can go here."]
    }
  ]
}
```

If any of the `puts` are invalid, then it responds with a **400** none of the values will be retrievable.
Assuming that all of the values are well-formed, then the server will respond with IDs which can be used to
fetch the values later.

```json
{
  "responses": [
    {"uuid": "279971e4-70f0-4b18-bd65-5c6e7aa75d40"},
    {"uuid": "147c9934-894b-4c1f-9a32-e7bb9cd15376"}
  ]
}
```


### GET /cache?uuid={id}

Retrieves a single value from the cache. If the `id` isn't recognized, then it will return a 404.

Assuming the above POST calls have been made, here are some sample GET responses.

---

**GET** */cache?uuid=279971e4-70f0-4b18-bd65-5c6e7aa75d40*

```
HTTP/1.1 200 OK
Content-Type: application/xml

<tag>Your XML content goes here.</tag>
```

---

**GET** */cache?uuid=147c9934-894b-4c1f-9a32-e7bb9cd15376*

```
HTTP/1.1 200 OK
Content-Type: application/json

[1, true, "JSON value of any type can go here."]
```

### Limitations

This section does not describe permanent API contracts; it just describes limitations on the current implementation.

- This application does *not* validate XML. If users `POST` malformed XML, they'll `GET` a bad response too.
- The host company can set a max length on payload size limits in the application config. This limit will vary from vendor to vendor.

## Development

### Prerequisites

[Golang](https://golang.org/doc/install) 1.9.1 or greater and [Dep](https://github.com/golang/dep#installation) must be installed on your system.

### Automated tests

`./validate.sh` runs the unit tests and reformats your code with [gofmt](https://golang.org/cmd/gofmt/).
`./validate.sh --nofmt` runs the unit tests, but will _not_ reformat your code.

### Manual testing

Run `prebid-cache` locally with:

```bash
go build .
./prebid-cache
```

The service will respond to requests on `localhost:2424`, and the admin data will be available on `localhost:2525`

### Configuration

Configuration is handled by [Viper](https://github.com/spf13/viper#putting-values-into-viper).
The easiest way to set config during development is by editing the [config.yaml](./config.yaml) file.

### Docker

Prebid Cache works in Docker out of the box.

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
docker build -t prebid-cache
docker run -p 2424:2424 -t prebid-cache .
```

### Profiling

[pprof stats](http://artem.krylysov.com/blog/2017/03/13/profiling-and-optimizing-go-web-applications/) can be accessed from a running app on `localhost:2525`
