package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"golang.org/x/sync/errgroup"
)

// Constants for ACL permissions
const (
	PermRead = 1 << iota
	PermWrite
	PermCreate
	PermDelete
	PermAdmin
	PermAll = 0x1f
)

type ServiceNode struct {
	Name string `json:"name"` // 服务名称，这里是user
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SdClient struct {
	zkServers []string // 多个节点地址
	zkRoot    string   // 服务根节点，这里是/api
	conn      *zk.Conn // zk的客户端连接
}

func NewClient(zkServers []string, zkRoot string, timeout int) (*SdClient, error) {
	client := new(SdClient)
	client.zkServers = zkServers
	client.zkRoot = zkRoot
	conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	client.conn = conn
	if err := client.ensureRoot(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

func (s *SdClient) Close() {
	s.conn.Close()
}

func (s *SdClient) ensureRoot() error {
	exists, _, err := s.conn.Exists(s.zkRoot)
	if err != nil {
		return err
	}
	if !exists {
		// Access Control List
		_, err := s.conn.Create(s.zkRoot, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *SdClient) Register(node *ServiceNode) error {
	if err := s.ensureName(node.Name); err != nil {
		return err
	}
	fmt.Println("11")

	path := s.zkRoot + "/" + node.Name + "/" + node.Host
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	fmt.Println("22:", path)
	_, err = s.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	fmt.Println("33")
	return err
}

func (s *SdClient) ensureName(name string) error {
	path := s.zkRoot + "/" + name
	exists, _, err := s.conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := s.conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *SdClient) GetW(name, host string) ([]byte, <-chan zk.Event, error) {
	path := s.zkRoot + "/" + name + "/" + host
	data, stat, ch, err := s.conn.GetW(path)
	if err != nil {
		return []byte(""), ch, err
	}
	fmt.Println("11===== ", data, stat)
	return data, ch, nil
}

func main() {
	servers := []string{"localhost:2181"}
	client, err := NewClient(servers, "/api", 10)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	node1 := &ServiceNode{"user", "127.0.0.1", 4000}
	//node2 := &ServiceNode{"user", "127.0.0.1", 4001}
	//node3 := &ServiceNode{"user", "127.0.0.1", 4002}
	if err := client.Register(node1); err != nil {
		panic(err)
	}
	/*
		if err := client.Register(node2); err != nil {
			panic(err)
		}
		if err := client.Register(node3); err != nil {
			panic(err)
		}
	*/

	ctx, _ := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		// GetW仅生效一次
		bb, ch, err := client.GetW("user", "127.0.0.1")
		if err != nil {
		} else {
			for {
				select {
				case msg, ok := <-ch:
					if ok {
						fmt.Println(bb, "||||||||| ", msg)
					} else {
						fmt.Println("----------------------- ok=false")
						time.Sleep(time.Second * 3)
					}
				}
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("wait --:", err)
	}
}
