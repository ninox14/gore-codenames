import AuthDialog from '@/components/AuthDialog';
import { useContext } from 'react';
import { Outlet } from 'react-router';
import { AuthContext } from '../AuthContext';
import { getLSToken } from '@/lib/utils';
import { GameContextProvider } from './GameContext';
import Header from './Header';

function GameLayout() {
  const { user } = useContext(AuthContext);
  const token = getLSToken();

  console.log('User', user);
  return (
    <>
      <AuthDialog open={!user && !token} />
      <GameContextProvider user={user}>
        <Header />
        <main className="min-h-screen bg-black py-12 px-4 sm:px-6 lg:px-8">
          <Outlet />
        </main>
      </GameContextProvider>
    </>
  );
}

export default GameLayout;
