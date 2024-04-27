package main

import (
	"bufio"
	"container/heap"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	inputFolder  = "./input"
	outputFolder = "./output"

	encodeMode = "encode"
	decodeMode = "decode"

	mapBinName  = "-map.bin"
	dataBinName = "-data.bin"
)

var wg sync.WaitGroup

func binaryStringToInt(binStr string) int {
	res := 0

	for i, ch := range binStr {
		if ch == '1' {
			res += int(math.Pow(float64(2), float64(7-i)))
		}
	}

	return res
}

func intToBinaryString(i int) string {
	str := strconv.FormatInt(int64(i), 2)

	if len(str) != 8 {
		for {
			str = "0" + str
			if len(str) == 8 {
				break
			}
		}
	}

	return str
}

func getCharCodes(node *Node, code string, m map[string]string) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		m[node.Character] = code
	}

	getCharCodes(node.Left, code+"0", m)
	getCharCodes(node.Right, code+"1", m)
}

func encodeFile(fileName string) {
	defer wg.Done()

	letters := make(map[string]int)
	stream, err := os.ReadFile(inputFolder + "/" + fileName)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, b := range stream {
		str := string(b)
		if letters[str] == 0 {
			letters[str] = 1
		} else {
			letters[str] += 1
		}
	}

	pq := mapToHeap(letters)

	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		combined := &Node{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}

		heap.Push(&pq, combined)
	}

	codes := make(map[string]string)
	getCharCodes(heap.Pop(&pq).(*Node), "", codes)

	fmap, err := os.Create(outputFolder + "/" + strings.Split(fileName, ".")[0] + mapBinName)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fmap.Close()

	reversedCodes := make(map[string]string)
	for key, value := range codes {
		reversedCodes[value] = key
	}

	enc := gob.NewEncoder(fmap)
	err = enc.Encode(reversedCodes)
	if err != nil {
		log.Println(err.Error())
		return
	}

	binaryString := ""

	for _, b := range stream {
		str := string(b)
		binaryString += codes[str]
	}

	fdata, err := os.Create(outputFolder + "/" + strings.Split(fileName, ".")[0] + dataBinName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fdata.Close()

	writer := bufio.NewWriter(fdata)
	byteString := ""
	test := ""

	for _, bit := range binaryString {
		byteString += string(bit)
		test += string(bit)
		if len(byteString) > 7 {
			b := []byte{byte(binaryStringToInt(byteString))}
			_, err = writer.Write(b)
			if err != nil {
				log.Println(err.Error())
				return
			}
			byteString = ""
		}
	}

	if len(byteString) != 0 {
		for {
			byteString += "0"
			if len(byteString) == 8 {
				b := []byte{byte(binaryStringToInt(byteString))}
				_, err = writer.Write(b)
				if err != nil {
					log.Println(err.Error())
					return
				}
				break
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func decodeFile(fileName string) {
	pathToData := fileName + dataBinName
	pathToMap := fileName + mapBinName

	fmap, err := os.OpenFile(pathToMap, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fmap.Close()

	dec := gob.NewDecoder(fmap)

	var codes map[string]string
	dec.Decode(&codes)

	dat, err := os.ReadFile(pathToData)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	binaryString := ""
	for _, b := range dat {
		str := intToBinaryString(int(b))
		binaryString += str
	}

	decodedString := ""
	temp := ""

	for _, ch := range binaryString {
		if codes[temp] == "" {
			temp += string(ch)
		} else {
			decodedString += codes[temp]
			temp = string(ch)
		}
	}

	name := strings.Split(fileName, "\\")
	foutput, err := os.Create(outputFolder + "\\" + name[len(name)-1] + ".txt")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer foutput.Close()

	writer := bufio.NewWriter(foutput)
	_, err = writer.Write([]byte(decodedString))
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = writer.Flush()
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func main() {
	defer func() {
		if recover() != nil {
			log.Fatal("No mode was provided. Please choose either encode or decode mode. If you are trying to decode, please provide path to file as second argument")
		}
	}()

	mode := os.Args[1]

	entries, err := os.ReadDir(inputFolder)

	if err != nil {
		log.Fatal(err.Error())
	}

	_ = os.Mkdir(outputFolder, 0755)

	if mode == encodeMode {
		for _, entry := range entries {
			wg.Add(1)
			go encodeFile(entry.Name())
		}
	} else if mode == decodeMode {
		decodeFile(os.Args[2])
	} else {
		log.Println("Selected incorrect mode. Should be either encode or decode")
	}

	wg.Wait()
}
