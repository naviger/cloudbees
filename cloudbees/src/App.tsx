import React from 'react'
import AdminPage from './pages/admin'
import ErrorPage, {ErrorPageProps} from './pages/error'
import DefaultPage from './pages/default'
import LoadingPage from './pages/loading'

import { useAuth } from "react-oidc-context"

import './App.css';

function App() {

  const auth = useAuth()

  switch (auth.activeNavigator) {
      case "signinSilent":
          return <div>Signing you in...</div>
      case "signoutRedirect":
          return <div>Signing you out...</div>
  }

  if (auth.isLoading) {
      return <div><LoadingPage /></div>
  }

  if (auth.error) {
      const props:ErrorPageProps = {
          source:"auth",
          message:auth.error.message
      }
      switch(auth.error.message) {
          case "Token is not active":
          case "No matching state found in storage":
          case "Session not active":
              window.location.href="http://isperience.web:3006"
              break;
          default:
              return <div><ErrorPage {...props} /></div>
      }
  }

  
    if (auth.isAuthenticated) {
      return (
          <div className="main">
            <AdminPage></AdminPage>
          </div>
      );
  }

  return <div className="splash">
      <div> </div>
      <div className="login">
          <div className="logo">
             <button className="login-button" onClick={() => void auth.signinRedirect()}>Log in</button>
          </div>
     
      </div>
  </div>
  
}

export default App;
