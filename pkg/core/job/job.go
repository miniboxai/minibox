package job

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"os"

	"github.com/fatih/color"
	"minibox.ai/pkg/core/volume"
)

type Job struct {
	ID       string
	Name     string
	RootPath string
	Status   int
	// Project *Project
	Tasks  []*TaskSpec
	events []Event
	chExit chan int
	out    io.ReadCloser
}

type TaskSpec struct {
	Name   string
	Image  string
	Labels []Label
	// Ports      []PortMap
	Volumes    []VolumeMap
	Links      []*TaskSpec
	WorkingDir string
	Cmd        []string
	Env        []string
}

type Label struct {
	Name  string
	Value string
}

type Event struct {
}

type PortMap struct {
	Src  int
	Dest int
}

type VolumeMap struct {
	Src  volume.Sourcer
	Dest string
}

type Status int

const (
	StateExit Status = iota
)

var colors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgWhite,
}

func (j *Job) SetOutput(rd io.ReadCloser) {
	r, w := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(rd)
		for scanner.Scan() {
			// ci := hash(j.Name) % uint32(len(colors))
			// clr := color.New(colors[ci]).SprintFunc()
			clr := color.New(color.FgCyan).SprintFunc()
			fmt.Fprintf(w, "%s %s\n", clr(j.Name+"\t|"), scanner.Text())
			// fmt.Println(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "log output:", err)
			r.Close()
		}
	}()
	j.out = r
}

func (j *Job) Output() io.ReadCloser {
	return j.out
}

func (j *Job) SetState(state Status, ch chan int) {
	switch state {
	case StateExit:
		j.chExit = ch
	}
}

func (j *Job) Wait() error {
	j.Status = <-j.chExit
	return nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
