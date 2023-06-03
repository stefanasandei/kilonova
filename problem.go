package kilonova

import (
	"time"
)

var (
	ErrAttachmentExists = Statusf(400, "Attachment with that name already exists!")
)

type ScoringType string

const (
	ScoringTypeNone        ScoringType = ""
	ScoringTypeMaxSub      ScoringType = "max_submission"
	ScoringTypeSumSubtasks ScoringType = "sum_subtasks"
)

type Problem struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Name          string    `json:"name"`
	TestName      string    `json:"test_name"`
	DefaultPoints int       `json:"default_points"`

	Visible bool `json:"visible"`

	// Limit stuff
	TimeLimit   float64 `json:"time_limit"`
	MemoryLimit int     `json:"memory_limit"`
	SourceSize  int     `json:"source_size"`

	SourceCredits string `json:"source_credits"`

	// Eval stuff
	ConsoleInput bool `json:"console_input"`

	ScoringStrategy ScoringType `json:"scoring_strategy"`
}

type StatementVariant struct {
	// Language, ie. ro/en
	Language string `json:"lang"`
	// Format, ie. pdf/md/etc.
	Format string `json:"format"`
	// Private is true if the attachment for this statement variant is private.
	// it may be private if it's currently being worked on.
	Private bool `json:"public"`
}

type ScoredProblem struct {
	Problem
	ScoreUserID *int `json:"score_user_id"`
	MaxScore    *int `json:"max_score"`
	// For showing the published/unpublished label on front page
	IsEditor bool `json:"is_editor"`
}

// ProblemFilter is the struct with all filterable fields on the problem
// It also provides a Limit and Offset field, for pagination
// This list might be expanded as time goes on
type ProblemFilter struct {
	ID           *int    `json:"id"`
	IDs          []int   `json:"ids"`
	ConsoleInput *bool   `json:"console_input"`
	Visible      *bool   `json:"visible"`
	Name         *string `json:"name"`

	FuzzyName *string `json:"name_fuzzy"`

	// EditorUserID filter marks if the user is part of the *editors* of the problem
	// Note that it excludes factors like admin or contest editor, it's just the editors in the access section.
	EditorUserID *int `json:"editor_user_id"`

	Tags []*TagGroup `json:"tags"`

	Look        bool       `json:"-"`
	LookEditor  bool       `json:"-"`
	LookingUser *UserBrief `json:"-"`

	// Unassociated filter ensures that all returned problems are not "bound" to a problem list
	Unassociated bool `json:"-"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ProblemUpdate struct {
	Name          *string `json:"name"`
	TestName      *string `json:"test_name"`
	DefaultPoints *int    `json:"default_points"`

	TimeLimit   *float64 `json:"time_limit"`
	MemoryLimit *int     `json:"memory_limit"`
	SourceSize  *int     `json:"source_size"`

	SourceCredits *string `json:"source_credits"`

	ConsoleInput *bool `json:"console_input"`
	Visible      *bool `json:"visible"`

	ScoringStrategy ScoringType `json:"scoring_strategy"`
}

type Attachment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Visible   bool      `json:"visible"`
	Private   bool      `json:"private"`
	Exec      bool      `json:"exec"`

	LastUpdatedAt time.Time `json:"last_updated_at"`
	LastUpdatedBy *int      `json:"last_updated_by"`

	Name string `json:"name"`
	// Data []byte `json:"data,omitempty"`
	Size int `json:"data_size"`
}

type AttachmentFilter struct {
	ID *int `json:"id"`
	// ProblemID *int    `json:"problem_id"`
	Visible *bool   `json:"visible"`
	Private *bool   `json:"private"`
	Exec    *bool   `json:"exec"`
	Name    *string `json:"name"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type AttachmentUpdate struct {
	Visible *bool   `json:"visible"`
	Private *bool   `json:"private"`
	Exec    *bool   `json:"exec"`
	Name    *string `json:"name"`
}

type ProblemEvalSettings struct {
	// If header/grader files are found, this is turned on to True
	OnlyCPP bool `json:"only_cpp"`
	// Files to be included during compilation
	HeaderFiles []string `json:"header_files"`
	// Files to be included in both the
	GraderFiles []string `json:"grader_files"`
	// If problem has custom checker, this is non-empty
	CheckerName string `json:"has_checker"`
	// If problem has custom checker that is marked as legacy
	LegacyChecker bool `json:"legacy_checker"`
	// If problem has ".output_only" attachment, show only outputOnly language as option
	OutputOnly bool `bool:"output_only"`
}
