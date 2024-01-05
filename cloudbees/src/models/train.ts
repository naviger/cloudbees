import { Car } from "./car"

export interface Train {
  id: string
  name:string
  cars:Array<Car>
}