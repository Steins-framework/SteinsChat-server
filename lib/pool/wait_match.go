package pool

import (
	"chat/models"
	"container/list"
)

var WaitMatchPool = WaitMatch{pool: list.New()}

type WaitMatch struct {
	pool *list.List
}

func (w *WaitMatch) Add(user *models.User) {
	w.pool.PushBack(user)
}

func (w *WaitMatch) Pop() *models.User {
	if w.pool.Len() == 0 {
		return nil
	}
	return w.pool.Remove(w.pool.Front()).(*models.User)
}

func (w *WaitMatch) Len() int {
	return w.pool.Len()
}