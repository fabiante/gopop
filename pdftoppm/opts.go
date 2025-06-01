package pdftoppm

import "strconv"

type ConvertOption = func(command *Command)

func First(page int) ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-f", strconv.Itoa(page))
	}
}

func Last(page int) ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-l", strconv.Itoa(page))
	}
}

func Resolution(dpi int) ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-r", strconv.Itoa(dpi))
	}
}

func ScaleTo(scale int) ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-scale-to", strconv.Itoa(scale))
	}
}

func JPG() ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-jpeg")
	}
}

func PNG() ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-png")
	}
}

func TIFF() ConvertOption {
	return func(command *Command) {
		command.args = append(command.args, "-tiff")
	}
}
