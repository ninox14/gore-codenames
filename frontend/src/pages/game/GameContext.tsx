import type { GameState, UserResponse } from '@/types';

import React, { createContext, useEffect, useState } from 'react';
import { useParams } from 'react-router';
import {
  useWebSocket,
  type ClientMessage,
  type ServerMessage,
} from '../useWebsocket';

type UseWebsockets = ReturnType<
  typeof useWebSocket<ClientMessage, ServerMessage>
>;

type GameContextType = {
  gameState?: GameState;
  send: UseWebsockets['send'];
  isConnected: boolean;
  close: UseWebsockets['close'];
};
export const GameContext = createContext<GameContextType>({
  isConnected: false,
  send: () => {},
  close: () => {},
});

function OnSocketOpen(socket: WebSocket, msg: ClientMessage) {
  socket.send(JSON.stringify(msg));
}

export function GameContextProvider({
  user,
  children,
}: {
  user: UserResponse | null;
  children: React.ReactNode;
}) {
  // TODO: handle case when you didnt get game state
  const [gameState, setGameState] = useState<GameState>();

  let { gameId } = useParams();
  // FIXME: better gameId parameter validation
  const { lastMessage, isConnected, close, send } = useWebSocket<
    ClientMessage,
    ServerMessage
  >({
    user,
    gameId,
    reconnect: false,
    onOpenCallback: (socket) =>
      OnSocketOpen(socket, {
        type: 'join_game',
        data: undefined,
        game_id: gameId,
      }),
  });

  useEffect(() => {
    if (!lastMessage) return;

    switch (lastMessage.type) {
      case 'game_state':
        setGameState(lastMessage.data);
        console.log('Game state update:', lastMessage.data);
        break;
      case 'error':
        console.error('Error:', lastMessage.data);
        break;
    }
  }, [lastMessage]);

  return (
    <GameContext.Provider value={{ gameState, isConnected, send, close }}>
      {children}
    </GameContext.Provider>
  );
}
