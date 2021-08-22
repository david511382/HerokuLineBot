import React, { useEffect, useState } from 'react'
import Accordions from '../../components/accordion/Accordions'
import {AccordionData} from '../../components/accordion/Domain'
import {
    GetRentalCourtsDayCourtsInfo as GetRentalCourtsDayCourtsInfoResp,
    RentalCourtsStatus,
} from '../../models/resp/rental-courts'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import {
    Time,
    FORMATE_DATE,
} from '../../models/util/Time'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    cancel: {
      backgroundColor: 'lightgray',
    },
    notPay: {
        backgroundColor: 'lightyellow',
    },
    notRefund: {
        backgroundColor: 'LightCoral',
    },
    ok: {
        backgroundColor: 'lightgreen',
    },
  }),
)

interface Props {
    datas: GetRentalCourtsDayCourtsInfoResp[]
}

export default function TotalDayCourtCalendarCard(props: Props) {
    const classes = useStyles()

    const  [courts, setCourts]= useState<GetRentalCourtsDayCourtsInfoResp[]>()

    const format = "hh:mm"

    useEffect(() => {
        setCourts(props.datas)
    },[courts])

    return (
        <div>
            {
               courts?.map((v,index)=>{
                    const fromTime = new Time(v.from_time, format)
                    const toTime = new Time(v.to_time, format)
            
                    const summarys = [
                        `${fromTime.Format()}~${toTime.Format()}`,
                        v.count.toString(),
                    ]
                    const details :AccordionData[]= [{
                        summarys:[
                            v.cost.toString(),
                            v.reason_message,
                        ]
                    }]
                    if (v.refund_time){
                        const refundTime = new Time(v.refund_time, FORMATE_DATE)
                        details[0].summarys.push(refundTime.Format())
                    }
                    const data =  {
                        summarys:summarys,
                        details:details,
                    }
                    
                    let className = ""
                    switch(v.status){
                        case RentalCourtsStatus.RENTAL_COURTS_STATUS_CANCEL:
                            className = classes.cancel
                            break
                        case RentalCourtsStatus.RENTAL_COURTS_STATUS_NOT_PAY:
                            className = classes.notPay
                            break
                        case RentalCourtsStatus.RENTAL_COURTS_STATUS_NOT_REFUND:
                            className = classes.notRefund
                            break
                        case RentalCourtsStatus.RENTAL_COURTS_STATUS_OK:
                            className = classes.ok
                            break
                    }
                    
                    const key = `TotalDayCourtCalendarCard${index}`
                    return (
                        <Accordions className={className} key={key} data={data}></Accordions>
                    )
                })
            }
        </div>
    )
}