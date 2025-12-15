Exercise A: Implement a toy peer-to-peer network
This first exercise asks you to program in a toy example of a peer-to-peer flooding network for sending strings around. The peer-to-peer network should then be used to build a distributed chat room. The chat room client should work as follows:

It runs as a command line program.
When it starts up it asks for the IP address and port number of an existing peer on the network. If the IP address or port is invalid or no peer is found at the address, the client starts its own new network with only itself as member.
Then the client prints its own IP address and the port on which it waits for connections.
Then it will iteratively prompt the user for text strings.
When the user types a text string at any connected client, then it will eventually be printed at all other clients.
Only the text string should be printed, no information about who sent it.
The system should be implemented as follows:

When a client connects to an existing peer, it will keep a TCP connection to that peer.
Then the client opens its own port where it waits for incoming TCP connections.
All the connections will be treated the same, they will be used for both sending and receiving strings.
It keeps a set of messages that it already sent.
When a string is typed by the user or a string arrives on any of its connections, the client checks if it is already sent. If so, it does nothing. Otherwise it adds it to MessagesSent and then sends it on all its connections. (Remember concurrency control. Probably several routines will access the set at the same time. Make sure that does not give problems.)
Ponder the following questions:

Does your system have eventual consistency in the sense that if all clients stop typing, then eventually all clients will print the same set of strings?
Exercise B: Implement a Simple Peer-to-Peer Ledger
Modify your code from Exercise A to add the following features:

The system now no longer broadcasts strings and prints them. Instead it implements a distributed ledger. Each client keeps a Ledger that keeps track of Accounts and their balancers.
Each client can make Transactions. When they do all other peers eventually update their ledger with the transaction.
The system should ensure eventual consistency, i.e., if all clients stop sending transactions, then all ledgers will eventually be in the same correct state.
Your system only has to work in a setting with two phases: first all the peers connect, then after a sufficiently long break they start making transactions. However, if you want to accommodate for late comers a way to do it is to let each client keep a list of all the transactions it saw and then forward them to clients that join the system later. You can assume that peers join the system one at a time with sufficient time between them. This way you do not have to worry what happens if peers hears about join events in different order.
Implement along these lines:

Keep a list of peers in the order in which their joined the network, with the latest peer to arrive being at the end.
When connecting to a peer, ask for its list of peers.
Then add yourself to the end of your own list
Then connect to the ten peers before you on the list. If the list has length less than 11 then just connect to all peers but yourself.
Then broadcast your own presence.
When a new presence is broadcast, add it to the end of your list of peers.
When a transaction is made, broadcast the Transaction object.
When a transaction is received, update the local Ledger object.
Ponder the following questions:

Discuss whether connection to the previous ten peers is a good strategy with respect to connectivity. In particular, if the network has 1000 peers, how many connections need to break to partition the network?
Argue that your system has eventual consistency if all parties are correct and the system is run in two-phase mode.
Assume we made the following change to the system: When a transaction arrives, it is rejected if the sending account goes below 0. Does your system still have eventual consistency? Why or why not?
Exercise C: Implement a Simple Peer-to-Peer Ledger
Modify your code from Exercise B to add the following features:

The system still keeps a Ledger.
Each client can make SignedTransactions, i.e., what is broadcast is now objects of the type SignedTransaction.
The sender and receive of a transaction are now RSA public keys encoded as strings. The client can only make a transaction if it knows the secret key corresponding to the sending account. This ensure that only the owner of the account can take money from the account. In a bit more detail, you have to find a way to encode and decode RSA public keys into the string type. If we call the encoding of pk by the name enc(pk), then the amount that "belongs" to pk is Accounts[enc(pk)]. To transfer money from pk one makes a SignedTransaction where pk is encoded and put in the From-field. An encoding of the RSA public key to receive the amount is placed in the To-field. All the fields (save Signature) are then signed under pk (using the corresponding secret key) and the signature is placed in the Signature-field. A SignedTransaction is valid if the signature is valid. Only valid transactions are executed. The invalid transactions are simply ignored
Implement as in Exercise B with these additions:

When a transaction is made, broadcast the SignedTransaction object.
When a transaction is received, update the local Ledger object if the SignedTransaction is has a valid signature and the amount is non-negative
You do not have to:

