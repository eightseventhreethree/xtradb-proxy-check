package output

// Responses stores body values returned by api
type Responses struct {
	Disabled   string
	ReadOnly   string
	Synced     string
	Unsynced   string
	Donor      string
	Error      string
	FullStatus string
}

// Messages returns Responses
func Messages() *Responses {
	baseMsg := "Percona XtraDB Cluster Node is "
	return &Responses{
		Disabled: " manually disabled.\n",
		ReadOnly: baseMsg + "read-only.\n",
		Synced:   baseMsg + "synced.\n",
		Unsynced: baseMsg + "not synced.\n",
		Donor:    baseMsg + "a donor.\n",
		Error: baseMsg + "unavailable OR in unknown state.\n" +
			"Verify connectivity and review logs.\n",
		FullStatus: "Offline: %t\nSynced: %t\nRead-Only: %t\nDonor: %t\nError: %s\n",
	}
}
