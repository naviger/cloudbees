import React from 'react'
import { useAuth } from "react-oidc-context"
import './nav.css'

const Nav = () => {
  const auth = useAuth()

  return (
    <div className="nav">
      <div className="logo">
      </div>
      
      <div className="center">CloudBees Train Service</div>
      <div className="auth">
      { auth.isAuthenticated && <div><span>{ auth.user?.profile.name}{" "}</span>
        <button onClick={() => void auth.signoutSilent()}>Log out</button></div> }
      { !auth.isAuthenticated && 
        <button className="login-button" onClick={() => void auth.signinRedirect()}>Log in</button>
        
      }
      </div>
      <div className="underlay"></div>
    </div>
  )
}

export default Nav