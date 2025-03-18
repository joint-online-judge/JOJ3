package stage

// copied from https://github.com/criyle/go-judge/blob/master/cmd/go-judge/model/model.go
// CmdFile defines file from multiple source including local / memory / cached or pipe collector
type CmdFile struct {
	Src       *string `json:"src"`
	Content   *string `json:"content"`
	FileID    *string `json:"fileId"`
	Name      *string `json:"name"`
	Max       *int64  `json:"max"`
	Symlink   *string `json:"symlink"`
	StreamIn  bool    `json:"streamIn"`
	StreamOut bool    `json:"streamOut"`
	Pipe      bool    `json:"pipe"`
}

// Cmd defines command and limits to start a program using in envexec
type Cmd struct {
	Args   []string `json:"args"`
	Env    []string `json:"env,omitempty"`
	Stdin  *CmdFile `json:"stdin,omitempty"`
	Stdout *CmdFile `json:"stdout,omitempty"`
	Stderr *CmdFile `json:"stderr,omitempty"`

	CPULimit     uint64 `json:"cpuLimit"`     // ns
	RealCPULimit uint64 `json:"realCpuLimit"` // deprecated: use clock limit instead (still working)
	ClockLimit   uint64 `json:"clockLimit"`   // ns
	MemoryLimit  uint64 `json:"memoryLimit"`  // byte
	StackLimit   uint64 `json:"stackLimit"`   // byte
	ProcLimit    uint64 `json:"procLimit"`
	CPURateLimit uint64 `json:"cpuRateLimit"` // limit cpu usage (1000 equals 1 cpu)
	CPUSetLimit  string `json:"cpuSetLimit"`  // set the cpuSet for cgroup

	CopyIn       map[string]CmdFile `json:"copyIn"`
	CopyInCached map[string]string  `json:"copyInCached"`
	CopyInDir    string             `json:"copyInDir"`

	CopyOut       []string `json:"copyOut"`
	CopyOutCached []string `json:"copyOutCached"`
	CopyOutMax    uint64   `json:"copyOutMax"`
	CopyOutDir    string   `json:"copyOutDir"`

	TTY               bool `json:"tty,omitempty"`
	StrictMemoryLimit bool `json:"strictMemoryLimit"`
	DataSegmentLimit  bool `json:"dataSegmentLimit"`
	AddressSpaceLimit bool `json:"addressSpaceLimit"`
}
