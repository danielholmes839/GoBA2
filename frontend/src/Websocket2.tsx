import { useCallback, useEffect, useState } from 'react';

// https://github.com/itays123/partydeck/blob/main/game/src/game/websocketUtils.ts

export type ConnectFN = (
  url: string,
) => void;

export interface Contextable {
  context: string;
}

export type SessionConnectHandler = (ev: Event) => any;
export type SessionMessageHanlder = (ev: MessageEvent<any>) => any;
export type SessionDisconnectHandler = (ev: Event) => any;

export type SessionHooks = {
  connect: (url: string) => void,
  disconnect: () => void,
  send: <T extends Contextable>(args: T) => void,
}

export type SessionHandlers = {
  open: SessionConnectHandler,
  message: SessionMessageHanlder,
  close: SessionDisconnectHandler
}

export function useSession({ open, close, message }: SessionHandlers): SessionHooks {
  const [session, setSession] = useState(null as unknown as WebSocket);

  const updateOpenHandler = () => {
    if (!session) return;
    session.addEventListener('open', open);
    return () => {
      session.removeEventListener('open', open);
    };
  };

  const updateMessageHandler = () => {
    if (!session) return;
    session.addEventListener('message', message);
    return () => {
      session.removeEventListener('message', message);
    };
  };

  const updateCloseHandler = () => {
    if (!session) return;
    session.addEventListener('close', close);
    return () => {
      session.removeEventListener('close', close);
    };
  };

  useEffect(updateOpenHandler, [session, open]);
  useEffect(updateMessageHandler, [session, message]);
  useEffect(updateCloseHandler, [session, close]);

  const connect = useCallback(
    (url: string) => {
      const ws = new WebSocket(url);
      setSession(ws);
    },
    []
  );

  const send = <T extends Contextable>(args: T) => {
    session.send(JSON.stringify(args));
  };

  const disconnect = useCallback(() => {
    if (session.readyState === session.OPEN) session.close(1000);
  }, [session]);

  return {
    connect: connect,
    disconnect: disconnect,
    send: send
  };
}