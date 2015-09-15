package main

import (
	"fmt"
	"github.com/heketi/heketi/utils"
)

type Node []*Device
type Zone []Node

type Topology struct {
	tmap  map[int]map[string][]*Device
	zones []Zone
	len   int
}

func NewTopology() *Topology {
	t := &Topology{}
	t.tmap = make(map[int]map[string][]*Device)

	return t
}

func (t *Topology) crateSlices() {
	t.zones = make([]Zone, 0)

	for _, n := range t.tmap {
		zone := make([]Node, 0)
		for _, d := range n {
			zone = append(zone, d)
		}
		t.zones = append(t.zones, zone)
	}

}

func (t *Topology) Add(d *Device) {
	if nodes, ok := t.tmap[d.zone]; ok {
		if _, ok := nodes[d.nodeid]; ok {
			nodes[d.nodeid] = append(nodes[d.nodeid], d)
		} else {
			nodes[d.nodeid] = []*Device{d}
		}
	} else {
		t.tmap[d.zone] = make(map[string][]*Device)
		t.tmap[d.zone][d.nodeid] = []*Device{d}
	}
}

func (t *Topology) Rebalance() []*Device {

	t.crateSlices()

	list := make([]*Device, 0)

	var device *Device
	for i := 0; len(t.zones) != 0; i++ {
		zone := i % len(t.zones)
		node := i % len(t.zones[zone])

		// pop device
		device, t.zones[zone][node] = t.zones[zone][node][len(t.zones[zone][node])-1], t.zones[zone][node][:len(t.zones[zone][node])-1]
		list = append(list, device)

		// delete node
		if len(t.zones[zone][node]) == 0 {
			t.zones[zone] = append(t.zones[zone][:node], t.zones[zone][node+1:]...)

			// delete zone
			if len(t.zones[zone]) == 0 {
				t.zones = append(t.zones[:zone], t.zones[zone+1:]...)
			}
		}
	}

	return list
}

type Device struct {
	zone             int
	nodeid, deviceid string
}

func (d *Device) String() string {
	return fmt.Sprintf("{Z:%v N:%v D:%v}",
		d.zone,
		d.nodeid,
		d.deviceid)
}

func (d *Device) Val() string {
	return fmt.Sprintf("%v%v%v",
		d.deviceid,
		d.nodeid,
		d.zone)
}

func main() {
	zones, nodes, drives := 10, 100, 48

	list := []*Device{}

	index := 0
	for z := 0; z < zones; z++ {
		for n := 0; n < nodes; n++ {
			nid := utils.GenUUID()[:4]
			for d := 0; d < drives; d++ {
				did := utils.GenUUID()[:4]
				dev := &Device{
					deviceid: did,
					nodeid:   nid,
					zone:     z,
				}
				list = append(list, dev)
				index++
				//s = append(s, fmt.Sprintf("d%v:n%v:z%v", utils.GenUUID()[:4], nid, z))
			}
		}
	}
	fmt.Println(list)
	fmt.Println("-------")

	t := NewTopology()
	for _, d := range list {
		t.Add(d)
	}
	l := t.Rebalance()
	fmt.Println(l)

	fmt.Println(len(l))
}
