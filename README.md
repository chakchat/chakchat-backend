[![Deploy to VPS](https://github.com/chakchat/chakchat-backend/actions/workflows/deploy.yaml/badge.svg)](https://github.com/chakchat/chakchat-backend/actions/workflows/deploy.yaml)

# About 
This is a messenger with secret mode where synchoronous e2e encryption is used.

# Architecture
![Architecture](./img/architecture.png)

# Observability
We make our services observable with **OpenTelemetry** **traces** and **metrics**. \
Logs are replaced with trace events and attributes in order to force using tracing.

We use OpenTelemetry collector configured in [otel-collector-config.yaml](otel-collector-config.yaml). \
Traces are visualized with **Jaeger** on `16686` port

# Deployment
For now all services are hosted on a single machine only for **Development** needs, so single docker-compose.yml file is used.

**Note**: This setup is intended for development only and is not suitable for production use. For security reasons, access to the development environment is restricted. \
In the future this all will work in k8s.

# Run
[Makefile](Makefile) is used for some frequent scenarios. \
If you want to run the backend on your local machine, you should firstly generate some keys and self-signed SSL certificates by `make gen` \
After that you can just use `make run` to start and `make down` to stop. \
(You'll also need a working docker daemon)

`make test` runs all existing tests including unit, service and intergation tests.

# Contributing
If you want to contribute to this project, please read [CONTRIBUTING.md](CONTRIBUTING.md) first.