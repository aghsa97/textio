package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var mode = 0
var source_file = ""
var text_buffer = [][]rune{}

func read_file(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		text_buffer = append(text_buffer, []rune{})
		os.Exit(1)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		text_buffer = append(text_buffer, []rune(line))
	}
}

func draw_text_buffer() {
	for row, line := range text_buffer {
		for col, ch := range line {
			termbox.SetCell(col, row, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func print_message(col, row int, fg, bg termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}

func get_event_key() termbox.Event {
	event := termbox.PollEvent()
	return event
}

func handle_event_key() {
	key_event := get_event_key()

	if key_event.Type == termbox.EventKey && key_event.Key == termbox.KeyEsc {
		mode = 0
	}

	if mode == 0 {
		switch key_event.Ch {
		case 'q':
			termbox.Close()
			os.Exit(0)
		case 'e':
			mode = 1
		}
	}
}

func print_status_bar() {
	var mode_status string
	_, row := termbox.Size()
	if mode == 0 {
		mode_status = "VIEW: "
	} else {
		mode_status = "EDIT: "
	}
	message := mode_status + source_file + " lines: " + fmt.Sprintf("%d", len(text_buffer)) + " Press e to edit or q to quit."
	print_message(0, row-1, termbox.ColorBlack, termbox.ColorWhite, message)
}

func run_editor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		source_file = os.Args[1]
	} else {
		source_file = "out.txt"
		file, err := os.Create(source_file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.WriteFile(source_file, []byte("Created Successfully!"), 0644)
		defer file.Close()
	}
	read_file(source_file)

	for {
		print_status_bar()
		draw_text_buffer()
		termbox.Flush()
		handle_event_key()
	}
}

func main() {
	run_editor()
}
