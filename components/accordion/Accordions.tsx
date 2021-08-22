import React from 'react';
import {
  makeStyles,
  Theme,
  createStyles,
} from '@material-ui/core/styles';
import {
  AccordionData,
  AccordionDataViewData,
} from './Domain'
import ChildAccordions from './ChildAccordions'

const useStyle =  makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: '100%',
    },
    expanded:{
      margin:'0px',
    }
  }),
);

interface Props {
  data?: AccordionData | AccordionData[] | AccordionDataViewData[]
  className? : string
}

export default function Accordions(props: Props) {
  const classes = useStyle()
  
  if (!props.data){
    return (null)
  }
  if (Array.isArray(props.data) &&
    props.data.length === 0){
    return (null)
  }
  
   let accordionDataViewDatas : AccordionDataViewData[] = []
  if (instanceOfAccordionData(props.data) ||
    !instanceOfAccordionDataViewData(props.data)){
    const accordionDatas : AccordionData[] = []
    if (instanceOfAccordionData(props.data)){
      accordionDatas.push(props.data as AccordionData)
    }else if (instanceOfAccordionData(props.data[0])) {
      const propsData = props.data as AccordionData[]
      propsData.forEach((accordionData)=>{
        if (!accordionData ||
          !accordionData.summarys ||
          accordionData.summarys.length === 0){
          return
        }

        accordionDatas.push(accordionData)
      })
    }

    const convertResult = convert(accordionDatas, 0)
    accordionDataViewDatas = convertResult.accordionDataViewDatas
  }else{
    accordionDataViewDatas = props.data as AccordionDataViewData[]
  }

  return(
    <div className={classes.root}>
      {
        accordionDataViewDatas.map((accordionDataViewData)=>{
          const accordionKey = `accordion${accordionDataViewData.start_id}`
          return (
            <ChildAccordions className={props.className} key={accordionKey} data={accordionDataViewData} ></ChildAccordions>
          )
        })
      }
    </div>
  )
}

interface ConvertResult {
  accordionDataViewDatas : AccordionDataViewData[],
  lastID : number
}

function convert(accordionDatas: AccordionData[], startID : number) :ConvertResult {
  const resultAccordionDataViewDatas : AccordionDataViewData[] = []
  let lastID = startID
  accordionDatas.forEach((accordionData)=>{
    let details : AccordionDataViewData[] = []
    if (accordionData.details){
      const convertResult = convert(accordionData.details, startID + 1)
      details = convertResult.accordionDataViewDatas
      lastID = convertResult.lastID
    }

    resultAccordionDataViewDatas.push(
      {
        summarys: accordionData.summarys,
        details : details,
        start_id: startID,
      },
    )

    startID = lastID
  })

  return {
    accordionDataViewDatas : resultAccordionDataViewDatas,
    lastID: lastID,
  }
}


function instanceOfAccordionDataViewData(object: any): object is AccordionDataViewData {
  return (object as AccordionDataViewData).start_id !== undefined
}

function instanceOfAccordionData(object: any): object is AccordionData {
  return (object as AccordionData).summarys !== undefined
}
