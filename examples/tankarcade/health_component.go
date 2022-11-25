package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
)

type healthComponent struct {
	hp           *float64
	hpPercentage float64
	body         *physics.Body
	maxHP        float64
	isBuilding   bool
	shieldLevel  int
}

func newHealthComponent(hp *float64, body *physics.Body) *healthComponent {
	return &healthComponent{hp: hp, maxHP: *hp, hpPercentage: 1, body: body}
}

func (c *healthComponent) IsAlive() bool {
	return *c.hp > 0
}

func (c *healthComponent) ApplyDamage(damage float64) bool {
	if damage < 0 && c.isBuilding {
		damage *= 2
	}
	*c.hp -= damage
	if damage < 0 && *c.hp > c.maxHP {
		*c.hp = c.maxHP
	}
	if *c.hp <= 0 {
		c.hpPercentage = 0
		*c.hp = 0
		return false
	}
	c.hpPercentage = *c.hp / c.maxHP
	return true
}

func (c *healthComponent) CheckProjectileCollisions(scene *ge.Scene) bool {
	for _, collision := range scene.GetCollisions(c.body) {
		switch obj := collision.Body.Object.(type) {
		case *battleMine:
			obj.Destroy()
			c.ApplyDamage(mineLayerDesign.extra.damage)
			scene.Audio().PlaySound(mineLayerDesign.extra.hitSound)
		case *projectile:
			if obj.CanHit() {
				switch c.shieldLevel {
				case 0:
					obj.Destroy()
					c.ApplyDamage(obj.config.design.damage)
				case 1:
					scene.Audio().PlaySound(AudioShieldAbsorb)
					obj.Destroy()
				case 2:
					scene.Audio().PlaySound(AudioShieldAbsorb)
					if obj.config.design.canReflect {
						obj.Reflect()
					} else {
						obj.Destroy()
					}
				}
			}
		}
	}
	return c.IsAlive()
}
