package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	inputFile := "input.txt"
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		inputFile = os.Args[1]
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	connections := make([]string, 0)
	for scanner.Scan() {
		connections = append(connections, scanner.Text())
	}

	startTime := time.Now()
	result := part1(connections)
	fmt.Printf("Part 1: %d\n", result)
	executionTime := float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 1 in %f seconds\n\n", executionTime)
}

func part1(connections []string) int {

	adjacency := make(map[string][]string)
	allEdges := make(map[string]bool)
	testHosts := make(map[string]bool)
	foundParties := make(map[string]bool)

	for _, conn := range connections {
		hosts := strings.Split(conn, "-")
		sort.Strings(hosts)
		edge := sortedHostsString(hosts)
		allEdges[edge] = true

		adj := adjacency[hosts[0]]
		adj = append(adj, hosts[1])
		adjacency[hosts[0]] = adj
		adj = adjacency[hosts[1]]
		adj = append(adj, hosts[0])
		adjacency[hosts[1]] = adj

		if hosts[0][0] == 't' {
			testHosts[hosts[0]] = true
		}
		if hosts[1][0] == 't' {
			testHosts[hosts[1]] = true
		}

		for host := range testHosts {
			for _, neighbor := range adjacency[host] {
				for _, candidate := range adjacency[neighbor] {
					if candidate != host {
						edge := []string{host, candidate}
						sort.Strings(edge)
						edgeString := sortedHostsString(edge)
						if allEdges[edgeString] {
							party := []string{host, neighbor, candidate}
							foundParties[sortedHostsString(party)] = true
						}
					}
				}
			}
		}

	}

	return len(foundParties)
}

func sortedHostsString(hosts []string) string {
	sorted := hosts
	sort.Strings(sorted)

	var sb strings.Builder
	for i, h := range hosts {
		if i > 0 {
			sb.WriteString("-")
		}
		sb.WriteString(h)
	}
	return sb.String()
}
