"use client";

import { Suspense } from "react";
import { app_config } from "../lib/app_config";
import useSSE, { ConnectionStatus } from "../lib/hooks/use_sse";

export default function SSEStatusIndicator() {
  const {
    connectionState,
    error: errorMsg,
    eventQueue,
  } = useSSE("http://localhost:8081" + "/events", {
    onMessage(event) {
      console.log(event.data);
    },
  });
  let bg_color = "bg-black"
  let shadow_color = "shadow-gray-500"
  if (connectionState == ConnectionStatus.ERROR) {
    console.log("patladiksss");
    console.log(errorMsg);
    bg_color = "bg-red-700"
    shadow_color = "shadow-red-500"
    return;
  }

  if (connectionState == ConnectionStatus.OPEN) {
    bg_color = "bg-teal-700"
    shadow_color= "shadow-green-500"
  }

    return (
      <Suspense fallback={<p className="text-teal-300"> Loading la </p>}>
        <div className="flex flex-row justify-center gap-x-1 items-bottom px-3 py-2">
          <div className={`rounded-full w-5 h-5 ${bg_color} ${shadow_color} shadow-md`}>  </div>
          <p className="font-medium  text-white align-middle gap-x-3">SSE Server Connection: <span className={`font-bold ${bg_color} rounded-xl px-1`}>{connectionState}</span></p>

      </div>
    </Suspense >
  ) 
  }
