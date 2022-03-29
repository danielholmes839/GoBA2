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
  write: <T extends Contextable>(args: T) => void,
}

export function useSession(
  onOpen: SessionConnectHandler,
  onMessage: SessionMessageHanlder,
  onClose: SessionDisconnectHandler
): SessionHooks {
  const [session, setSession] = useState(null as unknown as WebSocket);

  const updateOpenHandler = () => {
    if (!session) return;
    session.addEventListener('open', onOpen);
    return () => {
      session.removeEventListener('open', onOpen);
    };
  };

  const updateMessageHandler = () => {
    if (!session) return;
    session.addEventListener('message', onMessage);
    return () => {
      session.removeEventListener('message', onMessage);
    };
  };

  const updateCloseHandler = () => {
    if (!session) return;
    session.addEventListener('close', onClose);
    return () => {
      session.removeEventListener('close', onClose);
    };
  };

  useEffect(updateOpenHandler, [session, onOpen]);
  useEffect(updateMessageHandler, [session, onMessage]);
  useEffect(updateCloseHandler, [session, onClose]);

  const connect = useCallback(
    (url: string) => {
      const ws = new WebSocket(url);
      setSession(ws);
    },
    []
  );

  const sendMessage = <T extends Contextable>(args: T) => {
    session.send(JSON.stringify(args));
  };

  const close = useCallback(() => {
    if (session.readyState === session.OPEN) session.close(1000);
  }, [session]);

  return {
    connect: connect,
    disconnect: close,
    write: sendMessage
  };
}