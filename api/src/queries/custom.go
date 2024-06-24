package queries

import "context"

type QueryVectorStoreResponse struct {
	Vectors      []*VectorStore
	Documents    []*Document
	WebsitePages []*WebsitePage
}

func (q *Queries) CUSTOMQueryVectorStore(ctx context.Context, arg *QueryVectorStoreParams) (*QueryVectorStoreResponse, error) {
	rows, err := q.db.Query(ctx, queryVectorStore,
		arg.CustomerID,
		arg.Limit,
		arg.Embeddings,
		arg.Column4,
		arg.Column5,
		arg.Column6,
		arg.Column7,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vectors := make([]*VectorStore, 0)
	documents := make([]*Document, 0)
	pages := make([]*WebsitePage, 0)

	for rows.Next() {
		var v VectorStore
		var d Document
		var p WebsitePage
		if err := rows.Scan(
			&v.ID,
			&v.CustomerID,
			&v.Raw,
			&v.Embeddings,
			&v.ContentType,
			&v.ObjectID,
			&v.ObjectParentID,
			&v.Metadata,
			&v.CreatedAt,
			&d.ID,
			&d.ParentID,
			&d.CustomerID,
			&d.Filename,
			&d.Type,
			&d.SizeBytes,
			&d.Sha256,
			&d.Validated,
			&d.DatastoreType,
			&d.DatastoreID,
			&d.Summary,
			&d.SummarySha256,
			&d.VectorSha256,
			&d.CreatedAt,
			&d.UpdatedAt,
			&p.ID,
			&p.CustomerID,
			&p.WebsiteID,
			&p.Url,
			&p.Sha256,
			&p.IsValid,
			&p.Metadata,
			&p.Summary,
			&p.SummarySha256,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vectors = append(vectors, &v)
		documents = append(documents, &d)
		pages = append(pages, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &QueryVectorStoreResponse{
		Vectors:      vectors,
		Documents:    documents,
		WebsitePages: pages,
	}, nil
}
