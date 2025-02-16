# Intro

This project is a Stoik Take Home project. It's based on bazel, using go toolchain.

## Setup

Install Bazel (>=7.4.1). Instructions for doing so can be found [here](https://docs.bazel.build/versions/master/install.html).

Also consider using `bazelisk` to automatically install the proper bazel version

## Stoik Take Home Project

This project is a take home project for Stoik. It's a simple URL shortener service that uses a SQL storage.
More info in the [original README](apps/shortener/README.md).

### How to run

```sh
bazel run //apps/shortener:main
```

### How to test

```sh
bazel test //...
```

## Containerization

The builtin `stoikth_go_image` macros make it easy to containerize
any binary target. By default, these are built for the same architecture as the current host.

You can manually override this and build for a different platform using the
`--platforms=//tools/platforms:container_{arch}_linux` build flag.

```sh
# Build default image for current host architecture and load it into docker
bazel run //apps/shortener:shortener_img_load_docker
# Build default image for specific architecture
bazel run --platforms=//tools/platforms:container_x86_64_linux //apps/shortener:shortener_img_load_docker
```

## Project Vision & Technical Considerations

### Decision Process
- Chose Go for its performance, simplicity, and robust tooling.
- SQL storage was selected to ensure transactional consistency and ease of deployment.
- Dockerization allows uniform environments across development and production.
- Edge cases (e.g. invalid URLs, database errors, duplicate entries) are handled with input validations and error logging.

### Production Readiness
- Enhance logging, monitoring, and alerting (e.g., integration with systems like Prometheus and Grafana or Datadog).
- Consider switching to a more scalable database (e.g., PostgreSQL) along with migrations and backups.
- Enforce stricter input sanitation and implement rate-limiting and authentication for production security.
- Use orchestration tools (e.g., Kubernetes) for managing deployment, scaling, and fault tolerance, along with CI/CD pipelines.

### Future Enhancements
- Provide an admin dashboard for monitoring and management.
- Improve logging and error handling for better debugging and troubleshooting.