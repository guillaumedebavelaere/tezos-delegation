# tezos-delegation 

tezos delegation is composed of a delegation aggregation cronjob which aggregates the new delegation 
from the tezos blockchain and an api to expose the data.

## Requirements
- [go](https://go.dev/)
- [mage](https://magefile.org/)
- [docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/install/)

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

