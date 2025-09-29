import { useNavigate, useParams } from 'react-router';
import { type ClientEventMap } from '../useWebsocket';
import { Button } from '@/components/ui/button';
import { useContext } from 'react';
import { GameContext } from './GameContext';

function Game() {
  let { gameId } = useParams();
  const navigate = useNavigate();

  const { send, gameState } = useContext(GameContext);
  // FIXME: better gameId parameter validation
  if (!gameId) {
    navigate('/game');
    return null;
  }

  function handleMovePlayer(data: ClientEventMap['change_team']) {
    send({ type: 'change_team', data, game_id: gameId });
  }

  return (
    <div className="max-w-md mx-auto space-y-8">
      <div className="flex flex-col">
        <div className="flex flex-col">
          <div className="flex space-x-2">
            <p>Spectators</p>
            <Button
              onClick={() => handleMovePlayer({ destination: 'spectators' })}
            >
              Join
            </Button>
          </div>
          <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
            {gameState?.spectators.map((player) => `${player.name}, `)}
          </pre>
        </div>
        <div className="flex flex-col">
          <div className="flex space-x-2">
            <p>Team blue</p>
            <Button
              onClick={() =>
                handleMovePlayer({ destination: 'teams.blue.players' })
              }
            >
              Join
            </Button>
          </div>
          <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
            {gameState?.teams.blue.players.map((player) => `${player.name}, `)}
          </pre>
        </div>
        <div className="flex flex-col">
          <div className="flex space-x-2">
            <p>Team red</p>
            <Button
              onClick={() =>
                handleMovePlayer({ destination: 'teams.red.players' })
              }
            >
              Join
            </Button>
          </div>
          <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
            {gameState?.teams.red.players.map((player) => `${player.name}, `)}
          </pre>
        </div>

        <p>Board:</p>
        <pre className="rounded-md border border-blue-900 bg-slate-950 p-3">
          {JSON.stringify(gameState?.board, null, 2)}
        </pre>
      </div>
    </div>
  );
}

export default Game;
