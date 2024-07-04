package project

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type LinkedInPost struct {
	*queries.LinkedinPost
}

func (post *LinkedInPost) GetConfig(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
) (*queries.LinkedinPostConfig, error) {
	dmodel := queries.New(db)
	configRaw, err := dmodel.GetLinkedInPostConfig(ctx, &queries.GetLinkedInPostConfigParams{
		ProjectID:      utils.GoogleUUIDToPGXUUID(post.ProjectID),
		LinkedinPostID: utils.GoogleUUIDToPGXUUID(post.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get the configuration: %s", err)
	}
	config := utils.ReflectStructs[*queries.GetLinkedInPostConfigRow, *queries.LinkedinPostConfig](configRaw)
	if config == nil {
		return nil, fmt.Errorf("failed to convert the config record: %s", err)
	}

	return config, err
}

/*
Creates a new linkedin post record in the database, and attaches the project library
This post will be empty, ready for the content generation to happen
*/
func (p *Project) NewLinkedInPost(
	ctx context.Context,
	db queries.DBTX,
	ideaId string,
	title string,
) (*LinkedInPost, error) {
	dmodel := queries.New(db)

	// create the project library record
	record, err := dmodel.CreateProjectLibraryRecord(ctx, &queries.CreateProjectLibraryRecordParams{
		ProjectID:   p.ID,
		Title:       title,
		ContentType: string(ContentType_LinkedInPost),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the project library record: %s", err)
	}

	// create the linkedin post
	post, err := dmodel.CreateLinkedInPost(ctx, &queries.CreateLinkedInPostParams{
		ProjectID:        p.ID,
		ProjectLibraryID: record.ID,
		ProjectIdeaID:    utils.PGXUUIDFromString(ideaId),
		Title:            title,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the post record: %s", err)
	}

	return &LinkedInPost{post}, nil
}

func createLinkedInPostConfig(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
) {
	// parse the request
	body, valid := request.Decode[createLinkedinPostConfigRequest](w, r, p.logger)
	if !valid {
		return
	}

	// create a tx
	tx, err := pool.Begin(r.Context())
	if err != nil {
		p.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the ideas
	response, err := p.CreateLinkedInPostConfig(r.Context(), tx, &body)
	if err != nil {
		p.logger.Error("failed to generate the linkedin post config", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// return to the user
	request.Encode(w, r, p.logger, http.StatusOK, response)
}

func (p *Project) CreateLinkedInPostConfig(
	ctx context.Context,
	db queries.DBTX,
	args *createLinkedinPostConfigRequest,
) (*queries.LinkedinPostConfig, error) {
	dmodel := queries.New(db)

	config, err := dmodel.CreateLinkedInPostConfig(ctx, &queries.CreateLinkedInPostConfigParams{
		ProjectID:                 utils.GoogleUUIDToPGXUUID(p.ID),
		LinkedinPostID:            utils.PGXUUIDFromString(args.LinkedInPostId),
		MinSections:               int32(args.MinSections),
		MaxSections:               int32(args.MaxSections),
		NumDocuments:              int32(args.NumDocuments),
		NumWebsitePages:           int32(args.NumWebsitePages),
		LlmContentGenerationID:    utils.PGXUUIDFromString(args.LlmContentCenerationId),
		LlmVectorSummarizationID:  utils.PGXUUIDFromString(args.LlmVectorSummarizationId),
		LlmWebsiteSummarizationID: utils.PGXUUIDFromString(args.LlmWebsiteSummarizationId),
		LlmProofReadingID:         utils.PGXUUIDFromString(args.LlmProofReadingId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the linkedin post config: %s", err)
	}

	return config, nil
}

func linkedInPostHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		p *Project,
		post *LinkedInPost,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		projectHandler(func(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, p *Project) {
			// scan the id into a uuid
			idRaw := chi.URLParam(r, "linkedInPostId")
			id, err := uuid.Parse(idRaw)
			if err != nil {
				p.logger.Error("Invalid documentId", "documentId", idRaw)
				http.Error(w, fmt.Sprintf("Invalid documentId: %s", idRaw), http.StatusBadRequest)
				return
			}

			// parse as a docstore doc
			dmodel := queries.New(pool)
			post, err := dmodel.GetLinkedInPost(r.Context(), id)
			if err != nil {
				p.logger.Error("Error getting the linkedin post", "error", err)
				http.Error(w, fmt.Sprintf("There was an internal issue: %s", err), http.StatusInternalServerError)
				return
			}

			// pass to the handler
			handler(w, r, pool, p, &LinkedInPost{post})
		}),
	)
}

func generateLinkedInPost(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
	post *LinkedInPost,
) {
	// parse the request
	body, valid := request.Decode[generateLinkedInPostRequest](w, r, p.logger)
	if !valid {
		return
	}

	// create a tx
	tx, err := pool.Begin(r.Context())
	if err != nil {
		p.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// generate
	response, err := p.GenerateLinkedInPost(r.Context(), tx, post, &body)
	if err != nil {
		p.logger.Error("failed to generate the linkedin post", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// return to the user
	request.Encode(w, r, p.logger, http.StatusOK, response)
}

func (p *Project) GenerateLinkedInPost(
	ctx context.Context,
	db queries.DBTX,
	post *LinkedInPost,
	args *generateLinkedInPostRequest,
) (*generateLinkedInPostResponse, error) {
	logger := p.logger.With("linkedinPostId", post.ID.String(), "args", *args)
	logger.InfoContext(ctx, "Fetching the configurations ...")

	// get the post config
	config, err := post.GetConfig(ctx, logger, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get linkedin post config: %s", err)
	}

	// get the generation model
	genModel, err := llm.GetLLM(ctx, db, p.CustomerID, config.LlmContentGenerationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the generation llm: %s", err)
	}

	logger.InfoContext(ctx, "Configuring the conversation ...")

	// get the conversation
	conv, err := conversation.AutoConversation(
		ctx, logger, db, p.CustomerID, args.ConversationId, prompts.LINKEDIN_POST_SYSTEM,
		fmt.Sprintf("LinkedIn Post: %s-conv", post.Title),
		"linkedin-post",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get the conversation: %s", err)
	}

	// create a prompt
	var prompt string
	if conv.New {
		// first user message
		prompt = fmt.Sprintf("Title: %s\nWhat the post should be about: %s", post.Title, args.Input)
	} else {
		// provide feedback
		prompt = fmt.Sprintf("This was not quite what I am looking for. Please try again with this feedback: %s", args.Input)
	}

	logger.InfoContext(ctx, "Sending completion request ...")

	// create the post
	response, err := conv.Completion(ctx, db, genModel, gollm.NewUserMessage(prompt), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send the completion: %s", err)
	}

	return &generateLinkedInPostResponse{
		ConversationId: conv.ID,
		Messages:       conv.GetMessages(),
		LatestMessage:  response.Message.Message,
	}, nil
}
