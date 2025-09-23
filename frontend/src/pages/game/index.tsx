import { createNewGame } from '@/api';
import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router';
import { toast } from 'sonner';

function GameIndex() {
  const navigate = useNavigate();

  async function handleCreateNewGame() {
    const resp = await createNewGame();
    if (!resp) {
      toast('Failed to create new game');
      return;
    }

    navigate(`/game/${resp.game_id}`);
  }
  return (
    <>
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-100 mb-2">Game</h1>
          <p className="text-gray-50">Get started by creating new game</p>
          <Button onClick={handleCreateNewGame} className="mt-3">
            Create new game
          </Button>
        </div>
      </div>
    </>
  );
}

export default GameIndex;
