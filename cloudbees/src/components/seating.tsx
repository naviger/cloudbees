import React from 'react'
import './seating.css'
import { Car } from '../models/car'
import { Train } from '../models/train'

const Seating = () => {

  const cars:Array<Car> = []
  let a:Car = {
    name:"A",
    seats:[]
  }

  let b:Car = {
    name:"b",
    seats:[]
  }

  let r:number = 1;
  let p:string = "A"
  for(let i:number =0; i<10; i++) {
    a.seats.push({
      id: "A" + r + p,
      state: 'empty',
      passengerId: '',
      row: r,
      position: p
    })

    b.seats.push({
      id: "B" + r + p,
      state: 'empty',
      passengerId: '',
      row: r,
      position: p
    })

    switch(p) {
      case 'A':
        p='B'
        break
      case 'B':
        p = 'C'
        break
      case 'C':
        p = 'D'
        break
      case 'D':
        p = 'A'
        r = r + 1
        break
    }
  }


  let train:Train = {
    name:"London to Paris Cloud Bees Express",
    cars: [a, b]
  }
  return (
    <div className='content-frame'>
        <div className="train">
          { train.cars.map((c) => {
            return (<div key={"c"+ c.name} className="car">
              { c.seats.map((s, i)=> {
                return (
                  <div key={s.id} className={"seat" + (s.position === 'B' ? " aisle" : "")} >
                    <label>{s.id}</label>
                  </div>
                )
              })}
                <div className='lav'></div>
              </div>)
          })}
         </div>
    </div>
  )
}

export default Seating