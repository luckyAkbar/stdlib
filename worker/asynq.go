package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

var mux = asynq.NewServeMux()

type worker struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
}

// Client is the worker client
type Client interface {
	EnqueueTask(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error)
}

// Server is the worker server
type Server interface {
	Start(errch chan error) error
	Stop()
	RegisterTaskHandler([]TaskHandler)
	RegisterScheduler(task *asynq.Task, cronspec string) error
}

// Priority is worker priority
type Priority string

// list worker priority
var (
	PriorityHigh    Priority = "high"
	PriorityDefault Priority = "default"
	PriorityLow     Priority = "low"
)

// DefaultQueue is the default queue for worker. If you want to use this value
// you must use the defined priority above.
var DefaultQueue = map[string]int{
	string(PriorityHigh):    7,
	string(PriorityDefault): 2,
	string(PriorityLow):     1,
}

// DefaultHealtCheckFn is the default health check function for worker.
// This will only log the error
func DefaultHealtCheckFn(err error) {
	if err != nil {
		logrus.Errorf("unhealthy: %+v", err)
	}
}

// WorkerRateLimitError is used to indicate that the task is error because rate limited
type WorkerRateLimitError struct {
	RetryIn time.Duration
}

// Error return string representation of error
func (wrle *WorkerRateLimitError) Error() string {
	return fmt.Sprintf("rate limited (retry in  %v)", wrle.RetryIn)
}

// IsRateLimitError check if error is caused of rate limited
func IsRateLimitError(err error) bool {
	_, ok := err.(*WorkerRateLimitError)
	return ok
}

// NewWorkerRateLimitError create new WorkerRateLimitError based on supplied interval
func NewWorkerRateLimitError(interval time.Duration) *WorkerRateLimitError {
	return &WorkerRateLimitError{
		RetryIn: interval,
	}
}

// DefaultRetryDelayFn is the default retry delay function for worker. Will utilize rate limiter
var DefaultRetryDelayFn = func(n int, err error, task *asynq.Task) time.Duration {
	var rateLimiterErr *WorkerRateLimitError
	if errors.As(err, &rateLimiterErr) {
		return rateLimiterErr.RetryIn
	}

	return asynq.DefaultRetryDelayFunc(n, err, task)
}

// DefaultIsFailureCheckerFn check if the error is due to rate limitting. If not, don't mark it as failure
// so it will be retried again
var DefaultIsFailureCheckerFn = func(err error) bool {
	return !IsRateLimitError(err)
}

// DefaultEnqueueTaskFailureHandler is the default enqueue task failure handler. Will only log the error
var DefaultEnqueueTaskFailureHandler = func(task *asynq.Task, opts []asynq.Option, err error) {
	logrus.WithError(err).Errorf("failed to enqueue task %s", task.Type())
}

// NewClient create a new worker client
func NewClient(redisHost string) (Client, error) {
	redisOpts, err := asynq.ParseRedisURI(redisHost)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	client := asynq.NewClient(redisOpts)

	logrus.Info("worker client created")

	return &worker{
		client: client,
	}, nil
}

// NewServer creates a new worker server
func NewServer(redisHost string, serverCfg asynq.Config, schedulerCfg *asynq.SchedulerOpts) (Server, error) {
	redisOpts, err := asynq.ParseRedisURI(redisHost)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	client := asynq.NewClient(redisOpts)
	server := asynq.NewServer(
		redisOpts,
		serverCfg,
	)

	scheduler := asynq.NewScheduler(redisOpts, schedulerCfg)

	return &worker{
		client:    client,
		server:    server,
		scheduler: scheduler,
	}, nil
}

// Start start worker server
func (w *worker) Start(errch chan error) error {
	logrus.Info("starting worker...")

	go func() {
		logrus.Info("start to run the scheduler")
		if err := w.scheduler.Run(); err != nil {
			logrus.Error(err)

			errch <- err
		}
	}()

	go func() {
		logrus.Info("worker running...")
		if err := w.server.Run(mux); err != nil {
			logrus.Error(err)
			errch <- err
		}
	}()

	return nil
}

// Stop stop worker server
func (w *worker) Stop() {
	logrus.Info("stopping worker...")
	if w.client != nil {
		helper.WrapCloser(w.client.Close)
	}

	if w.server != nil {
		logrus.Info("stopping worker server...")
		w.server.Stop()
	}

	if w.scheduler != nil {
		logrus.Info("stopping worker scheduler...")
		w.scheduler.Shutdown()
	}

	logrus.Info("worker stopped.")
}

// TaskHandler is the task handler
type TaskHandler struct {
	Type    string
	handler asynq.HandlerFunc
}

// RegisterTaskHandler register task handler based on task type. This will be used by worker server and should be used before calling Start()
func (w *worker) RegisterTaskHandler([]TaskHandler) {
	for _, th := range []TaskHandler{} {
		mux.HandleFunc(th.Type, th.handler)
	}
}

func (w *worker) RegisterScheduler(task *asynq.Task, cronspec string) error {
	entryID, err := w.scheduler.Register(cronspec, task)
	if err != nil {
		logrus.Error("failed to register scheduler: ", err)
		return err
	}

	logrus.Info("success to register scheduler. entry id: ", entryID)

	return nil
}

func (w *worker) EnqueueTask(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	info, err := w.client.EnqueueContext(ctx, task)
	if err != nil {
		return nil, err
	}

	return info, nil
}
