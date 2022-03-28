import React from "react";

type Identity = {
  provider: string;
  user_id: string;
  user_name: string;
  avatar_id: string;
  color: string;
}

type Auth = {
  identity: Identity
  authenticated: boolean
}

const useAuth = () => {
  const context = React.useContext()
}
const Auth: React.FC = () = {

}

const App: React.FC = () => {
  return (
    <div>
      <div className="container mx-auto py-5">
        <h1>test</h1>
      </div>
    </div>
  );
}

export default App;
