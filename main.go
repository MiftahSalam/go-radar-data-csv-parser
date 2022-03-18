package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func processOneBin(w io.Writer, colI, colQ string) error {
	intDataI, err := strconv.ParseInt(colI, 10, 16)
	if err != nil {
		return err
	} else {
		rawDataI := int16(intDataI)

		err = binary.Write(w, binary.BigEndian, rawDataI)
		if err != nil {
			return err
		}

		intDataQ, err := strconv.ParseInt(colQ, 10, 16)
		if err != nil {
			return err
		} else {
			rawDataQ := int16(intDataQ)

			err = binary.Write(w, binary.BigEndian, rawDataQ)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func processOneChirp(out *os.File, I, Q *csv.Reader) {
	rowSize := 0
	for {
		rowI, err := I.Read()
		if err == io.EOF {
			log.Println("EOF file I")
			break
		}
		rowQ, err := Q.Read()
		if err == io.EOF {
			log.Println("EOF file Q")
			break
		}

		fmt.Printf("Row %v -> file I len: %v, file Q len: %v\n", rowSize, len(rowI), len(rowQ))

		err = processOneBin(out, rowI[0], rowQ[0])
		if err != nil {
			log.Fatalf("Row %v -> cannot process data bin with error: %v", rowSize, err)
		}
		rowSize++
	}

}

func processAllChirp(out *os.File, I, Q *csv.Reader) {
	const OFFSET_INC = 10400 * 4
	rowSize := 0
	offset := 0

	for {
		rowI, err := I.Read()
		if err == io.EOF {
			log.Println("EOF file I")
			break
		}
		rowQ, err := Q.Read()
		if err == io.EOF {
			log.Println("EOF file Q")
			break
		}

		fmt.Printf("Row %v, offset %v -> file I len: %v, file Q len: %v\n", rowSize, offset, len(rowI), len(rowQ))

		max_iterate := int(math.Max(float64(len(rowQ)), float64(len(rowQ))))
		offset = rowSize * 4
		for i := 0; i < max_iterate; i++ {
			// fmt.Printf("Row %v, offset %v\n", rowSize, offset)
			out.Seek(int64(offset), 0)
			err = processOneBin(out, rowI[i], rowQ[i])
			if err != nil {
				log.Fatalf("Row %v -> cannot process data bin %v with error:%v", rowSize, i, err)
			}
			offset += OFFSET_INC
		}
		rowSize++
	}
}
func main() {
	fileInputI, err := os.Open("../raw-data/data stadiun diam 6 I.csv")
	if err != nil {
		log.Fatal("Cannot open file I with error:", err)
	}
	fileInputQ, err := os.Open("../raw-data/data stadiun diam 6 Q.csv")
	if err != nil {
		log.Fatal("Cannot open file Q with error:", err)
	}

	fileOutput, err := os.OpenFile("../raw-data/data stadiun diam 6.dat", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal("Cannot create file raw with error:", err)
	}

	defer func() {
		fileInputI.Close()
		fileInputQ.Close()
		fileOutput.Close()
	}()

	csvReaderI := csv.NewReader(fileInputI)
	csvReaderQ := csv.NewReader(fileInputQ)

	// processOneChirp(fileOutput, csvReaderI, csvReaderQ)
	processAllChirp(fileOutput, csvReaderI, csvReaderQ)
}
