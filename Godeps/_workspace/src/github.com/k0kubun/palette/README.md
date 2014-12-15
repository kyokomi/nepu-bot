# Palette

Simple string colorizing utility for Go language

## Example

```go
fmt.Println(palette.Colorize("palette", "black"))
fmt.Println(palette.Colorize("palette", "red"))
fmt.Println(palette.Colorize("palette", "green"))
fmt.Println(palette.Colorize("palette", "yellow"))
fmt.Println(palette.Colorize("palette", "blue"))
fmt.Println(palette.Colorize("palette", "magenta"))
fmt.Println(palette.Colorize("palette", "cyan"))
fmt.Println(palette.Colorize("palette", "white"))
fmt.Println(palette.Colorize("palette", "Black"))
fmt.Println(palette.Colorize("palette", "Red"))
fmt.Println(palette.Colorize("palette", "Green"))
fmt.Println(palette.Colorize("palette", "Yellow"))
fmt.Println(palette.Colorize("palette", "Blue"))
fmt.Println(palette.Colorize("palette", "Magenta"))
fmt.Println(palette.Colorize("palette", "Cyan"))
fmt.Println(palette.Colorize("palette", "White"))

fmt.Println(palette.Black("black"))
fmt.Println(palette.Red("red"))
fmt.Println(palette.Green("green"))
fmt.Println(palette.Yellow("yellow"))
fmt.Println(palette.Blue("blue"))
fmt.Println(palette.Magenta("magenta"))
fmt.Println(palette.Cyan("cyan"))
fmt.Println(palette.White("white"))
fmt.Println(palette.BoldBlack("Black"))
fmt.Println(palette.BoldRed("Red"))
fmt.Println(palette.BoldGreen("Green"))
fmt.Println(palette.BoldYellow("Yellow"))
fmt.Println(palette.BoldBlue("Blue"))
fmt.Println(palette.BoldMagenta("Magenta"))
fmt.Println(palette.BoldCyan("Cyan"))
fmt.Println(palette.BoldWhite("White"))
```

### Result
![](http://i.gyazo.com/689584c52bfa9dbee0c97b0f16fc93ef.png)

## Alternatives

- [mitchellh/colorstring](https://github.com/mitchellh/colorstring)
- [fatih/color](https://github.com/fatih/color)

## License

MIT License
