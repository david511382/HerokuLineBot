import React, { useEffect, useState } from 'react'
import Accordions from '../../components/accordion/Accordions'
import { AccordionData } from '../../components/accordion/Domain'
import Calendar from '../../components/calendar/Calendar'
import type {
    GetRentalCourtsDayCourts as GetRentalCourtsDayCourtsResp,
    GetRentalCourtsDayCourtsInfo as GetRentalCourtsDayCourtsInfoResp,
    GetRentalCourtsPayInfo as GetRentalCourtsResp,
    GetRentalCourtsPayInfoDay,
    GetRentalCourtsCourtInfo,
} from '../../models/resp/rental-courts'
import { GetRentalCourts } from '../../service/badminton/RentalCourts'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import TextField from '@material-ui/core/TextField'
import Button from '@material-ui/core/Button'
import {
    Time,
    FORMATE_DATE,
} from '../../models/util/Time'
import TotalDayCourtCalendarCard from './TotalDayCourtCalendarCard'
import TemporaryDrawer, { Anchor } from '../../components/temporaryDrawer/TemporaryDrawer'

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            display: 'flex',
            flexWrap: 'wrap',
        },
        textField: {
            marginLeft: theme.spacing(1),
            marginRight: theme.spacing(1),
            width: 200,
        },
    }),
)

export default function RentalCourt() {
    const classes = useStyles()

    const [totalDayCourts, setTotalDayCourts] = useState<Map<number, () => JSX.Element | (JSX.Element | undefined)[]>>(new Map())
    const [notPayAccordionDatas, setNotPayAccordionData] = useState<AccordionData>()
    const [notRefundaccordionDatas, setNotRefundAccordionData] = useState<AccordionData>()
    const [fromDate, setFromDate] = useState<Time>()
    const [toDate, setToDate] = useState<Time>()
    const [openTemporaryDrawerEvent, setOpenTemporaryDrawerEvent] = useState<React.KeyboardEvent<Element> | React.MouseEvent<Element, MouseEvent>>()

    const getRentalCourts = async () => {
        if (!fromDate ||
            !toDate) {
            return
        }

        const resp = await GetRentalCourts(fromDate, toDate)
        const respData = resp.data
        if (respData.not_pay_day_courts) {
            setNotPayAccordionData(parseGetRentalCourtsRespToAccordionData(respData.not_pay_day_courts))
        }
        if (respData.not_refund_day_courts) {
            setNotRefundAccordionData(parseGetRentalCourtsRespToAccordionData(respData.not_refund_day_courts))
        }
        if (respData.total_day_courts) {
            setTotalDayCourts(parseTotalDayCourtsToDateAccordionDatasMap(respData.total_day_courts))
        }
    }

    useEffect(() => {
        const now = new Time(undefined, FORMATE_DATE)
        if (!fromDate) {
            setFromDate(now)
        }
        if (!toDate) {
            setToDate(now)
        }
    }, [])

    return (
        <div>
            {
                fromDate &&
                <TextField
                    label="起始時間"
                    type="date"
                    defaultValue={fromDate.Format()}
                    className={classes.textField}
                    onChange={(e) => { setFromDate(new Time(e.target.value, FORMATE_DATE)) }}
                    InputLabelProps={{
                        shrink: true,
                    }}
                />
            }

            {
                toDate &&
                <TextField
                    label="截止時間"
                    type="date"
                    defaultValue={toDate.Format()}
                    className={classes.textField}
                    onChange={(e) => { setToDate(new Time(e.target.value, FORMATE_DATE)) }}
                    InputLabelProps={{
                        shrink: true,
                    }}
                />
            }
            <Button
                variant="contained"
                color="primary"
                onClick={getRentalCourts}
            >
                搜尋
            </Button>

            <Button
                variant="contained"
                color="inherit"
                onClick={(e) => { setOpenTemporaryDrawerEvent(e) }}
            >未清償款項</Button>
            <TemporaryDrawer
                show={openTemporaryDrawerEvent && { anchor: Anchor.RIGHT, event: openTemporaryDrawerEvent }}
            >
                <p>未付款</p>
                {(!notPayAccordionDatas) ? <p>無</p> : undefined}
                <Accordions data={notPayAccordionDatas}></Accordions>
                <p>未退款</p>
                {(!notPayAccordionDatas) ? <p>無</p> : undefined}
                <Accordions data={notRefundaccordionDatas}></Accordions>
            </TemporaryDrawer>

            <Calendar
                dateViewMap={totalDayCourts}
            ></Calendar>
        </div>
    )
}

function parseTotalDayCourtsToDateAccordionDatasMap(data: GetRentalCourtsDayCourtsResp[]): Map<number, () => JSX.Element | (JSX.Element | undefined)[]> {
    const resultDateAccordionDatasMap = new Map<number, () => JSX.Element>()

    data.forEach((v) => {
        const date = new Time(v.date, FORMATE_DATE)
        resultDateAccordionDatasMap.set(date.valueOf(), () => (<TotalDayCourtCalendarCard datas={v.courts}></TotalDayCourtCalendarCard>))
    })

    return resultDateAccordionDatasMap
}

function parseGetRentalCourtsRespToAccordionData(data: GetRentalCourtsResp): AccordionData {
    const summarys = [
        "總計",
        data.cost.toString(),
    ]
    const details: AccordionData[] = []
    data.courts.forEach((getRentalCourtsPayInfoDay) => {
        details.push(parseGetRentalCourtsPayInfoDayToAccordionData(getRentalCourtsPayInfoDay))
    })

    return {
        summarys: summarys,
        details: details,
    }
}

function parseGetRentalCourtsPayInfoDayToAccordionData(data: GetRentalCourtsPayInfoDay): AccordionData {
    const date = new Time(data.date, "yyyy/MM/dd (w)")
    const summarys = [
        date.Format(),
        data.cost.toString(),
    ]
    const details: AccordionData[] = []
    data.courts.forEach((getRentalCourtsCourtInfo) => {
        details.push(parseGetRentalCourtsCourtInfoToAccordionData(getRentalCourtsCourtInfo))
    })

    return {
        summarys: summarys,
        details: details,
    }
}

function parseGetRentalCourtsCourtInfoToAccordionData(data: GetRentalCourtsCourtInfo): AccordionData {
    const format = "hh:mm"
    const fromTime = new Time(data.from_time, format)
    const toTime = new Time(data.to_time, format)
    const summarys = [
        data.place,
        data.count.toString(),
        `${fromTime.Format()}~${toTime.Format()}`,
        data.cost.toString(),
    ]

    return {
        summarys: summarys,
    }
}
