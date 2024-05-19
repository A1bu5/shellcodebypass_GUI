package main

import (
	"fmt"
	"image/color"
	"os/exec"
	"sgngui"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildCommand(params *sgngui.Options) string {
	args := []string{
		"sgn.exe",
		"-i", params.Input,
		"-o", params.Output,
		"-a", fmt.Sprintf("%d", params.Arch),
		"-c", fmt.Sprintf("%d", params.EncCount),
		"-M", fmt.Sprintf("%d", params.ObsLevel),
	}

	if params.PlainDecoder {
		args = append(args, "-plain")
	}
	if params.AsciiPayload {
		args = append(args, "-ascii")
	}
	if params.Safe {
		args = append(args, "-S")
	}
	if params.BadChars != "" {
		args = append(args, "-badchars", params.BadChars)
	}
	if params.Verbose {
		args = append(args, "-v")
	}

	return strings.Join(args, " ")
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Sh311C0d3CryPt0")
	statusLabel := canvas.NewText("Empty*", color.RGBA{255, 0, 0, 255})
	outputPathLabel := canvas.NewText("", color.RGBA{0, 0, 255, 255})

	var params *sgngui.Options

	startButton := widget.NewButton("sgnCONFIG", func() {
		sgngui.ShowGUI(myWindow, statusLabel, func(updatedOptions *sgngui.Options) {
			params = updatedOptions
		})
	})

	buildButton := widget.NewButton("Build", func() {
		if params == nil {
			statusLabel.Text = "Error: Parameters not configured"
			statusLabel.Color = color.RGBA{255, 0, 0, 255}
			statusLabel.Refresh()
			return
		}
		statusLabel.Text = "Progressing"
		statusLabel.Color = color.RGBA{255, 255, 0, 255}
		statusLabel.Refresh()

		command := buildCommand(params)
		fmt.Println("Executing command:", command)
		cmd := exec.Command("cmd", "/C", command)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		err := cmd.Run()
		if err != nil {
			statusLabel.Text = "Error: Command failed"
			statusLabel.Color = color.RGBA{255, 0, 0, 255}
		} else {
			statusLabel.Text = "Success"
			statusLabel.Color = color.RGBA{0, 255, 0, 255}
			outputPathLabel.Text = "Generated file: " + params.Output
			outputPathLabel.Refresh()
		}
		statusLabel.Refresh()
	})

	content := container.NewVBox(
		container.NewHBox(startButton, statusLabel),
		buildButton,
		outputPathLabel,
	)
	myWindow.SetContent(content)
	defaultSize := myWindow.Canvas().Size()
	myWindow.Resize(defaultSize.Add(fyne.NewSize(defaultSize.Width*0.3, defaultSize.Height*0.3)))
	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}
