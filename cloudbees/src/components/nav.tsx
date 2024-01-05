import React from 'react'
import { useAuth } from "react-oidc-context"
import './nav.css'

const Nav = () => {
  const auth = useAuth()

  const handleLogout = () => {
    // auth.removeUser();
    // auth.signoutRedirect();
    const url = window.location.href
    auth.signoutSilent()
    window.location.href = url
  };

  return (
    <div className="nav">
      <div className="logo">
      </div>
      
      <div className="center">CloudBees Train Service</div>
      <div className="auth">
      { auth.isAuthenticated && <div><span>{ auth.user?.profile.name}{" "}</span>
        <button onClick={handleLogout}>Log out</button></div> }
        {/* //auth.signoutSilent() {() => void auth.signoutSilent() */}
      { !auth.isAuthenticated && 
        <button className="login-button" onClick={() => void auth.signinRedirect()}>Log in</button>
        
      }
      </div>
      <div className="underlay"></div>
    </div>
  )
}

export default Nav