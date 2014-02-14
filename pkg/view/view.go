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
	Members map[Process]bool
}

func newView() *View {
	v := View{}
	v.Members = make(map[Process]bool)
	return &v
}

func NewWithProcesses(processes ...Process) *View {
	v := newView()
	for _, process := range processes {
		v.Members[process] = true
	}
	return v
}

func (v *View) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{")

	first := true
	for process, _ := range v.Members {
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
	return v.Members[p]
}

func (v *View) GetMembers() []Process {
	var members []Process
	for process, _ := range v.Members {
		members = append(members, process)
	}
	return members
}

func (v *View) QuorumSize() int {
	return v.quorumSize()
}

func (v *View) quorumSize() int {
	membersTotal := len(v.Members)
	return (membersTotal+1)/2 + (membersTotal+1)%2
}

func (v *View) N() int {
	return len(v.Members)
}

func (v *View) F() int {
	membersTotal := len(v.Members)
	return membersTotal - v.quorumSize()
}
