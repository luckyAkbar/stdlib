package worker

import (
	"context"
	"os"

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

type Client interface {
	RegisterTask(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error)
}

type Server interface {
	Start() error
	Stop()
	RegisterTaskHandler([]TaskHandler)
	RegisterScheduler(task *asynq.Task, cronspec string) error
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
func (w *worker) Start() error {
	logrus.Info("starting worker...")

	go func() {
		logrus.Info("start to run the scheduler")
		if err := w.scheduler.Run(); err != nil {
			logrus.Error(err)

			os.Exit(1)
		}
	}()

	if err := w.server.Run(mux); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	logrus.Info("worker running...")

	return nil
}

// Stop stop worker server
func (w *worker) Stop() {
	logrus.Info("stopping worker...")
	if w.client != nil {
		helper.WrapCloser(w.client.Close)
	}

	if w.server != nil {
		w.server.Stop()
	}

	logrus.Info("worker stopped.")
}

type TaskHandler struct {
	Type    string
	handler asynq.HandlerFunc
}

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

func (w *worker) RegisterTask(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	info, err := w.client.EnqueueContext(ctx, task)
	if err != nil {
		return nil, err
	}

	return info, nil
}
