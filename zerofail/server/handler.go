package server

import (
	"context"
	"time"

	proto "github.com/highonsemicolon/experiments/zerofail/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordServiceServer struct {
	proto.UnimplementedRecordServiceServer
}

func (s *RecordServiceServer) UpsertRecords(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {

	// input validation
	col1Set := make(map[string]bool)
	col2Set := make(map[string]bool)
	for _, r := range req.Records {
		if col1Set[r.Col1] || col2Set[r.Col2] {
			return &proto.UpsertResponse{
				Message: "Duplicate col1 or col2 in request itself",
				Success: false,
			}, nil
		}
		col1Set[r.Col1] = true
		col2Set[r.Col2] = true
	}

	pairs := make([]bson.D, 0)
	for _, r := range req.Records {
		pair := bson.D{
			{Key: "col1", Value: r.Col1},
			{Key: "col2", Value: r.Col2},
		}
		pairs = append(pairs, pair)
	}

	now := time.Now()
	var models []mongo.WriteModel

	upsert := mongo.NewReplaceOneModel().
		SetFilter(bson.M{"_id": req.OrderID}).
		SetReplacement(bson.M{
			"pairs":     pairs,
			"createdAt": now,
		}).
		SetUpsert(true)

	models = append(models, upsert)

	opts := options.BulkWrite().SetOrdered(false)
	_, err := RecordCollection.BulkWrite(ctx, models, opts)
	if err != nil {
		return &proto.UpsertResponse{
			Success: false,
			Message: "Upsert failed: " + err.Error(),
		}, nil
	}

	return &proto.UpsertResponse{
		Success: true,
		Message: "Upsert successful",
	}, nil
}
