package nanomarkup // import "nanomarkup.go"
FUNCTIONS
func Marshal(data any) ([]byte, error)
    Marshal returns the encoding data of input data.
    It traverses the value data recursively.
func Unmarshal(data []byte, v any) error
    Unmarshal parses the NanoM-encoded data and stores the result in the
    value pointed to by v. If v is nil or not a pointer, Unmarshal returns an
    InvalidArgumentError.
    It uses the inverse of the encodings that Marshal uses, allocating maps,
    slices, and pointers as necessary.
TYPES
type InvalidArgumentError struct {
	Context string
	Err     error
}
    InvalidArgumentError describes an invalid argument.
func (e *InvalidArgumentError) Error() string
type InvalidEntityError struct {
	Context string
	Entity  string
	Err     error
}
    InvalidEntityError describes an invalid entity.
func (e *InvalidEntityError) Error() string
