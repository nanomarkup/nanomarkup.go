package nanomarkup // import "nanomarkup.go"
FUNCTIONS
func Marshal(data any) ([]byte, error)
    Marshal returns the encoding data for the input value.
    It traverses the value recursively.
func Unmarshal(data []byte, v any) error
    Unmarshal parses the encoded data and stores the result in v. If v is nil or
    not a pointer, Unmarshal returns an InvalidArgumentError.
    It uses the inverse of the encodings that Marshal uses, allocating maps,
    slices, and pointers as necessary.
TYPES
type InvalidArgumentError struct {
	Context string
	Err     error
}
    InvalidArgumentError describes an error that occurs when an invalid argument
    is provided.
func (e *InvalidArgumentError) Error() string
    Error returns a string representation of the InvalidArgumentError.
type InvalidEntityError struct {
	Context string
	Entity  string
	Err     error
}
    InvalidEntityError describes an error that occurs when an attempt is made
    with an invalid entity.
func (e *InvalidEntityError) Error() string
    Error returns a string representation of the InvalidEntityError.
