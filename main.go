package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Cost struct {
	Id   int
	Cost int
}

type CostPQ []Cost

func (pq *CostPQ) Len() int {
	return len((*pq))
}

func (pq *CostPQ) Less(i, j int) bool {
	return (*pq)[i].Cost < (*pq)[j].Cost
}

func (pq *CostPQ) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *CostPQ) Push(x any) {
	*pq = append(*pq, x.(Cost))
}

func (pq *CostPQ) Pop() any {
	r := (*pq)[(*pq).Len()-1]
	*pq = (*pq)[:len(*pq)-1]
	return r
}

type Node struct {
	Id        int    `json:"id"`
	LineType  int    `json:"lineType"`
	Name      string `json:"name"`
	Adjacency []Cost `json:"adjacency"`
	Route     map[int]*Route
}

func getData() []Node {
	f, err := os.Open("./data_v3.json")
	if err != nil {
		log.Fatalln(err)
	}
	b, err := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var data []Node
	json.Unmarshal(b, &data)
	return data
}

type Route struct {
	Route []string
	Step  int
}

func getFind(data []Node) func(begin, end string) *Route {
	idMap := make(map[int]*Node)
	nameMap := make(map[string]*Node)
	for i := range data {
		idMap[data[i].Id] = &data[i]
		nameMap[fmt.Sprintf("%d:%s", data[i].LineType, data[i].Name)] = &data[i]
	}
	find := func(beginId int) map[int]*Route {
		beginNode := idMap[beginId]
		costPQ := &CostPQ{}
		costPQ.Push(Cost{
			Id:   beginNode.Id,
			Cost: 0,
		})
		idSet := make(map[int]struct{})
		routeMap := make(map[int]Cost)
		for costPQ.Len() > 0 {
			currentCost := heap.Pop(costPQ).(Cost)
			currentNode := idMap[currentCost.Id]
			for _, v := range currentNode.Adjacency {
				if _, ok := idSet[v.Id]; !ok {
					cost := v.Cost + currentCost.Cost
					route, ok := routeMap[v.Id]
					if !ok || cost < route.Cost {
						routeMap[v.Id] = Cost{
							Id:   currentNode.Id,
							Cost: cost,
						}
					}
					heap.Push(costPQ, Cost{
						Id:   v.Id,
						Cost: cost,
					})
				}
			}
			idSet[currentCost.Id] = struct{}{}
		}
		r := make(map[int]*Route)
		for _, v := range data {
			route := &Route{}
			currentId := v.Id
			prevName := ""
			for {
				currentNode := idMap[currentId]
				if prevName != currentNode.Name {
					route.Route = append(route.Route, currentNode.Name)
					prevName = currentNode.Name
					route.Step++
				}
				if currentId == beginNode.Id {
					break
				}
				currentId = routeMap[currentId].Id
			}
			for i, j := 0, len(route.Route)-1; i < j; i, j = i+1, j-1 {
				route.Route[i], route.Route[j] = route.Route[j], route.Route[i]
			}
			r[v.Id] = route
		}
		return r
	}
	for i := range data {
		data[i].Route = find(data[i].Id)
	}
	return func(begin, end string) *Route {
		beginNode := nameMap[begin]
		endId := nameMap[end].Id
		return beginNode.Route[endId]
	}
}

func main() {
	data := getData()
	find := getFind(data)
	var route *Route
	route = find("2:宁波火车站", "5:曹隘")
	fmt.Printf("%+v\n", route)
	route = find("1:高桥西", "5:钱湖南路")
	fmt.Printf("%+v\n", route)
	route = find("1:高桥西", "1:霞浦")
	fmt.Printf("%+v\n", route)
	route = find("4:大卿桥", "5:钱湖南路")
	fmt.Printf("%+v\n", route)
}
