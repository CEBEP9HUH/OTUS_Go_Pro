package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	const (
		count1 = 50
		count2 = 150
	)

	testSuites := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
		tasksCount     int
	}{
		{
			"tasks without errors: tasks count is greater than workers count",
			count1,
			1,
			count2,
		},
		{
			"tasks without errors: workers count is greater than tasks count",
			count2,
			1,
			count1,
		},
	}

	for _, ts := range testSuites {
		t.Run(ts.name, func(t *testing.T) {
			tasks := make([]Task, 0, ts.tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			for i := 0; i < ts.tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				sumTime += taskSleep

				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			requiredTime := sumTime / 2

			require.EventuallyWithT(t, func(c *assert.CollectT) {
				err := Run(tasks, ts.workersCount, ts.maxErrorsCount)
				require.NoError(c, err)
				require.Equal(c, runTasksCount, int32(ts.tasksCount), "not all tasks were completed")
			}, requiredTime, time.Millisecond, "tasks were run sequentially?")
		})
	}
}

func TestRun_EarlyReturn(t *testing.T) {
	defer goleak.VerifyNone(t)
	const (
		defaultWorkersCount   = 5
		defaultMaxErrorsCount = 1
		defaultTasksCount     = 50
	)

	testSuites := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
		tasksCount     int
	}{
		{
			"0 workers is allowed",
			0,
			defaultMaxErrorsCount,
			defaultTasksCount,
		},
		{
			"0 task has been passed",
			defaultWorkersCount,
			defaultMaxErrorsCount,
			0,
		},
	}

	for _, ts := range testSuites {
		t.Run(ts.name, func(t *testing.T) {
			tasks := make([]Task, 0, ts.tasksCount)

			var runTasksCount int32

			for i := 0; i < ts.tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, ts.workersCount, ts.maxErrorsCount)
			require.NoError(t, err)
			require.Equal(t, runTasksCount, int32(0), "extra tasks were started")
		})
	}
}

func TestRun_WithErrors(t *testing.T) {
	defer goleak.VerifyNone(t)
	const (
		defaultWorkersCount   = 10
		defaultMaxErrorsCount = 23
		defaultTasksCount     = 50
	)

	testSuites := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
		tasksCount     int
		maxTasksCount  int
	}{
		{
			"if were errors in first M tasks, than finished not more N+M tasks",
			defaultWorkersCount,
			defaultMaxErrorsCount,
			defaultTasksCount,
			defaultMaxErrorsCount + defaultWorkersCount,
		},
		{
			"0 errors is allowed",
			defaultWorkersCount,
			0,
			defaultTasksCount,
			0,
		},
	}

	for _, ts := range testSuites {
		t.Run(ts.name, func(t *testing.T) {
			tasks := make([]Task, 0, ts.tasksCount)

			var runTasksCount int32

			for i := 0; i < ts.tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			err := Run(tasks, ts.workersCount, ts.maxErrorsCount)

			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			require.LessOrEqual(t, runTasksCount, int32(ts.maxTasksCount), "extra tasks were started")
		})
	}
}
