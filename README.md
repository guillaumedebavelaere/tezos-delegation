# tezos-delegation 

tezos delegation is composed of a delegation aggregation cronjob which aggregates the new delegation 
from the tezos blockchain and an api to expose the data.

## Requirements
- [go](https://go.dev/)
- [mage](https://magefile.org/)
- [docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/install/)

## Setup
In the dev-tools, cron.delegation_aggregation and delegation_api folders,
copy and paste the `.env.dist` file to `.env` and fill the missing values.

## Development

You can do a `mage` to get a list of available commands.

### Build all
```bash
mage build
```

### Lint all
```bash
mage lint
```

### Run all tests
```bash
mage test:unit
```

### Start a mongodb docker container
```bash
mage mongodb:start 
```

### Run delegation aggregation cron                                
```bash
cd cron.delegation_aggregation
mage run
```

### Run delegation api server
```bash
cd service.delegation_api
mage run
```

### Calling the api endpoint
```bash
curl --location 'http://localhost:8088/xtz/delegations' | jq
```
jq is a lightweight command-line JSON processor https://jqlang.github.io/jq/

Or load the `dev-tools/Tezos.postman_collection.json` file in postman.

