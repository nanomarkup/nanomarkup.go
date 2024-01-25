package nanomarkup // import "nanomarkup.go"
FUNCTIONS
func Marshal(data any) ([]byte, error)
    Marshal returns the Nano Markup encoding of data.
    It traverses the value data recursively.
