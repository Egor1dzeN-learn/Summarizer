package resources

import (
	"context"
	"summarizer/backend/http/config"
	pb "summarizer/backend/http/generated/protos"
	"summarizer/backend/http/internal/pkgs/roundrobin"
	"summarizer/backend/http/internal/pkgs/tasksq"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SummarizeResultConsumer func(result string, err error)

type WorkerNode interface {
	Summarize(text string, prompt string, accept SummarizeResultConsumer)

	CloseConnection()
}

type workerNodeDescriptor struct {
	addr string

	conn   *grpc.ClientConn
	client pb.NodeWorkerClient
}

type workerNodeOrchestrator struct {
	balancer roundrobin.RoundRobin[workerNodeDescriptor]

	tq tasksq.TaskQueue[pb.SummarizeRequest, pb.SummarizeReply]
	hn map[string]SummarizeResultConsumer
}

func NewWorkerNodeOrchestrator(cfg *config.WorkerNodesConfig) WorkerNode {
	log.Printf("wnc=%s", cfg)
	descriptors := make([]*workerNodeDescriptor, len(cfg.Addresses))
	for i, addr := range cfg.Addresses {
		descriptors[i] = &workerNodeDescriptor{
			addr: addr,
		}

		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}
		descriptors[i].conn = conn

		descriptors[i].client = pb.NewNodeWorkerClient(conn)

		log.Printf("Added worker node: %s", addr)
	}
	o := &workerNodeOrchestrator{
		balancer: roundrobin.New(descriptors...),
		hn:       map[string]SummarizeResultConsumer{},
	}
	o.tq = tasksq.NewTaskQueue(o, 1, 16)
	return o
}

func (o *workerNodeOrchestrator) Summarize(text string, prompt string, accept SummarizeResultConsumer) {
	o.tq.Enqueue(&pb.SummarizeRequest{
		Text:   text,
		Prompt: prompt,
	}, func(tid string) {
		o.hn[tid] = accept
	})
}

func (o *workerNodeOrchestrator) CloseConnection() {
	panic("TODO")
}

func (o *workerNodeOrchestrator) Process(task tasksq.Task[pb.SummarizeRequest]) (*pb.SummarizeReply, error) {
	node := o.balancer.Next()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := node.client.Summarize(ctx, task.Req)
	return r, err
}
func (o *workerNodeOrchestrator) Proceed(completedTask tasksq.CompletedTask[pb.SummarizeRequest, pb.SummarizeReply]) {
	f := o.hn[completedTask.ID]
	if completedTask.Err != nil {
		f("todo nil", completedTask.Err)
	} else {
		f(completedTask.Res.Text, nil)
	}
	delete(o.hn, completedTask.ID)
}
