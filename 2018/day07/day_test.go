package day_test

import (
	"sort"
	"strings"
	"testing"
)

type ByByte []byte

func (a ByByte) Len() int           { return len(a) }
func (a ByByte) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByByte) Less(i, j int) bool { return a[i] < a[j] }

type worker struct {
	task   byte
	doneAt int
}

func NewWorker(task byte, now int) *worker {
	return &worker{
		task:   task,
		doneAt: now + int(task-'A'+1) + penalty,
	}
}

func (w *worker) work(task byte, now int) {
	w.task = task
	w.doneAt = now + int(task-'A'+1) + penalty
}

func (w *worker) String() string {
	if w.task == 0 {
		return "-"
	}
	return string(w.task)
}

func nextWorkerAvailableAt(workers []*worker) int {
	future := -1
	for _, w := range workers {
		if w != nil && (w.doneAt < future || future == -1) {
			future = w.doneAt
		}
	}
	if future == -1 {
		return 0
	}
	return future
}

/*
const penalty = 0
const numberOfWorkers = 2

var in = inputTest
/*/
const penalty = 60
const numberOfWorkers = 5

var in = input

//*/

func TestOne(t *testing.T) {
	workers := make([]*worker, numberOfWorkers)

	dependencies := make(map[byte][]byte)

	toComplete := make(map[byte]struct{})
	assigned := make(map[byte]struct{})

	for _, line := range strings.Split(in, "\n") {
		parent := line[5]
		child := line[55-19]
		dependencies[child] = append(dependencies[child], parent)
		toComplete[child] = struct{}{}
		toComplete[parent] = struct{}{}
	}

	hasUncompletedDependencies := func(child byte) bool {
		for _, parent := range dependencies[child] {
			if _, ok := toComplete[parent]; ok {
				return true
			}
		}
		dependencies[child] = nil
		return false
	}

	now := 0

	solution := ""
	for len(toComplete) > 0 {
		now = nextWorkerAvailableAt(workers)
		for i, w := range workers {
			if w == nil {
				continue
			}
			if w.doneAt <= now {
				delete(toComplete, w.task)
				workers[i] = nil
			}
		}

		var completable ByByte
		for task := range toComplete {
			if hasUncompletedDependencies(task) {
				continue
			}
			if _, ok := assigned[task]; ok {
				continue
			}
			completable = append(completable, task)
		}
		sort.Sort(completable)

		for i, w := range workers {
			if len(completable) == 0 {
				break
			}
			if w == nil {
				assigned[completable[0]] = struct{}{}
				workers[i] = NewWorker(completable[0], now)
				completable = completable[1:]
			}
		}
		t.Log(now, workers)
	}
	t.Fatal(solution, now)
}

var inputTest = `Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.`

var input = `Step Q must be finished before step O can begin.
Step Z must be finished before step G can begin.
Step W must be finished before step V can begin.
Step C must be finished before step X can begin.
Step O must be finished before step E can begin.
Step K must be finished before step N can begin.
Step P must be finished before step I can begin.
Step X must be finished before step D can begin.
Step N must be finished before step E can begin.
Step F must be finished before step A can begin.
Step U must be finished before step Y can begin.
Step M must be finished before step H can begin.
Step J must be finished before step B can begin.
Step B must be finished before step E can begin.
Step S must be finished before step L can begin.
Step A must be finished before step L can begin.
Step E must be finished before step L can begin.
Step L must be finished before step G can begin.
Step D must be finished before step I can begin.
Step Y must be finished before step I can begin.
Step I must be finished before step G can begin.
Step G must be finished before step R can begin.
Step V must be finished before step T can begin.
Step R must be finished before step H can begin.
Step H must be finished before step T can begin.
Step S must be finished before step E can begin.
Step C must be finished before step E can begin.
Step P must be finished before step T can begin.
Step I must be finished before step H can begin.
Step O must be finished before step P can begin.
Step M must be finished before step L can begin.
Step S must be finished before step D can begin.
Step P must be finished before step D can begin.
Step P must be finished before step R can begin.
Step I must be finished before step R can begin.
Step Y must be finished before step G can begin.
Step Q must be finished before step L can begin.
Step N must be finished before step R can begin.
Step J must be finished before step E can begin.
Step N must be finished before step T can begin.
Step B must be finished before step V can begin.
Step Q must be finished before step B can begin.
Step J must be finished before step H can begin.
Step F must be finished before step B can begin.
Step W must be finished before step X can begin.
Step S must be finished before step T can begin.
Step J must be finished before step G can begin.
Step O must be finished before step R can begin.
Step K must be finished before step B can begin.
Step Z must be finished before step O can begin.
Step Q must be finished before step S can begin.
Step K must be finished before step V can begin.
Step B must be finished before step R can begin.
Step J must be finished before step T can begin.
Step E must be finished before step T can begin.
Step G must be finished before step V can begin.
Step D must be finished before step Y can begin.
Step M must be finished before step Y can begin.
Step F must be finished before step G can begin.
Step C must be finished before step P can begin.
Step V must be finished before step R can begin.
Step R must be finished before step T can begin.
Step J must be finished before step Y can begin.
Step U must be finished before step R can begin.
Step Z must be finished before step F can begin.
Step Q must be finished before step V can begin.
Step U must be finished before step M can begin.
Step J must be finished before step R can begin.
Step L must be finished before step V can begin.
Step W must be finished before step K can begin.
Step B must be finished before step Y can begin.
Step O must be finished before step N can begin.
Step D must be finished before step V can begin.
Step P must be finished before step B can begin.
Step U must be finished before step I can begin.
Step O must be finished before step T can begin.
Step S must be finished before step G can begin.
Step X must be finished before step A can begin.
Step U must be finished before step T can begin.
Step A must be finished before step I can begin.
Step B must be finished before step G can begin.
Step N must be finished before step Y can begin.
Step Z must be finished before step J can begin.
Step M must be finished before step D can begin.
Step U must be finished before step A can begin.
Step S must be finished before step R can begin.
Step Z must be finished before step A can begin.
Step Y must be finished before step R can begin.
Step E must be finished before step Y can begin.
Step N must be finished before step G can begin.
Step Z must be finished before step X can begin.
Step P must be finished before step X can begin.
Step Z must be finished before step T can begin.
Step Z must be finished before step P can begin.
Step V must be finished before step H can begin.
Step P must be finished before step L can begin.
Step L must be finished before step H can begin.
Step X must be finished before step V can begin.
Step W must be finished before step G can begin.
Step N must be finished before step D can begin.
Step Z must be finished before step U can begin.`
