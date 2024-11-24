package database

import (
	"context"
	"fmt"
	"time"

	"github.com/F24-CSE535/2pc/cli/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertShards creates a shards collection in the database.
func (d *Database) InsertShards(list []*models.ClientShard) error {
	records := make([]interface{}, 0)
	for _, item := range list {
		records = append(records, item)
	}

	_, err := d.shardsCollection.InsertMany(context.TODO(), records)

	return err
}

// InsertEvent creates a new event in the database.
func (d *Database) InsertEvent(event *models.Event) error {
	event.CreatedAt = time.Now().String()
	event.IsDone = false

	_, err := d.eventsCollection.InsertOne(context.TODO(), event)

	return err
}

// GetAggregation returns a list of suggested clients to be moved.
func (d *Database) GetAggregation(count int) ([]*models.AggregationResult, error) {
	// create the mongodb pipeline
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "type",
						Value: bson.D{
							{Key: "$in",
								Value: bson.A{"inter-shard", "cross-shard"},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "pair",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.D{
									{Key: "if",
										Value: bson.D{
											{Key: "$lt", Value: bson.A{"$sender", "$receiver"}},
										},
									},
									{Key: "then",
										Value: bson.D{
											{Key: "first", Value: "$sender"},
											{Key: "second", Value: "$receiver"},
										},
									},
									{Key: "else",
										Value: bson.D{
											{Key: "first", Value: "$receiver"},
											{Key: "second", Value: "$sender"},
										},
									},
								},
							},
						},
					},
					{Key: "amount", Value: 1},
					{Key: "type", Value: 1},
					{Key: "participants", Value: 1},
				},
			},
		},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$pair"},
					{Key: "transactionCount", Value: bson.D{{Key: "$sum", Value: 1}}},
					{Key: "totalAmount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
					{Key: "types", Value: bson.D{{Key: "$addToSet", Value: "$type"}}},
					{Key: "participants", Value: bson.D{{Key: "$addToSet", Value: "$participants"}}},
				},
			},
		},
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "transactionCount", Value: bson.D{{Key: "$gt", Value: count}}},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "account1", Value: "$_id.first"},
					{Key: "account2", Value: "$_id.second"},
					{Key: "transactionCount", Value: 1},
					{Key: "totalAmount", Value: 1},
					{Key: "types", Value: 1},
					{Key: "participants",
						Value: bson.D{
							{Key: "$reduce",
								Value: bson.D{
									{Key: "input", Value: "$participants"},
									{Key: "initialValue", Value: bson.A{}},
									{Key: "in",
										Value: bson.D{
											{Key: "$setUnion", Value: bson.A{"$$value", "$$this"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$sort",
				Value: bson.D{
					{Key: "transactionCount", Value: -1},
					{Key: "totalAmount", Value: -1},
				},
			},
		},
	}

	// run the aggregate query
	cursor, err := d.sessionsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("database failed to run aggregate: %v", err)
	}
	defer cursor.Close(context.TODO())

	// create a list of aggregation results
	list := make([]*models.AggregationResult, 0)

	// extract them
	for cursor.Next(context.TODO()) {
		var result models.AggregationResult
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %v", err)
		}

		// process each result
		list = append(list, &result)
	}

	return list, nil
}
