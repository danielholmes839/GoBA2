import { useEffect, useRef, useState } from "react";

function deleteAllCookies() {
  var cookies = document.cookie.split(";");

  for (var i = 0; i < cookies.length; i++) {
    var cookie = cookies[i];
    var eqPos = cookie.indexOf("=");
    var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }
}

function Profile() {
  const [identity, setIdentity] = useState(null)

  useEffect(() => {
    fetch("http://localhost:3000/me", {
      method: "GET",
      credentials: "include",
    }).then(response => response.json()).then(data => {
      console.log("data", data)
      setIdentity(data)
    }).catch(err => {
      console.log(err)
    })
  }, [])

  if (identity === null) {
    return <><h1>identity not set</h1></>
  }

  return <>
    <span>
      <img className="rounded-full inline" style={{ background: identity.color, width: 32, border: "0.5px solid", borderColor: identity.color }} src={`https://cdn.discordapp.com/avatars/${identity.user_id}/${identity.avatar_id}.png?size=64`} />
      <span className="ml-3 inline text-lg font-mono">{identity.user_name}</span>
    </span>
  </>
}

function App() {
  const canvasRef = useRef(null);
  const [latest, setLatest] = useState({})

  useEffect(() => {
    const canvas = canvasRef.current
    const ctx = canvas.getContext('2d')
    ctx.fillStyle = '#ff7700'
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height)

    const ws = new WebSocket("ws://localhost:3000/game/connect")

    ws.addEventListener('open', (message) => {
      console.log("opened websocket", message);
    })

    ws.addEventListener('close', (message) => {
      console.log("closed websocket", message);
    })

    ws.addEventListener('error', (err) => {
      console.log("on error", err);
    })

    ws.addEventListener('message', (message) => {
      console.log("on message", message);
      setLatest(JSON.parse(message.data));
    })
  }, [])

  return <div className="p-5">
    <div className="container mx-auto">
      <h1 className="text-3xl mb-3">GoBA2!!!</h1>
      <Profile />
      <canvas style={{ width: 400, height: 400 }} ref={canvasRef} />
      <pre>{JSON.stringify(latest, null, 4)}</pre>

      <div>
        <button onClick={() => {
          console.log("clicked")
          window.location.assign("http://localhost:3000/auth/discord")
        }}>Login with Discord</button>
      </div>

      <button onClick={() => {
        deleteAllCookies();
        window.location.reload(false);
      }}>Logout</button>
    </div>
  </div>
}

export default App;
