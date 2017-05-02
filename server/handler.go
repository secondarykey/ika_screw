package server

import (
	"fmt"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/pb"
	"github.com/secondarykey/ikascrew/video"
)

func init() {
}

func (i *IkascrewServer) startRPC() error {

	lis, err := net.Listen("tcp", ":55555")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterIkascrewServer(s, i)

	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			fmt.Println("failed to serve: %v", err)
			panic(err)
		}
	}()
	return nil
}

func (i *IkascrewServer) Sync(ctx context.Context, r *pb.SyncRequest) (*pb.SyncReply, error) {
	i.window.FullScreen()

	rep := &pb.SyncReply{
		Source:  "logo.png",
		Type:    "image",
		Project: ikascrew.ProjectName(),
	}

	return rep, nil
}

func (i *IkascrewServer) Effect(ctx context.Context, r *pb.EffectRequest) (*pb.EffectReply, error) {

	rep := &pb.EffectReply{
		Success: false,
	}

	v, err := video.Get(video.Type(r.Type), r.Name)
	if err != nil {
		return rep, err
	}

	e, err := effect.Get(effect.Type(r.Effect), v, i.window.GetEffect())
	if err != nil {
		return rep, err
	}

	i.window.SetEffect(e)

	rep.Success = true
	return rep, nil
}
