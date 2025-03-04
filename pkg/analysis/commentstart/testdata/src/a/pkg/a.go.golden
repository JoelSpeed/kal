package pkg

type A struct {
	NoComment string `json:"noComment"` // want "field NoComment is missing godoc comment"
}

// To embed the same struct multiple times, we need to rename the type.
type Embedded A
type EmbeddedPointer A
