package a

type NoPhaseTestStruct struct {
	// +optional
	Phase *string `json:"phase,omitempty"` // want "field Phase: phase fields are deprecated and conditions should be preferred, avoid phase like enum fields"

}

// DoNothing is used to check that the analyser doesn't report on methods.
func (NoPhaseTestStruct) DoNothing() {}

type NoSubPhaseTestStruct struct {
	// +optional
	FooPhase *string `json:"fooPhase,omitempty"` // want "field FooPhase: phase fields are deprecated and conditions should be preferred, avoid phase like enum fields"

}

type SerializedPhaseTeststruct struct {
	// +optional
	FooField *string `json:"fooPhase,omitempty"` // want "field FooField: phase fields are deprecated and conditions should be preferred, avoid phase like enum fields"

}
