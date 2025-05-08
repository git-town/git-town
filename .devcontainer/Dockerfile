FROM mcr.microsoft.com/devcontainers/go:1.23

# Generate Spanish locale for E2E tests
RUN sed -i 's/# es_ES.UTF-8/es_ES.UTF-8/' /etc/locale.gen && \
    locale-gen
