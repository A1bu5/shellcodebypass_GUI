package sgngui

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var storedOptions = &Options{
	Arch:     64,
	EncCount: 1,
	ObsLevel: 50,
}

func ShowGUI(mainWindow fyne.Window, statusLabel *canvas.Text, onSave func(*Options)) {
	opts := storedOptions

	// Create a new window for the form
	formWindow := fyne.CurrentApp().NewWindow("Configuration Options")

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Raw or EXE File Path")
	inputEntry.SetText(opts.Input)
	inputFileButton := widget.NewButton("Open File", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				inputEntry.SetText(reader.URI().Path())
			}
		}, formWindow).Show()
	})

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Output Filename")
	outputEntry.SetText(opts.Output)

	archSelect := widget.NewSelect([]string{"32-bit", "64-bit"}, func(value string) {
		if value == "32-bit" {
			opts.Arch = 32
		} else {
			opts.Arch = 64
		}
	})
	if opts.Arch == 32 {
		archSelect.SetSelected("32-bit")
	} else {
		archSelect.SetSelected("64-bit")
	}

	encCountSlider := widget.NewSlider(1, 50)
	encCountSlider.Step = 1
	encCountSlider.Value = 1
	encCountSlider.SetValue(float64(opts.EncCount))
	encCountLabel := widget.NewLabel(fmt.Sprintf("Encryption Count: %d", int(encCountSlider.Value)))
	encCountSlider.OnChanged = func(value float64) {
		encCountLabel.SetText(fmt.Sprintf("Encryption Count: %d", int(value)))
		opts.EncCount = int(value)
	}

	obsLevelEntry := widget.NewEntry()
	obsLevelEntry.SetPlaceHolder("Default:50")
	obsLevelEntry.SetText(strconv.Itoa(opts.ObsLevel))
	plainDecoderCheck := widget.NewCheck("Do not encode the decoder stub", nil)
	plainDecoderCheck.SetChecked(opts.PlainDecoder)
	asciiPayloadCheck := widget.NewCheck("Generates a full ASCII printable payload, may take very long time to bruteforce", nil)
	asciiPayloadCheck.SetChecked(opts.AsciiPayload)
	safeCheck := widget.NewCheck("Preserve all register values (a.k.a. no clobber)", nil)
	safeCheck.SetChecked(opts.Safe)
	badCharsEntry := widget.NewEntry()
	badCharsEntry.SetPlaceHolder("Specified bad characters given in hex format YOU DON'T WANT TO USE")
	badCharsEntry.SetText(opts.BadChars)
	verboseCheck := widget.NewCheck("BlueTeam Blinder", nil)
	verboseCheck.SetChecked(opts.Verbose)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Input", Widget: container.NewBorder(nil, nil, nil, inputFileButton, inputEntry)},
			{Text: "Output", Widget: outputEntry},
			{Text: "Architecture", Widget: archSelect},
			{Text: "Encode Count", Widget: encCountSlider},
			{Text: "", Widget: encCountLabel},
			{Text: "Obfuscation Level", Widget: obsLevelEntry},
			{Text: "Plain Decoder", Widget: plainDecoderCheck},
			{Text: "ASCII Payload", Widget: asciiPayloadCheck},
			{Text: "Safe Mode", Widget: safeCheck},
			{Text: "Bad Characters", Widget: badCharsEntry},
			{Text: "Verbose", Widget: verboseCheck},
		},
		OnSubmit: func() {
			opts.Input = inputEntry.Text
			opts.Output = outputEntry.Text
			opts.ObsLevel, _ = strconv.Atoi(obsLevelEntry.Text)
			opts.PlainDecoder = plainDecoderCheck.Checked
			opts.AsciiPayload = asciiPayloadCheck.Checked
			opts.Safe = safeCheck.Checked
			opts.BadChars = badCharsEntry.Text
			opts.Verbose = verboseCheck.Checked

			err := ValidateOptions(opts, formWindow)
			if err != nil {
				return
			}

			// Update the global stored options
			storedOptions = opts

			// Call your function with opts
			configureOptions(opts)
			statusLabel.Text = "Check"
			statusLabel.Color = color.RGBA{0, 255, 0, 255}
			statusLabel.Refresh()
			formWindow.Close()
			onSave(opts)
		},
	}

	formWindow.SetContent(container.New(layout.NewVBoxLayout(), form))
	formWindow.Resize(fyne.NewSize(400, 600))
	formWindow.SetFixedSize(true)
	formWindow.Show()
}

func configureOptions(opts *Options) {
	// Your existing logic here, adapted for GUI
	if opts.Verbose {
		// Set verbose mode in your utility
	}

	// Handle other options as needed
	// For example:
	println("Input: ", opts.Input)
	println("Output: ", opts.Output)
	println("Architecture: ", opts.Arch)
	println("Encode Count: ", opts.EncCount)
	println("Obfuscation Level: ", opts.ObsLevel)
	println("Plain Decoder: ", opts.PlainDecoder)
	println("ASCII Payload: ", opts.AsciiPayload)
	println("Safe Mode: ", opts.Safe)
	println("Bad Characters: ", opts.BadChars)
	println("Verbose: ", opts.Verbose)
}
