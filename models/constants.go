package models

const (
	// General error messages.
	errNoInsertHasPrimaryKey = "%s has a primary key and cannot be inserted"

	// ObjectForm error messages.
	errObjectFormMustHaveKind = "ObjectForm must have a kind"

	// ObjectShadow error messages.
	errObjectShadowMustHaveFormID = "ObjectShadow must have a form ID"

	// ObjectCommit error messages.
	errObjectCommitMustHaveFormID   = "ObjectCommit must have a form ID"
	errObjectCommitMustHaveShadowID = "ObjectCommit must have a shadow ID"

	// ObjectHead error messages.
	errObjectHeadMustHaveContextID = "ObjectHead must have a context ID"
	errObjectHeadMustHaveCommitID  = "ObjectHead must have a commit ID"

	// ObjectContext error messages.
	errObjectContextMustHaveName = "ObjectContext must have a name"
)
