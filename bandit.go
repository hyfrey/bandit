package bandit

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Bandit can select arm or update information
type Bandit interface {
	SelectArm() int
	Update(arm int, reward float64)
	Reset()
	Version() string
}

// EpsilonGreedyNew constructs an epsilon greedy bandit.
func EpsilonGreedyNew(arms int, epsilon float64) (Bandit, error) {
	if !(epsilon >= 0 && epsilon <= 1) {
		return &epsilonGreedy{}, fmt.Errorf("epsilon not in [0, 1]")
	}

	return &epsilonGreedy{
		counts:  make([]int, arms),
		values:  make([]float64, arms),
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		arms:    arms,
		epsilon: epsilon,
	}, nil
}

// epsilonGreedy randomly selects arms with a probability of ε. The rest of
// the time, epsilonGreedy selects the currently best known arm.
type epsilonGreedy struct {
	counts  []int
	values  []float64
	epsilon float64
	arms    int
	rand    *rand.Rand
}

// SelectArm according to EpsilonGreedy strategy
func (e *epsilonGreedy) SelectArm() int {
	arm := 0
	if e.rand.Float64() > e.epsilon {
		// best arm
		for i := range e.values {
			if e.values[i] > e.values[arm] {
				arm = i
			}
		}
	} else {
		// random arm
		arm = e.rand.Intn(e.arms)
	}

	e.counts[arm] = e.counts[arm] + 1
	return arm + 1
}

// Update the running average
func (e *epsilonGreedy) Update(arm int, reward float64) {
	arm = arm - 1
	e.counts[arm] = e.counts[arm] + 1
	count := e.counts[arm]
	e.values[arm] = ((e.values[arm] * float64(count-1)) + reward) / float64(count)
}

// Version returns information on this bandit
func (e *epsilonGreedy) Version() string {
	return fmt.Sprintf("EpsilonGreedy(epsilon=%.2f)", e.epsilon)
}

// Reset returns the bandit to it's newly constructed state
func (e *epsilonGreedy) Reset() {
	e.counts = make([]int, e.arms)
	e.values = make([]float64, e.arms)
	e.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SoftmaxNew constructs a softmax bandit. Softmax explores non randomly
func SoftmaxNew(arms int, τ float64) (Bandit, error) {
	if !(τ >= 0.0) {
		return &softmax{}, fmt.Errorf("τ not in [0, ∞]")
	}

	return &softmax{
		counts: make([]int, arms),
		values: make([]float64, arms),
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		arms:   arms,
		tau:    τ,
	}, nil
}

// softmax holds counts values and temperature τ 
type softmax struct {
	counts []int
	values []float64
	tau    float64
	arms   int
	rand   *rand.Rand
}

// SelectArm 
func (s *softmax) SelectArm() int {
	z := 0.0
	for _, value := range s.values {
		z = z + math.Exp(value/s.tau)
	}

	var distribution []float64
	for _, value := range s.values {
		distribution = append(distribution, math.Exp(value/s.tau)/z)
	}

	accum := 0.0
	for i, p := range distribution {
		accum = accum + p
		if accum > z {
			return i
		}
	}

	return len(distribution) - 1
}

// Update the running average
func (s *softmax) Update(arm int, reward float64) {
	arm = arm - 1
	s.counts[arm] = s.counts[arm] + 1
	count := s.counts[arm]
	s.values[arm] = ((s.values[arm] * float64(count-1)) + reward) / float64(count)
}

// Version returns information on this bandit
func (s *softmax) Version() string {
	return fmt.Sprintf("Softmax(tau=%.2f)", s.tau)
}

// Reset returns the bandit to it's newly constructed state
func (s *softmax) Reset() {
	s.counts = make([]int, s.arms)
	s.values = make([]float64, s.arms)
	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}
