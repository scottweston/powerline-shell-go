package config

type Configuration struct {
	ShowWritable   bool `json:"showWritable"`
	ShowVirtualEnv bool `json:"showVirtualEnv"`
	ShowCwd        bool `json:"showCwd"`
	CwdMaxLength   int  `json:"cwdMaxLength"`
	BatteryWarn    int  `json:"batteryWarn"`
	ShowGit        bool `json:"showGit"`
	ShowHg         bool `json:"showHg"`
	ShowReturnCode bool `json:"showReturnCode"`
	Icons          struct {
		Powerline struct {
			Added         string `json:"added"`
			Ahead         string `json:"ahead"`
			Behind        string `json:"behind"`
			Branch        string `json:"branch"`
			Conflicted    string `json:"conflicted"`
			Detached      string `json:"detached"`
			Ellipsis      string `json:"ellipsis"`
			Modified      string `json:"modified"`
			Phases        string `json:"phases"`
			ReadOnly      string `json:"readonly"`
			Removed       string `json:"removed"`
			SeparatorThin string `json:"separatorthin"`
			Separator     string `json:"separator"`
			Untracked     string `json:"untracked"`
		} `json:"powerline"`
		Plain struct {
			Added         string `json:"added"`
			Ahead         string `json:"ahead"`
			Behind        string `json:"behind"`
			Branch        string `json:"branch"`
			Conflicted    string `json:"conflicted"`
			Detached      string `json:"detached"`
			Ellipsis      string `json:"ellipsis"`
			Modified      string `json:"modified"`
			Phases        string `json:"phases"`
			ReadOnly      string `json:"readonly"`
			Removed       string `json:"removed"`
			SeparatorThin string `json:"separatorthin"`
			Separator     string `json:"separator"`
			Untracked     string `json:"untracked"`
		} `json:"plain"`
	} `json:"icons"`
	Colours struct {
		Hg struct {
			BackgroundDefault int `json:"backgroundDefault"`
			BackgroundChanges int `json:"backgroundChanges"`
			Text              int `json:"text"`
		} `json:"hg"`
		Git struct {
			BackgroundDefault int `json:"backgroundDefault"`
			BackgroundChanges int `json:"backgroundChanges"`
			Text              int `json:"text"`
		} `json:"git"`
		Cwd struct {
			Background     int `json:"background"`
			Text           int `json:"text"`
			HomeBackground int `json:"homeBackground"`
			HomeText       int `json:"homeText"`
		} `json:"cwd"`
		Virtualenv struct {
			Background int `json:"background"`
			Text       int `json:"text"`
		} `json:"virtualenv"`
		Returncode struct {
			Background int `json:"background"`
			Text       int `json:"text"`
		} `json:"returncode"`
		Lock struct {
			Background int `json:"background"`
			Text       int `json:"text"`
		} `json:"lock"`
		Dollar struct {
			Background int `json:"background"`
			Text       int `json:"text"`
		} `json:"dollar"`
		Battery struct {
			Background int `json:"background"`
			Text       int `json:"text"`
		} `json:"battery"`
	} `json:"colours"`
}

func (self *Configuration) SetDefaults() {
	self.ShowWritable = true
	self.ShowVirtualEnv = true
	self.ShowCwd = true
	self.CwdMaxLength = 0
	self.BatteryWarn = 0
	self.ShowGit = true
	self.ShowHg = true
	self.ShowReturnCode = true
	self.Colours.Hg.BackgroundDefault = 22
	self.Colours.Hg.BackgroundChanges = 64
	self.Colours.Hg.Text = 251
	self.Colours.Git.BackgroundDefault = 17
	self.Colours.Git.BackgroundChanges = 21
	self.Colours.Git.Text = 251
	self.Colours.Cwd.Background = 40
	self.Colours.Cwd.Text = 237
	self.Colours.Cwd.HomeBackground = 31
	self.Colours.Cwd.HomeText = 15
	self.Colours.Virtualenv.Background = 35
	self.Colours.Virtualenv.Text = 0
	self.Colours.Returncode.Background = 196
	self.Colours.Returncode.Text = 16
	self.Colours.Lock.Background = 124
	self.Colours.Lock.Text = 254
	self.Colours.Dollar.Background = 240
	self.Colours.Dollar.Text = 15
	self.Colours.Battery.Background = 196
	self.Colours.Battery.Text = 16
}
