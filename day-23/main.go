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

type Set[T comparable] struct {
	data []T
	hash map[T]int
}

func (s *Set[T]) isNotInitialized() bool {
	return (s.hash == nil)
}

func (s *Set[T]) IsEmpty() bool {
	if len(s.hash) == 0 {
		return true
	}
	return false
}

func (s *Set[T]) Add(item T) bool {
	if s.isNotInitialized() {
		s.hash = make(map[T]int)
	}
	if _, found := s.hash[item]; !found {
		s.data = append(s.data, item)
		s.hash[item] = len(s.data) - 1
		return true
	}
	return false
}

func (s *Set[T]) AddAll(items []T) int {
	added := 0
	for _, item := range items {
		if s.Add(item) {
			added++
		}
	}
	return added
}

func (s *Set[T]) Remove(item T) bool {
	if s.isNotInitialized() {
		s.hash = make(map[T]int)
	}
	if i, found := s.hash[item]; found {
		for k := i + 1; k < len(s.data); k++ {
			v := s.data[k]
			s.hash[v] = k - 1
		}
		s.data = append(s.data[:i], s.data[i+1:]...)
		delete(s.hash, item)
		return true
	}
	return false
}

func (s *Set[T]) RetainAll(items []T) {
	if s.isNotInitialized() {
		s.hash = make(map[T]int)
	}
	keys := make(map[T]bool)
	for _, itm := range items {
		keys[itm] = true
	}
	for _, e := range s.Elements() {
		if !keys[e] {
			s.Remove(e)
		}
	}
}

func (s *Set[T]) Contains(item T) bool {
	_, found := s.hash[item]
	return found
}

func (s *Set[T]) Size() int {
	return len(s.data)
}

func (s *Set[T]) Elements() []T {
	copy := make([]T, 0)
	copy = append(copy, s.data...)
	return copy
}

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

	startTime = time.Now()
	result2 := part2(connections)
	fmt.Printf("Part 2: %s\n", result2)
	executionTime = float32(time.Since(startTime).Milliseconds()) / float32(1000)
	fmt.Printf("Completed Part 2 in %f seconds\n\n", executionTime)
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

func part2(connections []string) string {

	adjacency := make(map[string][]string)
	for _, conn := range connections {
		hosts := strings.Split(conn, "-")
		adj := adjacency[hosts[0]]
		adj = append(adj, hosts[1])
		adjacency[hosts[0]] = adj
		adj = adjacency[hosts[1]]
		adj = append(adj, hosts[0])
		adjacency[hosts[1]] = adj
	}

	vertices := Set[string]{}
	for v := range adjacency {
		vertices.Add(v)
	}

	currentClique := Set[string]{}
	alreadyProcessed := Set[string]{}

	allCliques := bronKerbosch(currentClique, vertices, alreadyProcessed, adjacency)

	maxSize := 0
	maxSizeIndex := -1
	for n, clique := range allCliques {
		if clique.Size() > maxSize {
			maxSize = clique.Size()
			maxSizeIndex = n
		}
	}
	return sortedHostsString(allCliques[maxSizeIndex].Elements())
}

func bronKerbosch(currentClique, potentialNodes, alreadyProcessed Set[string], adjacency map[string][]string) []Set[string] {

	cliques := make([]Set[string], 0)
	if potentialNodes.IsEmpty() && alreadyProcessed.IsEmpty() {
		clique := Set[string]{}
		clique.AddAll(currentClique.Elements())
		cliques = append(cliques, clique)
	}

	for !potentialNodes.IsEmpty() {
		n := potentialNodes.Elements()[0]
		newCurrentClique := Set[string]{}
		newCurrentClique.AddAll(currentClique.Elements())
		newCurrentClique.Add(n)

		newPotentialNodes := Set[string]{}
		newPotentialNodes.AddAll(potentialNodes.Elements())
		newPotentialNodes.RetainAll(adjacency[n])

		newAlreadyProcessed := Set[string]{}
		newAlreadyProcessed.AddAll(alreadyProcessed.Elements())
		newAlreadyProcessed.RetainAll(adjacency[n])

		cliques = append(cliques, bronKerbosch(newCurrentClique, newPotentialNodes, newAlreadyProcessed, adjacency)...)
		potentialNodes.Remove(n)
		alreadyProcessed.Add(n)
	}
	return cliques
}

func sortedHostsString(hosts []string) string {
	sorted := hosts
	sort.Strings(sorted)

	var sb strings.Builder
	for i, h := range hosts {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(h)
	}
	return sb.String()
}
