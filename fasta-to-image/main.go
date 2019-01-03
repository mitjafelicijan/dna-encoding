package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/profile"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// Point ...
type Point struct {
	Row    int
	Column int
	Text   string
	Color  color.RGBA
}

// Data ...
type Data struct {
	Points []Point
}

var (
	appName = "dnae-png"
	version = "0.0.1"

	rowCount    int
	columnCount int

	app      = kingpin.New(appName, "A command-line application that encodes FASTA file into PNG image.")
	input    = app.Flag("input", "Input FASTA file which will be encoded into PNG image.").PlaceHolder("INPUT").Required().Short('i').String()
	output   = app.Flag("output", "Output file in PNG format that represents DNA sequence in graphical way.").Short('o').Default("out.png").String()
	rectSize = app.Flag("size", "Size of pairings of DNA bases on image in pixels (lower resolution lower file size).").PlaceHolder("10").Short('s').Default("10").Int()
	pprof    = app.Flag("pprof", "Generates pprof file for profiling and debugging purposes.").Short('p').Bool()
)

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func main() {

	app.Version(fmt.Sprintf("%s %s", appName, version))
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *pprof {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	if *input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()

	log.Println("Gathering input file's stats ...")

	data := Data{
		Points: make([]Point, 0),
	}

	log.Println("Deconstructing FASTA file ...")

	file, _ := os.Open(*input)
	fileScanner := bufio.NewScanner(file)
	lineNumber := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if !strings.HasPrefix(line, ">") {
			for k, v := range line {

				var block color.RGBA
				switch string(v) {
				case "A":
					block = color.RGBA{0, 0, 255, 255}
				case "G":
					block = color.RGBA{0, 100, 0, 255}
				case "C":
					block = color.RGBA{255, 0, 0, 255}
				case "T":
					block = color.RGBA{255, 255, 0, 255}
				}

				point := Point{
					Row:    lineNumber,
					Column: k,
					Text:   string(v),
					Color:  block,
				}
				data.Points = append(data.Points, point)
			}
			lineNumber++

			if len(line) > columnCount {
				columnCount = len(line)
			}
		}
	}

	log.Println("Compositing image file ...")

	bar := pb.StartNew(len(data.Points))
	rowCount = lineNumber
	dnaImage := image.NewRGBA(image.Rect(0, 0, columnCount**rectSize, rowCount**rectSize))

	for _, v := range data.Points {
		rect := image.Rect(v.Column**rectSize, v.Row**rectSize, v.Column**rectSize+*rectSize, v.Row**rectSize+*rectSize)
		draw.Draw(dnaImage, rect, &image.Uniform{v.Color}, image.ZP, draw.Src)
		bar.Increment()
	}

	bar.Finish()

	log.Println("Saving output file ...")

	file, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed create file: %s", err)
	}
	defer file.Close()

	png.Encode(file, dnaImage)

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Output image file length is", byteCountDecimal(fi.Size()))

	elapsed := time.Since(start)

	log.Printf("Process took %s\n", elapsed)
	log.Println("Done ...")
}
