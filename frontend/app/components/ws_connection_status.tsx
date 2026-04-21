"use client";
import { Suspense } from "react";
import { app_config } from "../lib/app_config";
import useWebSocket, { WsConnectionStatus } from "../lib/hooks/use_socket";
import {
  errorEventHandler,
  EventCardCreated,
  CardCreatedPayload,
  CardUpdatedPayload,
  EventCardUpdated,
  ReadBoardState,
} from "../lib/ws/events";

export default function WsConnectionStatusIndicator() {
  const { connectionState, errorMsg, sendMessage, closeConnection } =
    useWebSocket({
      url: app_config.PUBLIC_API + "/ws",
      messageListeners: [
        { type: "error", handler: errorEventHandler },
        {
          type: "board.card.create ",
          handler(event) {
            console.log(event);
          },
        },
        {
          type: "board.state",
          handler(event) {
            console.log(event);
          },
        },
      ],
    });

  let bg_color = "bg-black";
  let shadow_color = "shadow-gray-500";
  if (connectionState == WsConnectionStatus.ERROR) {
    bg_color = "bg-red-700";
    shadow_color = "shadow-red-500";

    return;
  }

  if (connectionState == WsConnectionStatus.OPEN) {
    bg_color = "bg-teal-700";
    shadow_color = "shadow-green-500";
  }

  return (
    <Suspense fallback={<p className="text-teal-300"> Loading la </p>}>
      <div className="flex flex-row justify-center gap-x-1 items-bottom px-3 py-2">
        <div
          className={`rounded-full w-5 h-5 ${bg_color} ${shadow_color} shadow-md`}
        >
          {" "}
        </div>
        <p className="font-medium  text-white align-middle gap-x-3">
          WS Server Connection:{" "}
          <span className={`font-bold ${bg_color} rounded-xl px-1`}>
            {connectionState}
          </span>
        </p>
      </div>
      <div className="text-center">
        <button
          className="p-2 py-1 bg-blue-800 rounded-md text-gray-600 text-center cursor-pointer hover:bg-blue-600 hover:text-white active:bg-blue-700"
          onClick={() =>
            sendMessage({
              type: "hello.broadcast",
              payload: { message: "HELLO WORLDz!" },
            })
          }
        >
          Say Hello to Ws
        </button>
        <button
          className="p-2 py-1 rounded-md text-gray-500 text-center cursor-pointer bg-gray-700/50 hover:bg-gray-700 hover:text-white active:bg-gray-600 ml-1"
          onClick={closeConnection}
        >
          Close Connection
        </button>
      </div>
      <div className="text-center">
        <button
          className="p-2 py-1 bg-blue-800 rounded-md text-gray-600 text-center cursor-pointer hover:bg-blue-600 hover:text-white active:bg-blue-700"
          onClick={() => {
            const payload: CardCreatedPayload = {
              column: "1",
              title: "ornektitle",
            };
            const event = new EventCardCreated(payload);
            sendMessage(event);
          }}
        >
          Create Test
        </button>
        <button
          className="p-2 py-1 bg-blue-800 rounded-md text-gray-600 text-center cursor-pointer hover:bg-blue-600 hover:text-white active:bg-blue-700"
          onClick={() => {
            const payload: CardUpdatedPayload = {
              card_id: "1",
              column: "1",
              title: "ornektitle",
            };
            const event = new EventCardUpdated(payload);
            sendMessage(event);
          }}
        >
          Update Test
        </button>
      </div>
      <div className="rounded-md bg-gray-700 text-center max-w-3xl py-5 my-3 flex justify-center self-center place-self-center">
        <h1 className="text-center text-white font-bold italic underline uppercase">
          Cards
        </h1>
        <button
          className="p-2 py-1 bg-blue-800 rounded-md text-gray-600 text-center cursor-pointer hover:bg-blue-600 hover:text-white active:bg-blue-700"
          onClick={() => {
            
            sendMessage(new ReadBoardState([]));
          }}
        >
          Get Cards
        </button>
      </div>
    </Suspense>
  );
}
