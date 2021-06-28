package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olvrng/rbot/be/com/flowexec/flowcore"
	"github.com/olvrng/rbot/be/pkg/dot"
	"github.com/olvrng/rbot/be/pkg/xerrors"
)

type FlowFile struct {
	Last    map[string]*flowcore.FlowState   `json:"last"`
	History map[string][]*flowcore.FlowState `json:"history"`
}

func NewFlowFile() *FlowFile {
	return &FlowFile{
		Last:    make(map[string]*flowcore.FlowState),
		History: make(map[string][]*flowcore.FlowState),
	}
}

type FlowStateStore struct {
	FilePath string
	Data     *FlowFile
}

func NewFlowStateStore(filePath string) (*FlowStateStore, error) {
	s := &FlowStateStore{
		FilePath: filePath,
	}

	var data []byte
	_, err := os.Stat(filePath)
	switch {
	case err == nil: // load from storage
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &s.Data)
		if err != nil {
			return nil, err
		}
		if s.Data.Last == nil {
			s.Data.Last = make(map[string]*flowcore.FlowState)
		}
		if s.Data.History == nil {
			s.Data.History = make(map[string][]*flowcore.FlowState)
		}
		return s, err

	case os.IsNotExist(err): // try creating one
		s.Data = NewFlowFile()
		err = storeFile(filePath, s.Data)
		return s, err

	default:
		return nil, err
	}
}

func (s *FlowStateStore) SaveState(state *flowcore.FlowState) error {
	runID := mockEncodeRunID(state.PageID, state.PSID)
	data := s.Data
	data.Last[runID] = state
	data.History[runID] = append(data.History[runID], state)

	return storeFile(s.FilePath, data)
}

func (s *FlowStateStore) LoadState(ctx context.Context, pageID, psid dot.IntID) (*flowcore.FlowState, error) {
	data := s.Data
	runID := mockEncodeRunID(pageID, psid)
	state := data.Last[runID]
	if state == nil {
		return nil, xerrors.Errorf(xerrors.NotFound, nil, "not found")
	}
	return state, nil
}

func mockEncodeRunID(pageID, psID dot.IntID) string {
	return fmt.Sprintf("run:%v_%v", pageID, psID)
}

func storeFile(filePath string, data *FlowFile) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, out, 0644)
}
