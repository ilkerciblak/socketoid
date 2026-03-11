"use client";

import { Suspense, useState } from "react";
import { app_config } from "../lib/app_config";
import useSSE, { ConnectionStatus } from "../lib/hooks/use_sse";

export default function SSEStatusIndicator() {
  const [userList, setUserList] = useState<string[]>([]);
  const { connectionState, error: errorMsg } = useSSE(
    app_config.PUBLIC_API + "/events",
    {
      handlers: [
        {
          type: "userjoined",
          handler(event) {
            const data = JSON.parse(event.data);
            setUserList((prev) => [data["user-id"], ...prev]);
          },
        },
      ],
      auto: true,
    },
  );

  let bg_color = "bg-black";
  let shadow_color = "shadow-gray-500";
  if (connectionState == ConnectionStatus.ERROR) {
    bg_color = "bg-red-700";
    shadow_color = "shadow-red-500";
    return;
  }

  if (connectionState == ConnectionStatus.OPEN) {
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
          SSE Server Connection:{" "}
          <span className={`font-bold ${bg_color} rounded-xl px-1`}>
            {connectionState}
          </span>
        </p>
      </div>
      <ul className="flex flex-col justify-center items-center">
        {userList.map((u) => (
          <li key={u} className="italic underline text-gray-200 self-center">
            user {u} joined to chat 🦆
          </li>
        ))}
      </ul>
    </Suspense>
  );
}
