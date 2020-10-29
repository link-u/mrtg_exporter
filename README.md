# mrtg_exporter - Prometheus exporter for MRTG

[![Go Report Card](https://goreportcard.com/badge/github.com/link-u/cradle_exporter)](https://goreportcard.com/report/github.com/link-u/cradle_exporter)

[![Build on Linux](https://github.com/link-u/mrtg_exporter/workflows/Build%20on%20Linux/badge.svg)](https://github.com/link-u/mrtg_exporter/actions?query=workflow%3A%22Build+on+Linux%22)
[![Build on macOS](https://github.com/link-u/mrtg_exporter/workflows/Build%20on%20macOS/badge.svg)](https://github.com/link-u/mrtg_exporter/actions?query=workflow%3A%22Build+on+macOS%22)
[![Build on Windows](https://github.com/link-u/mrtg_exporter/workflows/Build%20on%20Windows/badge.svg)](https://github.com/link-u/mrtg_exporter/actions?query=workflow%3A%22Build+on+Windows%22)  
[![Publish Docker image](https://github.com/link-u/mrtg_exporter/workflows/Publish%20Docker%20image/badge.svg)](https://github.com/link-u/mrtg_exporter/actions?query=workflow%3A%22Publish+Docker+image%22)

`mrtg_exporter` scrapes [MRTG](https://oss.oetiker.ch/mrtg/) 's index pages and converts to prometheus metrics.

# How to use?

## docker-compose.yml

```yaml
  mrtg-exporter:
    image: ghcr.io/link-u/mrtg_exporter
    expose:
      - 9230
    restart: always
```

## Prometheus config

```yaml
  - job_name: 'mrtg'
    metrics_path: '/probe'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: '<mrtg-exporter>:9230'
    static_configs:
      - targets:
        - 'nickname https://user:pass@mrtg.example.com/path/to/mrtg/'
```

# License

MIT
