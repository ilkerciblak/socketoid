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
} from "../lib/ws/events";

export default function WsConnectionStatusIndicator() {
  const { connectionState, errorMsg, sendMessage, closeConnection } =
    useWebSocket({
      url: app_config.PUBLIC_API + "/ws",
      messageListeners: [{ type: "error", handler: errorEventHandler }],
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
            const payload: CardUpdatedPayload = {
              title: "arbitrary",
              column: "1",
              card_id: "1",
            };
            const event = new EventCardUpdated(payload);
            sendMessage(event);
            //sendMessage({
            // type: "hello.broadcast",
            //payload: { message: "HELLO WORLDz!" },
            //});
          }}
        >
          Arbitrary Test
        </button>
      </div>
    </Suspense>
  );
}
