package a

type StructWithRequiredField struct {
	// tsring is a string field.
	// +required
	String string `json:"string"`
}
