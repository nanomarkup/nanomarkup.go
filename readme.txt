package nanomarkup // import "nanomarkup.go"
FUNCTIONS
func Compact(dst *bytes.Buffer, src []byte) error
    Compact appends the nano-encoded src to dst, eliminating insignificant space
    characters.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error
    Indent function appends to `dst` the nano-encoded source (`src`) in an
    indented format. The data appended to dst does not begin with the prefix
    nor any indentation, to make it easier to embed inside other formatted
    nano-encoded data.
func Marshal(data any, meta *nanometadata.Metadata) ([]byte, error)
    Marshal returns the encoding data for the input value.
    It traverses the value recursively.
func MarshalIndent(data any, prefix, indent string) ([]byte, error)
    MarshalIndent is like Marshal but applies Indent to format the output.
func Unmarshal(data []byte, v any, meta *nanometadata.Metadata) error
    Unmarshal parses the encoded data and stores the result in v. If v is nil or
    not a pointer, Unmarshal returns an InvalidArgumentError.
    It uses the inverse of the encodings that Marshal uses, allocating maps,
    slices, and pointers as necessary.
