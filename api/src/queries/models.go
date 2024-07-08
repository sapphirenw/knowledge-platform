// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type VectorizeJobStatus string

const (
	VectorizeJobStatusWaiting    VectorizeJobStatus = "waiting"
	VectorizeJobStatusInProgress VectorizeJobStatus = "in-progress"
	VectorizeJobStatusComplete   VectorizeJobStatus = "complete"
	VectorizeJobStatusError      VectorizeJobStatus = "error"
	VectorizeJobStatusUnknown    VectorizeJobStatus = "unknown"
	VectorizeJobStatusRejected   VectorizeJobStatus = "rejected"
)

func (e *VectorizeJobStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VectorizeJobStatus(s)
	case string:
		*e = VectorizeJobStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for VectorizeJobStatus: %T", src)
	}
	return nil
}

type NullVectorizeJobStatus struct {
	VectorizeJobStatus VectorizeJobStatus `json:"vectorizeJobStatus"`
	Valid              bool               `json:"valid"` // Valid is true if VectorizeJobStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullVectorizeJobStatus) Scan(value interface{}) error {
	if value == nil {
		ns.VectorizeJobStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.VectorizeJobStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullVectorizeJobStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.VectorizeJobStatus), nil
}

type AssetCatalog struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	CustomerID   pgtype.UUID        `db:"customer_id" json:"customerId"`
	Datastore    string             `db:"datastore" json:"datastore"`
	DatastoreKey uuid.UUID          `db:"datastore_key" json:"datastoreKey"`
	Filetype     string             `db:"filetype" json:"filetype"`
	SizeBytes    int64              `db:"size_bytes" json:"sizeBytes"`
	Sha256       string             `db:"sha_256" json:"sha256"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type AvailableModel struct {
	ID                         string             `db:"id" json:"id"`
	Provider                   string             `db:"provider" json:"provider"`
	DisplayName                string             `db:"display_name" json:"displayName"`
	Description                string             `db:"description" json:"description"`
	InputTokenLimit            int32              `db:"input_token_limit" json:"inputTokenLimit"`
	OutputTokenLimit           int32              `db:"output_token_limit" json:"outputTokenLimit"`
	Currency                   string             `db:"currency" json:"currency"`
	InputCostPerMillionTokens  pgtype.Numeric     `db:"input_cost_per_million_tokens" json:"inputCostPerMillionTokens"`
	OutputCostPerMillionTokens pgtype.Numeric     `db:"output_cost_per_million_tokens" json:"outputCostPerMillionTokens"`
	DepreciatedWarning         bool               `db:"depreciated_warning" json:"depreciatedWarning"`
	IsDepreciated              bool               `db:"is_depreciated" json:"isDepreciated"`
	CreatedAt                  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt                  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BetaApiKey struct {
	ID        uuid.UUID          `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	Expired   bool               `db:"expired" json:"expired"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogCategory struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	CustomerID   uuid.UUID          `db:"customer_id" json:"customerId"`
	ProjectID    uuid.UUID          `db:"project_id" json:"projectId"`
	Title        string             `db:"title" json:"title"`
	TextColorHex *string            `db:"text_color_hex" json:"textColorHex"`
	BgColorHex   *string            `db:"bg_color_hex" json:"bgColorHex"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPost struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	CustomerID       uuid.UUID          `db:"customer_id" json:"customerId"`
	ProjectLibraryID uuid.UUID          `db:"project_library_id" json:"projectLibraryId"`
	BlogCategoryID   pgtype.UUID        `db:"blog_category_id" json:"blogCategoryId"`
	Title            string             `db:"title" json:"title"`
	Description      string             `db:"description" json:"description"`
	Metadata         []byte             `db:"metadata" json:"metadata"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostConfig struct {
	ID                               uuid.UUID          `db:"id" json:"id"`
	CustomerID                       uuid.UUID          `db:"customer_id" json:"customerId"`
	ProjectID                        uuid.UUID          `db:"project_id" json:"projectId"`
	MainTopic                        string             `db:"main_topic" json:"mainTopic"`
	Url                              *string            `db:"url" json:"url"`
	Metadata                         []byte             `db:"metadata" json:"metadata"`
	MinSections                      int32              `db:"min_sections" json:"minSections"`
	MaxSections                      int32              `db:"max_sections" json:"maxSections"`
	DocumentsPerSection              int32              `db:"documents_per_section" json:"documentsPerSection"`
	WebsitePagesPerSection           int32              `db:"website_pages_per_section" json:"websitePagesPerSection"`
	LlmContentGenerationDefaultID    pgtype.UUID        `db:"llm_content_generation_default_id" json:"llmContentGenerationDefaultId"`
	LlmVectorSummarizationDefaultID  pgtype.UUID        `db:"llm_vector_summarization_default_id" json:"llmVectorSummarizationDefaultId"`
	LlmWebsiteSummarizationDefaultID pgtype.UUID        `db:"llm_website_summarization_default_id" json:"llmWebsiteSummarizationDefaultId"`
	LlmProofReadingDefaultID         pgtype.UUID        `db:"llm_proof_reading_default_id" json:"llmProofReadingDefaultId"`
	CreatedAt                        pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt                        pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostSection struct {
	ID                          uuid.UUID          `db:"id" json:"id"`
	BlogPostID                  uuid.UUID          `db:"blog_post_id" json:"blogPostId"`
	AdditionalInstructions      string             `db:"additional_instructions" json:"additionalInstructions"`
	Title                       string             `db:"title" json:"title"`
	Description                 string             `db:"description" json:"description"`
	AssetID                     pgtype.UUID        `db:"asset_id" json:"assetId"`
	Metadata                    []byte             `db:"metadata" json:"metadata"`
	ContentGenerationModelID    pgtype.UUID        `db:"content_generation_model_id" json:"contentGenerationModelId"`
	VectorSummarizationModelID  pgtype.UUID        `db:"vector_summarization_model_id" json:"vectorSummarizationModelId"`
	WebsiteSummarizationModelID pgtype.UUID        `db:"website_summarization_model_id" json:"websiteSummarizationModelId"`
	ProofReadingModelID         pgtype.UUID        `db:"proof_reading_model_id" json:"proofReadingModelId"`
	CreatedAt                   pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt                   pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostSectionContent struct {
	ID                uuid.UUID          `db:"id" json:"id"`
	BlogPostSectionID uuid.UUID          `db:"blog_post_section_id" json:"blogPostSectionId"`
	Content           string             `db:"content" json:"content"`
	Feedback          string             `db:"feedback" json:"feedback"`
	Index             int32              `db:"index" json:"index"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostSectionDocument struct {
	ID                uuid.UUID          `db:"id" json:"id"`
	BlogPostSectionID uuid.UUID          `db:"blog_post_section_id" json:"blogPostSectionId"`
	DocumentID        uuid.UUID          `db:"document_id" json:"documentId"`
	Query             string             `db:"query" json:"query"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostSectionWebsitePage struct {
	ID                uuid.UUID          `db:"id" json:"id"`
	BlogPostSectionID uuid.UUID          `db:"blog_post_section_id" json:"blogPostSectionId"`
	WebsitePageID     uuid.UUID          `db:"website_page_id" json:"websitePageId"`
	Query             string             `db:"query" json:"query"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt         pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type BlogPostTag struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	BlogPostID uuid.UUID          `db:"blog_post_id" json:"blogPostId"`
	Title      string             `db:"title" json:"title"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ContentType struct {
	Title     string             `db:"title" json:"title"`
	Parent    string             `db:"parent" json:"parent"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Conversation struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	CustomerID       uuid.UUID          `db:"customer_id" json:"customerId"`
	Title            string             `db:"title" json:"title"`
	ConversationType string             `db:"conversation_type" json:"conversationType"`
	SystemMessage    string             `db:"system_message" json:"systemMessage"`
	Metadata         []byte             `db:"metadata" json:"metadata"`
	HasError         bool               `db:"has_error" json:"hasError"`
	ErrorMessage     *string            `db:"error_message" json:"errorMessage"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ConversationMessage struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	ConversationID uuid.UUID          `db:"conversation_id" json:"conversationId"`
	LlmID          pgtype.UUID        `db:"llm_id" json:"llmId"`
	Model          string             `db:"model" json:"model"`
	Temperature    float64            `db:"temperature" json:"temperature"`
	Instructions   string             `db:"instructions" json:"instructions"`
	Role           string             `db:"role" json:"role"`
	Message        string             `db:"message" json:"message"`
	Index          int32              `db:"index" json:"index"`
	ToolUseID      string             `db:"tool_use_id" json:"toolUseId"`
	ToolName       string             `db:"tool_name" json:"toolName"`
	ToolArguments  []byte             `db:"tool_arguments" json:"toolArguments"`
	ToolResults    []byte             `db:"tool_results" json:"toolResults"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt      pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Customer struct {
	ID        uuid.UUID          `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	Datastore string             `db:"datastore" json:"datastore"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type CustomerLlmConfiguration struct {
	CustomerID   uuid.UUID          `db:"customer_id" json:"customerId"`
	SummaryLlmID pgtype.UUID        `db:"summary_llm_id" json:"summaryLlmId"`
	ChatLlmID    pgtype.UUID        `db:"chat_llm_id" json:"chatLlmId"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Document struct {
	ID            uuid.UUID          `db:"id" json:"id"`
	ParentID      pgtype.UUID        `db:"parent_id" json:"parentId"`
	CustomerID    uuid.UUID          `db:"customer_id" json:"customerId"`
	Filename      string             `db:"filename" json:"filename"`
	Type          string             `db:"type" json:"type"`
	SizeBytes     int64              `db:"size_bytes" json:"sizeBytes"`
	Sha256        string             `db:"sha_256" json:"sha256"`
	Validated     bool               `db:"validated" json:"validated"`
	DatastoreType string             `db:"datastore_type" json:"datastoreType"`
	DatastoreID   string             `db:"datastore_id" json:"datastoreId"`
	Summary       string             `db:"summary" json:"summary"`
	SummarySha256 string             `db:"summary_sha_256" json:"summarySha256"`
	VectorSha256  string             `db:"vector_sha_256" json:"vectorSha256"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type DocumentVector struct {
	ID            uuid.UUID          `db:"id" json:"id"`
	DocumentID    uuid.UUID          `db:"document_id" json:"documentId"`
	VectorStoreID uuid.UUID          `db:"vector_store_id" json:"vectorStoreId"`
	CustomerID    uuid.UUID          `db:"customer_id" json:"customerId"`
	Index         int32              `db:"index" json:"index"`
	Metadata      []byte             `db:"metadata" json:"metadata"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"createdAt"`
}

type Folder struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	ParentID   pgtype.UUID        `db:"parent_id" json:"parentId"`
	CustomerID uuid.UUID          `db:"customer_id" json:"customerId"`
	Title      string             `db:"title" json:"title"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type LinkedinPost struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	ProjectID        uuid.UUID          `db:"project_id" json:"projectId"`
	ProjectLibraryID uuid.UUID          `db:"project_library_id" json:"projectLibraryId"`
	ProjectIdeaID    pgtype.UUID        `db:"project_idea_id" json:"projectIdeaId"`
	Title            string             `db:"title" json:"title"`
	AssetID          pgtype.UUID        `db:"asset_id" json:"assetId"`
	Metadata         []byte             `db:"metadata" json:"metadata"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type LinkedinPostConfig struct {
	ID                        uuid.UUID          `db:"id" json:"id"`
	ProjectID                 pgtype.UUID        `db:"project_id" json:"projectId"`
	LinkedinPostID            pgtype.UUID        `db:"linkedin_post_id" json:"linkedinPostId"`
	MinSections               int32              `db:"min_sections" json:"minSections"`
	MaxSections               int32              `db:"max_sections" json:"maxSections"`
	NumDocuments              int32              `db:"num_documents" json:"numDocuments"`
	NumWebsitePages           int32              `db:"num_website_pages" json:"numWebsitePages"`
	LlmContentGenerationID    pgtype.UUID        `db:"llm_content_generation_id" json:"llmContentGenerationId"`
	LlmVectorSummarizationID  pgtype.UUID        `db:"llm_vector_summarization_id" json:"llmVectorSummarizationId"`
	LlmWebsiteSummarizationID pgtype.UUID        `db:"llm_website_summarization_id" json:"llmWebsiteSummarizationId"`
	LlmProofReadingID         pgtype.UUID        `db:"llm_proof_reading_id" json:"llmProofReadingId"`
	CreatedAt                 pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt                 pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type LinkedinPostConversation struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	LinkedinPostID uuid.UUID          `db:"linkedin_post_id" json:"linkedinPostId"`
	ConversationID uuid.UUID          `db:"conversation_id" json:"conversationId"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt      pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Llm struct {
	ID           uuid.UUID          `db:"id" json:"id"`
	CustomerID   pgtype.UUID        `db:"customer_id" json:"customerId"`
	Title        string             `db:"title" json:"title"`
	Color        *string            `db:"color" json:"color"`
	Model        string             `db:"model" json:"model"`
	Temperature  float64            `db:"temperature" json:"temperature"`
	Instructions string             `db:"instructions" json:"instructions"`
	IsDefault    bool               `db:"is_default" json:"isDefault"`
	Public       bool               `db:"public" json:"public"`
	CreatedAt    pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Project struct {
	ID                    uuid.UUID          `db:"id" json:"id"`
	CustomerID            uuid.UUID          `db:"customer_id" json:"customerId"`
	Title                 string             `db:"title" json:"title"`
	Topic                 string             `db:"topic" json:"topic"`
	IdeaGenerationModelID pgtype.UUID        `db:"idea_generation_model_id" json:"ideaGenerationModelId"`
	CreatedAt             pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt             pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ProjectFolder struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	CustomerID uuid.UUID          `db:"customer_id" json:"customerId"`
	ProjectID  uuid.UUID          `db:"project_id" json:"projectId"`
	FolderID   uuid.UUID          `db:"folder_id" json:"folderId"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ProjectIdea struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	ProjectID      uuid.UUID          `db:"project_id" json:"projectId"`
	ConversationID pgtype.UUID        `db:"conversation_id" json:"conversationId"`
	Title          string             `db:"title" json:"title"`
	Used           bool               `db:"used" json:"used"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt      pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ProjectLibrary struct {
	ID          uuid.UUID          `db:"id" json:"id"`
	ProjectID   uuid.UUID          `db:"project_id" json:"projectId"`
	Title       string             `db:"title" json:"title"`
	ContentType string             `db:"content_type" json:"contentType"`
	Draft       bool               `db:"draft" json:"draft"`
	Published   bool               `db:"published" json:"published"`
	CreatedAt   pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type ProjectWebsite struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	CustomerID uuid.UUID          `db:"customer_id" json:"customerId"`
	ProjectID  uuid.UUID          `db:"project_id" json:"projectId"`
	WebsiteID  uuid.UUID          `db:"website_id" json:"websiteId"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type TokenUsage struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	CustomerID     uuid.UUID          `db:"customer_id" json:"customerId"`
	ConversationID pgtype.UUID        `db:"conversation_id" json:"conversationId"`
	Model          string             `db:"model" json:"model"`
	InputTokens    int32              `db:"input_tokens" json:"inputTokens"`
	OutputTokens   int32              `db:"output_tokens" json:"outputTokens"`
	TotalTokens    int32              `db:"total_tokens" json:"totalTokens"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
}

type VectorStore struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	CustomerID     uuid.UUID          `db:"customer_id" json:"customerId"`
	Raw            string             `db:"raw" json:"raw"`
	Embeddings     *pgvector.Vector   `db:"embeddings" json:"embeddings"`
	ContentType    string             `db:"content_type" json:"contentType"`
	ObjectID       uuid.UUID          `db:"object_id" json:"objectId"`
	ObjectParentID pgtype.UUID        `db:"object_parent_id" json:"objectParentId"`
	Metadata       []byte             `db:"metadata" json:"metadata"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
}

type VectorStoreDefault struct {
	ID             uuid.UUID          `db:"id" json:"id"`
	CustomerID     uuid.UUID          `db:"customer_id" json:"customerId"`
	Raw            string             `db:"raw" json:"raw"`
	Embeddings     *pgvector.Vector   `db:"embeddings" json:"embeddings"`
	ContentType    string             `db:"content_type" json:"contentType"`
	ObjectID       uuid.UUID          `db:"object_id" json:"objectId"`
	ObjectParentID pgtype.UUID        `db:"object_parent_id" json:"objectParentId"`
	Metadata       []byte             `db:"metadata" json:"metadata"`
	CreatedAt      pgtype.Timestamptz `db:"created_at" json:"createdAt"`
}

type VectorizeJob struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	CustomerID uuid.UUID          `db:"customer_id" json:"customerId"`
	Documents  bool               `db:"documents" json:"documents"`
	Websites   bool               `db:"websites" json:"websites"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type VectorizeJobItem struct {
	ID        uuid.UUID          `db:"id" json:"id"`
	JobID     uuid.UUID          `db:"job_id" json:"jobId"`
	Status    VectorizeJobStatus `db:"status" json:"status"`
	Message   string             `db:"message" json:"message"`
	Error     string             `db:"error" json:"error"`
	CreatedAt pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type Website struct {
	ID         uuid.UUID          `db:"id" json:"id"`
	CustomerID uuid.UUID          `db:"customer_id" json:"customerId"`
	Protocol   string             `db:"protocol" json:"protocol"`
	Domain     string             `db:"domain" json:"domain"`
	Blacklist  []string           `db:"blacklist" json:"blacklist"`
	Whitelist  []string           `db:"whitelist" json:"whitelist"`
	CreatedAt  pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type WebsitePage struct {
	ID            uuid.UUID          `db:"id" json:"id"`
	CustomerID    uuid.UUID          `db:"customer_id" json:"customerId"`
	WebsiteID     uuid.UUID          `db:"website_id" json:"websiteId"`
	Url           string             `db:"url" json:"url"`
	Sha256        string             `db:"sha_256" json:"sha256"`
	IsValid       bool               `db:"is_valid" json:"isValid"`
	Metadata      []byte             `db:"metadata" json:"metadata"`
	Summary       string             `db:"summary" json:"summary"`
	SummarySha256 string             `db:"summary_sha_256" json:"summarySha256"`
	VectorSha256  string             `db:"vector_sha_256" json:"vectorSha256"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at" json:"updatedAt"`
}

type WebsitePageVector struct {
	ID            uuid.UUID          `db:"id" json:"id"`
	WebsitePageID uuid.UUID          `db:"website_page_id" json:"websitePageId"`
	VectorStoreID uuid.UUID          `db:"vector_store_id" json:"vectorStoreId"`
	CustomerID    uuid.UUID          `db:"customer_id" json:"customerId"`
	Index         int32              `db:"index" json:"index"`
	Metadata      []byte             `db:"metadata" json:"metadata"`
	CreatedAt     pgtype.Timestamptz `db:"created_at" json:"createdAt"`
}
