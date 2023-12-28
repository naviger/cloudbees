import React from 'react'
import { useAuth } from "react-oidc-context"
import './common.css'
import Seating from '../components/seating'

const AdminPage = () => {
  let auth = useAuth()

  return (
    <div className="login">
        {/* <div className="logo">
            <button className="login-button" onClick={() => void auth.signinRedirect()}>Log in</button> */}
            <Seating />
        {/* </div> */}
    </div>
  )
}

export default AdminPage
