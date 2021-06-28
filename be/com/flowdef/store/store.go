package store

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/olvrng/rbot/be/com/flowdef/types"
	"github.com/olvrng/rbot/be/pkg/dot"
	"github.com/olvrng/rbot/be/pkg/l"
)

var ll = l.New()

type FlowFile struct {
	Flows []*types.Flow `json:"flows"`
}

func (ff *FlowFile) GetByID(id dot.IntID) *types.Flow {
	for _, f := range ff.Flows {
		if f.ID == id {
			return f
		}
	}
	return nil
}

func (ff *FlowFile) GetByPageID(pageID dot.IntID) *types.Flow {
	for _, flow := range ff.Flows {
		for _, _pageID := range flow.PageIDs {
			if _pageID == pageID {
				return flow
			}
		}
	}
	return nil
}

var _ FlowStore = (*FlowFileStore)(nil)

type FlowFileStore struct {
	FilePath  string
	FlowsData *FlowFile
}

func NewFlowFileStore(filePath string) (*FlowFileStore, error) {
	s := &FlowFileStore{
		FilePath: filePath,
	}

	var flowFile *FlowFile
	_, err := os.Stat(filePath)
	switch {
	case err == nil: // load from storage
		flowFile, err = loadJson(filePath)
		s.FlowsData = flowFile
		if ce := ll.Check(l.DebugLevel, "debug"); ce != nil {
			out, _ := json.Marshal(flowFile)
			ll.Sugar().Debugf("flow-data: %s", out)
		}
		return s, err

	case os.IsNotExist(err): // try creating one
		flowFile = &FlowFile{}
		err = saveJson(filePath, flowFile)
		return s, err

	default:
		return nil, err
	}
}

func (s *FlowFileStore) SaveFlow(flow *types.Flow) *types.Flow {
	if flow.ID == 0 {
		flow.ID = dot.NewIntID()
	}
	existingFlow := s.FlowsData.GetByID(flow.ID)
	if existingFlow != nil {
		*existingFlow = *flow // overwrite
	} else {
		s.FlowsData.Flows = append(s.FlowsData.Flows, flow)
	}
	return flow
}

func (s *FlowFileStore) LoadFlowByID(id dot.IntID) *types.Flow {
	return s.FlowsData.GetByID(id)
}

func (s *FlowFileStore) LoadFlowByPageID(pageID dot.IntID) *types.Flow {
	// mock: only implement page_id
	for _, flow := range s.FlowsData.Flows {
		for _, _pageID := range flow.PageIDs {
			if _pageID == pageID {
				return flow
			}
		}
	}
	return nil
}

func loadJson(filePath string) (out *FlowFile, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &out)
	return out, err
}

func saveJson(filePath string, in *FlowFile) error {
	data, err := json.MarshalIndent(in, "", "\t")
	if err != nil {
		return nil
	}
	return ioutil.WriteFile(filePath, data, 0644)
}
