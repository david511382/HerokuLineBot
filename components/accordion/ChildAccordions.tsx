import React from 'react';
import {
  makeStyles,
  Theme,
  createStyles,
  ThemeProvider,
  createTheme,
} from '@material-ui/core/styles';
import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import {
  AccordionDataViewData,
} from './Domain'

export const theme = createTheme({
  overrides: {
    MuiAccordion : {
      "root": {
        "&$expanded": {
          "margin": "0px 0",
        }
      }
    },
    MuiAccordionSummary : {
      "root": {
        "padding":"",
      }
    },
    MuiAccordionDetails : {
      "root": {
        "flex-direction": "column",
        "padding":"",
      }
    },
  },
});

export const useStyles = makeStyles<Theme, accordionDataViewDataStylesStyleProps>((theme: Theme) =>
  createStyles({
    root: {
      width: '100%',
    },
    heading: {
      fontSize: theme.typography.pxToRem(15),
      flexBasis: props => props.flexBasisStr,
      flexShrink: 0,
    },
    minLevel: {
      textAlign: "center",
    },
}));

interface Props  {
  data: AccordionDataViewData 
  className? : string
}

export default function ChildAccordions(props: Props) {
  const accordionDataViewData = props.data
  const flexBasis = 100 / accordionDataViewData.summarys.length
  const flexBasisStr = `${flexBasis.toString()}%`
  const classes = useStyles({flexBasisStr})

  const startID = accordionDataViewData.start_id
  const isContainChildern = accordionDataViewData.details.length > 0
  const ariaControls = `panel${startID}bh-content`
  const id = `panel${startID}bh-header`
  const accordionKey = `accordion${startID}`
  const typographyBaseKey = `typography${startID}-`
  
  return(
    <ThemeProvider theme={theme}>
      <Accordion 
        className={`${classes.root} ${props.className}`}
        expanded={
          (isContainChildern)?
            undefined:
            false
        }
      >
        <AccordionSummary
          expandIcon={
            isContainChildern &&
            <ExpandMoreIcon />
          }
          aria-controls={ariaControls}
          id={id}
        >
          {
            accordionDataViewData.summarys.map(
              (summary,i)=> 
                <Typography 
                  key={`${typographyBaseKey}${i.toString()}`}
                  className={
                    (isContainChildern)?
                      classes.heading:
                      `${classes.heading} ${classes.minLevel}`
                  }
                >
                  {summary}
                </Typography>
            )
          }
        </AccordionSummary>
        <AccordionDetails>
          {
            accordionDataViewData.details.map((accordionDataViewData,i)=>{
              return (
                <ChildAccordions key={`${accordionKey}-${i.toString()}`} data={accordionDataViewData} ></ChildAccordions>
              )
            })
          }
        </AccordionDetails>
      </Accordion>
    </ThemeProvider>
  )
}

interface accordionDataViewDataStylesStyleProps {
  flexBasisStr : string
}