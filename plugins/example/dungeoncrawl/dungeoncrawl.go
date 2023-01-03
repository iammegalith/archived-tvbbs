package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

type Module interface {
	Execute() error
}

type DungeonModule struct {
	width       int
	height      int
	numRooms    int
	roomMinSize int
	roomMaxSize int
	mapString   string
	playerPos   [2]int
	enemies     []Enemy
}

type Enemy struct {
	pos      [2]int
	hp       int
	maxHp    int
	attack   int
	defense  int
	dead     bool
	deadTurn int
}

func (m *DungeonModule) generateMap() {
	// Create a random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize the map
	m.mapString = strings.Repeat("#", m.width*m.height)

	// Generate the rooms
	for i := 0; i < m.numRooms; i++ {
		// Choose random width and height for the room
		roomWidth := r.Intn(m.roomMaxSize-m.roomMinSize+1) + m.roomMinSize
		roomHeight := r.Intn(m.roomMaxSize-m.roomMinSize+1) + m.roomMinSize

		// Choose random position for the top-left corner of the room
		roomX := r.Intn(m.width - roomWidth - 1)
		roomY := r.Intn(m.height - roomHeight - 1)

		// Draw the room on the map
		for y := roomY; y < roomY+roomHeight; y++ {
			for x := roomX; x < roomX+roomWidth; x++ {
				m.mapString = m.mapString[:y*m.width+x] + "." + m.mapString[y*m.width+x+1:]
			}
		}
	}

	// Place the player in a random room
	for {
		x := r.Intn(m.width)
		y := r.Intn(m.height)
		if m.mapString[y*m.width+x] == '.' {
			m.mapString = m.mapString[:y*m.width+x] + "P" + m.mapString[y*m.width+x+1:]
			m.playerPos = [2]int{x, y}
			break
		}
	}

	// Place the enemies on the map
	for i := 0; i < 5; i++ {
		for {
			x := r.Intn(m.width)
			y := r.Intn(m.height)
			if m.mapString[y*m.width+x] == '.' {
				m.mapString = m.mapString[:y*m.width+x] + "E" + m.mapString[y*m.width+x+1:]
				m.enemies = append(m.enemies, Enemy{
					pos:      [2]int{x, y},
					hp:       10,
					maxHp:    10,
					attack:   3,
					defense:  1,
					dead:     false,
					deadTurn: 0,
				})
				break
			}
		}
	}
}

func (m *DungeonModule) drawMap() {
	// Clear the screen
	fmt.Print("\033[2J")

	// Draw the map
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			fmt.Print(string(m.mapString[y*m.width+x]))
		}
		fmt.Println()
	}

	// Draw the enemy health bars
	for i, enemy := range m.enemies {
		if !enemy.dead {
			fmt.Printf("Enemy %d: %d/%d\n", i+1, enemy.hp, enemy.maxHp)
		}
	}
}

func (m *DungeonModule) move(dx, dy int) {
	// Check if the new position is within the bounds of the map
	if m.playerPos[0]+dx < 0 || m.playerPos[0]+dx >= m.width || m.playerPos[1]+dy < 0 || m.playerPos[1]+dy >= m.height {
		return
	}

	// Check if the new position is a wall
	if m.mapString[(m.playerPos[1]+dy)*m.width+m.playerPos[0]+dx] == '#' {
		return
	}

	// Update the player's position on the map
	m.mapString = m.mapString[:m.playerPos[1]*m.width+m.playerPos[0]] + "." + m.mapString[m.playerPos[1]*m.width+m.playerPos[0]+1:]
	m.playerPos[0] += dx
	m.playerPos[1] += dy
	m.mapString = m.mapString[:m.playerPos[1]*m.width+m.playerPos[0]] + "P" + m.mapString[m.playerPos[1]*m.width+m.playerPos[0]+1:]
}

func (m *DungeonModule) attack(i int) {
	// Get the enemy and player stats
	enemy := m.enemies[i]
	playerAttack := rand.Intn(5) + 1
	playerDefense := rand.Intn(2)

	// Calculate the damage
	damage := playerAttack - enemy.defense
	if damage < 0 {
		damage = 0
	}
	enemy.hp -= damage

	// Check if the enemy is dead
	if enemy.hp <= 0 {
		enemy.dead = true
		enemy.deadTurn = 1
	}
}

func (m *DungeonModule) enemy
if dx == 0 && dy == 0 {
	continue
}
if m.mapString[(enemy.pos[1]+dy)*m.width+enemy.pos[0]+dx] == '#' {
	continue
}
if m.mapString[(enemy.pos[1]+dy)*m.width+enemy.pos[0]+dx] == 'P' {
	// Player is in range, attack them
	enemyAttack := rand.Intn(5)
	playerHp := m.playerHp - (enemyAttack - playerDefense)
	if playerHp <= 0 {
		// Player is dead
		fmt.Println("You have been killed by an enemy!")
		fmt.Print("Press Enter to continue: ")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	m.playerHp = playerHp
} else {
	// Move the enemy
	m.mapString = m.mapString[:enemy.pos[1]*m.width+enemy.pos[0]] + "." + m.mapString[enemy.pos[1]*m.width+enemy.pos[0]+1:]
	enemy.pos[0] += dx
	enemy.pos[1] += dy
	m.mapString = m.mapString[:enemy.pos[1]*m.width+enemy.pos[0]] + "E
}
}
}
}
}

func (m *DungeonModule) Execute() error {
// Initialize the player's hit points
m.playerHp = 10

// Generate the map
m.generateMap()

// Main game loop
for {
// Draw the map
m.drawMap()

// Get the player's input
fmt.Print("> ")
input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
input = strings.TrimSpace(input)

// Parse the input
if input == "w" {
m.move(0, -1)
} else if input == "a" {
m.move(-1, 0)
} else if input == "s" {
m.move(0, 1)
} else if input == "d" {
m.move(1, 0)
} else if strings.HasPrefix(input, "attack ") {
i, err := strconv.Atoi(strings.TrimPrefix(input, "attack "))
if err != nil || i < 1 || i > len(m.enemies) {
	fmt.Println("Invalid enemy number")
	continue
}
m.attack(i - 1)
} else if input == "exit" {
break
} else {
fmt.Println("Invalid command")
continue
}

// Move the enemies
m.moveEnemies()
}

return nil
}
