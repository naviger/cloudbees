import { User } from "oidc-client-ts"
import { Seat } from "./seat"

export interface GetSeatResult {
  id:string
  status:string 
  customerId:string 
  seat:Seat
  user:User 
}