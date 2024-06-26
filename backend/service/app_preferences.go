package service

import (
	"context"
	"devbox/backend/utils/container"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/adrg/sysfont"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"

	"devbox/backend/consts"
	storage2 "devbox/backend/storage"
	"devbox/backend/types"
)

type preferencesService struct {
	pref          *storage2.PreferencesStorage
	clientVersion string
	ctx           context.Context
}

var preferences *preferencesService
var oncePreferences sync.Once

func Preferences() *preferencesService {
	if preferences == nil {
		oncePreferences.Do(func() {
			preferences = &preferencesService{
				pref:          storage2.NewPreferences(),
				clientVersion: "",
			}
		})
	}
	return preferences
}

func (p *preferencesService) Startup(ctx context.Context) {
	p.ctx = ctx
	preferences.ctx = ctx
	preferences.pref.SetContext(ctx)
}

func (p *preferencesService) GetPreferences() (resp types.JSRsp) {
	resp.Data = p.pref.GetPreferences()
	resp.Code = types.CodeOK
	return
}

func (p *preferencesService) SetPreferences(pf types.Preferences) (resp types.JSRsp) {
	err := p.pref.SetPreferences(&pf)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Code = types.CodeOK
	return
}

func (p *preferencesService) UpdatePreferences(value map[string]any) (resp types.JSRsp) {
	err := p.pref.UpdatePreferences(value)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Code = types.CodeOK
	return
}

func (p *preferencesService) RestorePreferences() (resp types.JSRsp) {
	defaultPref := p.pref.RestoreDefault()
	resp.Data = map[string]any{
		"pref": defaultPref,
	}
	resp.Code = types.CodeOK
	return
}

type FontItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (p *preferencesService) GetFontList() (resp types.JSRsp) {
	finder := sysfont.NewFinder(nil)
	fontSet := container.NewSet[string]()
	var fontList []FontItem
	for _, font := range finder.List() {
		if len(font.Family) > 0 && !strings.HasPrefix(font.Family, ".") && fontSet.Add(font.Family) {
			fontList = append(fontList, FontItem{
				Name: font.Family,
				Path: font.Filename,
			})
		}
	}
	sort.Slice(fontList, func(i, j int) bool {
		return fontList[i].Name < fontList[j].Name
	})
	resp.Data = map[string]any{
		"fonts": fontList,
	}
	resp.Code = types.CodeOK
	return
}

func (p *preferencesService) SetAppVersion(ver string) {
	if !strings.HasPrefix(ver, "v") {
		p.clientVersion = "v" + ver
	} else {
		p.clientVersion = ver
	}
}

func (p *preferencesService) GetAppVersion() (resp types.JSRsp) {
	resp.Code = types.CodeOK
	resp.Data = map[string]any{
		"version": p.clientVersion,
	}
	return
}

func (p *preferencesService) SaveWindowSize(width, height int, maximised bool) {
	if maximised {
		// do not update window size if maximised state
		p.UpdatePreferences(map[string]any{
			"behavior.windowMaximised": true,
		})
	} else if width >= consts.MinWindowWidth && height >= consts.MinWindowHeight {
		p.UpdatePreferences(map[string]any{
			"behavior.windowWidth":     width,
			"behavior.windowHeight":    height,
			"behavior.windowMaximised": false,
		})
	}
}

func (p *preferencesService) GetWindowSize() (width, height int, maximised bool) {
	data := p.pref.GetPreferences()
	width, height, maximised = data.Behavior.WindowWidth, data.Behavior.WindowHeight, data.Behavior.WindowMaximised
	if width <= 0 {
		width = consts.DefaultWindowWidth
	}
	if height <= 0 {
		height = consts.DefaultWindowHeight
	}
	return
}

func (p *preferencesService) GetWindowPosition(ctx context.Context) (x, y int) {
	data := p.pref.GetPreferences()
	x, y = data.Behavior.WindowPosX, data.Behavior.WindowPosY
	width, height := data.Behavior.WindowWidth, data.Behavior.WindowHeight
	var screenWidth, screenHeight int
	if screens, err := runtime2.ScreenGetAll(ctx); err == nil {
		for _, screen := range screens {
			if screen.IsCurrent {
				screenWidth, screenHeight = screen.Size.Width, screen.Size.Height
				break
			}
		}
	}
	if screenWidth <= 0 || screenHeight <= 0 {
		screenWidth, screenHeight = consts.DefaultWindowWidth, consts.DefaultWindowHeight
	}
	if x <= 0 || x+width > screenWidth || y <= 0 || y+height > screenHeight {
		// out of screen, reset to center
		x, y = (screenWidth-width)/2, (screenHeight-height)/2
	}
	return
}

func (p *preferencesService) SaveWindowPosition(x, y int) {
	if x > 0 || y > 0 {
		p.UpdatePreferences(map[string]any{
			"behavior.windowPosX": x,
			"behavior.windowPosY": y,
		})
	}
}

type latestRelease struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

func (p *preferencesService) CheckForUpdate() (resp types.JSRsp) {
	// request latest version
	res, err := http.Get("https://api.github.com/repos/eyotang/dev-box/releases/latest")
	if err != nil || res.StatusCode != http.StatusOK {
		resp.Msg = "network error"
		return
	}

	var respObj latestRelease
	err = json.NewDecoder(res.Body).Decode(&respObj)
	if err != nil {
		resp.Msg = "invalid content"
		return
	}

	// compare with current version
	resp.Code = types.CodeOK
	resp.Data = map[string]any{
		"version":  p.clientVersion,
		"latest":   respObj.TagName,
		"page_url": respObj.HtmlUrl,
	}
	return
}
