package elf

type Toplevel struct {
	Title   string   `json:"title"`
	Modules []Module `json:"modules"`
}

type Module struct {
	Entries   []Entry `json:"entries"`
	DebugInfo string  `json:"debug_info"`
}

type Entry []any

type Report struct {
	File  string `json:"file"`
	Name  string `json:"name"`
	Cases []Case `json:"cases" mapstructure:"cases"`
}

type Case struct {
	Binders        []Binder `mapstructure:"binders"`
	Context        string   `mapstructure:"context"`
	Depths         int      `mapstructure:"depths"`
	Code           string   `mapstructure:"code"`
	Plain          int      `mapstructure:"plain"`
	Weighed        float64  `mapstructure:"weighed"`
	Detail         string   `mapstructure:"detail"`
	SimilarityRate float64  `mapstructure:"similarity_rate"`
	Sources        []Source `mapstructure:"srcs"`
}

type Binder struct {
	Binder string `json:"binder"`
	Pos    string `json:"pos"`
}

type Source struct {
	Context string `json:"context"`
	Code    string `json:"code"`
}
