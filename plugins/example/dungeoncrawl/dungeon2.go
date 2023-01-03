type Leaderboard struct {
	Scores []Score
}

type Score struct {
	Name  string
	Score int
}

func (lb *Leaderboard) Load(filename string) error {
	// Read the leaderboard data from the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the data into the Scores field
	err = json.Unmarshal(data, &lb.Scores)
	if err != nil {
		return err
	}

	return nil
}

func (lb *Leaderboard) Save(filename string) error {
	// Marshal the Scores field into JSON data
	data, err := json.Marshal(lb.Scores)
	if err != nil {
		return err
	}

	// Write the data to the file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (lb *Leaderboard) AddScore(name string, score int) {
	// Add the new score to the leaderboard
	lb.Scores = append(lb.Scores, Score{Name: name, Score: score})

	// Sort the scores in descending order
	sort.Slice(lb.Scores, func(i, j int) bool {
		return lb.Scores[i].Score > lb.Scores[j].Score
	})

	// Truncate the leaderboard if it has more than 10 entries
	if len(lb.Scores) > 10 {
		lb.Scores = lb.Scores[:10]
	}
}

// // Modify the Execute:
func (m *DungeonModule) Execute() error {
	// Load the leaderboard
	err := m.leaderboard.Load("leaderboard.json")
	if err != nil {
		return err
	}

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

		// Check if the player is dead
		if m.playerHp <= 0 {
			// Prompt the player for their name
			fmt.Print("Enter your name: ")
			name, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			name = strings.TrimSpace(name)

			// Add the player's score to the leaderboard
			m.leaderboard.AddScore(name, m.turn)

			// Display the leaderboard
			fmt.Println("Leaderboard:")
			for i, score := range m.leaderboard.Scores {
				fmt.Printf("%d. %s: %d\n", i+1, score.Name, score.Score)
			}
			fmt.Print("Press Enter to continue: ")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			break
		}
	}

	// Save the leaderboard
	err = m.leaderboard.Save("leaderboard.json")
	if err != nil {
		return err
	}

	return nil
}

