package opendream

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type handler struct {
	provider Provider
}

func NewHandler(provider Provider) *handler {
	return &handler{provider: provider}
}

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

func (h *handler) ServeCSV(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rows, err := QueryBigQuery(ctx, "your-gcp-project-id")
	if err != nil {
		http.Error(w, "Unable to query BigQuery", http.StatusInternalServerError)
		return
	}

	var records []EEGRecord

	for {
		var record EEGRecord
		err := rows.Next(&record)
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error reading from BigQuery", http.StatusInternalServerError)
			return
		}
		records = append(records, record)
	}

	// Convert records to JSON
	jsonData, err := json.Marshal(records)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func QueryBigQuery(ctx context.Context, projectID string) (*bigquery.RowIterator, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	query := client.Query(`SELECT * FROM skillful-flow-399108.texteeg.eeg2text LIMIT 1000`)
	return query.Read(ctx)
}
