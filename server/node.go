package server

import "os"

type Node struct {}

func (n *Node) Version(request *string, version *string) error {

	*version = os.Getenv("VERSION")

	return nil
}