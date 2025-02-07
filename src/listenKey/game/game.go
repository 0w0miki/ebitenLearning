package game

import (
	"ebitenLearning/src/utils"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PREPARE GameStatus = iota
	RUNNING
	FAILURE
	SUCCESS
)

type GameStatus int

type Game struct {
	cfg               *config
	p                 *plane
	enemies           map[*enemy]struct{}
	lastLoadEnemy     time.Time
	bg                *background
	point             int
	status            GameStatus
	menu              *Menu
	restartButton     *Button
	lastClickInterval time.Time
}

const (
	resourcePath = "resource"
)

func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.Width, cfg.Hight)
	ebiten.SetWindowTitle(cfg.Title)
	return &Game{
		cfg:           cfg,
		p:             loadPlane(resourcePath+"/airplane/plane/plane1.png", cfg),
		enemies:       make(map[*enemy]struct{}),
		bg:            loadBackground(resourcePath+"/background/bg_plain.jpg", 1),
		menu:          loadMenu(),
		restartButton: loadIcon("resource/icon/restart.png", 200, float64(cfg.Hight)/2, 1),
		status:        PREPARE,
	}
}

// update the running data
func (g *Game) Update() error {
	switch g.status {
	case RUNNING:
		g.bg.update()
		g.p.update(g.cfg)
		for enemy := range g.enemies {
			enemy.update(g.cfg)
		}
		g.GenerateEnemy()
		g.CollisionDetect()
	case PREPARE:
		g.menu.update(g)
	case FAILURE:
		g.restartButton.update(g, PREPARE)
		if g.status == PREPARE {
			g.p.bullets = make(map[*bullet]struct{})
			g.enemies = make(map[*enemy]struct{})
			g.p.x = float64(g.cfg.Width-g.p.image.Bounds().Dx()) / 2
			g.p.y = float64(g.cfg.Hight - g.p.image.Bounds().Dy())
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.status {
	case FAILURE:
		g.bg.draw(screen, g.cfg)
		g.p.Draw(screen, g.cfg)
		for enemy := range g.enemies {
			enemy.draw(screen, g.cfg)
		}
		g.restartButton.draw(screen)
	case RUNNING:
		g.bg.draw(screen, g.cfg)
		g.p.Draw(screen, g.cfg)
		for enemy := range g.enemies {
			enemy.draw(screen, g.cfg)
		}
	case PREPARE:
		g.menu.draw(screen, g.cfg)
	}
}

// logic size, use to zoom in/out the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.Width, g.cfg.Hight
}

func (g *Game) GenerateEnemy() {
	if time.Since(g.lastLoadEnemy).Milliseconds() > g.cfg.EnemyInterval {
		g.enemies[loadEnemy(g.cfg)] = struct{}{}
		g.lastLoadEnemy = time.Now()
	}
}

func (g *Game) CollisionDetect() {
	g.killEnemy()
	g.survival()
}

func (g *Game) killEnemy() {
	deadEnemies := make([]*enemy, 0, 10)
	deadBullets := make([]*bullet, 0, 10)
	for bullect := range g.p.bullets {
		r1 := utils.Rectangle{
			Left:  utils.Point{X: bullect.x, Y: bullect.y + float64(bullect.image.Bounds().Dy())*0.5},
			Right: utils.Point{X: bullect.x + float64(bullect.image.Bounds().Dx()), Y: bullect.y + float64(bullect.image.Bounds().Dy())},
		}
		for enemy := range g.enemies {
			r2 := utils.Rectangle{
				Left:  utils.Point{X: enemy.x + float64(enemy.image.Bounds().Dx())*0.25, Y: enemy.y + float64(enemy.image.Bounds().Dy())*0.25},
				Right: utils.Point{X: enemy.x + float64(enemy.image.Bounds().Dx())*0.75, Y: enemy.y + float64(enemy.image.Bounds().Dy())*0.75},
			}
			if utils.IsOverlappingPoint(r1, r2) {
				deadBullets = append(deadBullets, bullect)
				deadEnemies = append(deadEnemies, enemy)
				break
			}
		}
	}
	for _, bullet := range deadBullets {
		delete(g.p.bullets, bullet)
	}
	for _, enemy := range deadEnemies {
		g.point += 500
		delete(g.enemies, enemy)
	}
}

func (g *Game) survival() {
	r1 := utils.Rectangle{
		Left:  utils.Point{X: g.p.x + float64(g.p.image.Bounds().Dx())*0.5, Y: g.p.y + float64(g.p.image.Bounds().Dy())*0.5},
		Right: utils.Point{X: g.p.x + float64(g.p.image.Bounds().Dx()), Y: g.p.y + float64(g.p.image.Bounds().Dy())},
	}
	for enemy := range g.enemies {
		r3 := utils.Rectangle{
			Left:  utils.Point{X: enemy.x, Y: enemy.y},
			Right: utils.Point{X: enemy.x + float64(enemy.image.Bounds().Dx()), Y: enemy.y + float64(enemy.image.Bounds().Dy())},
		}
		if utils.IsOverlappingPoint(r1, r3) {
			g.status = FAILURE
		}
		for _, bullet := range enemy.bullets {
			r2 := utils.Rectangle{
				Left:  utils.Point{X: bullet.x, Y: bullet.y},
				Right: utils.Point{X: bullet.x + float64(bullet.image.Bounds().Dx()), Y: bullet.y + float64(bullet.image.Bounds().Dy())},
			}
			if utils.IsOverlappingPoint(r1, r2) {
				g.status = FAILURE
			}
		}
	}
}
