import React, { useEffect, useState } from "react"
import { useAuth, AuthProvider } from "./Auth";
import { useSession } from "./Websocket2";

const Avatar: React.FC = () => {
  const { identity, authenticated, login, logout } = useAuth();
  if (!authenticated || identity === null) {
    return <></>
  }

  return <span>
    <img className="rounded-full inline" style={{ background: identity.color, height: "100%", border: "0.5px solid", borderColor: identity.color }} src={`https://cdn.discordapp.com/avatars/${identity.user_id}/${identity.avatar_id}.png?size=32`} />
  </span>
}

const Main: React.FC = () => {
  const { identity, authenticated, login, logout } = useAuth();
  const { connect, disconnect, send } = useSession({
    open: (e) => {
      console.log("open", e)
    },
    close: (e) => {
      console.log("close", e)
    },
    message: (e) => {
      console.log("message", e)
    }
  })

  const [counter, setCounter] = useState<any>({})
  return (
    <div className="bg-slate-800 text-white h-100">
      <div className="sm:container mx-auto p-10">
        <div>
          <span className="mr-3"><Avatar /></span>
          <button className="bg-blue-500 hover:bg-blue-600 text-sm text-white font-mono font-semibold px-3 py-1 mr-3 rounded-full" onClick={logout}>Logout</button>
          <button className="bg-blue-500 hover:bg-blue-600 text-sm text-white font-mono font-semibold px-3 py-1 mr-3 rounded-full" onClick={login}>Login</button>
          <button className="bg-blue-500 hover:bg-blue-600 text-sm text-white font-mono font-semibold px-3 py-1 mr-3 rounded-full" onClick={() => connect("ws://localhost:3000/connect")}>connect</button>
          <button className="bg-red-500 hover:bg-red-600 text-sm text-white font-mono font-semibold px-3 py-1 mr-3 rounded-full" onClick={disconnect}>disconnect</button>
        </div>

        <pre>{JSON.stringify(counter, null, 4)}</pre>
        <pre>{JSON.stringify(identity, null, 4)}</pre>

      </div>
    </div>
  );
}
const App: React.FC = () => {
  return <AuthProvider><Main /></AuthProvider>
}

export default App;
