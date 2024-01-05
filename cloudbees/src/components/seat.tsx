import React from 'react'
import { Seat } from '../models/seat'
import { useAuth } from "react-oidc-context"
import './seat.css'
import { Dropdown, DropdownProps, Option, makeStyles, shorthands, tokens } from '@fluentui/react-components'
import TrainController from '../controllers/trainController'
import { User } from 'oidc-client-ts'
import { GetSeatResult } from '../models/getSeatResult'
import { Receipt } from '../models/receipt'
import { CancelResult } from '../models/cancelResult'
import { ChangeResult } from '../models/changeResult'
export type SeatProps = {
  trainId: string,
  seat: Seat | undefined
  openSeats: Array<Seat>
  controller:TrainController,
  getSeats:Function,
}

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
      // need padding to see the background color for filled variants
      ...shorthands.padding("5px", "20px", "10px"),
    },
  },
  // filledLighter and filledDarker appearances depend on particular background colors
  filledLighter: {
    backgroundColor: tokens.colorNeutralBackgroundInverted,
    "> label": {
      color: tokens.colorNeutralForegroundInverted2,
    },
    "> h3": {
      color: tokens.colorNeutralForegroundInverted2,
    },
  },
  filledDarker: {
    backgroundColor: tokens.colorNeutralBackgroundInverted,
    "> label": {
      color: tokens.colorNeutralForegroundInverted2,
    },
    "> h3": {
      color: tokens.colorNeutralForegroundInverted2,
    },
  },
});

type ReserveForm = {
  firstName: string
  lastName: string
  email: string
}

