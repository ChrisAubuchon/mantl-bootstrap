package listen

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/asteris-llc/mantl-bootstrap/common"
	"github.com/asteris-llc/mantl-bootstrap/structs"

	"github.com/spf13/cobra"
)

func Init(root *cobra.Command) {
	listenCmd := &cobra.Command{
		Use: "listen",
		Short: "Listen for bootstrap information",
		Long: "Listen for bootstrap information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Listen(args)
		},
	}

	root.AddCommand(listenCmd)
}

func Listen(args []string) error {
	c := common.DefaultConfig()

	l, err := net.Listen(c.Type, fmt.Sprintf(":%d", c.Port))
	if err != nil {
		return err
	}

	var bs structs.Bootstrap

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		d := json.NewDecoder(conn)

		if err := d.Decode(&bs); err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}

		common.SendResponse(common.MSG_SUCCESS, "Success", conn)
		conn.Close()

		break
	}


	fmt.Printf("%+v\n", bs)

	return nil
}
