# CED (consul-ext-dns)

Project aims to act as an load balancer by pointing the A records in external providers like porkbun to the ip address of node of healthy load balancer.

Reasoning: If both the load balancer and ced are run in a nomad + consul cluster, then consul keeps track of all healthy instances of load balancer. ced will collect the IP address of nodes containing healthy load balancer from consul. It then sets those IP addresses as A records in a DNS provider. Since nomad + consul will ensure this job is running with health checks (to be done), we can make sure the A records are **almost** always correct. This is similar to DNS failover provided by many DNS providers, but in a selfhosted form run in your datacenter.

## Docker

```sh
docker run -v /path/to/properties:/etc/ced/ced.properties blmhemu/ced
```

The reference properties file can be found in https://github.com/blmhemu/ced/blob/main/ced.properties

## Build

See `newrelease.sh`

## Acknowledgments

Much of the *design and helper code* is inspired and taken from [fabio](https://github.com/fabiolb/fabio)
