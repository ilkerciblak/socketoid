import { useCallback, useEffect, useRef, useState } from "react";

export const enum ConnectionStatus {
  CONNECTING = "connecting",
  OPEN = "open",
  CLOSED = "closed",
  ERROR = "error",
}

type onMessage = { type: string; handler: (event: MessageEvent) => void };

export default function useSSE(
  url: string,
  options?: {
    handlers: onMessage[];
  },
) {
  const [connectionState, setConnectionState] = useState<ConnectionStatus>(
    ConnectionStatus.CONNECTING,
  );
  const [error, setErrorMsg] = useState<string | null>(null);

  const [eventQueue, setEventQueue] = useState<string[]>([]);
  const eventSourceRef = useRef<EventSource | null>(null);
  const optionsRef = useRef(options);

  // const maxReconnectAttemps = 5;

  const connect = useCallback(
    function connect() {
      if (eventSourceRef.current) eventSourceRef.current.close();
      const eventSource = new EventSource(url, { withCredentials: false });
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        setConnectionState(ConnectionStatus.OPEN);
        setErrorMsg(null);
      };

      if (optionsRef?.current?.handlers.length) {
        optionsRef?.current?.handlers.forEach((handler) => {
          eventSource.addEventListener(handler.type, (e) => {
            try {
              handler.handler(e);
            } catch (error) {
              setErrorMsg(`failed to process message: ${error}`);
            }
          });
        });
      } else {
        eventSource.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            setEventQueue((prev) => [...prev, data]);
          } catch (error) {
            setErrorMsg(`${error}`);
          }
        };
      }

      eventSource.onerror = (error) => {
        error.preventDefault();
        setConnectionState(ConnectionStatus.ERROR);
        setErrorMsg(`Failed to establish connection:\n ${error}`);
      };
    },
    [url],
  );

  useEffect(() => {
    connect();
    return () => {
      eventSourceRef.current?.close();
    };
  }, [connect, url]);

  return { connectionState, error, eventQueue, connect };
}
