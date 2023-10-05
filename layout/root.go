package layout

type Root struct {
	Parameters Parameters `json:"parameters"`
	Chains     Chains     `json:"chains"`
	Layouts    Layouts    `json:"layouts"`
}
