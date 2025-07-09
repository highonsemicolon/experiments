package server

import (
	"context"
	"time"

	proto "github.com/highonsemicolon/experiments/zerofail/proto"
	"go.mongodb.org/mongo-driver/bson"
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

	pairs := make([]bson.D, len(req.Records))
	for i, r := range req.Records {
		pairs[i] = bson.D{
			{Key: "col1", Value: r.Col1},
			{Key: "col2", Value: r.Col2},
		}
	}

	replacement := bson.M{
		"pairs":     pairs,
		"createdAt": time.Now(),
	}

	filter := bson.M{"_id": req.OrderID}
	opts := options.Replace().SetUpsert(true)
	_, err := RecordCollection.ReplaceOne(ctx, filter, replacement, opts)
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
