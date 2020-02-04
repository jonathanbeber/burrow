# Burrow

## What is it?

Burrow is a proxy that accepts simple (conventional) DNS requests and proxy it to DNS servers running with [DNS over TLS (DoT)](https://en.wikipedia.org/wiki/DNS_over_TLS). Once you start it, a Gopher burrows your plain-text DNS traffic to a cryptographed server.

DoT provides privacy and security improvements taking advantage of encrypted DNS traffic. To allow clients that don't support DoT, basically, Burrow accepts plain text DNS requests and upstream it to servers using TLS.

## Getting started

### Docker (recommended)

**Note**: A Docker image is ready on [DockerHub](https://hub.docker.com/r/jonathanbeber/burrow)

Run the following command to build this image:
```
docker build -t jonathanbeber/burrow .
```

Here's an example running its image and mapping the 53 TCP and UDP port on the same ports on the docker host:
```
docker run -itp 53:53 -p 53:53/UDP burrow
```
### Running manually

This projects requires a [Go Lang](https://golang.org/) installation and uses [dep](https://github.com/golang/dep) to manage dependencies:
```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
```

After that, start the server:

```
go run main.go
```

### Using it

A valid test can be achieved using `dig`:
```
$ dig +short cloudflare.com @127.0.0.1
198.41.215.162
198.41.214.162
```

Also, if it's accepting TCP connections:
```
$ dig +short +tcp cloudflare.com @127.0.0.1
198.41.215.162
198.41.214.162
```

Configuring the system's DNS resolver to the localhost address (for example, configuring `/etc/resol.conf` to set `127.0.0.1` as your main `nameserver`) will make all the host's DNS traffic to use DoT.

## Implementation

Burrow is writen in [Go Lang](https://golang.org/) and relies on the [miekg/dns](https://github.com/miekg/dns) library. [Miekg/dns](https://github.com/miekg/dns) library is used for great projects as [coredns](https://github.com/coredns/coredns), [godns cache](https://github.com/kenshinx/godns) and [mesos-dns](https://mesosphere.github.io/mesos-dns/).

## Configuration

Burrow supports the following environment variable as configuration:

| Env var name               | Default value | Description |
|----------------------------|---------------|-------------|
| BURROW_UPSTREAM_TIMEOUT | `2000ms`      | Duration period that Burrow waits for an upstream server response before canceling the request |
| BURROW_UPSTREAM_SERVER  | `1.1.1.1`     | The upstream server that will receive the encrypted connection |
| BURROW_UPSTREAM_PORT    | `853`         | Port that the upstream DoT server is running |
| BURROW_ENABLE_TCP       | `true`        | Enable or disable the TCP incoming requests |
| BURROW_ENABLE_UDP       | `true`        | Enable or disable the UDP incoming requests |

## Security concerns

Burrow enables encrypted connection to upstream DoT servers, but all the traffic until this service, including its responses to clients, still not secure. When using it, you will have to ensure that all the communication between your client and this service is secure. For example, if you host this service in a public address and your DNS client points to it over public internet access, you can be a victim of a [man in the middle attack](https://en.wikipedia.org/wiki/Man-in-the-middle_attack). The usage of this service on a controlled network environment increases the security level.

## Usage examples

### Localhost

It's possible to secure all DNS traffic that relies on the system DNS resolver configuration. Usually, all the process running on a system follow the configuration defined by this file to find the DNS server to make requests. Running this service locally and configuring the system's DNS config to localhost all the DNS operations will be encrypted once it leaves the local machine network layer. Look at the following diagram:

```
 +---------------------------------+
 | Local machine                   |
 |-----------------------------    |                  +---------------+
 | +-------+                       |                  |DoT External   |
 | |Exp APP+---+                   |                  |---------------|
 | +-------+   |                   |                  |               |
 | +-------+   |                   |                  |               |
 | |Browser+---+                   |                  |               |
 | +-------+   |       +---------+ |                  |               |
 | +-------+   |     +>|Burrow   +------------------->853/TCP         |
 | |  DIG  +---+     | +---------+ |                  |               |
 | +-------+   |     |             |                  |               |
 |             |     |             |                  |               |
 |             v     +             |                  +---------------+
 +-----+53/TCP&UDP port+-----------+
```

### Docker

This scenario assumes that your docker images are running in the same docker network.

Start the burrow docker image:
```
docker run -d --rm --name=burrow jonathanbeber/burrow
```

Now, start a docker image setting it DNS to be the address from the previous image:

```
docker run -it --dns=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' burrow) YOUR_APPLICATION_IMAGE
```

All the DNS traffic from your application's container will be proxied by Burrow.

### Kubernetes


Run the Burrow container as a daemonset (so that it runs on every node) with `hostNetwork: true`. Check the [deployment file](kubernetes/deployment.yaml). Apply it to your cluster with the following command:

```
kubectl apply -f kubernetes/deployment.yaml
```

That'll open the capability that each Kubernetes Node to uses the localhost address as its own NameServer. Kubernetes uses DNS service that resolves the cluster internal names and just the Pods of this DNS service will talk with Burrow daemonset directly. Take a look at the following diagram:

```
+---------------------------------------+
| Kubernetes Node                       |
|---------------------------------------|
|  +------+     +------+    +------+    |
|  |      |     |      |    |      |    |
|  |POD 1 |     |POD 2 |    |POD 3 |    |
|  |      |     |      |    |      |    |
|  +--+---+     +---+--+    +---+--+    |
|     |             |           |       |
| +---v-------------v-----------v-----+ |
| |      Kubernetes DNS service       | |                 +------------------+
| +----------------+------------------+ |                 | DoT External     |
|           +------+                    |                 |------------------|
|           |                           |                 |                  |
|        +--v---+      +------+         |                 |                  |
|        | K8s  |      |BURROW|         |                 |                  |
|        | DNS  |      |      +-------------------------->853/TCP port       |
|        | POD  |      | POD  |         |                 |                  |
|        +---+--+      +--^---+         |                 |                  |
|            |            |             |                 |                  |
|            v            +             |                 +------------------+
+----------+53-TCP&UDP port+------------+
```

Note that, changing the `resolv.conf` file of a running Kubernetes Node after its Kubernetes DNS service POD's be already running will not generate any effect. If you use a tool like [launch configuration](https://docs.aws.amazon.com/autoscaling/ec2/userguide/LaunchConfiguration.html) like, change it, and recreate your nodes to apply your changes and make Burrow to take control of your DNS traffic. A [kops rolling-update](https://github.com/kubernetes/kops/blob/master/docs/cli/kops_rolling-update.md) like action is the most indicated in this case.

# Next steps

- Implementing cache
- Integration tests
- Better understanding in Cloud environments how Kubernetes acts without the internal names resolution. For example: `*.eu-west-1.compute.internal` addresses on AWS
