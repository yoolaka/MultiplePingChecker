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
	TempDB struct {
		PingHosts    map[string]*PingHostEntry
		DefaultCount int
		DefaultName  string
	}
)

var (
	PrimaryDB *TempDB = nil
	SUCCESS           = 1
	FAILED            = 0
)

func (t *TempDB) InsertPingHost(key string, value *PingHostEntry) int {
	if _, exist := t.PingHosts[key]; exist {
		return FAILED
	}
	t.PingHosts[key] = value
	return SUCCESS
}
func (t *TempDB) DeletePingHost(key string) int {
	delete(t.PingHosts, key)
	return SUCCESS
}
func (t *TempDB) SearchHost(key string) (*PingHostEntry, bool) {
	if host_entry, exist := t.PingHosts[key]; exist {
		return host_entry, exist
	} else {
		return nil, exist
	}
}
func (t *TempDB) GetPingHost() *map[string]*PingHostEntry {
	return &t.PingHosts
}

func InitDB(default_cnt int, default_name string) *TempDB {
	if PrimaryDB == nil {
		PrimaryDB = &TempDB{map[string]*PingHostEntry{}, default_cnt, default_name}
	}

	return PrimaryDB
}
