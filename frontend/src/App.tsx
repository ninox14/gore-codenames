import { Route, Routes } from 'react-router';
import Home from './pages/Home';
import GameIndex from './pages/game';
import GameLayout from './pages/game/GameLayout';
import RootLayout from './pages/RootLayout';
import Game from './pages/game/Game';
import { AuthContextProvider } from './pages/AuthContext';

function App() {
  return (
    <Routes>
      <Route path="/" element={<RootLayout />}>
        <Route index element={<Home />} />
        <Route
          path="game"
          element={
            <AuthContextProvider>
              <GameLayout />
            </AuthContextProvider>
          }
        >
          <Route index element={<GameIndex />} />
          <Route path=":gameId" element={<Game />} />
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
