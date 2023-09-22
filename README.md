# toll-calculator
```
docker compose up -d
```
```
docker run -d --net="host" -p 9090:9090 -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml:z prom/prometheus
```