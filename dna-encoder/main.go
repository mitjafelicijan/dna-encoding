package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/profile"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	input   = flag.String("i", "", "Input file")
	output  = flag.String("o", "out.fa", "Output file")
	seqName = flag.String("seq", "SEQ1", "Sequence name")
	columns = flag.Int("col", 60, "Characters by column")
	pprof   = flag.Bool("profile", false, "Output profile pprof")
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

	flag.Parse()

	if *pprof {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	start := time.Now()

	log.Println("Gathering input file's stats")

	file, err := os.Open(*input)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

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
			} else if binary[element] == 48 && binary[element+1] == 49 {
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
