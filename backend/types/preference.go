package types

import "devbox/backend/consts"

type Preferences struct {
	Behavior PreferencesBehavior `json:"behavior" yaml:"behavior"`
	General  PreferencesGeneral  `json:"general" yaml:"general"`
	Editor   PreferencesEditor   `json:"editor" yaml:"editor"`
}

func NewPreferences() Preferences {
	return Preferences{
		Behavior: PreferencesBehavior{
			AsideWidth:   consts.DefaultAsideWidth,
			WindowWidth:  consts.DefaultWindowWidth,
			WindowHeight: consts.DefaultWindowHeight,
		},
		General: PreferencesGeneral{
			Theme:       "auto",
			Language:    "auto",
			FontSize:    consts.DefaultFontSize,
			CheckUpdate: true,
			AllowTrack:  true,
		},
		Editor: PreferencesEditor{
			FontSize:    consts.DefaultFontSize,
			ShowLineNum: true,
			ShowFolding: true,
			DropText:    true,
			Links:       true,
		},
	}
}

type PreferencesBehavior struct {
	Welcomed        bool `json:"welcomed" yaml:"welcomed"`
	AsideWidth      int  `json:"asideWidth" yaml:"aside_width"`
	WindowWidth     int  `json:"windowWidth" yaml:"window_width"`
	WindowHeight    int  `json:"windowHeight" yaml:"window_height"`
	WindowMaximised bool `json:"windowMaximised" yaml:"window_maximised"`
	WindowPosX      int  `json:"windowPosX" yaml:"window_pos_x"`
	WindowPosY      int  `json:"windowPosY" yaml:"window_pos_y"`
}

type PreferencesGeneral struct {
	Theme           string   `json:"theme" yaml:"theme"`
	Language        string   `json:"language" yaml:"language"`
	Font            string   `json:"font" yaml:"font,omitempty"`
	FontFamily      []string `json:"fontFamily" yaml:"font_family,omitempty"`
	FontSize        int      `json:"fontSize" yaml:"font_size"`
	UseSysProxy     bool     `json:"useSysProxy" yaml:"use_sys_proxy,omitempty"`
	UseSysProxyHttp bool     `json:"useSysProxyHttp" yaml:"use_sys_proxy_http,omitempty"`
	CheckUpdate     bool     `json:"checkUpdate" yaml:"check_update"`
	SkipVersion     string   `json:"skipVersion" yaml:"skip_version,omitempty"`
	AllowTrack      bool     `json:"allowTrack" yaml:"allow_track"`
}

type PreferencesEditor struct {
	Font        string   `json:"font" yaml:"font,omitempty"`
	FontFamily  []string `json:"fontFamily" yaml:"font_family,omitempty"`
	FontSize    int      `json:"fontSize" yaml:"font_size"`
	ShowLineNum bool     `json:"showLineNum" yaml:"show_line_num"`
	ShowFolding bool     `json:"showFolding" yaml:"show_folding"`
	DropText    bool     `json:"dropText" yaml:"drop_text"`
	Links       bool     `json:"links" yaml:"links"`
}
