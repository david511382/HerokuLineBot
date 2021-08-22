import React, { useEffect, useState } from 'react'
import {
  makeStyles,
  createStyles,
} from '@material-ui/core/styles';
import {
  Time,
} from '../../models/util/Time'

const useStyle =  makeStyles(theme =>
  createStyles({
    today:{
      backgroundColor: "yellow",
    },

    outDay:{
      backgroundColor: "lightgray",
    },

    day:{
      textAlign: "right",
      padding: "14px 0px",
      fontSize: "14px",
      borderBottom: "1px solid",
      borderRight: "1px solid",
      borderColor: "lightgray",
      display: "flex",
      flexDirection: "column",
    },
    innerDay:{
      display: "flex",
      flexDirection: "column",
      width: "100%",
    },
  }),
);

interface Props  {
  isToday: boolean 
  isCurrent: boolean 
  time: Time
  children?: JSX.Element | (JSX.Element | undefined)[] 
}

export default function CalendarDates(props: Props) {
  const classes = useStyle()

  const  [date, setDate]= useState<number>(props.time.getDate())
  const  [className, setClassName]= useState<string>()
  const  [children, setChildren]= useState<JSX.Element | (JSX.Element | undefined)[]>()

  useEffect(() => {
    setDate(props.time.getDate())
  },[props.time])

  useEffect(() => {
    let className = classes.day
    if (!props.isCurrent){
      className = className + ` ${classes.outDay}`
    }
    if (props.isToday){
      className = className + ` ${classes.today}`
    }
    setClassName(className)
  },[props.isCurrent, props.isToday])

  useEffect(() => {
    setChildren(props.children)
  },[props.children])
 
  return(
    <div className={className}>
      <span
        style={{ paddingRight: "6px" }}
      >
        {date}
      </span>
      <div className={classes.innerDay}>{
        children
      }</div>
    </div>
  )
}

