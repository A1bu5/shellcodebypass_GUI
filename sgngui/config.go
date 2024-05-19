package sgngui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var (
	Version = "?"
)

type Options struct {
	Input        string
	Output       string
	Arch         int
	EncCount     int
	ObsLevel     int
	PlainDecoder bool
	AsciiPayload bool
	Safe         bool
	BadChars     string
	Verbose      bool
}

func ValidateOptions(opts *Options, window fyne.Window) error {
	if opts.Input == "" {
		dialog.ShowError(errors.New("input file parameter is mandatory"), window)
		return errors.New("input file parameter is mandatory")
	}
	if opts.Output == "" {
		dialog.ShowError(errors.New("output file parameter is mandatory"), window)
		return errors.New("output file parameter is mandatory")
	}
	return nil
}
