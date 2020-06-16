package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mailru_task/analysis"
	"log"
	"os"
)

//Читает из Stdin строки и записывает в слайс
func readInput() ([]string, error) {
	var paths []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

func main() {
	goMax := flag.Int("k", 5, "Максимальное кол-во горутин-обработчиков")
	word := flag.String("word", "Go", "Искомое слово")
	flag.Parse()

	paths, err := readInput()
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	result, dataErrors := analysis.RunWordAnalysis(paths, *goMax, *word)

	for _, err := range dataErrors {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	var sum int64
	for _, elem := range result {
		if elem.Err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v, source: %v \n", elem.Err, elem.Path)
			continue
		}
		fmt.Printf("Count for %v: %v\n", elem.Path, elem.Value)
		sum += elem.Value
	}
	fmt.Printf("Total: %v\n", sum)
}
