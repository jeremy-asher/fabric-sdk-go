/*
Copyright SecureKey Technologies Inc. All Rights Reserved.


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

package fabricclient

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	rwsetutil "github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	kvrwset "github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"

	pb "github.com/hyperledger/fabric/protos/peer"
)

// MockEndorserServer mock endoreser server to process endorsement proposals
type MockEndorserServer struct {
	ProposalError error
	AddkvWrite    bool
}

// ProcessProposal mock implementation that returns success if error is not set
// error if it is
func (m *MockEndorserServer) ProcessProposal(context context.Context,
	proposal *pb.SignedProposal) (*pb.ProposalResponse, error) {
	if m.ProposalError == nil {
		return &pb.ProposalResponse{Response: &pb.Response{
			Status: 200,
		}, Endorsement: &pb.Endorsement{Endorser: []byte("endorser"), Signature: []byte("signature")},
			Payload: m.createProposalResponsePayload()}, nil
	}
	return &pb.ProposalResponse{Response: &pb.Response{
		Status:  500,
		Message: m.ProposalError.Error(),
	}}, m.ProposalError
}

func (m *MockEndorserServer) createProposalResponsePayload() []byte {

	prp := &pb.ProposalResponsePayload{}
	ccAction := &pb.ChaincodeAction{}
	txRwSet := &rwsetutil.TxRwSet{}

	if m.AddkvWrite {
		txRwSet.NsRwSets = []*rwsetutil.NsRwSet{
			&rwsetutil.NsRwSet{NameSpace: "ns1", KvRwSet: &kvrwset.KVRWSet{
				Reads:  []*kvrwset.KVRead{&kvrwset.KVRead{Key: "key1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}}},
				Writes: []*kvrwset.KVWrite{&kvrwset.KVWrite{Key: "key2", IsDelete: false, Value: []byte("value2")}},
			}}}
	}

	txRWSetBytes, err := txRwSet.ToProtoBytes()
	if err != nil {
		return nil
	}
	ccAction.Results = txRWSetBytes
	ccActionBytes, err := proto.Marshal(ccAction)
	if err != nil {
		return nil
	}
	prp.Extension = ccActionBytes
	prpBytes, err := proto.Marshal(prp)
	if err != nil {
		return nil
	}
	return prpBytes
}