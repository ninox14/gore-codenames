import AuthDialog from "@/components/AuthDialog";
import { useContext } from "react";
import { Outlet } from "react-router";
import { AuthContext } from "../AuthContext";
import { getLSToken } from "@/lib/utils";
import { GameContextProvider } from "./GameContext";
import Header from "./components/Header";

function GameLayout() {
  const { user } = useContext(AuthContext);
  const token = getLSToken();

  console.log("User", user);
  return (
    <>
      <AuthDialog open={!user && !token} />
      <GameContextProvider user={user}>
        <Header />
        <main className="min-h-[calc(100vh-48px)] bg-black">
          <Outlet />
        </main>
      </GameContextProvider>
    </>
  );
}

export default GameLayout;
