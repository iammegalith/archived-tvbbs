type ChatModule struct {
	conn      net.Conn
	clients   []net.Conn
	commands  chan string
	clientIDs map[net.Conn]int
	nextID    int
}

func (m *ChatModule) handleConnection(conn net.Conn) {
	// Read messages from the connection
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		m.commands <- line
	}

	// Remove the client from the list of clients
	delete(m.clientIDs, conn)
	for i, c := range m.clients {
		if c == conn {
			m.clients = append(m.clients[:i], m.clients[i+1:]...)
			break
		}
	}
}

func (m *ChatModule) broadcast(message string) {
	// Send the message to all clients
	for _, conn := range m.clients {
		conn.Write([]byte(message))
	}
}

func (m *ChatModule) Execute() error {
	// Add the existing connection to the list of clients
	m.clients = append(m.clients, m.conn)
	m.clientIDs[m.conn] = m.nextID
	m.nextID++

	// Start listening for incoming messages
	go m.handleConnection(m.conn)

	// Main loop
	for {
		select {
		case line := <-m.commands:
			// Broadcast the message to all clients
			m.broadcast(line)
		}
	}
}
/// To use this modified ChatModule with the existing connection of an ANSI BBS, you can create a new ChatModule instance and pass it the net.Conn object as an argument. For example:
chatModule := &ChatModule{conn: conn}
///This ChatModule instance will use the provided net.Conn object to send and receive data, and will allow other clients to connect and participate in the chat.

/*
To add the previous code to a menus.yaml file, you will need to specify the name of the plugin and the command that should be used to execute it. You can also specify any additional options or parameters that the plugin requires.

Here's an example of how you could add the previous code to a menus.yaml file:
- name: My Plugin
  command: myplugin
  options:
    conn: "{{.Conn}}"

This will create a menu item called "My Plugin" that, when selected, will execute the myplugin command with the conn option set to the current net.Conn object.

You can then use the conn option in the plugin to access the existing connection of the ANSI BBS. For example:

type MyPlugin struct {
	conn net.Conn
	// Other fields and methods
}

func (p *MyPlugin) Execute() error {
	// Use the conn field to send and receive data over the existing connection
	// ...
}
*/
