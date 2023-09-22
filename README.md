# toll-calculator
```
docker compose up -d
```
```
docker run -d -p 9093:9090 --net="host"  -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml:z prom/prometheus
```