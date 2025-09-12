import { Outlet } from 'react-router';
import { AuthContextProvider } from '../AuthContext';

function GameLayout() {
  return (
    <div className="min-h-screen bg-black py-12 px-4 sm:px-6 lg:px-8">
      <AuthContextProvider>
        <Outlet />
      </AuthContextProvider>
    </div>
  );
}

export default GameLayout;
