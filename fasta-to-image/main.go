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
	rowCount    int
	columnCount int
	input       = flag.String("i", "", "Fasta file")
	output      = flag.String("o", "out.png", "Output image file")
	rectSize    = flag.Int("size", 5, "Block size")
	pprof       = flag.Bool("profile", false, "Output profile pprof")
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

	flag.Parse()

	if *pprof {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
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
