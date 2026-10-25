# konzatsu

> [!NOTE]
> This project was written with AI assistance. Don't expect this to work well.

Prometheus exporter for KONAMI Sports Club congestion levels.

Scrapes the club's info page and exports `congestion_threshold` (gauge, 1–3) on `/metrics`.

```
docker build -f Dockerfile addons/konzatsu -t konzatsu
docker run -p 2112:2112 konzatsu -url <URL> -listen :2112
```

Also available as a Home Assistant add-on — see [`addons/konzatsu/DOCS.md`](addons/konzatsu/DOCS.md).
