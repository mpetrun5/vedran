package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/NodeFactoryIo/vedran/internal/rpc"
	log "github.com/sirupsen/logrus"
)

func (c ApiController) RPCHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	reqBody, _ := ioutil.ReadAll(r.Body)
	isBatch := rpc.IsBatch(reqBody)
	var reqRPCBody rpc.RPCRequest
	var reqRPCBodies []rpc.RPCRequest
	var err error
	if isBatch {
		err = json.Unmarshal(reqBody, &reqRPCBodies)
	} else {
		err = json.Unmarshal(reqBody, &reqRPCBody)
	}

	if err != nil {
		log.Errorf("Request failed because of: %v", err)
		_ = json.NewEncoder(w).Encode(
			rpc.CreateRPCError(isBatch, reqRPCBody, reqRPCBodies, rpc.ParseError, "Parse error"))
		return
	}

	nodes, err := c.nodeRepo.GetActiveNodes()
	if err != nil || len(*nodes) == 0 {
		log.Error("Request failed because vedran has no available nodes")
		_ = json.NewEncoder(w).Encode(
			rpc.CreateRPCError(isBatch, reqRPCBody, reqRPCBodies, rpc.InternalServerError, "No available nodes"))
		return
	}

	// @TODO: Peer selection code

	for _, node := range *nodes {
		rpcResponse, err := rpc.SendRequestToNode(isBatch, node, reqBody)
		if err != nil {
			log.Errorf("Request failed to node %s because of: %v", node.ID, err)
			continue
		}

		_ = json.NewEncoder(w).Encode(rpcResponse)
		return
	}

	log.Error("Request failed because all nodes returned invalid rpc response")
	_ = json.NewEncoder(w).Encode(
		rpc.CreateRPCError(isBatch, reqRPCBody, reqRPCBodies, rpc.InternalServerError, "Internal Server Error"))
}
