### Exercise A: Peer-to-Peer Chat System

In this exercise, you will implement a simple peer-to-peer network that functions as a distributed chat room. The key features include:

1. **Command Line Interface:** The chat client runs as a command-line program.
2. **Peer Connection:** Upon startup, the client requests the IP address and port of an existing peer. If the peer is not found, it starts its own network.
3. **Client Details:** The client prints its IP address and port for incoming connections.
4. **Messaging:** The client prompts the user for text strings. When a string is typed, it is broadcast to all connected clients.
5. **Message Handling:** The system ensures that messages are printed at all clients without revealing the sender's information.

### Exercise B: Distributed Ledger

This exercise extends the previous chat system into a distributed ledger with the following features:

1. **Ledger Management:** Each client maintains a ledger that tracks accounts and balances.
2. **Transactions:** Clients can make transactions, and all peers update their ledgers accordingly.
3. **Eventual Consistency:** The system ensures that if all clients stop sending transactions, all ledgers will eventually converge to the same state.
4. **Peer Connections:** Clients connect to a list of peers and forward transactions to maintain ledger consistency.

### Exercise C: Signed Transactions

Enhance the ledger system with signed transactions:

1. **Signed Transactions:** Transactions are now signed with RSA keys, ensuring only account owners can authorize transactions.
2. **Validation:** Transactions are validated based on their signature and non-negative amounts.
3. **RSA Integration:** Encode and decode RSA public keys for transaction processing.

### Exercise D: Total Order with Sequencer

Modify the system to achieve total order for transactions:

1. **Sequencer Role:** A designated sequencer establishes the order of transactions.
2. **Block Creation:** The sequencer creates blocks containing transaction IDs and signs them.
3. **Order Enforcement:** Clients process transactions based on the sequencer's block order.

### Exercise E: Proof-of-Stake Blockchain

Implement a proof-of-stake based blockchain with total order:

1. **Genesis Block:** Includes special public keys with initial balances.
2. **Proof-of-Stake:** A tree-based blockchain where block creation is based on the stake (account balance).
3. **Transaction Fees:** Implement transaction fees and rewards for block creators.
4. **Testing:** Validate the system's performance, including handling invalid transactions, rollbacks, and resistance to malicious peers.
