package d2asset

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
)

const (
	animationBudget = 64
)

type animationManager struct {
	cache    d2interface.Cache
	renderer d2interface.Renderer
}

func (am *animationManager) ClearCache() {
	panic("implement me")
}

func (am *animationManager) GetCache() d2interface.Cache {
	panic("implement me")
}

func createAnimationManager(renderer d2interface.Renderer) *animationManager {
	return &animationManager{
		renderer: renderer,
		cache:    d2common.CreateCache(animationBudget),
	}
}

func (am *animationManager) LoadAnimation(
	animationPath, palettePath string,
	transparency int ) (d2interface.Animation, error) {
	cachePath := fmt.Sprintf("%s;%s;%d", animationPath, palettePath, transparency)
	if animation, found := am.cache.Retrieve(cachePath); found {
		return animation.(*Animation).Clone(), nil
	}

	var animation d2interface.Animation

	ext := strings.ToLower(filepath.Ext(animationPath))
	switch ext {
	case ".dc6":
		dc6, err := loadDC6(animationPath)
		if err != nil {
			return nil, err
		}

		palette, err := LoadPalette(palettePath)
		if err != nil {
			return nil, err
		}

		animation, err = CreateAnimationFromDC6(am.renderer, dc6, palette)
		if err != nil {
			return nil, err
		}
	case ".dcc":
		dcc, err := loadDCC(animationPath)
		if err != nil {
			return nil, err
		}

		palette, err := LoadPalette(palettePath)
		if err != nil {
			return nil, err
		}

		animation, err = CreateAnimationFromDCC(am.renderer, dcc, palette, transparency)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown animation format: %s", ext)
	}

	if err := am.cache.Insert(cachePath, animation.Clone(), 1); err != nil {
		return nil, err
	}

	return animation, nil
}
