package models

import (
	"context"
	"errors"
	"github.com/koko990/logcollector/models/dao"
	"github.com/koko990/logcollector/models/merge"
	"sync"
)

//
func (StoreJob) Run(m *sync.Mutex) error {
	rootCtx := context.Background()
	mergeModelCtx, notifyMerge := context.WithCancel(rootCtx)
	pushModelCtx, notifyPush := context.WithCancel(rootCtx)
	go merge.RunMergeModel(mergeModelCtx, m, notifyMerge)
	go dao.PushMultiModelToDB(pushModelCtx, notifyPush)
	select {
	case <-mergeModelCtx.Done():
		return errors.New("merge job is stop!!")
	case <-pushModelCtx.Done():
		return errors.New("push db is done")
	}
	return nil
}
