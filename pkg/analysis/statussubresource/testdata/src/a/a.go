package a

// +kubebuilder:object:root:=true
type Foo struct {
	Spec FooSpec `json:"spec"`
}

type FooSpec struct {
	Name string `json:"name"`
}

// +kubebuilder:object:root:=true
// +kubebuilder:subresource:status
type Bar struct { // want "root object type \"Bar\" is marked to enable the status subresource with marker \"kubebuilder:subresource:status\" but has no status field"
	Spec BarSpec `json:"spec"`
}

type BarSpec struct {
	Name string `json:"name"`
}

// +kubebuilder:object:root:=true
type Baz struct { // want "root object type \"Baz\" has a status field but does not have the marker \"kubebuilder:subresource:status\" to enable the status subresource"
	Spec   BazSpec   `json:"spec"`
	Status BazStatus `json:"status,omitempty"`
}

type BazSpec struct {
	Name string `json:"name"`
}

type BazStatus struct {
	Name string `json:"name"`
}

// +kubebuilder:object:root:=true
// +kubebuilder:subresource:status
type FooBar struct {
	Spec   FooBarSpec   `json:"spec"`
	Status FooBarStatus `json:"status"`
}

type FooBarSpec struct {
	Name string `json:"name"`
}

type FooBarStatus struct {
	Name string `json:"name"`
}

// Test that it works with 'root=true' as well
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type FooBarBaz struct { // want "root object type \"FooBarBaz\" is marked to enable the status subresource with marker \"kubebuilder:subresource:status\" but has no status field"
    Spec FooBarBazSpec `json:"spec"`
}

type FooBarBazSpec struct {
    Name string `json:"name"`
}
