#!/usr/bin/bashio

URL=$(bashio::config 'url')
PORT=$(bashio::config 'listen_port')

bashio::log.info "Starting konzatsu on port ${PORT} (scraping ${URL})"

exec /usr/bin/konzatsu -url "${URL}" -listen ":${PORT}"
