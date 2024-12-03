package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type measurement struct {
	max   int
	min   int
	sum   int
	count int
}

func renderResults(m map[string]measurement) string {
	var ret = ""
	var min float64
	var max float64
	var mean float64
	for city, measure := range m {
		min = float64(measure.min) / 10
		max = float64(measure.max) / 10
		mean = float64(measure.sum) / (float64(measure.count) * 10)
		ret += fmt.Sprintf("%s -> min:%.1f max:%.1f mean:%.1f\n", city, min, max, mean)
	}
	return ret
}

func main() {

	m := make(map[string]measurement)

	if len(os.Args) != 2 {
		log.Fatalf("Missing measurements filename")
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	ch := make(chan string)

	go func() {
		fmt.Println("spawned")
		for {
			tmp, ok := <-ch
			if !ok {
				return
			}
			data := strings.Split(tmp, ";")
			data[1] = strings.Replace(data[1], ".", "", -1)
			temp, err := strconv.Atoi(data[1])
			if err != nil {
				panic(err)
			}
			value, present := m[data[0]]
			if !present {
				m[data[0]] = measurement{max: temp, min: temp, sum: temp, count: 1}
			} else {
				if value.max < temp {
					value.max = temp
				}
				if value.min > temp {
					value.min = temp
				}
				value.sum += temp
				value.count += 1
				m[data[0]] = value
			}
		}
	}()
	for fileScanner.Scan() {
		ch <- fileScanner.Text()
	}
	fmt.Println("done")
	close(ch)
	fmt.Print(renderResults(m))

}
