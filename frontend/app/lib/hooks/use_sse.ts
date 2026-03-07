import { useCallback, useEffect, useRef, useState } from "react";

export const enum ConnectionStatus {
  CONNECTING = "connecting",
  OPEN = "open",
  CLOSED = "closed",
  ERROR = "error",
}

export default function useSSE(
  url: string,
  options?: { onMessage?: (event: MessageEvent) => void },
) {
  const [connectionState, setConnectionState] = useState<ConnectionStatus>(
    ConnectionStatus.CONNECTING,
  );
  const [error, setErrorMsg] = useState<string | null>(null);

  const [eventQueue, setEventQueue] = useState<string[]>([]);
  const eventSourceRef = useRef<EventSource | null>(null);

  // const maxReconnectAttemps = 5;

  const connect = useCallback(
    function connect() {
      if (eventSourceRef.current) eventSourceRef.current.close();
      const eventSource = new EventSource(url, { withCredentials: false });
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        setConnectionState(ConnectionStatus.OPEN);
        console.log("open oldu")
        setErrorMsg(null);
      };

      eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (options?.onMessage) {
            options.onMessage?.(data);
          }
          setEventQueue((prev) => [...prev, data]);
          console.log("event", data)
          console.log("que", eventQueue)
        } catch (error) {
          console.log("failed to parse event message:\n", error);
          setErrorMsg(`${error}`);
        }
      };

      eventSource.onerror = (error) => {
        error.preventDefault();
        setConnectionState(ConnectionStatus.ERROR);
        setErrorMsg(`Failed to establish connection:\n ${error}`);
        eventSource.close();
      };
    },
    [url, options],
  );

  useEffect(() => {
    connect();
    console.log("hello there");
    return () => {
      eventSourceRef.current?.close();
    };
  }, [connect, url]);

  return { connectionState, error, eventQueue, connect };
}
