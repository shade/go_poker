import {
    ITable,
    Player,
    TableState
} from './types';


export class Table implements ITable {
    players: Player[];

    state: TableState = TableState.PREFLOP;

    bigBlind: number;
    maxSeats: number;
    lastBet: number = 0;

    constructor(
        maxSeats: number,
        bigBlind: number) {
            this.bigBlind = bigBlind;
            this.maxSeats = maxSeats;
    }
}