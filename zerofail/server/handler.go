package server

import (
	"context"
	"time"

	proto "github.com/highonsemicolon/experiments/zerofail/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecordServiceServer struct {
	proto.UnimplementedRecordServiceServer
}

func (s *RecordServiceServer) InsertRecords(ctx context.Context, req *proto.InsertRequest) (*proto.InsertResponse, error) {

	session, err := Client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		var docs []any
		for _, r := range req.Records {
			docs = append(docs, bson.M{"col1": r.Col1, "col2": r.Col2, "createdAt": time.Now(), "deleted": false})
		}

		if _, err := RecordCollection.InsertMany(sc, docs); err != nil {
			_ = session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return &proto.InsertResponse{Message: "Insert failed: " + err.Error(), Success: false}, nil
	}

	return &proto.InsertResponse{Message: "Insert successful", Success: true}, nil
}

func (s *RecordServiceServer) DeleteRecords(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	session, err := Client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		for _, r := range req.Records {
			filter := bson.M{"col1": r.Col1, "col2": r.Col2}
			var record bson.M
			err := RecordCollection.FindOneAndDelete(sc, filter).Decode(&record)
			if err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}

			record["deletedAt"] = time.Now()
			record["deletedBy"] = req.DeletedBy

			if _, err := DeletedCollection.InsertOne(sc, record); err != nil {
				_ = session.AbortTransaction(sc)
				return err
			}
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return &proto.DeleteResponse{Message: "Delete failed: " + err.Error(), Success: false}, nil
	}

	return &proto.DeleteResponse{Message: "Delete successful", Success: true}, nil
}
