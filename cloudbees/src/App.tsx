import React from 'react'
import AdminPage from './pages/admin'
import ErrorPage, {ErrorPageProps} from './pages/error'
import DefaultPage from './pages/default'
import LoadingPage from './pages/loading'

import { useAuth } from "react-oidc-context"

import './App.css';
import Nav from './components/nav'

const parseJwt = (token:string) => {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch (e) {
    return null;
  }
}

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
            window.location.href="https://cloudbees.dev:3443"
            break;
        default:
            return <div><ErrorPage {...props} /></div>
    }
  }

  let isAdmin = false
  if(auth.isAuthenticated) {
    const jwt = parseJwt(auth.user?.access_token as string)
    console.log("SECURE DATA:", auth.user, jwt)
    
    if(jwt["resource_access"]) {
      const ra = jwt["resource_access"]
      if (ra["cloudbees-client"]) {
        const roles = ra["cloudbees-client"]["roles"]
        if(roles.indexOf("travel_admin") > -1){
          isAdmin=true
        }
        console.log(roles)
      }
    }
  }

  
  if (auth.isAuthenticated && isAdmin) {
    return (
        <div className="main">
           <Nav></Nav>
          <AdminPage></AdminPage>
        </div>
    )
  }

  return (
    <div className="main">
      <Nav></Nav>
      <DefaultPage></DefaultPage>
    </div>
  )
  
}

export default App;
