package interfaces


type IDeck interface {	
	Shuffle(seed uint32)
	GetCard(amount int)	
}
type ICard interface {

}

type IDealer interface {
	GetHand(player IPlayer) [2]ICard

	GetFlop() [3]ICard
	GetTurn() ICard
	GetRiver() ICard
}

type ITable interface {}

type IPlayer interface {}

