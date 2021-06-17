# CED (consul-ext-dns)

Project aims to act as an load balancer by pointing the A records for external providers like cloudflare to the ip address of node of healthy load balancer.

## Build

```
export CED_VERSION="0.1.1"
# Create git tags
git tag -a v$CED_VERSION
git push origin v$CED_VERSION
# Run go releaser for binary releases
goreleaser release --rm-dist
# Run docker build push
docker buildx build --platform linux/arm64,linux/arm/v7,linux/arm/v6,linux/amd64 --build-arg CED_VERSION=$CED_VERSION --push -t blmhemu/ced:$CED_VERSION -t blmhemu/fabio:latest .
```

## Acknowledgments

Much of the *design and helper code* is inspired and taken from [fabio](https://github.com/fabiolb/fabio)
