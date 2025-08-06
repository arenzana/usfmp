# USFM Parser Docker Image
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 usfmp && \
    adduser -u 1000 -G usfmp -s /bin/sh -D usfmp

# Set working directory
WORKDIR /app

# Copy binary from goreleaser
COPY usfmp /usr/local/bin/usfmp

# Copy sample data and documentation
COPY bsb_usfm/ /app/bsb_usfm/
COPY README.md LICENSE /app/

# Ensure binary is executable
RUN chmod +x /usr/local/bin/usfmp

# Create directories for input/output with proper permissions
RUN mkdir -p /app/input /app/output && \
    chown -R usfmp:usfmp /app

# Switch to non-root user
USER usfmp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD usfmp --help > /dev/null || exit 1

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/usfmp"]

# Default command shows help
CMD ["--help"]

# Labels for metadata
LABEL org.opencontainers.image.title="USFM Parser" \
      org.opencontainers.image.description="A comprehensive Go parser for USFM (Unified Standard Format Marker) files" \
      org.opencontainers.image.url="https://github.com/arenzana/usfmp" \
      org.opencontainers.image.documentation="https://github.com/arenzana/usfmp/blob/main/README.md" \
      org.opencontainers.image.source="https://github.com/arenzana/usfmp" \
      org.opencontainers.image.vendor="Ismael Arenzana" \
      org.opencontainers.image.licenses="MIT"