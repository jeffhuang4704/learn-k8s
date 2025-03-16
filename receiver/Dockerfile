# Use Ubuntu as the base image
FROM ubuntu:latest

# Install necessary dependencies (SSL/TLS, etc.)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the precompiled binary (receiver) and cert files (cert.pem and key.pem)
COPY receiver /root/
COPY cert.pem /root/
COPY key.pem /root/

# Expose the port the app will run on
EXPOSE 8443

# Run the receiver binary
CMD ["./receiver"]

