import React from "react"
import { useAuth, AuthProvider } from "./Auth";
import useWebSocket, { ReadyState } from 'react-use-websocket';

const someText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam imperdiet neque nec feugiat blandit. Mauris in est volutpat, ultricies felis a, rutrum erat. Aliquam non augue"

const Container: React.FC = ({ children }) => {
  return (
    <div>
      <div className="sm:container mx-auto p-10">
        {children}
      </div>
    </div>
  )
}
const Avatar: React.FC = () => {
  const { identity, authenticated, login, logout } = useAuth();
  if (!authenticated || identity === null) {
    return <></>
  }

  return <span>
    <img className="rounded-full inline" style={{ background: identity.color, height: "32px", border: "1px solid", borderColor: "#bbbbbb" }} src={`https://cdn.discordapp.com/avatars/${identity.user_id}/${identity.avatar_id}.png?size=512`} />
  </span>
}

const Game: React.FC = () => {
  const { readyState, lastJsonMessage } = useWebSocket("ws://localhost:3000/connect", {
    onClose: (e) => {
      console.log(e)
    },
    onError: (e) => {
      console.log(e)
    },
    onMessage: (e) => {
      console.log(e)
    },
    onOpen: (e) => {
      console.log(e)
    }
  })

  if (readyState === ReadyState.CONNECTING) {
    return <>connecting</>
  }

  if (readyState === ReadyState.CLOSED) {
    return <>disconnected</>
  }

  return <>
    <pre>{JSON.stringify(lastJsonMessage, null, 4)}</pre>
  </>
}


const LoggedOut: React.FC = () => {
  const { identity, login } = useAuth();
  return <>
    <button className="bg-blue-500 hover:bg-blue-600 text-sm text-white font-mono font-semibold px-3 py-1 mr-3 rounded-full" onClick={login}>Login</button>
  </>
}
const LoggedIn: React.FC = () => {
  const { identity, logout } = useAuth();
  const games = [
    {
      id: "connect4",
      title: "Connect 4",
      description: someText
    },
    {
      id: "goba2",
      title: "GoBA 2",
      description: someText
    }
  ]
  return (
    <Container>
      <div><Avatar /><h1 className="inline ml-3 text-xl font-mono">{identity?.user_name}</h1></div> 
      {/* <button className="bg-gray-100 hover:bg-gray-200 text-sm text-gray-900 font-mono px-3 py-1 mr-3 rounded-sm mt-3" onClick={logout}>Logout</button> */}
      <div className="grid lg:grid-cols-4 gap-4 mt-5">
        {games.map(({ title, description }) => <div className="shadow p-3 rounded-sm">
          <h1 className="text-lg font-bold">{title}</h1>
          <p className="text-sm">{description}</p>
          <button className="bg-blue-500 hover:bg-blue-600 text-sm text-white font-mono font-semibold px-5 py-1 mr-3 rounded-sm shadow mt-2" onClick={logout}>Play</button>
        </div>)}
      </div>
      <Game />
    </Container>
  )
}

const Main: React.FC = () => {
  const { authenticated, loading } = useAuth();
  if (loading) {
    return <>loading</>
  }

  if (authenticated) {
    return <LoggedIn/>
  }

  return <LoggedOut/>
}
const App: React.FC = () => {
  return <AuthProvider><Main /></AuthProvider>
}

export default App;
