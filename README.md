# 2PC

## Cluster

### Load shard data

Input file:

```csv
Client, Balance
S1, 10
S2, 20
S3, 40
```

Shards input:

```csv
Shard, Cluster, Range
D1, C1, 1-1000
```

Output needs to be sharded like this:

```json
{
    "client": "S1",
    "cluster": "C1",
    "shard": "D1"
}
```

Store these JSON objects inside `shards` collection in a MongoDB cluster.

## Client

### PrintBalance

Accept a client, find the cluster, send a request to get the client balance.

### Performance

Calculate the system's performance when sending a request.
