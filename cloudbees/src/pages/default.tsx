import React from 'react'
import { useAuth } from "react-oidc-context"
import './default.css'
import Seating from '../components/seating'
import {
  Dropdown,
  DropdownProps,
  Label,
  makeStyles,
  Option,
  shorthands,
  useId,
} from "@fluentui/react-components";
import TrainController from '../controllers/trainController';
import { Seat } from '../models/seat';
import { Train } from '../models/train';
import { Car } from '../models/car';
import { Stats } from '../models/stats';
import SeatDetails from '../components/seat';


const useStyles = makeStyles({
  root: {
    // Stack the label above the field with a gap
    display: "grid",
    gridTemplateRows: "repeat(1fr)",
    justifyItems: "start",
    ...shorthands.gap("20px"),
    maxWidth: "400px",
    "> div": {
      display: "grid",
      gridTemplateRows: "repeat(1fr)",
      justifyItems: "start",
      ...shorthands.gap("2px"),
    },
  },
});

const DefaultPage = () => {
  let auth = useAuth()
  const [trains, setTrains] = React.useState<Array<string>>([])
  const [currentTrain, setCurrentTrain] = React.useState<Train>({ id:"", name: "", cars:[]})
  const [stats, setStats] = React.useState<Stats>({ total: 0, occupied: 0})
  const [currentSeat, setCurrentSeat] = React.useState<Seat>()
  
  let token = ""
  if(auth.isAuthenticated) {
    token = auth.user?.access_token as string
  }
  
  const trainController = new TrainController({
    token: token,
    protocol: "https",
    host: "cloudbees.dev",
    port: 5007,
  })

  const setCurrentTrainData = (rawtrain:any) => {
    let train = JSON.parse(rawtrain)
    console.log("TRAIN: ", train)
    let t:Train ={
      id: train.trainId,
      name: "Cloudbees London to Paris Express: " + train.trainId,
      cars: []
    } 
    let carA:Car =  {
      name: 'Car A',
      seats: []
    }

    let carB:Car =  {
      name: 'Car B',
      seats: []
    }

    let cnt:number = 0
    for(var i:number=0; i < train["Seats"].length; i++) {
      let s = train.Seats[i]
      s.id = s.id.split(":")[1]
      if(s.status === 'occupied') {
        cnt++;
      }
      if(i < 10)  {
        carA.seats.push(s)
      } else {
        carB.seats.push(s)
      }
    }

    setStats({ total:  train["Seats"].length, occupied: cnt})
    t.cars.push(carA)
    t.cars.push(carB)

    console.log("Ordered train:", t)
    setCurrentTrain(t)
  }

  const getSeats = (trainId:string) => {
    trainController.GetTrain(trainId, setCurrentTrainData)
  }

  React.useEffect(() => {
    trainController.GetTrains(setTrains)
  }, [])
  
  const selectTrain : DropdownProps["onOptionSelect"] = (ev, data) => { 
    trainController.GetTrain(data.optionText as string, setCurrentTrainData)
  }

  let open:Array<Seat> = []
  if(currentTrain && currentTrain.cars.length === 2) {
    var seats:Array<Seat> = currentTrain.cars[0].seats.concat(currentTrain.cars[1].seats)
    open = seats.filter((s)=> { return s.status==="vacant"})
  }

  return (
    <div className="content">  
      <div className="form">
        <div className="formline">
          <div className="label">Train:</div>
          <Dropdown 
            placeholder="Select a train"
            onOptionSelect={selectTrain}
            size='small'
            appearance='underline'
          >
              {trains.map((t) => (
                <Option key={t} style={{ backgroundColor:'white'}} >{t}</Option>
              ))}
            
          </Dropdown>
        </div>
        { stats.total > 0 && <div className='stats'>
                Current occupancy is at {( stats.occupied/stats.total * 100).toFixed(0)}%. There are {stats.total - stats.occupied} open seats.
          </div>}
          <SeatDetails trainId={currentTrain.id} seat={currentSeat} openSeats={open} controller={trainController} getSeats={getSeats} />
      </div>
      <Seating train={currentTrain} setCurrentSeat={setCurrentSeat} />
    </div>
  )
}

export default DefaultPage
