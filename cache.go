package unity_admin_observer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var Cache map[string]int

func ReadCache() map[string]int {
	bytes, err := ioutil.ReadFile("/home/siruba_2050/unity_admin_observer/cache/cache.txt")
	if err != nil {
		panic(err)
	}

	cache := make(map[string]int)

	var data []string
	line := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	for _, l := range line {
		data = strings.Split(l, ",")
	}
	for i := 0; i < len(data); i += 2 {
		n, err := strconv.Atoi(data[i+1])
		if err != nil {
			panic(err)
		}
		cache[data[i]] = n
	}

	return cache
}

func WriteCache(packages []string, nowSales []int) {
	if len(packages) != len(nowSales) {
		panic("len(packages) != len(nowSales)")
	}

	file, err := os.OpenFile("/home/siruba_2050/unity_admin_observer/cache/cache.txt", os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < len(packages); i++ {
		_, err := fmt.Fprintf(file, "%s,%d\n", packages[i], nowSales[i])
		if err != nil {
			panic(err)
		}
	}
}
