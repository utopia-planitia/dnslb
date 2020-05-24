# DNSLB

DNS Load Balancer runs as a sidecar for ingress and apiserver pods.

Its purpose is to manage DNS entries, during node drains and failures.

## Life cycle

While running, it checks the local port and adds or removes (on failure) its own IP to the subdomain.

Process termination removes the local IP entries from the subdomain and waits 120 seconds for DNS TTL.

In case the node shut down without warning, a cronjob checks all DNS endpoints and removes unhealthy IPs.

## Config

### System variables

| Name         | Default value | Required | Purpose                                                       |
|--------------|---------------|----------|---------------------------------------------------------------|
| CF_API_TOKEN | ""            |          | Used to log in to cloudflares API.                            |
| CF_API_KEY   | ""            |          | Used together with CF_API_EMAIL to log in to cloudflares API. |
| CF_API_EMAIL | ""            |          | Used together with CF_API_KEY to log in to cloudflares API.   |
| CF_ZONE      |               | yes      | Defines the Domain to use.                                    |

### Arguments

| Name         | Default value | Required | Purpose                                                       |
|--------------|---------------|----------|---------------------------------------------------------------|
| --subdomain  |               | yes      |                                                               |
| --ports      |               | yes      | Ports required to be open. Can be used multiple times.        |
| --ipv4       | true          | no       | Enable or disable IPv4 for networking.                        |
| --ipv6       | true          | no       | Enable or disable IPv6 for networking.                        |
