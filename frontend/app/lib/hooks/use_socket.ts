import { useCallback, useEffect, useRef, useState } from "react";

export const enum WsConnectionStatus {
  IDLE = "idle",
  CONNECTING = "connecting",
  OPEN = "open",
  CLOSED = "closed",
  ERROR = "error",
}

type onMessage = {
  type: string;
  handler: (event: MessageEvent) => void;
};

type onEvent = (event: Event) => void;

export default function useWebSocket(g: {
  url: string;
  messageListeners: onMessage[];
  onOpen?: onEvent;
  onDisconnect?: onEvent;
  onError?: onEvent;
}) {
  const [connectionState, setConnectionState] = useState<WsConnectionStatus>(
    WsConnectionStatus.IDLE,
  );

  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const webSocketRef = useRef<WebSocket | null>(null);
  const listeners = useRef<onMessage[]>(g.messageListeners);

  const connect = useCallback(() => {
    if (webSocketRef) webSocketRef.current = null;

    try {
      const ws = new WebSocket(g.url);
      setConnectionState(WsConnectionStatus.CONNECTING);
      webSocketRef.current = ws;
      webSocketRef.current.onopen = (event) => {
        setConnectionState(WsConnectionStatus.OPEN);
        if (g.onOpen) g.onOpen(event);
      };

      webSocketRef.current.onmessage = (event) => {
        try {
          listeners.current.forEach((listener) => listener.handler(event));
        } catch (e) {
          setErrorMsg(`${e}`);
        }
      };

      webSocketRef.current.onclose = (event) => {
        setConnectionState(WsConnectionStatus.CLOSED);
        if (g.onDisconnect) g.onDisconnect(event);
      };

      webSocketRef.current.onerror = (event) => {
        setConnectionState(WsConnectionStatus.ERROR);
        if (g.onError) g.onError(event);
      };
    } catch (error) {
      console.log(`${error}`);
    }
  }, [g.url]);

  const sendMessage = useCallback(
    (data: { type: string; payload: unknown }) => {
      try {
        if (
          webSocketRef.current == null ||
          (webSocketRef.current &&
            webSocketRef.current.readyState != webSocketRef.current.OPEN)
        ) {
          throw new Error(
            `connection is not open: ${webSocketRef.current?.readyState}`,
          );
        }
        const serialized = JSON.stringify(data);
        webSocketRef.current!.send(serialized);
      } catch (error) {
        setErrorMsg("serialization error");
        console.log(`serialization error ${error}`);
      }
    },
    [webSocketRef],
  );

  const closeConnection = useCallback(() => {
    webSocketRef.current?.close();
  }, [webSocketRef]);

  useEffect(() => {
    connect();

    return () => {
      if (webSocketRef?.current?.readyState != WebSocket.CLOSED) {
        webSocketRef?.current?.close();
      }
    };
  }, [g.url]);

  return { connectionState, errorMsg, sendMessage, closeConnection };
}
