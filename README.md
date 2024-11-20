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

### PrintLogs

Accept a target and prints its logs.

### PrintDatastore

Accept a target and prints its committed transactions.

### 2PC

In order to process a transaction, the client reads transactions from a `csv` file. For each transaction, the client triggers an inter-shard or cross-shard procedure based on the clusters that are needed in processing the transaction.

1. The client sends `prepare` to all participates.
2. The client receives `ack` from all participates.
3. The client sends `commit` or `abort` based on the `ack` values.

For concurrent transaction handling, the client assigns a `sessionId` which will be used for each transaction.

## Cluster

When a cluster is running, it clones its data from the `global` database. For each of its nodes, it runs a new process manager. A process manager starts a node, and waits for input commands from the cluster manager. For each node, there is collection in the cluster's database.

The cluster manager looks a collection called `events` in the database. In an interval, it get's all events that are belong to it from that collection and performs a logic (`scale-up`, `scale-down`, `reshard`). After that, it will mark that event as done.

### Node

Each node has a `gRPC` interface that accepts `RPC` calls from both client and other nodes. A list of these `RPC` calls are as follow.

- `PrintBalance` : accepts a client name and returns its balance.
- `PrintLogs` : returns all write-ahead logs inside this node.
- `PrintDatastore` : returns all committed transactions inside this node.

### 2PL

Each server has a local lock for input transactions. If a lock is being set, then the input transaction will be aborted. Otherwise, the server accepts the input transaction. The lock table is a key-value map that tells wheter a record is being used by the processor or not.

### WAL

For each operation, the server logs the operations in the following manner.

```
<T1, Start>
<T1, Update, Record, New>
<T1, Commit/Abort>
```
