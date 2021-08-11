package scanner

import (
	"fmt"
	"net"
	"sort"
	
)

// Worker function runs concurrently to make a connection to a url with specific port
func worker(ports chan int, results chan int, url string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", url, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

// Scanner function scans how many ports are open for a given url
func Scanner() []int {
	
	url := "scanme.nmap.org"
	ports := make(chan int, 100)
	results := make(chan int)
	var openPorts []int
	
	for i:=0; i < cap(ports); i++ {
		go worker(ports, results, url)
	}
	
	go func() {
		for i:=0; i < 150; i++ {
			ports <- i
		}
	}()
	
	for i:=0; i < 150; i++ {
		port := <- results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}
	
	close(ports)
	close(results)
	sort.Ints(openPorts)
	return openPorts
}
