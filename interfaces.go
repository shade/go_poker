
interface IDeck {	
	Shuffle(seed uint32)
	GetCard(amount int)	
}

interface IDealer {
	GetHand(player Player) [2]Card

	GetFlop() [3]Card
	GetTurn() Card
	GetRiver() Card
}

interface ITable {
	CheckTable(Player player)
	CallTable(Player player)
	RaiseTable(Player player)
	FoldTable(Player player)
}

interface IPlayer {}

