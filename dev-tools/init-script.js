db = new Mongo().getDB("tezos_delegation");

db.delegations.createIndex({ "timestamp": 1 }, { unique: true });