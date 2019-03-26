package temp_db

import (
	"os/exec"
)

type (
	PingHostEntry struct {
		HostName string
		Count    int
		Channel  chan string
		Buffer   string
		PingCmd  *exec.Cmd
	}
	tempDB struct {
		PingHosts    map[string]*PingHostEntry
		DefaultCount int
		DefaultName  string
	}
	TempDB interface {
		InsertPingHost(key string, value *PingHostEntry) int
		DeletePingHost(key string) int
		SearchHost(key string) (*PingHostEntry, bool)
		GetPingHost() *map[string]*PingHostEntry
	}
)

var (
	PrimaryDB *tempDB = nil
	SUCCESS           = 1
	FAILED            = 0
)

func (t tempDB) InsertPingHost(key string, value *PingHostEntry) int {
	if _, exist := t.PingHosts[key]; exist {
		return FAILED
	}
	t.PingHosts[key] = value
	return SUCCESS
}
func (t tempDB) DeletePingHost(key string) int {
	delete(t.PingHosts, key)
	return SUCCESS
}
func (t tempDB) SearchHost(key string) (*PingHostEntry, bool) {
	if host_entry, exist := t.PingHosts[key]; exist {
		return host_entry, exist
	} else {
		return nil, exist
	}
}
func (t tempDB) GetPingHost() *map[string]*PingHostEntry {
	return &t.PingHosts
}

func InitDB(default_cnt int, default_name string) tempDB {
	if PrimaryDB == nil {
		PrimaryDB = &tempDB{map[string]*PingHostEntry{}, default_cnt, default_name}
	}

	return *PrimaryDB
}
