# Ipfs Document Manager Service

The Ipfs Document Manager Service is a RESTful API service designed to manage documents in IPFS. It provides functionalities for creating, retrieving, and deleting documents.

## API Documentation
The API documentation is written in Swagger [Swagger Web UI](docs/swagger.json)

It can be accessed in local docker-compose environment [swagger_ui](http://localhost:8000/swagger/index.html).

It is generated using [swagger generation tool](https://github.com/swaggo/gin-swagger)

In case of changes in api definitions, it should be updated. Please run:
```cmd
go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency
```
## Available Endpoints

- POST `/v1/tenants/tenant_space/api/ipfs/create`: Creates a new document in IPFS.
- GET `/v1/tenants/tenant_space/api/ipfs/:id`: Retrieves an existing document from IPFS by its CID.
- DELETE `/v1/tenants/tenant_space/api/ipfs/:id`: Deletes an existing document from IPFS by its CID.

## Dependency services

[docker-compose](./deployment/docker/docker-compose.yml)

- IPFS node - [Kubo](https://docs.ipfs.tech/install/command-line/)

## Running the Service

### In [Docker container](https://docs.docker.com/engine/install/) (recommended)

1. Setup 
Add following to file [.env](.env)
```
PROJECT_NAME=ipfs_doc_manager
IPFS_HOST=ipfs-node
IPFS_API_GATEWAY_URL=http://ipfs-node:8080/ipfs
IPFS_RPC_API_PORT=5001
IPFS_LOG_LEVEL=debug
IPFSMANAGER_SERVERMODE=debug
IPFSMANAGER_LISTEN_PORT=8000
```
2. Build and run service

If necessary, copy to terminal corresponding command from [Makefile](makefile)

```make docker-compose-run```

### Locally

1. Setup
Define environment variables in terminal
```
EXPORT IPFS_HOST=0.0.0.0
EXPORT IPFS_API_GATEWAY_URL=http://0.0.0.0:8080/ipfs
EXPORT IPFS_RPC_API_PORT=5001
EXPORT IPFS_LOG_LEVEL=debug
EXPORT IPFSMANAGER_SERVERMODE=debug
EXPORT IPFSMANAGER_LISTEN_PORT=8000
EXPORT VDR_TYPE=ipfs
```

2. Build plugin

```cmd
mkdir etc
cd etc
git clone https://gitlab.eclipse.org/eclipse/xfsc/libraries/ssi/verifiable-data-registry/plugins/vdr-ipfs.git .
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=mod -buildmode=plugin -o plugins/ipfs
```
3. Build and run service

```cmd
go mod download
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /microservice
./microservice
```

## Deploying the service

Deployment is managed using [Helm](https://helm.sh)

Charts are defined in [./deployment/helm](deployment/helm)

From root

```cmd
cd deployment
helm install ipfs-document-manager-service ./helm -n <your namespace | default> --kubeconfig <path to kubeconfig of necessary cluster>
```