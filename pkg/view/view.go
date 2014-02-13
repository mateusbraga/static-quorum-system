package view

import (
	"bytes"
	"fmt"
)

type Process struct {
	Addr string
}

func (thisProcess Process) Less(otherProcess Process) bool {
	return thisProcess.Addr < otherProcess.Addr
}

// -------- View type -----------

type View struct {
	members map[Process]bool
}

func newView() *View {
	v := View{}
	v.members = make(map[Process]bool)
	return &v
}

func NewWithProcesses(processes ...Process) *View {
	v := newView()
	for _, process := range processes {
		v.members[process] = true
	}
	return v
}

func (v *View) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{")

	first := true
	for process, _ := range v.members {
		if !first {
			fmt.Fprintf(&b, ", ")
		}
		fmt.Fprintf(&b, "%v", process.Addr)
		first = false
	}

	fmt.Fprintf(&b, "}")
	return b.String()
}

func (v *View) HasMember(p Process) bool {
	return v.members[p]
}

func (v *View) GetMembers() []Process {
	var members []Process
	for process, _ := range v.members {
		members = append(members, process)
	}
	return members
}

func (v *View) QuorumSize() int {
	return v.quorumSize()
}

func (v *View) quorumSize() int {
	membersTotal := len(v.members)
	return (membersTotal+1)/2 + (membersTotal+1)%2
}

func (v *View) N() int {
	return len(v.members)
}

func (v *View) F() int {
	membersTotal := len(v.members)
	return membersTotal - v.quorumSize()
}
