// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.


package rpc

import (
	"context"
	"log"
	"sync"
)

// handler handles JSON-RPC messages. There is one handler per connection. Note that
// handler is not safe for concurrent use. Message handling never blocks indefinitely
// because RPCs are processed on background goroutines launched by handler.
//
// The entry points for incoming messages are:
//
//	h.handleMsg(message)
//	h.handleBatch(message)
//
// Outgoing calls use the requestOp struct. Register the request before sending it
// on the connection:
//
//	op := &requestOp{ids: ...}
//	h.addRequestOp(op)
//
// Now send the request, then wait for the reply to be delivered through handleMsg:
//
//	if err := op.wait(...); err != nil {
//		h.removeRequestOp(op) // timeout, etc.
//	}
type handler struct {
	respWait             map[string]*requestOp          // active client requests
	callWG               sync.WaitGroup                 // pending call goroutines
	rootCtx              context.Context                // canceled by close()
	cancelRoot           func()                         // cancel function for rootCtx
	conn                 jsonWriter                     // where responses will be sent
	log                  log.Logger
	allowSubscribe       bool
	batchRequestLimit    int
	batchResponseMaxSize int

	subLock    sync.Mutex
}

type callProc struct {
	ctx       context.Context
}

func newHandler(connCtx context.Context, conn jsonWriter, batchRequestLimit, batchResponseMaxSize int) *handler {
	rootCtx, cancelRoot := context.WithCancel(connCtx)
	h := &handler{
		conn:                 conn,
		rootCtx:              rootCtx,
		cancelRoot:           cancelRoot,
		allowSubscribe:       true,
		batchRequestLimit:    batchRequestLimit,
		batchResponseMaxSize: batchResponseMaxSize,
	}
	if conn.remoteAddr() != "" {
		// todo log
	}
	//h.unsubscribeCb = newCallback(reflect.Value{}, reflect.ValueOf(h.unsubscribe))
	return h
}

// batchCallBuffer manages in progress call messages and their responses during a batch
// call. Calls need to be synchronized between the processing and timeout-triggering
// goroutines.
type batchCallBuffer struct {
	mutex sync.Mutex
	calls []*jsonrpcMessage
	resp  []*jsonrpcMessage
	wrote bool
}




// close cancels all requests except for inflightReq and waits for
// call goroutines to shut down.
func (h *handler) close(err error, inflightReq *requestOp) {
	h.cancelAllRequests(err, inflightReq)
	h.callWG.Wait()
	h.cancelRoot()
	//h.cancelServerSubscriptions(err)
}

// addRequestOp registers a request operation.
func (h *handler) addRequestOp(op *requestOp) {
	for _, id := range op.ids {
		h.respWait[string(id)] = op
	}
}

// removeRequestOps stops waiting for the given request IDs.
func (h *handler) removeRequestOp(op *requestOp) {
	for _, id := range op.ids {
		delete(h.respWait, string(id))
	}
}

// cancelAllRequests unblocks and removes pending requests and active subscriptions.
func (h *handler) cancelAllRequests(err error, inflightReq *requestOp) {
	didClose := make(map[*requestOp]bool)
	if inflightReq != nil {
		didClose[inflightReq] = true
	}

	for id, op := range h.respWait {
		// Remove the op so that later calls will not close op.resp again.
		delete(h.respWait, id)

		if !didClose[op] {
			op.err = err
			close(op.resp)
			didClose[op] = true
		}
	}
	//for id, sub := range h.clientSubs {
	//	delete(h.clientSubs, id)
	//	sub.close(err)
	//}
}
