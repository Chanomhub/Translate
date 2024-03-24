package Animations

type EnemyData struct {
	ID             int           `json:"id"`
	Animation1Hue  int           `json:"animation1Hue"`
	Animation1Name string        `json:"animation1Name"`
	Animation2Hue  int           `json:"animation2Hue"`
	Animation2Name string        `json:"animation2Name"`
	Frames         [][][]float64 `json:"frames"`
	Name           string        `json:"name"`
	Position       int           `json:"position"`
	Timings        []Timings     `json:"timings"`
}

type Timings struct {
	FlashColor    []int `json:"flashColor"`
	FlashDuration int   `json:"flashDuration"`
	FlashScope    int   `json:"flashScope"`
	Frame         int   `json:"frame"`
	Se            Se    `json:"se"`
}

type Se struct {
	Name   string `json:"name"`
	Pan    int    `json:"pan"`
	Pitch  int    `json:"pitch"`
	Volume int    `json:"volume"`
}
