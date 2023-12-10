# Cron Delegation Aggregation 
Delegation aggregation cronjob which aggregates the new delegations
from the tezos blockchain.

## Commands

### Build
```bash
mage build
```

### Clean
```bash
mage clean
```

### Generate datastore and tezos API mocks
```bash
mage gen
```

### Run the cron (requires mongodb docker container to be started: `mage mongodb:start` from tezos-delegation folder)
```bash
mage run
```