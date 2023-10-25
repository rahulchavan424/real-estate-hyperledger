## Code Changes Require Self-Compilation Before Docker Deployment

Backend:

1. Navigate to the 'server' directory and run `./build.sh`.
2. In the 'docker-compose.yml' file, configure the local image as follows: `fabric-realty.server:latest`.

Frontend:

1. Go to the 'web' directory and execute `./build.sh`.
2. In the 'docker-compose.yml' file, configure the local image as follows: `fabric-realty.web:latest`.

## Support for Local Development Mode

Backend:

1. Change the configuration file path in `server/blockchain/sdk.go` to `configPath = "config-local-dev.yaml"`.
2. Execute `go run main.go`.

Frontend:

1. Modify the backend API address in `web/vue.config.js` to `http://127.0.0.1:8888`.
2. Run `yarn install` to download dependencies.
3. Execute `yarn run dev` to start the development server.