​A simple real-time strategy game optimized for gamepads and hot seat battles against other players.

Connect up to 4 gamepads and play the battles with your friends using a single computer!

## How to play: a quick guide

* Send the starting builder to some nearby sector
* Build builders to expand, you need many sectors
* Order some tanks to defend your bases and to attack the enemy
* When under attack, consider to build a defensive turret (E on keyboard)

## Game basics overview

Players start with a base (HQ) and one builder. A builder can be sent to some free sector and capture it by building a new base. Bases (including HQ) can be used to build new units.

Every unit (tanks) consist of two parts: a hull and a turret. Depending on what resources and situation you have, different combinations can be considered optimal. You should always try to utilize your resources in the best way possible.

There are 3 types of resources:

* Iron (metal nugget icon)
* Gold (yellow sphere icon)
* Oil (blue barrel icon)

Normally, a sector has one resource assigned to it. When you build a base in that sector, you will get +1 of that resource per 5 seconds. Starting sectors have a "combined" resource, they produce 1 unit of each resource kind.

When you build a tank, it will stay near the base that produced it. You can select all tanks in the sector by pressing Space, then select a target sector and press Space again to send the units. After you send the units, you can't control them.

If your squad had a builder, it will capture a sector if it is not yet occupied. Note that if manage to build a base, all escort tanks will stay near that newly established base.

## Controls

In menu:

* W, S - selects the focused menu button
* A, D - loop through the focused button options
* Enter - activate button / toggle option
* Escape - to the previous screen

In game:

* W, A, S, D - move the sector selector
* Space - select units / send units
* Enter - open sector menu (when owned sector is selected)
* E - build a defensive turret (when owned sector is selected)
* Q - cancel selection / action
* Escape - back to the menu

In sector menu (also in game):

* W, S - select turret/hull entry
* A, S - loop through the selected item (turret or hull) options
* Space - order selected design
* Enter, Q - close the menu

### Gamepad controls

![](https://user-images.githubusercontent.com/580022/45268303-10a03e80-b4ce-11e8-883c-1f586566c040.png)

> We use Xbox controller button names in this section.
> Up, Down, Left, Right are the buttons from the D-pad.

In menu:

* Up, Down - selects the focused menu button
* Left, Right - loop through the focused button options

In game:

* Up, Down, Left, Right - move the sector selector
* A - select units / send units
* X - open sector menu (when owned sector is selected)
* Y - build a defensive turret (when owned sector is selected)
* B - cancel selection / action

In sector menu:

* Up, Down - select turret/hull entry
* Left, Right - loop through the selected item (turret or hull) options
* A - order selected design
* X, B - close the menu

## New game configuration

Extra game rules (can enable more than one):

* Close combat: make starting locations closer to each other.
* Barren center: central sectors have no resources.
* Doubled income: all income is doubled (2 per sector instead of 1).
* Quick start: players start with two bases instead of one.
* Balanced Resources: generate a map with fair resources distribution
* HQ siege: losing an HQ base casuses the immediate defeat.
* No fortifications: building battle post fortifications (turrets) is prohibited.
* Mud terrain: causes tanks to move and turn slower.

Team modes:

* 2 vs 2: slot1+slot2 versus slot3+slot4
* 1 vs 3: slot1 versus slot2+slot3+slot4
* deathmatch: everyone versus everyone
* vs leader: dynamically changing alliances

In `vs leader` teams mode, you start like in `deathmatch`, but every few seconds a "leader" is selected, making everyone except that leader allied against that leader. A leader is a player that has the most sectors captured.
