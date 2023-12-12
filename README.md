# tezos-delegation 

tezos delegation is a Golang project composed of a delegation aggregation cron which aggregates the new delegation 
from the tezos blockchain and an api service to expose the data.

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
curl --location 'http://localhost:8088/xtz/delegations?page=1&size=100' | jq
```
jq is a lightweight command-line JSON processor https://jqlang.github.io/jq/

Or load the `dev-tools/Tezos.postman_collection.json` file in postman.

## Architecture choices

### Project structure and build tool
It is probably not the simplest structure, but I choose to use a monorepo structure and mage as a build tool because 
it is what I know from my current experience and comfortable with (I have been working with Golang for around 2 years).
Moreover, it is a structure that can be easily scaled and deployed on kubernetes for example
with a specific deployment for each service.

### Delegation aggregation cron
The delegation aggregation cron is a Golang program which aggregates the new delegation from the tezos blockchain, 
storing it in a mongodb database.

I choose to use a mongodb database because it is a NoSQL database, which is a good fit 
for delegation data which are just a collection of key-value pairs.
Moreover, it can handle a large amount of data and can be scaled horizontally.

I choose to separate it from the api service, so it could be configured as a cron job 
and a specific schedule on kubernetes for example.
We would have to adjust the cron schedule related to the number of delegations we expect to have, and the limit parameter
of the tezos api (currently set to 100).

### Delegation api service
The delegation api service is a Golang program which exposes the delegation data stored by the cron.
It is a REST api which exposes the data in a paginated way to limit the amount of data returned.
I choose to separate it from the delegation cron, so it could be scaled easily and independently, 
for example on kubernetes with autoscaler.
                                          
## Improvements
- Add more unit tests
- Add integration tests
- Add a swagger documentation
- Add a CI/CD pipeline
- Add a kubernetes deployment
- Tezos has a websocket api, it could be interesting to use it to get the new delegations in real time

