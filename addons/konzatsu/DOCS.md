# Konzatsu

Prometheus exporter for KONAMI Sports Club congestion levels.

## Installation

### Add this repository to Home Assistant

1. In the Home Assistant Supervisor, go to **Add-ons → Add-on Store**.
2. Click the three dots in the top right and select **Repositories**.
3. Add `https://github.com/tobschli/konzatsu` and click **Save**.

The add-on will then appear in the store and can be installed and started.

### Local installation

1. Copy the `addons/konzatsu` folder to the device (e.g. via Samba share to `/addons`).
2. In **Add-ons → Add-on Store**, click **Check for updates** / refresh, then the add-on from the `local` repository appears.
3. Install and start it.

## Configuration

- `url`: The KONAMI Sports Club information page to scrape.
- `listen_port`: Port the Prometheus metrics endpoint listens on (default `2112`).

## Usage

Scrape `http://<home-assistant-host>:2112/metrics` from Prometheus (or use the
`prometheus` integration's manual scrape targets). Exports
`congestion_threshold` as a gauge (or only the documented metric) with the
current congestion threshold level.
