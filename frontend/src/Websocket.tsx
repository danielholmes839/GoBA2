// import React, { useEffect, useState } from "react";
// import { useAuth } from "./Auth";

// type Game = "game1" | "game2" | "game3"


// type WebsocketContextValue = {
//   connected: boolean;
//   error: string;
//   socket: WebSocket;
// }

// const WebsocketContext = React.createContext<WebsocketContextValue>({
//   connected: false,
//   error: "",
//   socket: new WebSocket("")
// })

// export const useWebsocket = (): WebsocketContextValue => {
//   return React.useContext(WebsocketContext)
// }


// // type WebsocketProviderProps = {
// //   socket: WebSocket;
// //   onOpen: () => void;
// //   onClose: () => void;
// //   onMessage: () => void;
// // }

// export const WebsocketProvider: React.FC = ({ children }) => {
//   const { authenticated } = useAuth();
//   const [state, setState] = useState({
//     connected: false,
//     connecting: false,
//     error: ""
//   });

//   const connect = () => {
//     const ws = new WebSocket("ws://localhost:3000/connect")
//     setState({...state, connecting: true})

//     ws.addEventListener('open', (message) => {
//       setState({...state, connecting: false, connected: true})
//       console.log("opened websocket", message);
//     })

//     ws.addEventListener('close', (message) => {
//       setState({...state, connecting: false, connected: true})
//       console.log("closed websocket", message);
//     })

//     ws.addEventListener('error', (err) => {
//       setState({...state, connecting: false, connected: false, error})
//       console.log("on error", err);
//     })

//     ws.addEventListener('message', (message) => {
//       console.log("message", message)
//     })
//   }

//   if (!authenticated) {
//     return <WebsocketContext.Provider value={{
//       connected: false,
//       error: "please login",
//     }}>{children}</WebsocketContext.Provider>
//   }

//   return <WebsocketContext.Provider value={{
//       connected: false,
//       error: "please login",
//       socket: socket
//   }}>

//   </WebsocketContext.Provider>
// }

export const data = {}