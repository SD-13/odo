package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/redhat-developer/odo/pkg/api"
	"github.com/redhat-developer/odo/pkg/testingutil/filesystem"
)

type State struct {
	content Content
	fs      filesystem.Filesystem
}

var _ Client = (*State)(nil)

func NewStateClient(fs filesystem.Filesystem) *State {
	return &State{
		fs: fs,
	}
}

func (o *State) SetForwardedPorts(fwPorts []api.ForwardedPort) error {
	// TODO(feloy) When other data is persisted into the state file, it will be needed to read the file first
	o.content.ForwardedPorts = fwPorts
	return o.save()
}

func (o *State) GetForwardedPorts() ([]api.ForwardedPort, error) {
	err := o.read()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil // if the state file does not exist, no ports are forwarded
		}
		return nil, err
	}
	return o.content.ForwardedPorts, err
}

func (o *State) SaveExit() error {
	o.content.ForwardedPorts = nil
	return o.save()
}

// save writes the content structure in json format in file
func (o *State) save() error {
	jsonContent, err := json.MarshalIndent(o.content, "", " ")
	if err != nil {
		return err
	}
	// .odo directory is supposed to exist, don't create it
	dir := filepath.Dir(_filepath)
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}
	return o.fs.WriteFile(_filepath, jsonContent, 0644)
}

func (o *State) read() error {
	jsonContent, err := o.fs.ReadFile(_filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonContent, &o.content)
}