const SeatDetails = (props: SeatProps) => {
  const styles = useStyles();
  const auth = useAuth()
  
  var frm: ReserveForm = {
    firstName: '',
    lastName: '',
    email: ''
  }
  
  let isCustomer:boolean = false
  let isAdmin:boolean = false

  if(auth.isAuthenticated) {
    const resources:any = auth.user?.profile["resource_access"]
    const roles:string[] = resources["cloudbees-client"]["roles"]
    isCustomer = roles.indexOf("travel_customer") > -1
    isAdmin = roles.indexOf("travel_admin") > -1
    frm.firstName = auth.user?.profile.given_name as string
    frm.lastName = auth.user?.profile.family_name as string
    frm.email = auth.user?.profile.email as string
  }
  
  const [mode, setMode] = React.useState<string>("")
  const [resultMode, setResultMode] = React.useState<string>("")
  const [sourceSeat, setSourceSeat] = React.useState<Seat|undefined>(props.seat)
  const [destinationSeat, setDestinationSeat] = React.useState<Seat>()
  const [reserveForm, setReserveForm] = React.useState<ReserveForm>(frm)
  const [resultData, setResultData] = React.useState<any|string|undefined>()
  //GetSeatResult|Receipt|CancelResult|ChangeResult

  const orderSeat = () => {
    setMode("reserve")
  }

  React.useEffect(() => {
    if(sourceSeat!= props.seat) {
      setResultMode("")
      setResultData("")
    }
  }, [])

  const cancelSeat = () => {
    //console.log("CANCEL SEAT: " + props.seat?.id)
    props.controller.cancelSeat(props.trainId, props.seat?.id as string, auth.user?.profile.preferred_username as string, setCancelled)
  }

  const changeSeat = () => {
    setMode("change")
  }

  const getReceipt = () => {
    props.controller.GetReceipt(props.trainId, props.seat?.id as string, auth.user?.profile?.preferred_username as string, showReceipt)
  }

  const selectDestinationSeat  : DropdownProps["onOptionSelect"] = (ev, data) => { 
    const s:Seat|undefined = props.openSeats.find((s:Seat) => {return s.id === data.optionValue})
    console.log(s, data.optionValue)
    if(s) {
      setDestinationSeat(s)
      document.getElementById("btnChangeSeat")?.removeAttribute("disabled")
    }
  }

  const setSeatChange = () => {
    console.log("CHANGE SEAT:", props.seat?.id, destinationSeat?.id)
    document.getElementById("btnChangeSeat")?.setAttribute("disabled", "true")
    props.controller.changeSeat(props.trainId, props.seat?.id as string, destinationSeat?.id as string, auth.user?.profile.preferred_username as string, setChanged)
  }

  const setSeatReserve = (event:any) => {
    console.log("RESERVE SEAT:", props.seat?.id)
    document.getElementById("btnChangeSeat")?.setAttribute("disabled", "true")
    props.controller.GetSeat(props.trainId, props.seat?.id as string, reserveForm.email, reserveForm.firstName, reserveForm.lastName, seatResults)
  }

  const validateReserveForm = () => {
    if(reserveForm.firstName.length > 2 && reserveForm.lastName.length > 2 && validateEmail(reserveForm.email)) {
      document.getElementById("btnSeatReserve")?.removeAttribute("disabled")
    } else {
      document.getElementById("btnSeatReserve")?.setAttribute("disabled", "true")
    }
  }

  const changeInput = (event:any) => {
    let value = event.target.value;
    let frm:ReserveForm = structuredClone(reserveForm)
    console.log("id:", event.target.id, "value:", value)
    switch (event.target.id) {
      case "frmFN":
        frm.firstName = value
        break
      case "frmLN":
        frm.lastName = value
        break;
      case "frmEM":
        frm.email = value
        break;
    }
    console.log("FORM:", frm)
    validateReserveForm()
    setReserveForm(frm)
  }

  const seatResults = (data:any) => {
    const rslt = JSON.parse(data) 
    console.log("SEAT RESULTS:", rslt)
    setResultData(rslt)
    props.getSeats(props.trainId)
    setResultMode("seatResult")
    setMode("")
  }

  const showReceipt = (data:any) => {
    const rslt =  data 
    console.log("RECEIPT:", rslt)
    setResultData(data)
    setResultMode("showReceipt")
    setMode("")
  }

  const setCancelled = (data:any) => {
    console.log("seat results:", JSON.parse(data))
    const rslt = JSON.parse(data) 
    setResultData(rslt)
    setResultMode("cancelResult")
    props.getSeats(props.trainId)
    setMode("")
  }

  const setChanged = (data:any) => {
    const rslt = JSON.parse(data) 
    console.log("seat results:", JSON.parse(data))
    setResultData(rslt)
    setResultMode("changeResult")
    props.getSeats(props.trainId)
    setMode("")
  }

 
  const validateEmail = (email:string):boolean => {
    const em = String(email)
      .toLowerCase()
      .match(
        /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      );
    console.log(email, em)
    return em ? true:false;
  };

  
  
  return (
    <div className={ 'seat-details ' }>
      { props.seat != undefined && 
      <div>
        <div className='seat-name'>{props.seat.id}</div>
        <div className='seat-details-inner'>
          <div className={"seat-operations "  }>
            { props.seat?.status === "occupied" && isCustomer && <div></div> }
            { props.seat?.status === "occupied" && isAdmin && <div className="button-container">
              <button onClick={cancelSeat}>Cancel Seat</button>
              <button onClick={changeSeat}>Change Seat</button>
            </div> }

            { props.seat?.status === "vacant" && (isCustomer || !auth.isAuthenticated )&& <button onClick={orderSeat}>Reserve Seat</button>  }
            { props.seat?.status === "occupied" && isCustomer && props.seat.passengerId === auth.user?.profile.preferred_username &&
              <div className="button-container">
                { isCustomer &&  <button onClick={getReceipt}> Get Receipt</button>}
                <button onClick={cancelSeat}>Cancel Seat</button>
              <button onClick={changeSeat}>Change Seat</button>
              </div>}
            { mode==="change" && <div className="change-seat">
                <label>Change Seat</label>
                <Dropdown className="seat-dd"
                style={{minWidth: '160px', maxWidth: '160px', fontSize:'0.8em'}}
                  placeholder="Select a new seat"
                  onOptionSelect={selectDestinationSeat}
                  size='small'
                  appearance='underline'
                >
                    { props.openSeats.map((s) => (
                      <Option key={s.id} style={{ backgroundColor:'white'}}  >{s.id}</Option>
                  ))} 
                </Dropdown>
                <div className="change-details">
                  { !destinationSeat && <span>Please select a destination seat.</span> }
                  { destinationSeat && <span>Change seat from {props.seat.id} to {destinationSeat.id}.</span> }
                </div> 
                <button id="btnChangeSeat" onClick={setSeatChange}>Go</button>
              </div>
            }
             { mode==="reserve" && <div className="reserve-seat">
              <label>Reserve Seat</label>
                <div className="formline-uo">
                  <label className="lbl-offset">First Name</label>
                  <input id="frmFN" type="text" defaultValue={reserveForm.firstName} onChange={changeInput}  />
                </div>
                <div className="formline-uo">
                  <label className="lbl-offset">Last Name</label>
                  <input  id="frmLN" type="text" defaultValue={reserveForm.lastName}  onChange={changeInput}  />
                </div>
                <div className="formline-uo">
                  <label  className="lbl-offset">Email Address</label>
                  <input  id="frmEM" type="text" defaultValue={reserveForm.email}  onChange={changeInput}  />
                </div>
                <button id="btnSeatReserve" onClick={(e:any)=>{console.log('A');setSeatReserve(e)}}>Go</button>
              </div>
            }
          </div>
          <div className='seat-data'>
            { props.seat?.status === "occupied" && isAdmin && <div>
              <label>Passenger</label><input type='text' value={props.seat.passengerId} /> 
            </div> }
            { props.seat?.status === "occupied" && isCustomer && props.seat.passengerId === auth.user?.profile.preferred_username && <div>
              <label>Passenger</label><input type='text' value={props.seat.passengerId} /> 
            </div> }
           
            { resultMode === "seatResult" && <div>
              {/* <div>{getSeatResultText()}</div> */}
              <div className="formline-uo">
                <label>Id</label>
                <input readOnly type="text" value ={resultData?.id} />
              </div>
              <div className="formline-uo">
                <label>Status</label>
                <input readOnly type="text" value ={(resultData as GetSeatResult)?.status } />
              </div>
              { auth.isAuthenticated && 
                <div className="formline-uo">
                  <button onClick={getReceipt}> Get Receipt</button>
                </div>
              } 
              {!auth.isAuthenticated && 
                <div className="formline-uo">
                  <span>To get receipt, please log in using <i>&lt;firstname&gt;.&lt;lastname&gt;</i> as a username and <i>Test123!</i> as a password</span>
                </div>
              }


            </div> }
            { resultMode==="showReceipt" && 
              <div className="receipt">
                <pre>
                  {/* { getReceiptText() } */}
                  {resultData}
                </pre>
              </div>
            }
            { resultMode==="cancelResult" && 
              <div className="receipt">
                <div className="formline-uo">
                  <label>Id</label>
                  <input readOnly type="text" value ={resultData?.id} />
                </div>
                <div className="formline-uo">
                  <label >Status</label>
                  <input readOnly type="text" value ={(resultData as CancelResult)?.status + ": Reservation has been cancelled" } />
                </div>
              </div>
            }
            { resultMode==="changeResult" && 
              <div className="receipt">
                <div className="formline-uo">
                  <label>Id</label>
                  <input readOnly type="text" value ={resultData?.id} />
                </div>
                <div className="formline-uo">
                  <label>Status</label>
                  <input readOnly type="text" value ={((resultData as ChangeResult).sourceSeat as Seat).id + " changed to " + ((resultData as ChangeResult).destSeat as Seat).id } />
                </div>
              </div>
            }
          </div>
        </div>
      </div>
    }
    { props.seat === undefined && <div className="select-seat-msg">Please select a seat   &gt;&gt;&gt;&gt;&gt;&gt;------------&gt;</div> }
    </div>
  )
}

export default SeatDetails