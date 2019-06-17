package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/vardius/go-api-boilerplate/cmd/auth/infrastructure/persistence"
	oauth2 "gopkg.in/oauth2.v3"
	oauth2_models "gopkg.in/oauth2.v3/models"
)

// NewClientStore create client store
func NewClientStore(repository persistence.ClientRepository) *ClientStore {
	return &ClientStore{
		repository: repository,
		internal:   make(map[string]oauth2.ClientInfo),
	}
}

// ClientStore client information store
type ClientStore struct {
	sync.RWMutex
	internal   map[string]oauth2.ClientInfo
	repository persistence.ClientRepository
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	i, err := cs.Internal(id)
	if err == nil {
		return i, nil
	}

	c, err := cs.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return cs.toClientInfo(c.GetData())
}

// Internal according to the ID for the internal client information
func (cs *ClientStore) Internal(id string) (cli oauth2.ClientInfo, err error) {
	cs.RLock()
	defer cs.RUnlock()
	if c, ok := cs.internal[id]; ok {
		cli = c
		return
	}
	err = errors.New("not found")
	return
}

// SetInternal set internal system client information
func (cs *ClientStore) SetInternal(id string, cli oauth2.ClientInfo) (err error) {
	cs.Lock()
	defer cs.Unlock()
	cs.internal[id] = cli
	return
}

func (cs *ClientStore) toClientInfo(data []byte) (oauth2.ClientInfo, error) {
	info := oauth2_models.Client{}
	err := json.Unmarshal(data, &info)

	return &info, err
}
