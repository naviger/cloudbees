import { Seat } from "./seat"

export interface ChangeResult {
  id:string
  customerId:string 
  sourceSeat:Seat 
  destSeat:Seat 
}