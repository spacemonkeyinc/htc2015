# Future Tank

## Game

Outlive the other tank! Employ futuristic weaponry to defeat your enemy!

Your tank lives on a grid of cells. Tanks can be oriented north, south, east,
or west. Tanks can move forward, turn 90 degrees left or right, and fire.

Tanks start off with a certain amount of health and energy. As the game
progresses, tanks lose health naturally, and health is only replensished
in small amounts by obtaining batteries.

If a tank gets shot, it loses a large amount of health. Tank energy is used
during weapon fire, but can be replenished by driving to cells that contain
batteries. Batteries are dropped at random times and locations throughout the
course of a game.

The game proceeds in ticks. Every two game ticks, tanks are asked to make a
decision about what to do next. The tank's next action is performed in the
next game tick. The other game tick is used for making lasers move faster
than tanks can move.

If your tank outlives the other tank, you win!

## API

### Starting the game

Since multiple games can run at the same time on the same game server, you must
agree with your opponent on a game id. Let's assume your game id is `tankyou`.

To join the `tankyou` game, you must do a `POST` to
`http://gameserver:8080/game/tankyou/join`, with header
`X-Sm-Playermoniker: yourname`. If your game is being televised, your player
moniker will show up on the scoreboard. With `curl`, this looks like:

```
curl -X POST -H 'X-Sm-Playermoniker: yourname' http://gameserver:8080/game/tankyou
```

This request will not return until the game has started and its your turn to
move.

The response will include the `X-Sm-Playerid` header, which you will need to
save and include in all future action requests.

The response will include a JSON object (the one described below) with an
additional field called `config`. `config` itself has the following fields:

```
  * `turn_timeout` - how long you have each turn to take your turn, in
   nanoseconds. If you take longer than this time, then you default to a
   noop action.
  * `connect_back_timeout` - a timeout value in seconds in which you have to respond
   before we assume you are no longer playing and you self destruct.
  * `max_health` - how much health you start with
  * `max_energy` - the maximum amount of energy possible, per player
  * `health_loss` - how much health you automatically lose each turn
  * `laser_damage` - how much health is subtracted when hit by a laser
  * `laser_distance` - how many cells a laser travels before fizzing out
  * `laser_energy` - how much energy it takes to fire a laser
  * `battery_power` - how much energy is restored by picking up a battery, up
   to the `maximum_energy` limit
  * `battery_health` - how much health is restored by picking up a battery, up
   to the `maximum_health` limit
```

### Turns

Your turn begins when your previous `POST` request returns the current game
state to you. Game state will be a JSON object like the following example:

```
{
	"status": "running",
	"health": 200,
	"energy": 10,
	"orientation": "north",
	"grid": <grid>
}
```

 * `status` - can either be `running`, `won`, `lost`, or `draw`. When the
  status is not `running`, you are not expected to make future requests. The
  game is over.
 * `health` - An integer, starts off at the max health possible and decreases
  over time, possibly rapidly if you're getting shot.
 * `energy` - An integer, does not necessarily start at the max energy
  possible. Decreases whenever you fire your weapon. You cannot fire if you
  don't have enough energy.
 * `orientation` - The direction you are currently facing on the baord.
 * `grid` - A string, detailing the current state of the board.


The grid will be a string containing something like the following contents:

```
________________________
___W_____WWWWWWWW_______
___W_W__________________
___W_W_______B__________
___W_W__________________
___W_W__________________
_WWWWWWWWW___L____O_____
_____W__________________
_____W_WWWWW____________
_________WWWWWWWW_______
________________________
___________WWWW_________
__X_____________________
________________________
____WWW_________________
________________________
```

Empty cells are `_`, walls are `W`, your tank is `X`, the other tank is `O`,
batteries are `B`, and lasers are `L`.

Once you have computed your next action, you must make an HTTP request with
that action. Your action can be `move`, `left`, `right`, `fire`, or `noop`, and
you should make a `POST` request to
`http://gameserver:8080/game/tankyou/action`, except replacing `tankyou` and
`action` with appropriate values. You should send the `X-Sm-Playerid` header
with each action request.

## Other notes

 * You have a fixed amount of time to make your move. If you take longer than
   the turn timeout, you miss your turn.
 * If you take longer than the connect-back timeout, you forfeit the game.
 * The board wraps around in both directions, like Pac-Man. Laser fire also
   wraps.
 * You can shoot batteries.
 * Colliding lasers nullify each other.
