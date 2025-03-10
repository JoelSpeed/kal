package a

type (
	DefinedType  map[string]Component
	DefinedType2 *map[string]Component
	DefinedType3 map[string]int
)

type (
	TypeAlias  = map[string]Component
	TypeAlias2 = map[string]*Component
	TypeAlias3 = DefinedType
	TypeAlias4 = *DefinedType
)

type NoMapsTestStruct struct {
	Primitive     int32                `json:"primitive"`
	Components    []Component          `json:"components"`
	MapComponents map[string]Component `json:"mapComponents"` // want "MapComponents should not use a map type, use list type instead"
	DefinedType   DefinedType          `json:"definedType"`   // want "DefinedType should not use a map type, use list type instead"
	DefinedType2  DefinedType2         `json:"definedType2"`  // want "DefinedType2 should not use a map type, use list type instead"
	DefinedType3  DefinedType3         `json:"definedType3"`  // want "DefinedType3 should not use a map type, use list type instead"
	TypeAlias     TypeAlias            `json:"typeAlias"`     // want "TypeAlias should not use a map type, use list type instead"
	TypeAlias2    TypeAlias2           `json:"typeAlias2"`    // want "TypeAlias2 should not use a map type, use list type instead"
	TypeAlias3    TypeAlias3           `json:"typeAlias3"`    // want "TypeAlias3 should not use a map type, use list type instead"
	TypeAlias4    TypeAlias4           `json:"typeAlias4"`    // want "TypeAlias4 should not use a map type, use list type instead"
	MapStringInt  map[string]int       `json:"mapStringInt"`  // want "MapStringInt should not use a map type, use list type instead"
	Labels        map[string]string    `json:"specialCase"`
}

type Component struct {
	Key   string `json:"key"`
	Value int32  `json:"value"`
}
