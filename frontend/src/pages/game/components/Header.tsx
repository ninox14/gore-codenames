import { useContext } from "react";
import { GameContext } from "../GameContext";
import { Button } from "@/components/ui/button";
import { Link, useParams } from "react-router";

function Header() {
  const { isConnected, close } = useContext(GameContext);

  const { gameId } = useParams();

  return (
    <header className="flex h-12 items-center justify-between px-5">
      <Link className="block text-2xl font-bold" to={"/game"}>
        Game
      </Link>
      {gameId && (
        <>
          <p className="grow text-center">
            Status: {isConnected ? "🟢 Connected" : "🔴 Disconnected"}
          </p>
          <Button className="block cursor-pointer" onClick={close}>
            CLOSE CONNECTION
          </Button>
        </>
      )}
    </header>
  );
}

export default Header;
