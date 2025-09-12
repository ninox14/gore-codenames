import { useContext } from 'react';
import { AuthContext } from '../AuthContext';
import AuthDialog from '@/components/AuthDialog';
import { getLSToken } from '@/lib/utils';

function Game() {
  const { user } = useContext(AuthContext);
  console.log('USER', user);

  const token = getLSToken();
  return (
    <>
      <AuthDialog open={!user && !token} />
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-100 mb-2">Game</h1>
          <p className="text-gray-50">Get started by creating new game</p>
        </div>
      </div>
    </>
  );
}

export default Game;
