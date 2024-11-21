package memory

import (
	"fmt"
	"strings"
)

func (s *SharedMemory) SetClusterIPs() {
	// split the cluster endpoints by ':'
	parts := strings.Split(s.iptable[fmt.Sprintf("E%s", s.clusterName)], ":")

	// ip list
	list := make([]string, 0)
	for _, key := range parts {
		if key != s.nodeName {
			list = append(list, s.iptable[key])
		}
	}

	// set cluster IPs
	s.clusterIPs = list
}

// SetLeader updates the current leader id.
func (s *SharedMemory) SetLeader(leader string) {
	s.leader = leader
}
