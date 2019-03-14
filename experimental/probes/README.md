# Probes

## Probe Types

### Ping Server

This probe doesn't export any metrics, it's required for the **Ping Client** to run. 

#### Example
```
go run cmd/main.go --mode=ping-server --metric-bind-address=:8070 --ping-server-bind-address=0.0.0.0:8081
```

  
### Ping Client

This probe exports the `probes_in_cluster_network_latency` metric.

#### Example

```
go run cmd/main.go --mode=ping-client --metric-bind-address=:8071 --ping-server-address=127.0.0.1:8081
```

  

## Go Modules

This project uses [Go Modules].


[Go Modules]: https://github.com/golang/go/wiki/Modules
