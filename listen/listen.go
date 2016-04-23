package listen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"

	"github.com/asteris-llc/mantl-bootstrap/cert"
	"github.com/asteris-llc/mantl-bootstrap/common"
	pb "github.com/asteris-llc/mantl-bootstrap/proto"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Init(root *cobra.Command) {
	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Listen for bootstrap information",
		Long:  "Listen for bootstrap information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Listen(args)
		},
	}

	root.AddCommand(listenCmd)
}

type server struct{
	Listener net.Listener
}

func (s *server) ConfigureConsul(ctx context.Context, in *pb.ConsulConfig) (*pb.Response, error) {
	var cconfig common.ConsulStatic

	if _, err := os.Stat(common.ConsulStaticPath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist. Creating\n", common.ConsulStaticPath)
	}

	// in.Data is a common.ConsulStatic structure. Unmarshal.
	if err := json.Unmarshal(in.Data, &cconfig); err != nil {
		return failure(err), err
	}

	// Read CA certificate and key
	cacert, err := ioutil.ReadFile(cconfig.CaFile)
	if err != nil {
		return failure(err), err
	}

	cakey, err := ioutil.ReadFile(common.CaKeyPath)
	if err != nil {
		return failure(err), err
	}

	// Generate host certificate
	c, err := cert.GenerateCert(cacert, cakey, cconfig.AdvertiseAddr, cconfig.Domain)
	if err != nil {
		return failure(err), err
	}

	// Write Host cert and set cconfig member
	if err := ioutil.WriteFile(common.HostCertPath, c.Cert, 0644); err != nil {
		return failure(err), err
	}
	cconfig.CertFile = common.HostCertPath
		
	// Write Host key and set cconfig member
	if err := ioutil.WriteFile(common.HostKeyPath, c.Key, 0644); err != nil {
		return failure(err), err
	}
	cconfig.KeyFile = common.HostKeyPath

	// Set BootstrapExpect member if a server
	if cconfig.Server {
		be := len(cconfig.RetryJoin)
		if be > 3 {
			be = 3
		}

		cconfig.BootstrapExpect = be
	}

	consulJson, err := json.MarshalIndent(cconfig, "", "  ")
	if err != nil {
		return failure(err), err
	}
	ioutil.WriteFile(common.ConsulStaticPath, consulJson, 0664)
	cmd := exec.Command("chown", "consul:consul", common.ConsulStaticPath)
	if err := cmd.Run(); err != nil {
		return failure(err), err
	}

	return success(), nil
}

func (s *server) WriteFile(ctx context.Context, in *pb.FileData) (*pb.Response, error) {
	fmt.Printf("Mode: %o\n", in.Mode)
	if err := ioutil.WriteFile(in.Path, in.Data, os.FileMode(in.Mode)); err != nil {
		return failure(err), err
	}

	return success(), nil
}

func (s *server) Shutdown(ctx context.Context, in *pb.ShutdownMsg) (*pb.Response, error) {
	s.Listener.Close()
	return success(), nil
}

func failure(err error) *pb.Response {
	return &pb.Response{
		Code: pb.Response_Failure,
		Mesg: err.Error(),
	}
}

func success() *pb.Response {
	return &pb.Response{
		Code: pb.Response_Success,
		Mesg: "success",
	}
}
	

func Listen(args []string) error {
	c := common.DefaultConfig()

	l, err := net.Listen(c.Type, fmt.Sprintf(":%d", c.Port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterBootstrapRPCServer(s, &server{Listener: l})
	s.Serve(l)

	return nil
}
