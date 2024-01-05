import { SeatReturn } from "./seatreturn";
import { User } from "./user";

export type Receipt = {
  id: string,
  trainId: string, 
  seat: SeatReturn,
  from: string,
  to: string,
  price:number,
  user: User
}