Handle overdraft, i.e., we allow that accounts become negative.
Protection against cheating parties (neither Byzantine errors nor crash errors).
Exercise D: Total Order by Sequencer
Start from your solution in Exercise C and make it into a system with using the following idea:

Your system runs in two phases. In phase 1 the peers connect to the network. In phase 2 they can send signed transactions.
The peer that started the network is a designated sequencer.
The sequencer creates a special RSA key pair called the sequencer key pair.
When connecting to a network the new client is informed who is the sequencer.
It is the order in which the sequencer received the transactions that counts. This is communicated to the other peers as follows: Every 10 seconds the sequencer will take the transactions that it saw, but which have so far not been sequenced. Then it puts the IDs of those transactions into a block. A block has a block number and an ordered list of IDs, string[]. It numbers the blocks 0, 1, ... in the order they are sent. The sequencer signs the block and sends the block on the network.
A client will accept a block if and only if it has the next block number it has not seen yet and the block is signed by the sequencer.
All clients process the transactions they receive in the order chosen by the sequencer.
A transaction is ignored if it would make the sending account negative.
Exercise E: Static Proof-of-Stake
Start from your solution in Exercise C and use parts of your code from Exercise D. The code in Exercise C should already be a distributed ledger with authenticated transactions. However, it does not have total order. Change it such that it gives a total order of all transactions and rejects transactions that would bring an account into minus. Do it by adding a proof-of-stake, tree-based, totally-ordered broadcast. Implement total-order using a tree-based blockchain based on proof-of-stake. In a bit more detail, implement it as follows:

The initial seed Seed is picked by you and hardcoded into the genesis block.
Transactions are conducted in the unit DKK.
The genesis block contains ten special public keys which by definition have 106 DKKs on them. All other accounts have 0 DKK on them initially. The special accounts are generated by you, and you know the secret keys.
Transactions are in integral DKKs.
A transaction must send at least 1 DKK to be valid.
A block can contain any number of transactions, and might even contain as little as no transactions if there are none to add to the block.
SlotLength is 1 second. You might set it larger if your signatures are very slow to compute. Recall that you need to compute one signature per slot.
To take part in the lottery and making blocks you need an account in the ledger with a positive balance. Your number of tickets is the balance of the account. Throughout the system your number of tickets is the balance in the genesis block, so only the ten accounts you created can be part of running the system.
The signature keys used in the lottery is the same as those used in the ledger system.
Make the system run with 10 peers.
Set the hardness such that your system creates a new block about every 10 seconds. If this is too often for you system to grow a longest chain, then make it longer.
A block is not added to the tree unless all transactions are correctly signed and valid (they make no account go below 0 at any point).
When a transaction is made, the receiver gets 1 DKK less than what was sent. This is a transaction fee.
When a new block is made, then the account of the block creator gets 10 DKK plus one DKK for each transaction in the block.
Ponder the following questions:

When the system is not under attack, how many transactions per second can the sys- tem handle. If you compute this number of different values of BlockSize, which value of BlockSize is the best for throughout? (A transaction is not counted as done until it has been ordered and the balance of the accounts have been updated with that transaction.)
Test the following:

You should test your code. You should do the same testing as in Exercise B and Exercise C. In addition you should now send transactions (for instance 25% of them) which are invalid, for instance negative amount, 0 amount, invalid signature, or which would bring the sending account into overdraft. Test that they are rejected.
During test, try to set your block time so low that you provoke rollbacks now and then to make sure your system can tolerate this. If you simulate your network, you might have to insert a simulated network delay to see this.
Also do transfers to and from accounts not being the initial ones.
Try to run with some fraction 
ϕ
 of bad peers trying to destroy the system. For instance, always let them build on the second longest chain and the longest side branch (a branch in the tree such that the distance up to hitting the longest chain is maximal). Note that a corrupt peer can use a single winning ticket to extend both. Run the system for some fixed amount of time, like 10 or 30 minutes. Keep track of the longest rollback that an honest peer saw. Try to do it for different values of corruption. For instance, how does it evolve for 
ϕ
 = 10%, 25%, ... , 50%, 60%. Is this the behaviour we expected to see?
