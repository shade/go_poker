

export interface IPeer {
    send(data: Uint8Array): void;
}

export interface IPeerMan {
    peers: IPeer[];

    broadcast(data: Uint8Array): void;
}


export enum TableState {
    PREFLOP,
    FLOP,
    TURN,
    RIVER,
    SHOWDOWN
};

export enum HandState {
    CALL,
    CHECK,
    BET,
    FOLD,
    SHOVE
};

export interface Player {
    chips: number;
    state: HandState;

    isDealer(): boolean;
    updateState(state: HandState): void;
}

export interface ITable {
    players: Player[];

    state: TableState;

    bigBlind: number;
    maxSeats: number;

    lastBet: number;
}



