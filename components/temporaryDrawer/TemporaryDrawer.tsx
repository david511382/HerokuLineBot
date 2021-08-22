import React, { useEffect } from 'react'
import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';

const useStyles = makeStyles({
  list: {
    width: 250,
  },
  fullList: {
    width: 'auto',
  },
});

export enum Anchor  {
  TOP = 'top',
  LEFT = 'left',
  BOTTOM = 'bottom',
  RIGHT = 'right',
}

interface Props  {
  children?: JSX.Element | (JSX.Element | undefined)[] 
  clickOnce?: boolean
  show?: {
    anchor : Anchor,
    event : React.KeyboardEvent | React.MouseEvent,
  },
}

export default function TemporaryDrawer(props: Props) {
  const classes = useStyles();
  const [state, setState] = React.useState({
    top: false,
    left: false,
    bottom: false,
    right: false,
  });

  useEffect(() => {
    console.log(props.show)
    if (props.show)
      toggleDrawer(props.show.anchor, true)(props.show.event)
  },[props.show])

  const toggleDrawer = (anchor: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }

    setState({ ...state, [anchor]: open });
  };

  const list = (anchor: Anchor) => (
    <div
      className={clsx(classes.list, {
        [classes.fullList]: anchor === 'top' || anchor === 'bottom',
      })}
      role="presentation"
      onClick={(props.clickOnce)? toggleDrawer(anchor, false): undefined} 
      onKeyDown={(props.clickOnce)? toggleDrawer(anchor, false): undefined} 
    >
      {props.children}
    </div>
  );

  return (
    <div>
      {(['left', 'right', 'top', 'bottom'] as Anchor[]).map((anchor) => (
        <React.Fragment key={anchor}>
          <Drawer anchor={anchor} open={state[anchor]} onClose={toggleDrawer(anchor, false)}>
            {list(anchor)}
          </Drawer>
        </React.Fragment>
      ))}
    </div>
  );
}
