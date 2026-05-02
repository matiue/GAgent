# System Monitor Agent

A lightweight monitoring and logging agent that runs as a service, reads Linux `/proc` data, and broadcasts it over gRPC.

## Current Status

- **Agent** – written in Go, compiles to a static binary for easy deployment (e.g., as a `systemd` service).  
- **Receiver** – **not yet implemented**; planned to be written in Python.

## Future Features (for the Receiver)

- [ ] Real‑time anomaly detection  
- [ ] Webhook‑based alerting  
- [ ] Minimal dashboard  
- [ ] Optional logging
