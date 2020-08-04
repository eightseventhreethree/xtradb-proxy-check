package output

// Responses stores body values returned by api
type Responses struct {
	Disabled   string
	ReadOnly   string
	Synced     string
	Unsynced   string
	FullStatus string
}

// Messages returns Responses
func Messages() *Responses {
	return &Responses{
		Disabled:   "Percona XtraDB Cluster Node is manually disabled.\n",
		ReadOnly:   "Percona XtraDB Cluster Node is read-only.\n",
		Synced:     "Percona XtraDB Cluster Node is synced.\n",
		Unsynced:   "Percona XtraDB Cluster Node is not synced.\n",
		FullStatus: "Offline: %t\nSynced: %t\nRead-Only: %t\n",
	}
}
