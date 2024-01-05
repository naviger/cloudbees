import axios, { Axios, AxiosResponse, HttpStatusCode } from 'axios'

export type TrainControllerProps = {
  token:string
  protocol: string
  host:string
  port:number
}


export default class TrainController {
  token:string
  protocol: string
  host: string
  port?: number
  axios:Axios
  
  constructor(props:TrainControllerProps) {
    this.token = props.token
    this.protocol = props.protocol
    this.host = props.host
    this.port = props.port
    if (this.token.length > 0 ) {
      this.axios = axios.create({
        baseURL: this.protocol + "://" + this.host + (this.port ? ":" + this.port : ""),
        headers: {
          'Authorization': "Bearer " + this.token,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
         }
      })
    } else  {
      this.axios = axios.create({
        baseURL: props.protocol + "://" + props.host + (props.port ? ":" + props.port : ""),
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      })
    }
  }

  GetTrains = async (setTrains:Function) => {
    this.axios.get('/train')
    .then((response:AxiosResponse)=> {
      console.log(response.data)
      setTrains(response.data.data)
    })
  }

  GetTrain = async(trainId:string, setCurrentTrain:Function) => {
    this.axios.get('/train/' + trainId)
    .then((response:AxiosResponse) => {
      console.log(response.data)
      setCurrentTrain(response.data.data)
    })
  }

  GetSeat = async(trainId: string, seatId: string, email:string, firstName:string, lastName:string, setSeatResult:Function) => {
    this.axios.post('/train/' + trainId + "/seat", {
      seat: seatId,
      firstName: firstName,
      lastName: lastName,
      email: email
    })
    .then((response:AxiosResponse) => {
      console.log(response.data)
      setSeatResult(response.data.data)
    })
  }

  GetReceipt = async(trainId:string, seatId:string, customerId:string, showReceipt:Function) => {
    this.axios.get('/train/' + trainId + "/seat/" + seatId + "/receipt", { params: { customerId: customerId
    }})
    .then((response:AxiosResponse) => {
      console.log(response.data)
      showReceipt(response.data.data)
    })
  }

  changeSeat = async(trainId:string, source:string, dest:string, customerId:string, changeResult:Function) => {
    this.axios.patch('/train/' + trainId + "/seat/" + source + "/change", { 
      customerId: customerId,
      source: source,
      dest: dest
    })
    .then((response:AxiosResponse) => {
      console.log(response.data)
      changeResult(response.data.data)
    })
  }

  cancelSeat = async(trainId:string, seatId:string,customerId:string, cancelResult:Function) => {
    this.axios.patch('/train/' + trainId + "/seat/" + seatId + "/cancel", { 
      customerId: customerId,
      seatId: seatId,
    })
    .then((response:AxiosResponse) => {
      console.log(response.data)
      cancelResult(response.data.data)
    })
  }
}