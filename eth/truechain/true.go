/*
Copyright (c) 2018 TrueChain Foundation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package truechain

import (
	"net"
    "math/big"
    "errors"
    
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/common"
)

type HybridConsensusHelp struct {
}
func (s *HybridConsensusHelp) PutBlock(ctx context.Context, block *TruePbftBlock) (*CommonReply, error) {
    // do something
    return &CommonReply{Message: "success "}, nil
}
func (s *HybridConsensusHelp) ViewChange(ctx context.Context, in *EmptyParam) (*CommonReply, error) {
    // do something
    return &CommonReply{Message: "success "}, nil
}

type PyHybConsensus struct {
}

type TrueHybrid struct {
    quit        bool
    address     string
    curCmm      []*CommitteeMember
    oldCmm      []*CommitteeMember

    sdmsize     int
    sdm         []*StandbyInfo
    crpmsg      []*TrueCryptoMsg        // authenticated msg by block comfirm
    crptmp      []*TrueCryptoMsg        // unauthenticated msg by block comfirm
}

func (t *TrueHybrid) HybridConsensusHelpInit() {
    port = 17546
    lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterHybridConsensusHelpServer(s, &HybridConsensusHelp{})
	// Register reflection service on gRPC server.
	// reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (t *TrueHybrid) MembersNodes(nodes []*TruePbftNode) error{
    // Set up a connection to the server.
    conn, err := grpc.Dial(t.address, grpc.WithInsecure())
    if err != nil {
        return err
    }
    defer conn.Close()   
    c := NewPyHybConsensusClient(conn)
 
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    // pbNodes := make([]*TruePbftNode,0,0)
    // for _,v := range nodes {
    //     pbNodes = append(pbNodes,&TruePbftNode{
    //         Addr:       v.Addr,
    //         Pubkey:     v.pubkey,
    //         Privkey:    v.Privkey,
    //     })
    // }
    _, err1 := c.MembersNodes(ctx, &Nodes{Nodes:nodes})
    if err1 != nil {
        return err1
    }
    return nil
}
func (t *TrueHybrid) SetTransactions(txs []*types.Transaction){
    // Set up a connection to the server.
    conn, err := grpc.Dial(t.address, grpc.WithInsecure())
    if err != nil {
        return err
    }
    defer conn.Close()   

    c := NewPyHybConsensusClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    pbTxs := make([]*Transaction,0,0)
    for _,v := range txs {
        to := make([]byte,0,0)
        if t := v.To(); t != nil {
            to = t.Bytes()
        }
        v,r,s := v.RawSignatureValues()
        pbTxs = append(pbTxs,&Transaction{
            Data:       &TxData{
                AccountNonce:       v.Nonce(),
                Price:              v.GasPrice().Int64(),
                GasLimit:           v.Gas().Int64(),
                Recipient:          to,
                Amount:             v.Value().Int64(),
                Payload:            v.Data(),
                V:                  v.Int64(),
                R:                  r.Int64(),
                S:                  s.Int64(),
            },
        })
    }
    _, err1 := c.SetTransactions(ctx, &Transactions{Txs:pbTxs})
    if err1 != nil {
        return err1
    }
    return nil     
}
func (t *TrueHybrid) Start() error{
    // Set up a connection to the server.
    conn, err := grpc.Dial(t.address, grpc.WithInsecure())
    if err != nil {
        return err
    }
    defer conn.Close()   
    c := NewPyHybConsensusClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    _, err1 := c.Start(ctx, &EmptyParam{})
    if err1 != nil {
        return err1
    }
    return nil    
}
func (t *TrueHybrid) Stop() error{
    // Set up a connection to the server.
    conn, err := grpc.Dial(t.address, grpc.WithInsecure())
    if err != nil {
        return err
    }
    defer conn.Close()   
    c := NewPyHybConsensusClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    _, err1 := c.Stop(ctx, &EmptyParam{})
    if err1 != nil {
        return err1
    }   
    return nil
}