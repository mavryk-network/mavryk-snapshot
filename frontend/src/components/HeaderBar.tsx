import React, { useEffect, useState } from 'react'
import AppBar from '@mui/material/AppBar'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { useTheme } from '@mui/material/styles'
import SnapshotLink from './SnapshotLink'
import Separator from './Separator'
import Link from '@mui/material/Link'
import IconButton from '@mui/material/IconButton'
import Brightness4Icon from '@mui/icons-material/Brightness4'
import Brightness7Icon from '@mui/icons-material/Brightness7'
import { ColorModeContext } from '../ThemeContext'

export default function HeaderBar () {
  const [width, setWidth] = useState<number>(window.innerWidth)

  function handleWindowSizeChange () {
    setWidth(window.innerWidth)
  }
  useEffect(() => {
    window.addEventListener('resize', handleWindowSizeChange)
    return () => {
      window.removeEventListener('resize', handleWindowSizeChange)
    }
  }, [])

  const isMobile = () => width <= 860
  const theme = useTheme()
  const colorMode = React.useContext(ColorModeContext)

  return (<AppBar position="relative" style={{
    height: '66px',
    borderColor: 'white',
    border: 'solid',
    borderWidth: '1px',
    justifyContent: 'center'
  }}>
    <Toolbar>
        <Link style={{ color: theme.palette.text.primary, flexGrow: isMobile() ? 1 : 0, fontFamily: '"Roboto","Helvetica","Arial",sans-serif', display: 'flex', alignItems: 'center' }} href="https://mavrykdynamics.com" underline="none">
          <img style={{ marginRight: '10px' }} src="img/mavryk-small-light.svg" alt="Marigold Logo" width="24" height="24"></img>

        <Typography style={{ marginRight: isMobile() ? '0px' : '24px' }} variant="h6" color="inherit" noWrap>
          MAVRYK {isMobile() && <span> SNAPSHOTS </span>}
          </Typography>
        </Link>

      {!isMobile() && <Separator></Separator> }

      {!isMobile() && <Box style={{
        paddingLeft: '10px', justifyContent: 'left'
      }} sx={{ flexGrow: 1 }}>
        <Typography style={{ color: theme.palette.text.primary, marginLeft: '25px' }} variant="h6" color="inherit" noWrap>
          MAVRYK SNAPSHOTS
        </Typography>
      </Box>
      }

      {!isMobile() &&
        // <span style={{ display: 'flex', alignItems: 'center' }}>
        //   <Separator></Separator>
        //   <SnapshotLink url="https://snapshots.api.mavryk.network">
        //     API
        //   </SnapshotLink>
        //   <Separator></Separator>
        //   <SnapshotLink url="https://snapshots.api.mavryk.network/mainnet/full">
        //     FULL MAINNET
        //   </SnapshotLink>
        //   <Separator></Separator>
        //   <SnapshotLink url="https://snapshots.api.mavryk.network/mainnet/rolling">
        //     ROLLING MAINNET
        //   </SnapshotLink>
        //   <Separator></Separator>
        //   <SnapshotLink url="https://snapshots.api.mavryk.network/basenet/full">
        //     FULL BASE
        //   </SnapshotLink>
        //   <Separator></Separator>
        //   <SnapshotLink url="https://snapshots.api.mavryk.network/basenet/rolling">
        //     ROLLING BASE
        //   </SnapshotLink>
        //   <Separator></Separator>
        // </span>}
        <span style={{ display: 'flex', alignItems: 'center' }}>
          <Separator></Separator>
          <SnapshotLink url="https://snapshots.api.mavryk.network">
            API
          </SnapshotLink>
          <Separator></Separator>
          <SnapshotLink url="https://snapshots.api.mavryk.network/basenet/archive">
            ARCHIVE ATLASNET
          </SnapshotLink>
          <Separator></Separator>
          <SnapshotLink url="https://snapshots.api.mavryk.network/basenet/full">
            FULL ATLASNET
          </SnapshotLink>
          <Separator></Separator>
          <SnapshotLink url="https://snapshots.api.mavryk.network/basenet/rolling">
            ROLLING ATLASNET
          </SnapshotLink>
          <Separator></Separator>
        </span>}

      <IconButton sx={{ ml: 1, marginLeft: '24px' }} onClick={colorMode.toggleColorMode} color="inherit">
        {theme.palette.mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
      </IconButton>
    </Toolbar>
  </AppBar>)
}
