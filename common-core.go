package main

import (
	"fmt"
)

// 2 checked main
func main() {
	n := 3
	majority := n/2 + 1
	//majority := (2*n + 1) / 3
	maxR := 3

	ack := 0
	nack := 0

	for q := 0; q < 100000; q++ {
		finalValues := make([][]int, n)

		for i := 0; i < n; i++ {
			finalValues[i] = make([]int, 0)
		}

		complete := make(chan bool, n)

		chanArr := make([][][]chan []int, maxR)

		for r := 0; r < maxR; r++ {
			chanArr[r] = make([][]chan []int, n)
			for i := 0; i < n; i++ {
				chanArr[r][i] = make([]chan []int, n)
				for j := 0; j < n; j++ {
					chanArr[r][i][j] = make(chan []int, n*n)
				}
			}
		}

		for i := 0; i < n; i++ {
			go func(id int) {
				initV := id
				values := []int{initV}
				for r := 0; r < maxR-1; r++ {

					//fmt.Printf("id: %v, r: %v, values: %v\n", id, r, values)
					for j := 0; j < n; j++ {
						if id == j {
							continue
						} else {
							chanArr[r][id][j] <- values
						}
					}
					count := 1
					for count < majority {
						for j := 0; j < n; j++ {
							if j == id {
								continue
							} else {
								select {
								case v := <-chanArr[r][j][id]:
									count++
									values = addMissingitems(values, v)
									if count == majority {
										j = n + 1
									}
									break
								default:
									// do nothing
									break
								}
							}
						}
					}
					// round i complete
					if len(values) < majority {
						panic("should not happen")
					}
					//rndTime = rand.Intn(2 + id)
					//time.Sleep(time.Duration(rndTime) * time.Microsecond)

				}
				//randSleep := rand.Intn(10)
				//time.Sleep(time.Duration(randSleep))
				//fmt.Printf("id: %v, r: %v, values: %v\n", id, maxR-1, values)
				finalValues[id] = values
				complete <- true
			}(i)
		}

		for i := 0; i < n; i++ {
			<-complete
		}

		if checkCommonCore(q, finalValues, majority) {
			ack++
		} else {
			nack++
		}
	}
	fmt.Printf("Ack: %v, NACK: %v\n", ack, nack)

}

// checked  2 common-core
func checkCommonCore(q int, values [][]int, majority int) bool {
	for i := 0; i < len(values); i++ {
		if len(values[i]) < majority {
			panic("should not happen")
		}
	}
	// check if there is common set of n/2+1 values in all the arrays
	// if yes, then print the common core
	frequency := make(map[int]int)
	for _, arr := range values {
		for _, val := range arr {
			_, ok := frequency[val]
			if ok {
				frequency[val] = frequency[val] + 1
			} else {
				frequency[val] = 1
			}
		}
	}
	core := make([]int, 0)
	count := 0
	for k, v := range frequency {
		if v == len(values) {
			core = append(core, k)
			count++
		}
	}
	if count >= majority {
		//fmt.Printf("\n Common Core: %v \n", core)
		return true
	} else {
		//fmt.Printf("\n No common core found in %v \n", q)
		//fmt.Printf("\n Values: %v \n", values)
		return false
	}
}

// checked 2 isAvailable
func isAvailable(arr []int, v int) bool {
	for _, val := range arr {
		if val == v {
			return true
		}
	}
	return false
}

// checked 2 addMissingitems
func addMissingitems(arr1 []int, arr2 []int) []int {
	for _, val := range arr2 {
		if !isAvailable(arr1, val) {
			arr1 = append(arr1, val)
		}
	}
	return arr1
}
