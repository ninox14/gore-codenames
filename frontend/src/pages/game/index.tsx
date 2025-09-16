import { Button } from '@/components/ui/button';
import { Link } from 'react-router';

function GameIndex() {
  // FIXME: handle game creation on server
  const uuid = self.crypto.randomUUID();
  return (
    <>
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-100 mb-2">Game</h1>
          <p className="text-gray-50">Get started by creating new game</p>
          <Button asChild>
            <Link className="mt-3" to={`/game/${uuid}`}>
              Create new game
            </Link>
          </Button>
        </div>
      </div>
    </>
  );
}

export default GameIndex;
