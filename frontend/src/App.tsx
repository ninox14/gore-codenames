import { Route, Routes } from 'react-router';
import Home from './pages/Home';
import Game from './pages/game';
import GameLayout from './pages/game/GameLayout';
import RootLayout from './pages/RootLayout';

function App() {
  return (
    <Routes>
      <Route path="/" element={<RootLayout />}>
        <Route index element={<Home />} />
        <Route path="game" element={<GameLayout />}>
          <Route index element={<Game />} />
          {/* <Route path=":gameId" element={<Register />} /> */}
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
