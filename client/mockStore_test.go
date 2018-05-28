package client_test

import (
	"gitlab.33.cn/chain33/chain33/queue"
	"gitlab.33.cn/chain33/chain33/types"
)

type mockStore struct {
}

func (m *mockStore) SetQueueClient(q queue.Queue) {
	go func() {
		client := q.Client()
		client.Sub("store")
		for msg := range client.Recv() {
			switch msg.Ty {
			case types.EventStoreGet:
				msg.Reply(client.NewMessage("store", types.EventStoreGetReply, &types.StoreReplyValue{}))
			case types.EventStoreGetTotalCoins:
				if req, ok := msg.GetData().(*types.IterateRangeByStateHash); ok {
					if req.Count == 10 {
						msg.Reply(client.NewMessage("store", types.EventStoreGetReply, &types.Transaction{}))
					} else {
						msg.Reply(client.NewMessage("store", types.EventStoreGetReply, &types.ReplyGetTotalCoins{}))
					}
				} else {
					msg.ReplyErr("Do not support", types.ErrInvalidParam)
				}
			default:
				msg.ReplyErr("Do not support", types.ErrNotSupport)
			}
		}
	}()
}

func (m *mockStore) Close() {
}