import React, { useEffect, useState } from "react";

function deleteAllCookies() {
  var cookies = document.cookie.split(";");

  for (var i = 0; i < cookies.length; i++) {
    var cookie = cookies[i];
    var eqPos = cookie.indexOf("=");
    var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }
}

type Identity = {
  provider: string;
  user_id: string;
  user_name: string;
  avatar_id: string;
  color: string;
}

export type AuthValue = {
  identity: Identity | null;
  loading: boolean,
  authenticated: boolean;
}

export type AuthContextValue = AuthValue & {
  login: () => void,
  logout: () => void,
}

const defaultAuth = {
  loading: true,
  authenticated: false,
  identity: null,
  login: () => { },
  logout: () => { }
}

const AuthContext = React.createContext<AuthContextValue>(defaultAuth)

export const useAuth = () => {
  return React.useContext(AuthContext)
}

export const AuthProvider: React.FC = ({ children }) => {
  const [state, setState] = useState<AuthValue>(defaultAuth);

  useEffect(() => {
    loadIdentity()
  }, [])

  const loadIdentity = async () => {
    try {
      setState({ ...state, loading: true })
      const response = await fetch("http://localhost:3000/@me", { credentials: "include" });
      const identity: Identity = await response.json();
      setState({ ...state, identity: identity, authenticated: true, loading: false })
    } catch {
      console.log("user not logged in")
      setState({ ...state, identity: null, authenticated: false, loading: false })
    }
  }

  const login = () => {
    window.location.assign("http://localhost:3000/auth/discord")
  }

  const logout = async () => {
    deleteAllCookies();
    setState({ authenticated: false, identity: null, loading: false });
  }

  return <AuthContext.Provider value={{
    ...state, login, logout
  }}>
    {children}
  </AuthContext.Provider>
}