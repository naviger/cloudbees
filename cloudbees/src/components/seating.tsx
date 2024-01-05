import React, { MouseEventHandler } from 'react'
import './seating.css'
import { Car } from '../models/car'
import { Train } from '../models/train'
import { Seat } from '../models/seat'
import { Toilet } from './toilet'
import { useAuth } from "react-oidc-context"

export type TrainProps = {
  train: Train
  setCurrentSeat: Function
}
const Seating = (props: TrainProps) => {
  const auth = useAuth();
  const selectSeat = (sid:string) => {
    console.log(sid)
   
    const seats:Array<Seat> = props.train.cars[0].seats.concat(props.train.cars[1].seats)
    let st:Seat|undefined = seats.find((s)=> {return s.id === sid})
    props.setCurrentSeat(st)
  }

  return (
    <div className='content-frame'>
        <div className="train-name">{props.train.name}</div>
        <div className="train">
          { props.train.cars.map((c) => {
            return (<div key={"c"+ c.name} className="car">
              { c.seats.map((s, i)=> {
                return (
                  <div 
                    key={s.id} 
                    id={s.id} 
                    className={"seat" + (s.status==="occupied"? " occupied": " vacant") + (s.position === 'B' ? " aisle" : "" ) + (s.passengerId === auth.user?.profile.preferred_username ? " current-user" : "")} 
                    onClick={() => { selectSeat(s.id)}}>
                    <label>{s.id}</label>
                  </div>
                )
              })}
                <div className='lav'>
                  <Toilet />
                </div>
              </div>)
          })}
         </div>
    </div>
  )
}



export default Seating