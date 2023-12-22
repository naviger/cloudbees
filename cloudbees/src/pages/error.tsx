import React from 'react'
import './error.css'

export type ErrorPageProps = {
  source: string
  message: string
}

const ErrorPage = (props:ErrorPageProps) => {

  return (
    <div className='content-frame'>
      {props.source==="auth" && 
        <div>
          <div className="auth-error">AUTHENTICATION ERROR</div>
          <div className="error-label">Error:</div><div className="error-message">{props.message}</div>
          <div>An authentication error occured. Please try to <a href="/" title="log in again">log in again</a>.</div>
        </div>
      } 
    </div>
  )
}

export default ErrorPage
