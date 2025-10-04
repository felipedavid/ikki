package scheduler

type Scheduler interface {
	SelectCandiateNodes()
	Score()
	Pick()
}