package opendream

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// EEGRecord represents the structure of the data retrieved from BigQuery.
type EEGRecord struct {
	InputEmbeddings     string  `json:"input_embeddings"`
	SeqLen              int     `json:"seq_len"`
	InputAttnMask       string  `json:"input_attn_mask"`
	InputAttnMaskInvert string  `json:"input_attn_mask_invert"`
	TargetIds           float64 `json:"target_ids"`
	TargetMask          float64 `json:"target_mask"`
	SentimentLabel      int     `json:"sentiment_label"`
	SentLevelEEG        string  `json:"sent_level_EEG"`
}

// DataStore provides methods to query BigQuery
type DataStore interface {
	QueryEEGData(ctx context.Context, projectID, rowLimit string) ([]EEGRecord, error)
}

type bigQueryStore struct{}

func NewDataStore() DataStore {
	return &bigQueryStore{}
}

// QueryEEGData queries EEG data from BigQuery
func (store *bigQueryStore) QueryEEGData(ctx context.Context, projectID, rowLimit string) ([]EEGRecord, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	queryString := fmt.Sprintf(`SELECT input_embeddings, seq_len, input_attn_mask, input_attn_mask_invert, target_ids, target_mask, sentiment_label, sent_level_EEG FROM skillful-flow-399108.texteeg.all LIMIT %s`, rowLimit)
	query := client.Query(queryString)

	it, err := query.Read(ctx)
	if err != nil {
		return nil, err
	}

	var records []EEGRecord
	for {
		var record EEGRecord
		err := it.Next(&record)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
