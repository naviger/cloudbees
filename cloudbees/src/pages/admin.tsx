import React from 'react'
import './common.css'
import { useAuth } from "react-oidc-context"


const AdminPage = () => {

  const auth = useAuth()
  console.log(auth.user)
  return (
    <div className='content-frame'>
        <div>
          <div className="header">Admin</div>
          
        </div>
    </div>
  )
}

export default AdminPage
