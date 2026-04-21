export interface BoardEvent {
  type: string;
  payload: unknown;
}

export type CardCreatedPayload = {
  title: string;
  column: string;
};



export class EventCardCreated implements BoardEvent {
  type: string;
  payload: CardCreatedPayload;

  constructor(payload: CardCreatedPayload) {
    this.type = "board.card.created";
    this.payload = payload;
  }
}

export type CardDeletedPayload = {
  card_id: string;
};
export class EventCardDeleted implements BoardEvent {
  type: string;
  payload: CardDeletedPayload;

  constructor(payload: CardDeletedPayload) {
    this.type = "board.card.deleted";
    this.payload = payload;
  }
}

export type CardUpdatedPayload = {
  title: string;
  column: string;
  card_id: string;
};
export class EventCardUpdated implements BoardEvent {
  type: string;
  payload: CardUpdatedPayload;

  constructor(payload: CardUpdatedPayload) {
    this.type = "board.card.updated";
    this.payload = payload;
  }
}

export type CardMovedEventPayload = {
  card_id: string;
  column: string;
};

export class EventCardMoved implements BoardEvent {
  type: string;
  payload: unknown;

  constructor(payload: CardMovedEventPayload) {
    this.type = "board.card.moved";
    this.payload = payload;
  }
}

export class ReadBoardState implements BoardEvent{
    type: string;
    payload:  CardUpdatedPayload[];

    constructor(payload: CardUpdatedPayload[]){
      this.type="board.state";
      this.payload = payload;
    }

}

export const errorEventHandler = (event: MessageEvent): void | Error => {
  const parsed = JSON.parse(event.data);

  if (parsed.type == "error") {
    console.log(JSON.parse(event.data));
    return;
  }

  return new Error(`mis-type used in error-handler: ${event.type}`);
};
