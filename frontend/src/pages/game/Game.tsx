import { useNavigate, useParams } from 'react-router';
import {
  useWebSocket,
  type ClientMessage,
  type ServerMessage,
} from '../useWebsocket';
import { useEffect, useState } from 'react';
import { makeWsUrlWithToken } from '@/lib/utils';
import type { GameState } from '@/types';
import { Button } from '@/components/ui/button';

function OnSocketOpen(socket: WebSocket, msg: ClientMessage) {
  socket.send(JSON.stringify(msg));
}

function Game() {
  let { gameId } = useParams();
  const navigate = useNavigate();

  // TODO: handle case when you didnt get game state
  const [gameState, setGameState] = useState<GameState>();

  // FIXME: better gameId parameter validation
  if (!gameId) {
    navigate('/game');
    return null;
  }
  const { lastMessage, isConnected, close } = useWebSocket<
    ClientMessage,
    ServerMessage
  >(makeWsUrlWithToken(gameId), {
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
    <>
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-100 mb-2">Game</h1>
          <p className="text-gray-50">We are in a game</p>
          <p>Status: {isConnected ? '🟢 Connected' : '🔴 Disconnected'}</p>
          <Button className="my-3" onClick={close}>
            CLOSE CONNECTION
          </Button>
          <div className="flex flex-col">
            <div className="">
              <p>Spectators:</p>
              <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
                {gameState?.spectators.map((player) => `${player.name}, `)}
              </pre>
            </div>
            <div className="">
              <p>Team blue:</p>
              <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
                {gameState?.teams.blue.players.map(
                  (player) => `${player.name}, `
                )}
              </pre>
            </div>
            <div className="">
              <p>Team red:</p>
              <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
                {gameState?.teams.red.players.map(
                  (player) => `${player.name}, `
                )}
              </pre>
            </div>

            <p>Board:</p>
            <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
              {JSON.stringify(gameState?.board, null, 2)}
            </pre>
          </div>
        </div>
      </div>
    </>
  );
}

export default Game;
