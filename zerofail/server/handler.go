package server

import (
	"context"
	"fmt"
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

	now := time.Now()
	var models []mongo.WriteModel
	var idsToCleanup []string

	for i, r := range req.Records {
		docID := fmt.Sprintf("%s_%02d", req.OrderID, i)
		idsToCleanup = append(idsToCleanup, docID)

		upsert := mongo.NewReplaceOneModel().
			SetFilter(bson.M{"_id": docID}).
			SetReplacement(bson.M{
				"_id":       docID,
				"col1":      r.Col1,
				"col2":      r.Col2,
				"createdAt": now,
			}).
			SetUpsert(true)

		models = append(models, upsert)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := RecordCollection.BulkWrite(ctx, models, opts)
	if err != nil {
		// on any failure, rollback all
		_, _ = RecordCollection.DeleteMany(ctx, bson.M{
			"_id": bson.M{"$in": idsToCleanup},
		})

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
