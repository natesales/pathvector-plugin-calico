package main

import (
	"context"
	"fmt"
	"os"

	"github.com/natesales/pathvector/pkg/config"
	"github.com/natesales/pathvector/pkg/plugin"
	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var version = "0.0.1"

type Plugin struct{}

var _ plugin.Plugin = (*Plugin)(nil)

func (Plugin) Description() string {
	return "Calico BGP neighbor integration"
}

func (*Plugin) Version() string {
	return version
}

func (*Plugin) Command() *cobra.Command {
	return nil
}

func (Plugin) Modify(c *config.Config) error {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}
	cclientset, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}
	kclientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	// List BGP peers
	peers, err := cclientset.ProjectcalicoV3().BGPPeers().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}
	if len(peers.Items) < 1 {
		return fmt.Errorf("calico has no configured BGP peers")
	}
	remoteAS := int(peers.Items[0].Spec.ASNumber)

	// List nodes
	nodes, err := kclientset.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	var nodeIPs []string
	for _, node := range nodes.Items {
		nodeIPs = append(nodeIPs, node.Status.Addresses[0].Address)
	}

	// Add the peer
	c.Peers["Calico"] = &config.Peer{
		ASN:         &remoteAS,
		NeighborIPs: &nodeIPs,
	}

	return nil
}

func init() {
	plugin.Register("calico", &Plugin{})
}
