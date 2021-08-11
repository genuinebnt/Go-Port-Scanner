package scanner

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
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
func Scanner(url string, portStrings string) []int {
	ports, _ := stringSplit(portStrings)
	portsChan := make(chan int, 100)
	results := make(chan int)
	var openPorts []int
	
	for i:=0; i < cap(portsChan); i++ {
		go worker(portsChan, results, url)
	}
	
	go func() {
		for _, port := range ports {
			portsChan <- port
		}
	}()
	
	for _ = range ports {
		portScanned := <- results
		if portScanned != 0 {
			openPorts = append(openPorts, portScanned)
		}
	}
	
	close(portsChan)
	close(results)
	sort.Ints(openPorts)
	return openPorts
}

func stringSplit(s string) ([]int, error) {
	var ports []int
	if strings.Contains(s, ",") && strings.Contains(s, "-") {
		sp := strings.Split(s, ",")
		for _, p := range sp {
			if strings.Contains(p, "-") {
				err := dashSplit(p, &ports)
				if err != nil {
					return ports, err
				}
			} else {
				if err := checkValidPort(p, &ports); err != nil {
					return ports, err
				}
			}
		}
	} else if strings.Contains(s, ",") {
		sp := strings.Split(s, ",")
		for _, p := range sp {
			if err := checkValidPort(p, &ports); err != nil {
				return ports, err
			}
		}
	} else if strings.Contains(s, "-") {
		if err := dashSplit(s, &ports); err != nil {
			return ports, err
		}
	}
	return ports, nil
}

func checkValidPort(port string, ports *[]int) error {
	port = strings.TrimSpace(port)
	p, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("invalid port")
	}
	if p < 1 || p > 65535 {
		return errors.New("invalid port")
	}
	*ports = append(*ports, p)
	return nil
}

func dashSplit(sp string, ports *[]int) error {
	ds := strings.Split(sp, "-")
	for i, port := range ds {
		ds[i] = strings.TrimSpace(port)
	}
	if len(ds) != 2 {
		return errors.New("invalid port range")
	}
	
	start, err := strconv.Atoi(ds[0])
	if err != nil {
		return errors.New("invalid port range")
	}
	end, err := strconv.Atoi(ds[1])
	if err != nil {
		return errors.New("invalid port range")
	}
	
	if (start < 1 || start > 65535) || (end < 1 || end > 65535) || start > end {
		return errors.New("invalid port range")
	}
	for i:=start; i <= end; i++ {
		*ports = append(*ports, i)
	}
	
	return nil
}