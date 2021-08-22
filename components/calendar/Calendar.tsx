import React, { useEffect, useState } from 'react'
import {
  makeStyles,
  createStyles,
} from '@material-ui/core/styles';
import {
  Time,
  FORMATE_DATE,
  TimeType,
} from '../../models/util/Time'
import CalendarDates from './CalendarDates'

const useStyle =  makeStyles(theme =>
  createStyles({
    root: {
      width: '100%',
    },

    dayName:{
      fontSize: "16px",
      textTransform: "uppercase",
      textAlign: "center",
      borderBottom: "1px solid",
      borderTop: "1px solid",
      lineHeight: "50px",
      fontWeight: 500,
    },

    calendarHeader:{
      display: "flex",
      justifyContent: "space-between",
      alignItems: "center",
    },
    calendarTitle:{
      marginTop: "0.83em",
      marginBottom: "0.83em",
    },
    calendarNavigate:{
      padding: "10px",
      opacity:  0.5,
      '&:hover': {
        cursor: "pointer",
        opacity: 0.9,
      },
    },

    calendarBody:{
      width: "100%",
      display: "grid",
      gridTemplateColumns: "repeat(7, minmax(40px, 1fr))",
    },
    calendarDates:{
      display: "grid",
      gridAutoRows: "minmax(120px, auto)",
    },

    unselectable:{
      WebkitUserSelect: "none",
      MozUserSelect: "none",
      msUserSelect: "none",
      userSelect: "none",
    },

    sticky: {
      position: 'sticky',
      top: 0,
      zIndex: 1,
      backgroundColor: 'inherit',
    },
  }),
);

interface ChildAccordionsProps  {
  isToday: boolean 
  isCurrent: boolean 
  time: Time
}

interface Props {
  defaultTime? : Time
  dayStart? : number
  dateViewMap?: Map<number,()=>JSX.Element | (JSX.Element | undefined)[]>
  disableHeadSticky?: boolean
}

export default function Calendar(props: Props) {
  const classes = useStyle()

  const  [currentTime, setCurrentTime]= useState<Time>()
  const  [today, setToday]= useState<Time>()
  const  [dates, setDates]= useState<ChildAccordionsProps[]>([])
  const  [dateViewMap, setDateViewMap]= useState<Map<number, ()=>JSX.Element | (JSX.Element | undefined)[]>>()

  const defaultDayStart = (props.dayStart) ?
    props.dayStart:
    0
  const weekdays = Time.Weekdays(defaultDayStart)
  const defaultTimeType = TimeType.MONTH

  useEffect(() => {
    setDateViewMap(props.dateViewMap)
  },[props.dateViewMap])

  useEffect(() => {
    if (!today || !currentTime){
      const today = new Time(undefined, FORMATE_DATE)
      setToday(today)

      if (currentTime){
        return
      }

      if (props.defaultTime){
        setCurrentTime(props.defaultTime.Of(defaultTimeType))
      }else{
        setCurrentTime(today.Of(defaultTimeType))
      }

      return
    }

    const newDates : ChildAccordionsProps[] = []
    const dayOfFirst = currentTime.getDay()
    const preDateCount = dayOfFirst - defaultDayStart
    const startDate = currentTime.Next(TimeType.DATE, -preDateCount)
    const nextMonth = currentTime.Next(defaultTimeType, 1)
    const currentLastDate = nextMonth.Next(TimeType.DATE, -1)
    const dayOfEnd =currentLastDate.getDay()
    const nextDateCount = defaultDayStart + 6 - dayOfEnd
    const endNextDate = currentLastDate.Next(TimeType.DATE, nextDateCount + 1)
    
    const todayDate = today.Of(TimeType.DATE)
    Time.Slice(
      startDate, currentTime,
      (t)=>t.Next(TimeType.DATE, 1),
      (runTime, next) =>{
        newDates.push({        
          isToday: todayDate.valueOf() === runTime.valueOf(),
          isCurrent: false, 
          time: runTime,
        })
        
        return true
      }
    )
    Time.Slice(
      currentTime, nextMonth,
      (t)=>t.Next(TimeType.DATE, 1),
      (runTime, next) =>{
        newDates.push({        
          isToday: todayDate.valueOf() === runTime.valueOf(),
          isCurrent: true, 
          time: runTime,
        })
        return true
      }
    )
    Time.Slice(
      nextMonth, endNextDate,
      (t)=>t.Next(TimeType.DATE, 1),
      (runTime, next) =>{
        newDates.push({        
          isToday: todayDate.valueOf() === runTime.valueOf(),
          isCurrent: false, 
          time: runTime,
        })
        return true
      }
    )
    setDates(newDates)
  },[currentTime, today])

  const renderDays = () =>{
    return weekdays.map((weekDay, i) => (
      <div
        className={classes.dayName}
        key={"day-of-week-" + i}
        style={{ borderColor: "LightGray" }}
      >
        {weekDay}
      </div>
    ));
  }

  const renderDates = () => {
    return [
      dates.map(dateTime => {
        const date = dateTime.time.getDate()
        const key = (dateTime.isCurrent)?
          "day-" + date :
          "out-day-" + date
        const view = dateViewMap?.get(dateTime.time.valueOf())
        
        return (
          <CalendarDates
            key={key}
            isToday={dateTime.isToday}
            isCurrent={dateTime.isCurrent}
            time={dateTime.time}
          >{
            view &&
            view()
          }</CalendarDates>
        )
      })
    ]
  }

  return(
    <div
      className="calendar"
      style={{
        fontSize: "18px",
        border: "1px solid",
        minWidth: "300px",
        position: "relative",
        borderColor: "LightGray",
        color: "#51565d",
        backgroundColor: "white",
      }}
    >
      <div className={
          (props.disableHeadSticky)?
            "":
            classes.sticky
        }
      >
        <div className={classes.calendarHeader}>
          <div
            className={`${classes.calendarNavigate} ${classes.unselectable}`}
            onClick={()=>{setCurrentTime(currentTime?.Next(TimeType.MONTH,-1))}}
          >
            &#10094;
          </div>
          <div>
            <h2 className={classes.calendarTitle}>
              {currentTime?.Format("yyyy MM月")}
            </h2>
          </div>
          <div
            className={`${classes.calendarNavigate} ${classes.unselectable}`}
            onClick={()=>{setCurrentTime(currentTime?.Next(TimeType.MONTH,1))}}
          >
            &#10095;
          </div>
        </div>
        <div className={classes.calendarBody}>
          {renderDays()}
        </div>
      </div>
      <div className={`${classes.calendarBody} ${classes.calendarDates}`}>
        {renderDates()}
      </div>
    </div>
  )
}

