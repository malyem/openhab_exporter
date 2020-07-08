# OpenHAB exporter for Prometheus

Implementation in Go ready to export [Openhab](https://www.openhab.org/) metrics to [Prometheus](https://prometheus.io/).

Exported are metrics from items types: `Number`, `Dimmer`, `Switch`,`Contact`.

## Example usage
Only one parameter is required `--apiurl`, which should be address of your Openhab installation.
```
-> % ./openhab_exporter --apiurl=http://192.168.0.116:8080
level=info ts=2020-07-08T15:34:45.713Z caller=main.go:42 msg="Starting openhab_exporter" version="(version=, branch=, revision=)"
level=info ts=2020-07-08T15:34:45.713Z caller=main.go:43 build_context="(go=go1.14.4, user=, date=)"
level=info ts=2020-07-08T15:34:45.879Z caller=main.go:47 msg="Listening on address" address=:9266
```

Open new terminal window and check results:
```
-> % curl -v localhost:9266/metrics
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 9266 (#0)
> GET /metrics HTTP/1.1
> Host: localhost:9266
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; version=0.0.4; charset=utf-8
< Date: Wed, 08 Jul 2020 15:35:24 GMT
< Transfer-Encoding: chunked
<
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0

...

# HELP openhab_item_state_current Openhab items current state
# TYPE openhab_item_state_current gauge
openhab_item_state_current{groupnames="",item="ActionsPower",label="Power on/off",tags="",type="Switch"} 0
openhab_item_state_current{groupnames="",item="AlarmAlert",label="Alarm alert",tags="",type="Switch"} 0
openhab_item_state_current{groupnames="",item="Gateway_AddDevice",label="",tags="",type="Switch"} 0

...

```
