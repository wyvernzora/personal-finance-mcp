package types

// Annotations maps internal keys to descriptive strings for metadata.
// Use this to store system-generated annotations alongside user data.
type Annotations map[string]string

// AnnotatedObject provides a user-settable description and a set of system descriptions.
// Embed this struct in domain types to include metadata annotations.
type AnnotatedObject struct {
	// Annotations holds internal system-generated descriptors for this item.
	Annotations Annotations `json:"annotations,omitempty"`
}

// NewAnnotatedObject creates a Describable with an empty user description
// and an initialized system descriptions map.
func NewAnnotatedObject() AnnotatedObject {
	return AnnotatedObject{make(Annotations)}
}

func (o *AnnotatedObject) Annotate(key, value string) {
	o.Annotations[key] = value
}
