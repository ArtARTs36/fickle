# fickle

**fickle** - is proxy for Docker containers, runs container on incoming requests

## Metrics

**fickle** scraps metrics when containers are stopped, and also when your prometheus asks for them and the container is running.

Paths:
- fickle-host/metrics - metrics of **fickle**
- fickle-host/metrics/<proxy-host> - metrics of proxy app
