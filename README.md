# Octopus

Simple and Tested HTTP Server in Go for CI Experimentations.

This repository contains a simple HTTP server implemented in Go, designed for CI experimentations. The server includes features such as Dockerfile configuration, end-to-end tests using Playwright, unit testing, an endpoint with a database connection, deployment with Kustomization in staging and production environments, and support for canary deployment.

## Features

- **HTTP Server in Go:** A basic HTTP server written in Go to serve as a foundation for CI experiments.

- **Dockerfile:** Includes a Dockerfile for containerizing the Go application.

- **End-to-End Tests with Playwright:** Provides end-to-end tests using Playwright to ensure the reliability of the application.

- **Unit Testing:** Includes unit tests for the Go codebase.

- **Database Connection:** Demonstrates an endpoint with a database connection for more realistic scenarios.

- **Deployment with Kustomization:**
  - Configurations for staging and production environments using Kustomize.
  - Secrets for sensitive information like `DATABASE_PASSWORD` and private registry credentials.

- **Canary Deployment:** Supports canary deployments for gradual rollouts.

## Prerequisites

Before you begin, ensure you have the following:

- [Docker](https://www.docker.com/get-started)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) for deploying to Kubernetes clusters.
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) for customizing Kubernetes manifests.

## Configuration

Before deploying, make sure to configure the necessary secrets:

### Secrets

1. **DATABASE_PASSWORD:**
   - Create a secret named `database-secret` in your Kubernetes cluster with the key `DATABASE_PASSWORD`. Update the value accordingly.

      ```bash
      kubectl create secret generic database-secret --from-literal=DATABASE_PASSWORD=your_password_here
      ```

2. **Private Registry Credentials:**
   - Create a secret named `registry-secret` in your Kubernetes cluster with the Docker registry credentials. You can create the `.dockerconfigjson` file and then create the secret using:

      ```bash
      kubectl create secret generic registry-secret --from-file=.dockerconfigjson=/path/to/your/.dockerconfigjson
      ```

   Ensure that your `.dockerconfigjson` contains credentials for your private Docker registry.

## Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/michaelahli/octopus.git
   cd octopus
   ```

2. Build the Docker image:

   ```bash
   docker build -t octopus:latest .
   ```

3. Run the Docker container:

   ```bash
   docker run -p 8080:8080 -e DATABASE_PASSWORD=your_database_password -d octopus:latest
   ```

   Replace `your_database_password` with the actual database password.

4. Access the application at [http://localhost:8080](http://localhost:8080).

## Deployment

### Staging Environment

Deploy to the staging environment using Kustomize:

```bash
cd deploy/overlays/staging
kustomize edit set image host.docker.internal:30500/octopus=host.docker.internal:5000/octopus:latest
kustomize build | kubectl apply -f -
```

### Production Environment

Deploy to the production environment using Kustomize:

```bash
cd deploy/overlays/production
kustomize edit set image host.docker.internal:30500/octopus=host.docker.internal:5000/octopus:latest
kustomize build | kubectl apply -f -
```

### Canary Deployment

For canary deployment, modify the Kustomization file in the `k8s/canary` directory as needed and apply the changes:

```bash
cd deploy/overlays/canary
kustomize edit set image host.docker.internal:30500/octopus=host.docker.internal:5000/octopus:latest
kustomize build | kubectl apply -f -
```