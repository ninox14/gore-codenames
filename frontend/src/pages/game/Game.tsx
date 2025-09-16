import { useNavigate, useParams } from 'react-router';
import useGameWS from '../useWebsocket';
import { useEffect } from 'react';
import { makeWsUrlWithToken } from '@/lib/utils';
import { Button } from '@/components/ui/button';

function Game() {
  let { gameId } = useParams();
  const navigate = useNavigate();
  const { ws, status } = useGameWS(makeWsUrlWithToken());

  // FIXME: better gameId parameter validation
  if (!gameId) {
    navigate('/game');
    return null;
  }

  useEffect(() => {
    if (!ws) {
      console.log('WEBSOCKET IS NOT CONNECTED', ws);
      return;
    }
  }, []);

  function sendSomething() {
    ws?.send(JSON.stringify({ type: 'bla-bla', data: 'ARBITRARY DATA ' }));
  }

  console.log(status);
  return (
    <>
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-100 mb-2">Game</h1>
          <p className="text-gray-50">We are in a game</p>
          <Button className="mt-3" onClick={sendSomething}>
            SEND SMTHING{' '}
          </Button>
        </div>
      </div>
    </>
  );
}

export default Game;
