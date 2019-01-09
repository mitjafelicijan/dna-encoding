package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	appName = "dnae-encode"
	version = "0.0.1"

	app     = kingpin.New(appName, "A command-line application that encodes file into DNA sequence.")
	input   = app.Flag("input", "Input file (ASCII or binary) which will be encoded into DNA sequence.").PlaceHolder("INPUT").Required().Short('i').String()
	output  = app.Flag("output", "Output file which stores DNA sequence in FASTA format.").Short('o').Default("out.fa").String()
	seqName = app.Flag("sequence", "The description line (defline) or header/identifier line, gives a name and/or a unique identifier for the sequence.").PlaceHolder("SEQ1").Short('s').Default("SEQ1").String()
	columns = app.Flag("columns", "Row characters length (no more than 120 characters). Devices preallocate fixed line sizes in software.").PlaceHolder("60").Short('c').Default("60").Int()
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

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

	if *input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()

	file, err := os.Open(*input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	log.Println("Gathering input file's stats")

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting encoding ...")

	bar := pb.StartNew(int(fi.Size()))
	bar.SetUnits(pb.U_BYTES)

	sequence := make([]byte, 0)
	buffer := make([]byte, 1)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		binary := fmt.Sprintf("%08b", buffer[0])
		for _, element := range []int{0, 2, 4, 6} {
			if binary[element] == 48 && binary[element+1] == 48 {
				sequence = append(sequence, 'A')
			} else if binary[element] == 48 && binary[element+1] == 49 {
				sequence = append(sequence, 'G')
			} else if binary[element] == 49 && binary[element+1] == 48 {
				sequence = append(sequence, 'C')
			} else if binary[element] == 49 && binary[element+1] == 49 {
				sequence = append(sequence, 'T')
			}
		}

		bar.Increment()
	}

	bar.Finish()

	log.Println("Saving to FASTA file ...")

	file, err = os.OpenFile(
		*output,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf(">%s\n", *seqName))
	for i := 0; i < len(sequence); i += *columns {
		byteSlice := sequence[i:min(i+*columns, len(sequence))]
		file.Write(byteSlice)
		file.WriteString("\n")
	}

	fi, err = file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Output FASTA file length is", byteCountDecimal(fi.Size()))

	elapsed := time.Since(start)

	log.Printf("Process took %s\n", elapsed)
	log.Println("Done ...")

}
